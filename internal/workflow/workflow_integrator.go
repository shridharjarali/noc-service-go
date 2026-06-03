package wfintegrator

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"digit-oss/noc-services/internal/config"
	"digit-oss/noc-services/internal/domain"
	"digit-oss/noc-services/internal/repository/postgres"
	"digit-oss/noc-services/internal/util"
)

// WorkflowIntegrator calls the workflow transition endpoint and
// updates the NOC's applicationStatus from the workflow response.
type WorkflowIntegrator struct {
	Cfg        *config.Config
	SvcRequest *postgres.ServiceRequestRepository
}

// CallWorkFlow posts a process-instance transition to the workflow service.
// On success it sets noc.ApplicationStatus from the workflow response.
func (w *WorkflowIntegrator) CallWorkFlow(nocReq *domain.NocRequest, businessServiceValue string) error {
	noc := nocReq.Noc

	// Build the ProcessInstance array as generic JSON.
	pi := map[string]interface{}{
		"businessId":      noc.ApplicationNo,
		"tenantId":        noc.TenantID,
		"businessService": businessServiceValue,
		"moduleName":      util.NOCModule,
		"action":          noc.Workflow.Action,
	}
	if noc.Workflow.Comment != "" {
		pi["comment"] = noc.Workflow.Comment
	}
	if len(noc.Workflow.Assignes) > 0 {
		uuidMaps := make([]map[string]string, 0, len(noc.Workflow.Assignes))
		for _, a := range noc.Workflow.Assignes {
			uuidMaps = append(uuidMaps, map[string]string{"uuid": a})
		}
		pi["assignes"] = uuidMaps
	}
	if len(noc.Workflow.Documents) > 0 {
		pi["documents"] = noc.Workflow.Documents
	}

	body := map[string]interface{}{
		"RequestInfo":      nocReq.RequestInfo,
		"ProcessInstances": []interface{}{pi},
	}

	url := w.Cfg.WfHost + w.Cfg.WfTransitionPath
	log.Printf("[WorkflowIntegrator] POST %s", url)

	var rawResp json.RawMessage
	if err := w.SvcRequest.FetchResult(url, body, &rawResp); err != nil {
		return util.NewCustomError("EG_WF_ERROR", "Exception while integrating with workflow: "+err.Error())
	}

	// Parse response to extract status
	var respObj struct {
		ProcessInstances []struct {
			BusinessID string `json:"businessId"`
			State      struct {
				ApplicationStatus string `json:"applicationStatus"`
			} `json:"state"`
		} `json:"ProcessInstances"`
	}
	if err := json.Unmarshal(rawResp, &respObj); err != nil {
		return util.NewCustomError("EG_WF_ERROR", "Failed to parse workflow response: "+err.Error())
	}

	for _, inst := range respObj.ProcessInstances {
		if strings.EqualFold(inst.BusinessID, noc.ApplicationNo) {
			noc.ApplicationStatus = inst.State.ApplicationStatus
			break
		}
	}

	return nil
}

// WorkflowService fetches business-service definitions and inspects states.
type WorkflowService struct {
	Cfg        *config.Config
	SvcRequest *postgres.ServiceRequestRepository
}

// GetBusinessService fetches the BusinessService definition from the workflow engine.
// Returns nil if no business service is found (matches Java null-return behaviour).
func (ws *WorkflowService) GetBusinessService(noc *domain.Noc, requestInfo *domain.RequestInfo, businessServiceValue string) (*BusinessServiceDef, error) {
	url := fmt.Sprintf("%s%s?tenantId=%s&businessServices=%s",
		ws.Cfg.WfHost, ws.Cfg.WfBusinessServiceSearchPath,
		noc.TenantID, businessServiceValue)

	reqBody := map[string]interface{}{
		"RequestInfo": requestInfo,
	}

	var resp BusinessServiceSearchResponse
	if err := ws.SvcRequest.FetchResult(url, reqBody, &resp); err != nil {
		return nil, util.NewCustomError("PARSING ERROR", "Failed to parse response: "+err.Error())
	}

	if len(resp.BusinessServices) == 0 {
		return nil, nil
	}
	return &resp.BusinessServices[0], nil
}

// GetCurrentState returns the State whose ApplicationStatus matches.
func (ws *WorkflowService) GetCurrentState(status string, bs *BusinessServiceDef) *StateDef {
	for _, s := range bs.States {
		if strings.EqualFold(s.ApplicationStatus, status) {
			return &s
		}
	}
	return nil
}

// IsStateUpdatable returns whether the state for the given status is updatable.
func (ws *WorkflowService) IsStateUpdatable(status string, bs *BusinessServiceDef) bool {
	for _, s := range bs.States {
		if strings.EqualFold(s.ApplicationStatus, status) {
			if s.IsStateUpdatable != nil {
				return *s.IsStateUpdatable
			}
			return false
		}
	}
	return false
}

// ──────────────────────────────────────────────────────────────────────────────
// Local types for parsing workflow responses (not exported to domain to avoid
// circular deps — these are API-response shapes, not domain entities).
// ──────────────────────────────────────────────────────────────────────────────

// BusinessServiceSearchResponse is the shape of the workflow businessService search response.
type BusinessServiceSearchResponse struct {
	BusinessServices []BusinessServiceDef `json:"BusinessServices"`
}

// BusinessServiceDef is a lightweight representation of a workflow BusinessService.
type BusinessServiceDef struct {
	TenantID        string     `json:"tenantId"`
	BusinessService string     `json:"businessService"`
	States          []StateDef `json:"states"`
}

// StateDef is a lightweight representation of a workflow State.
type StateDef struct {
	State             string `json:"state"`
	ApplicationStatus string `json:"applicationStatus"`
	IsStateUpdatable  *bool  `json:"isStateUpdatable"`
}
