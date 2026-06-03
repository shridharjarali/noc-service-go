package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"digit-oss/noc-services/internal/config"
	"digit-oss/noc-services/internal/domain"
	"digit-oss/noc-services/internal/repository/postgres"
	"digit-oss/noc-services/internal/util"
	"digit-oss/noc-services/internal/validator"
	"digit-oss/noc-services/internal/workflow"
)

// NOCServiceImpl implements domain.NOCService.
type NOCServiceImpl struct {
	Cfg          *config.Config
	Repo         domain.NOCRepository
	Enrichment   *EnrichmentService
	Validator    *validator.NOCValidator
	WfIntegrator *wfintegrator.WorkflowIntegrator
	WfService    *wfintegrator.WorkflowService
	SvcRequest   *postgres.ServiceRequestRepository
}

// Compile-time interface check
var _ domain.NOCService = (*NOCServiceImpl)(nil)

// Create validates and persists a new NOC application.
func (s *NOCServiceImpl) Create(ctx context.Context, nocReq *domain.NocRequest) (*domain.Noc, error) {
	noc := nocReq.Noc
	tenantID := rootTenant(noc.TenantID)

	mdmsData, err := s.mdmsCall(nocReq.RequestInfo, tenantID)
	if err != nil {
		return nil, err
	}

	additionalDetails, err := s.Validator.GetOrValidateBusinessService(noc, mdmsData)
	if err != nil {
		return nil, err
	}

	if err := s.Validator.ValidateCreate(nocReq, mdmsData); err != nil {
		return nil, err
	}

	if err := s.Enrichment.EnrichCreateRequest(nocReq); err != nil {
		return nil, err
	}

	// Workflow transition (if workflow action is provided)
	if noc.Workflow != nil && noc.Workflow.Action != "" {
		if err := s.WfIntegrator.CallWorkFlow(nocReq, additionalDetails[util.WorkflowCode]); err != nil {
			return nil, err
		}
	} else {
		noc.ApplicationStatus = util.CreatedStatus
	}

	// Push to Kafka save topic (persister pattern — same as Java)
	if err := s.Repo.Save(ctx, nocReq); err != nil {
		return nil, err
	}

	return noc, nil
}

// Update validates and applies changes to an existing NOC application.
func (s *NOCServiceImpl) Update(ctx context.Context, nocReq *domain.NocRequest) (*domain.Noc, error) {
	noc := nocReq.Noc
	tenantID := rootTenant(noc.TenantID)

	mdmsData, err := s.mdmsCall(nocReq.RequestInfo, tenantID)
	if err != nil {
		return nil, err
	}

	// Get workflow code from additionalDetails or MDMS
	var additionalDetails map[string]string
	if len(noc.AdditionalDetails) > 0 {
		_ = json.Unmarshal(noc.AdditionalDetails, &additionalDetails)
	}
	if len(additionalDetails) == 0 {
		additionalDetails, err = s.Validator.GetOrValidateBusinessService(noc, mdmsData)
		if err != nil {
			return nil, err
		}
	}

	// Fetch existing record
	searchResult, err := s.getNocForUpdate(ctx, nocReq)
	if err != nil {
		return nil, err
	}

	// Block AUTO_APPROVED → INPROGRESS
	if strings.EqualFold(searchResult.ApplicationStatus, "AUTO_APPROVED") &&
		strings.EqualFold(noc.ApplicationStatus, "INPROGRESS") {
		return nil, util.NewCustomError("AutoApproveException",
			"NOC_UPDATE_ERROR_AUTO_APPROVED_TO_INPROGRESS_NOTALLOWED")
	}

	if err := s.Validator.ValidateUpdate(nocReq, searchResult, additionalDetails[util.Mode], mdmsData); err != nil {
		return nil, err
	}

	s.Enrichment.EnrichNocUpdateRequest(nocReq, searchResult)

	if noc.Workflow != nil && noc.Workflow.Action != "" {
		wfCode := additionalDetails[util.WorkflowCode]

		if err := s.WfIntegrator.CallWorkFlow(nocReq, wfCode); err != nil {
			return nil, err
		}

		s.Enrichment.PostStatusEnrichment(nocReq, "")

		// Get business service to determine isStateUpdatable
		bs, err := s.WfService.GetBusinessService(noc, nocReq.RequestInfo, wfCode)
		if err != nil {
			return nil, err
		}

		if bs == nil {
			// No business service found — push with isStateUpdatable=true
			// (matches Java: nocRepository.update(nocRequest, true))
			if err := s.Repo.Update(ctx, nocReq, true); err != nil {
				return nil, err
			}
		} else {
			// Post-status enrichment with current state
			stateObj := s.WfService.GetCurrentState(noc.ApplicationStatus, bs)
			if stateObj != nil {
				if err := s.Enrichment.PostStatusEnrichment(nocReq, stateObj.State); err != nil {
					return nil, err
				}
			}
			isUpdatable := s.WfService.IsStateUpdatable(noc.ApplicationStatus, bs)
			if err := s.Repo.Update(ctx, nocReq, isUpdatable); err != nil {
				return nil, err
			}
		}
	} else {
		// No workflow action — push with isStateUpdatable=false
		// (matches Java: nocRepository.update(nocRequest, Boolean.FALSE))
		if err := s.Repo.Update(ctx, nocReq, false); err != nil {
			return nil, err
		}
	}

	return noc, nil
}

