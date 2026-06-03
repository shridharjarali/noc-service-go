package postgres

import (
	"context"
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"digit-oss/noc-services/internal/config"
	"digit-oss/noc-services/internal/domain"
	"digit-oss/noc-services/internal/repository/models"
)

// NocRepository implements domain.NOCRepository using the Kafka persister
// pattern for writes (Save/Update) and GORM for reads (Search).
type NocRepository struct {
	DB       *gorm.DB
	Cfg      *config.Config
	Producer domain.Producer
}

// Compile-time interface check.
var _ domain.NOCRepository = (*NocRepository)(nil)

// ──────────────────────────────────────────────────────────────────────────────
// Save pushes the NocRequest to the Kafka save topic.
// The egov-persister service listens on this topic and performs the actual
// SQL INSERT into eg_noc and eg_noc_document.
// ──────────────────────────────────────────────────────────────────────────────

func (r *NocRepository) Save(_ context.Context, req *domain.NocRequest) error {
	log.Printf("[NocRepository.Save] pushing to topic=%s applicationNo=%s",
		r.Cfg.SaveTopic, req.Noc.ApplicationNo)
	return r.Producer.Push(r.Cfg.SaveTopic, req)
}

// ──────────────────────────────────────────────────────────────────────────────
// Update pushes the NocRequest to the appropriate Kafka update topic.
// If isStateUpdatable is true → update topic (full record update).
// If isStateUpdatable is false → update-workflow topic (workflow-only update).
// ──────────────────────────────────────────────────────────────────────────────

func (r *NocRepository) Update(_ context.Context, req *domain.NocRequest, isStateUpdatable bool) error {
	log.Printf("[NocRepository.Update] applicationStatus=%s isStateUpdatable=%v",
		req.Noc.ApplicationStatus, isStateUpdatable)
	if isStateUpdatable {
		return r.Producer.Push(r.Cfg.UpdateTopic, req)
	}
	return r.Producer.Push(r.Cfg.UpdateWorkflowTopic, req)
}

// ──────────────────────────────────────────────────────────────────────────────
// Search finds NOC applications matching the given criteria using GORM.
// Documents are eager-loaded via Preload. Returns the list and a total count.
// ──────────────────────────────────────────────────────────────────────────────

func (r *NocRepository) Search(ctx context.Context, criteria domain.NocSearchCriteria) ([]domain.Noc, int, error) {
	db := r.DB.WithContext(ctx)

	// 1. Count query (distinct NOC IDs) — uses its own query chain.
	var count int64
	countQuery := r.applyFilters(db, criteria)
	if err := countQuery.Distinct("id").Count(&count).Error; err != nil {
		return nil, 0, fmt.Errorf("count query: %w", err)
	}

	// 2. Data query (paginated) with eager-loaded documents — separate chain.
	dataQuery := r.applyFilters(db, criteria)

	limit := r.Cfg.DefaultLimit
	offset := r.Cfg.DefaultOffset

	if criteria.Limit != nil {
		if *criteria.Limit <= r.Cfg.MaxSearchLimit {
			limit = *criteria.Limit
		} else {
			limit = r.Cfg.MaxSearchLimit
		}
	}
	if criteria.Offset != nil {
		offset = *criteria.Offset
	}

	// Order by lastmodifiedtime DESC to match original DENSE_RANK ordering.
	dataQuery = dataQuery.Order("lastmodifiedtime DESC").
		Preload("Documents")

	if limit != -1 {
		dataQuery = dataQuery.Offset(offset).Limit(limit)
	}

	var nocModels []models.NocModel
	if err := dataQuery.Find(&nocModels).Error; err != nil {
		return nil, 0, fmt.Errorf("search query: %w", err)
	}

	// 3. Convert GORM models to domain objects.
	nocList := make([]domain.Noc, 0, len(nocModels))
	for i := range nocModels {
		nocList = append(nocList, nocModels[i].ToDomain())
	}

	log.Printf("[NocRepository.Search] found %d results (total %d)", len(nocList), count)
	return nocList, int(count), nil
}

// ──────────────────────────────────────────────────────────────────────────────
// applyFilters chains GORM Where clauses for each non-empty search criterion.
// ──────────────────────────────────────────────────────────────────────────────

func (r *NocRepository) applyFilters(db *gorm.DB, c domain.NocSearchCriteria) *gorm.DB {
	query := db.Model(&models.NocModel{})

	// ── tenantId ─────────────────────────────────────────────────────────
	if c.TenantID != "" {
		query = query.Where("tenantid = ?", c.TenantID)
	}

	// ── ids ──────────────────────────────────────────────────────────────
	if len(c.IDs) > 0 {
		query = query.Where("id IN ?", c.IDs)
	}

	// ── applicationNo (supports comma-separated + optional fuzzy) ───────
	if c.ApplicationNo != "" {
		appNos := splitTrimmed(c.ApplicationNo)
		if r.Cfg.IsFuzzyEnabled {
			query = query.Where("applicationno LIKE ANY(?)", pgStringArray(wrapLike(appNos)))
		} else {
			query = query.Where("applicationno IN ?", appNos)
		}
	}

	// ── nocNo (supports comma-separated + optional fuzzy) ───────────────
	if c.NocNo != "" {
		nocNos := splitTrimmed(c.NocNo)
		if r.Cfg.IsFuzzyEnabled {
			query = query.Where("nocno LIKE ANY(?)", pgStringArray(wrapLike(nocNos)))
		} else {
			query = query.Where("nocno IN ?", nocNos)
		}
	}

	// ── source ───────────────────────────────────────────────────────────
	if c.Source != "" {
		query = query.Where("source = ?", c.Source)
	}

	// ── sourceRefId (supports comma-separated + optional fuzzy) ──────────
	if c.SourceRefID != "" {
		cleaned := strings.NewReplacer("[", "", "]", "").Replace(c.SourceRefID)
		refs := splitTrimmed(cleaned)
		if r.Cfg.IsFuzzyEnabled {
			query = query.Where("sourcerefid LIKE ANY(?)", pgStringArray(wrapLike(refs)))
		} else {
			query = query.Where("sourcerefid IN ?", refs)
		}
	}

	// ── nocType (comma-separated, exact match only) ─────────────────────
	if c.NocType != "" {
		nocTypes := splitTrimmed(c.NocType)
		query = query.Where("noctype IN ?", nocTypes)
	}

	// ── status (list) ───────────────────────────────────────────────────
	if len(c.Status) > 0 {
		query = query.Where("applicationstatus IN ?", c.Status)
	}

	return query
}

// ──────────────────────────────────────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────────────────────────────────────

// splitTrimmed splits s on "," and trims whitespace from each element.
func splitTrimmed(s string) []string {
	raw := strings.Split(s, ",")
	out := make([]string, 0, len(raw))
	for _, v := range raw {
		v = strings.TrimSpace(v)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

// wrapLike wraps each string with % for SQL LIKE fuzzy matching.
func wrapLike(vals []string) []string {
	out := make([]string, len(vals))
	for i, v := range vals {
		out[i] = "%" + v + "%"
	}
	return out
}

// pgStringArray builds a PostgreSQL ARRAY literal as a gorm.Expr for use
// with LIKE ANY(ARRAY[...]) clauses.
func pgStringArray(vals []string) clause.Expr {
	quoted := make([]string, len(vals))
	for i, v := range vals {
		escaped := strings.ReplaceAll(v, "'", "''")
		quoted[i] = "'" + escaped + "'"
	}
	return gorm.Expr("ARRAY[" + strings.Join(quoted, ",") + "]")
}
