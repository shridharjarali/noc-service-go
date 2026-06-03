package http

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"digit-oss/noc-services/internal/domain"
	"digit-oss/noc-services/internal/util"
)

// NOCHandler provides gin handlers for the NOC REST API.
// Translates NOCController.java — one handler per endpoint.
type NOCHandler struct {
	Service domain.NOCService
}

// ──────────────────────────────────────────────────────────────────────────────
// POST /v1/noc/_create
// ──────────────────────────────────────────────────────────────────────────────

func (h *NOCHandler) Create(c *gin.Context) {
	var nocReq domain.NocRequest
	if err := c.ShouldBindJSON(&nocReq); err != nil {
		writeErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST",
			"Failed to parse request body: "+err.Error(), nil)
		return
	}

	// Validate tenantId is not empty
	if nocReq.Noc.TenantID == "" {
		writeErrorResponse(c, http.StatusBadRequest, "MISSING_TENANT_ID",
			"tenantId is required in Noc object", nocReq.RequestInfo)
		return
	}

	noc, err := h.Service.Create(context.Background(), &nocReq)
	if err != nil {
		handleServiceError(c, err, nocReq.RequestInfo)
		return
	}

	resp := domain.NocResponse{
		ResponseInfo: createResponseInfo(nocReq.RequestInfo, true),
		Noc:          []domain.Noc{*noc},
	}
	c.JSON(http.StatusOK, resp)
}

// ──────────────────────────────────────────────────────────────────────────────
// POST /v1/noc/_update
// ──────────────────────────────────────────────────────────────────────────────

func (h *NOCHandler) Update(c *gin.Context) {
	var nocReq domain.NocRequest
	if err := c.ShouldBindJSON(&nocReq); err != nil {
		writeErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST",
			"Failed to parse request body: "+err.Error(), nil)
		return
	}

	// Validate tenantId is not empty
	if nocReq.Noc.TenantID == "" {
		writeErrorResponse(c, http.StatusBadRequest, "MISSING_TENANT_ID",
			"tenantId is required in Noc object", nocReq.RequestInfo)
		return
	}

	noc, err := h.Service.Update(context.Background(), &nocReq)
	if err != nil {
		handleServiceError(c, err, nocReq.RequestInfo)
		return
	}

	resp := domain.NocResponse{
		ResponseInfo: createResponseInfo(nocReq.RequestInfo, true),
		Noc:          []domain.Noc{*noc},
	}
	c.JSON(http.StatusOK, resp)
}

// ──────────────────────────────────────────────────────────────────────────────
// POST /v1/noc/_search
// Search takes RequestInfoWrapper in body + NocSearchCriteria from query params.
// ──────────────────────────────────────────────────────────────────────────────

func (h *NOCHandler) Search(c *gin.Context) {
	var wrapper domain.RequestInfoWrapper
	if err := c.ShouldBindJSON(&wrapper); err != nil {
		writeErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST",
			"Failed to parse request body: "+err.Error(), nil)
		return
	}

	criteria := parseCriteriaFromQuery(c)

	nocs, count, err := h.Service.Search(context.Background(), criteria, wrapper.RequestInfo)
	if err != nil {
		handleServiceError(c, err, wrapper.RequestInfo)
		return
	}

	resp := domain.NocResponse{
		ResponseInfo: createResponseInfo(wrapper.RequestInfo, true),
		Noc:          nocs,
		Count:        &count,
	}
	c.JSON(http.StatusOK, resp)
}

// ──────────────────────────────────────────────────────────────────────────────
// Helpers — ResponseInfo factory (mirrors ResponseInfoFactory.java)
// ──────────────────────────────────────────────────────────────────────────────

func createResponseInfo(ri *domain.RequestInfo, success bool) *domain.ResponseInfo {
	status := "successful"
	if !success {
		status = "failed"
	}

	apiId := ""
	ver := ""
	msgId := ""
	var ts *int64

	if ri != nil {
		apiId = ri.APIId
		ver = ri.Ver
		ts = ri.Ts
		msgId = ri.MsgId
	}

	return &domain.ResponseInfo{
		APIId:  apiId,
		Ver:    ver,
		Ts:     ts,
		ResMsg: "uief87324", // matches Java hard-coded value
		MsgId:  msgId,
		Status: status,
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// Error handling
// ──────────────────────────────────────────────────────────────────────────────

type errorResponse struct {
	ResponseInfo *domain.ResponseInfo `json:"ResponseInfo"`
	Errors       []errorDetail        `json:"Errors"`
}

type errorDetail struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func handleServiceError(c *gin.Context, err error, ri *domain.RequestInfo) {
	switch e := err.(type) {
	case *util.CustomError:
		writeErrorResponse(c, http.StatusBadRequest, e.Code, e.Message, ri)
	case *util.MultiError:
		for code, msg := range e.Errors {
			writeErrorResponse(c, http.StatusBadRequest, code, msg, ri)
			return
		}
	default:
		log.Printf("[NOCHandler] internal error: %v", err)
		writeErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), ri)
	}
}

func writeErrorResponse(c *gin.Context, statusCode int, code, message string, ri *domain.RequestInfo) {
	resp := errorResponse{
		ResponseInfo: createResponseInfo(ri, false),
		Errors: []errorDetail{
			{Code: code, Message: message, Description: message},
		},
	}
	c.JSON(statusCode, resp)
}

// ──────────────────────────────────────────────────────────────────────────────
// Query-param parsing for NocSearchCriteria (@ModelAttribute equivalent)
// ──────────────────────────────────────────────────────────────────────────────

func parseCriteriaFromQuery(c *gin.Context) domain.NocSearchCriteria {
	criteria := domain.NocSearchCriteria{
		TenantID:      c.Query("tenantId"),
		ApplicationNo: c.Query("applicationNo"),
		MobileNumber:  c.Query("mobileNumber"),
		NocNo:         c.Query("nocNo"),
		Source:        c.Query("source"),
		NocType:       c.Query("nocType"),
		SourceRefID:   c.Query("sourceRefId"),
	}

	if ids := c.Query("ids"); ids != "" {
		criteria.IDs = splitComma(ids)
	}
	if accountID := c.Query("accountId"); accountID != "" {
		criteria.AccountID = splitComma(accountID)
	}
	if status := c.Query("status"); status != "" {
		criteria.Status = splitComma(status)
	}

	if v := c.Query("offset"); v != "" {
		if n := parseInt(v); n >= 0 {
			criteria.Offset = &n
		}
	}
	if v := c.Query("limit"); v != "" {
		if n := parseInt(v); n > 0 {
			criteria.Limit = &n
		}
	}

	return criteria
}

func splitComma(s string) []string {
	result := make([]string, 0)
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			part := s[start:i]
			if len(part) > 0 {
				result = append(result, part)
			}
			start = i + 1
		}
	}
	return result
}

func parseInt(s string) int {
	n := 0
	neg := false
	i := 0
	if len(s) > 0 && s[0] == '-' {
		neg = true
		i = 1
	}
	for ; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return 0
		}
		n = n*10 + int(s[i]-'0')
	}
	if neg {
		n = -n
	}
	return n
}
