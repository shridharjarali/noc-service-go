package validator

import (
	"encoding/json"
	"fmt"
	"strings"

	"digit-oss/noc-services/internal/config"
	"digit-oss/noc-services/internal/domain"
	"digit-oss/noc-services/internal/util"
)

// NOCValidator translates NOCValidator.java + MDMSValidator.java.
type NOCValidator struct {
	Cfg *config.Config
}

// ValidateCreate validates a create request against MDMS data and documents.
func (v *NOCValidator) ValidateCreate(nocReq *domain.NocRequest, mdmsData json.RawMessage) error {
	if err := v.validateMdmsData(mdmsData); err != nil {
		return err
	}
	noc := nocReq.Noc

	errorMap := make(map[string]string)

	if noc.TenantID == "" {
		errorMap["EG_NOC_TENANTID_MANDATORY"] = "TenantId is mandatory"
	}
	if noc.NocType == "" {
		errorMap["EG_NOC_NOCTYPE_MANDATORY"] = "NocType is mandatory"
	}
	if noc.Source == "" {
		errorMap["EG_NOC_SOURCE_MANDATORY"] = "Source is mandatory"
	}
	if noc.SourceRefID == "" {
		errorMap["EG_NOC_SOURCEREFID_MANDATORY"] = "SourceRefId is mandatory"
	}
	if noc.Workflow == nil || noc.Workflow.Action == "" {
		errorMap["EG_NOC_WORKFLOW_ACTION_MANDATORY"] = "Workflow action is mandatory"
	}

	if len(errorMap) > 0 {
		var msgs []string
		for _, msg := range errorMap {
			msgs = append(msgs, msg)
		}
		return util.NewCustomError("EG_NOC_ERR", strings.Join(msgs, ", "))
	}

	if len(noc.Documents) > 0 {
		if err := v.validateDuplicateDocuments(noc); err != nil {
			return err
		}
	}
	return nil
}

// ValidateUpdate validates an update request.
func (v *NOCValidator) ValidateUpdate(nocReq *domain.NocRequest, searchResult *domain.Noc, mode string, mdmsData json.RawMessage) error {
	noc := nocReq.Noc
	if err := v.validateMdmsData(mdmsData); err != nil {
		return err
	}

	errorMap := make(map[string]string)

	if noc.ID == "" {
		errorMap["UPDATE ERROR"] = fmt.Sprintf("Application Not found in the System %+v", noc)
	}

	if noc.Workflow != nil && noc.Workflow.Action != "" {
		action := strings.ToUpper(noc.Workflow.Action)
		if action == util.ActionReject && noc.Workflow.Comment == "" {
			errorMap["NOC_UPDATE_ERROR_COMMENT_REQUIRED"] = "Comment is mandatory, please provide the comments"
		}
	}

	if err := util.NewMultiError(errorMap); err != nil {
		return err
	}
	if err := v.validateDuplicateDocuments(noc); err != nil {
		return err
	}
	return nil
}

// GetOrValidateBusinessService fetches the business-service code for the
// NOC type from MDMS data and stores it in additionalDetails.
// Returns a map with keys "mode" and "workflowCode".
func (v *NOCValidator) GetOrValidateBusinessService(noc *domain.Noc, mdmsData json.RawMessage) (map[string]string, error) {
	var mdms map[string]interface{}
	if err := json.Unmarshal(mdmsData, &mdms); err != nil {
		return nil, util.NewCustomError(util.InvalidTenantIDMdmsKey, util.InvalidTenantIDMdmsMsg)
	}

	nocTypes, err := extractNocTypes(mdms)
	if err != nil || len(nocTypes) == 0 {
		return nil, util.NewCustomError("MDMS DATA ERROR", "Unable to fetch NocType from MDMS")
	}

	var matched map[string]interface{}
	for _, entry := range nocTypes {
		m, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		code, _ := m["code"].(string)
		if strings.EqualFold(code, noc.NocType) {
			matched = m
			break
		}
	}
	if matched == nil {
		return nil, util.NewCustomError("MDMS DATA ERROR",
			fmt.Sprintf("Unable to fetch %s workflow mode from MDMS", noc.NocType))
	}

	businessValues := make(map[string]string)
	modeVal, _ := matched[util.Mode].(string)
	businessValues[util.Mode] = modeVal

	if modeVal == util.OnlineMode {
		wf, _ := matched[util.OnlineWF].(string)
		businessValues[util.WorkflowCode] = wf
	} else {
		wf, _ := matched[util.OfflineWF].(string)
		businessValues[util.WorkflowCode] = wf
	}

	if noc.Workflow != nil && noc.Workflow.Action == util.ActionInitiate {
		businessValues[util.InitiatedTime] = fmt.Sprintf("%d", util.TimeNowMillis())
	}

	detailsJSON, _ := json.Marshal(businessValues)
	noc.AdditionalDetails = detailsJSON

	return businessValues, nil
}

// ──────────────────────────────────────────────────────────────────────────────
// Internal helpers
// ──────────────────────────────────────────────────────────────────────────────

func (v *NOCValidator) validateMdmsData(mdmsData json.RawMessage) error {
	var mdms map[string]interface{}
	if err := json.Unmarshal(mdmsData, &mdms); err != nil {
		return util.NewCustomError(util.InvalidTenantIDMdmsKey, util.InvalidTenantIDMdmsMsg)
	}
	nocTypes, err := extractNocTypes(mdms)
	if err != nil || len(nocTypes) == 0 {
		return util.NewCustomError("MDMS DATA ERROR ", "Unable to fetch NocType codes from MDMS")
	}
	return nil
}

func (v *NOCValidator) validateDuplicateDocuments(noc *domain.Noc) error {
	if len(noc.Documents) == 0 {
		return nil
	}
	seen := make(map[string]bool)
	for _, doc := range noc.Documents {
		if seen[doc.FileStoreID] {
			return util.NewCustomError("NOC_DUPLICATE_DOCUMENT", "Same document cannot be used multiple times")
		}
		seen[doc.FileStoreID] = true
	}
	return nil
}

func extractNocTypes(mdms map[string]interface{}) ([]interface{}, error) {
	mdmsRes, ok := mdms["MdmsRes"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing MdmsRes")
	}
	nocModule, ok := mdmsRes["NOC"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing NOC module")
	}
	nocTypeList, ok := nocModule["NocType"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("missing NocType")
	}
	return nocTypeList, nil
}
