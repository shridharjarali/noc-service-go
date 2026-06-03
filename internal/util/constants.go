package util

// ─── Module / Master names ───────────────────────────────────────────────────

const (
	SearchModule       = "rainmaker-nocsrv"
	NOCModule          = "NOC"
	NOCType            = "NocType"
	CommonMastersModule = "common-masters"
)

// ─── MDMS JSONPath codes ─────────────────────────────────────────────────────

const (
	NOCJsonPathCode           = "$.MdmsRes.NOC"
	NOCTypeJsonPathCode       = "$.MdmsRes.NOC.NocType"
	CommonMasterJsonPathCode  = "$.MdmsRes.common-masters"
)

// ─── Error constants ─────────────────────────────────────────────────────────

const (
	InvalidTenantIDMdmsKey = "INVALID TENANTID"
	InvalidTenantIDMdmsMsg = "No data found for this tenentID"
	ParsingError           = "PARSING_ERROR"
)

// ─── Status / State constants ────────────────────────────────────────────────

const (
	ApprovedState     = "APPROVED"
	AutoApprovedState = "AUTO_APPROVED"
	CreatedStatus     = "CREATED"
	VoidedStatus      = "VOIDED"
)

// ─── Workflow action constants ───────────────────────────────────────────────

const (
	ActionApprove     = "APPROVE"
	ActionAutoApprove = "AUTO_APPROVE"
	ActionReject      = "REJECT"
	ActionVoid        = "VOID"
	ActionInitiate    = "INITIATE"
)

// ─── MDMS mode / workflow code keys ──────────────────────────────────────────

const (
	Mode         = "mode"
	OnlineMode   = "online"
	OfflineMode  = "offline"
	OnlineWF     = "onlineWF"
	OfflineWF    = "offlineWF"
	WorkflowCode = "workflowCode"
)

// ─── Document type MDMS keys ─────────────────────────────────────────────────

const (
	NOCDocTypeMapping = "DocumentTypeMapping"
	DocumentType      = "DocumentType"
)

// ─── Enrichment keys ─────────────────────────────────────────────────────────

const (
	InitiatedTime = "SubmittedOn"
)

// ─── SMS notification action_status codes ────────────────────────────────────

const (
	ActionStatusCreated   = "null_CREATED"
	ActionStatusInitiated = "INITIATE_INPROGRESS"
	ActionStatusRejected  = "REJECT_REJECTED"
	ActionStatusApproved  = "APPROVE_APPROVED"
)

// ─── NOC type labels ─────────────────────────────────────────────────────────

const (
	FireNOCType    = "FIRE_NOC"
	AirportNOCType = "AIRPORT_AUTHORITY"
)
