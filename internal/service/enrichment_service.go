package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"digit-oss/noc-services/internal/config"
	"digit-oss/noc-services/internal/domain"
	"digit-oss/noc-services/internal/repository/postgres"
	"digit-oss/noc-services/internal/util"
)

// EnrichmentService translates EnrichmentService.java.
type EnrichmentService struct {
	Cfg       *config.Config
	IdGenRepo *postgres.IdGenRepository
}

// EnrichCreateRequest populates IDs, audit details, and application number.
func (e *EnrichmentService) EnrichCreateRequest(nocReq *domain.NocRequest) error {
	// Ensure RequestInfo and UserInfo are present to avoid nil deref
	ri := nocReq.RequestInfo
	var userUUID string
	if ri == nil {
		ri = &domain.RequestInfo{}
		nocReq.RequestInfo = ri
	}
	if ri.UserInfo != nil {
		userUUID = ri.UserInfo.UUID
	}
	audit := util.GetAuditDetails(userUUID, true)
	noc := nocReq.Noc

	noc.AuditDetails = audit
	noc.ID = util.GenerateUUID()
	noc.AccountID = audit.CreatedBy

	// Generate application number via idgen
	ids, err := e.IdGenRepo.GetID(ri, noc.TenantID, e.Cfg.ApplicationNoIdgenName, 1)
	if err != nil {
		return util.NewCustomError("IDGEN ERROR", "No ids returned from idgen Service: "+err.Error())
	}
	if len(ids) == 0 {
		return util.NewCustomError("IDGEN ERROR", "No ids returned from idgen Service")
	}
	noc.ApplicationNo = ids[0]

	// Assign UUIDs to documents missing IDs
	for i := range noc.Documents {
		if noc.Documents[i].ID == "" {
			noc.Documents[i].ID = util.GenerateUUID()
		}
	}

	return nil
}

// EnrichNocUpdateRequest populates audit details and document IDs for update.
func (e *EnrichmentService) EnrichNocUpdateRequest(nocReq *domain.NocRequest, searchResult *domain.Noc) {
	ri := nocReq.RequestInfo
	var userUUID string
	if ri == nil {
		ri = &domain.RequestInfo{}
		nocReq.RequestInfo = ri
	}
	if ri.UserInfo != nil {
		userUUID = ri.UserInfo.UUID
	}
	audit := util.GetAuditDetails(userUUID, false)
	noc := nocReq.Noc

	noc.AuditDetails = audit
	// Preserve original createdBy / createdTime from search result
	if searchResult.AuditDetails != nil {
		noc.AuditDetails.CreatedBy = searchResult.AuditDetails.CreatedBy
		noc.AuditDetails.CreatedTime = searchResult.AuditDetails.CreatedTime
	}
	noc.ApplicationNo = searchResult.ApplicationNo

	// Assign UUIDs to new documents
	for i := range noc.Documents {
		if noc.Documents[i].ID == "" {
			noc.Documents[i].ID = util.GenerateUUID()
		}
	}
	// Workflow documents
	if noc.Workflow != nil {
		for i := range noc.Workflow.Documents {
			if noc.Workflow.Documents[i].ID == "" {
				noc.Workflow.Documents[i].ID = util.GenerateUUID()
			}
		}
	}
}

// PostStatusEnrichment enriches after workflow transition:
// generates nocNo on APPROVED/AUTO_APPROVED, sets INACTIVE on VOIDED,
// stamps initiatedTime on INITIATE.
func (e *EnrichmentService) PostStatusEnrichment(nocReq *domain.NocRequest, currentState string) error {
	noc := nocReq.Noc

	state := strings.ToUpper(currentState)

	if state == util.ApprovedState || state == util.AutoApprovedState {
		ids, err := e.IdGenRepo.GetID(nocReq.RequestInfo, noc.TenantID, e.Cfg.ApplicationNoIdgenName, 1)
		if err != nil || len(ids) == 0 {
			return util.NewCustomError("IDGEN ERROR", "Failed to generate NOC number")
		}
		noc.NocNo = ids[0]
	}

	if state == util.VoidedStatus {
		noc.Status = "INACTIVE"
	}

	if noc.Workflow != nil && noc.Workflow.Action == util.ActionInitiate {
		details := make(map[string]interface{})
		if len(noc.AdditionalDetails) > 0 {
			_ = json.Unmarshal(noc.AdditionalDetails, &details)
		}
		details[util.InitiatedTime] = fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
		raw, _ := json.Marshal(details)
		noc.AdditionalDetails = raw
	}

	return nil
}