// Search retrieves NOC applications matching the criteria.
func (s *NOCServiceImpl) Search(ctx context.Context, criteria domain.NocSearchCriteria, requestInfo *domain.RequestInfo) ([]domain.Noc, int, error) {
	log.Printf("[NOCService.Search] criteria=%+v", criteria)
	nocs, count, err := s.Repo.Search(ctx, criteria)
	if err != nil {
		return nil, 0, err
	}
	return nocs, count, nil
}

// getNocForUpdate fetches the existing NOC by ID for validation before update.
func (s *NOCServiceImpl) getNocForUpdate(ctx context.Context, nocReq *domain.NocRequest) (*domain.Noc, error) {
	criteria := domain.NocSearchCriteria{
		IDs: []string{nocReq.Noc.ID},
	}
	nocs, _, err := s.Search(ctx, criteria, nocReq.RequestInfo)
	if err != nil {
		return nil, err
	}
	if len(nocs) == 0 {
		return nil, util.NewCustomError("INVALID_NOC_SEARCH",
			fmt.Sprintf("Noc Application not found for: %s :ID", nocReq.Noc.ID))
	}
	if len(nocs) > 1 {
		return nil, util.NewCustomError("INVALID_NOC_SEARCH",
			fmt.Sprintf("Multiple Noc Application(s) found for: %s :ID", nocReq.Noc.ID))
	}
	return &nocs[0], nil
}

// mdmsCall fetches MDMS data for the given tenant.
func (s *NOCServiceImpl) mdmsCall(requestInfo *domain.RequestInfo, tenantID string) (json.RawMessage, error) {
	url := s.Cfg.MdmsHost + s.Cfg.MdmsEndPoint

	mdmsReq := map[string]interface{}{
		"RequestInfo": requestInfo,
		"MdmsCriteria": map[string]interface{}{
			"tenantId": tenantID,
			"moduleDetails": []map[string]interface{}{
				{
					"moduleName": util.NOCModule,
					"masterDetails": []map[string]interface{}{
						{"name": util.NOCType, "filter": "$.[?(@.isActive==true)]"},
						{"name": util.NOCDocTypeMapping},
					},
				},
				{
					"moduleName": util.CommonMastersModule,
					"masterDetails": []map[string]interface{}{
						{"name": util.DocumentType, "filter": "$.[?(@.active==true)]"},
					},
				},
			},
		},
	}

	var result json.RawMessage
	if err := s.SvcRequest.FetchResult(url, mdmsReq, &result); err != nil {
		return nil, util.NewCustomError(util.InvalidTenantIDMdmsKey, util.InvalidTenantIDMdmsMsg)
	}
	return result, nil
}

// rootTenant extracts "pb" from "pb.amritsar".
func rootTenant(tenantID string) string {
	parts := strings.SplitN(tenantID, ".", 2)
	return parts[0]
}
