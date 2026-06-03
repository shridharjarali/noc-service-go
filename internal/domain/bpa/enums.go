package bpa

// BPAStatus represents the active/inactive status for BPA entities.
type BPAStatus string

const (
	BPAStatusActive   BPAStatus = "ACTIVE"
	BPAStatusInactive BPAStatus = "INACTIVE"
)

// Channel represents the channel through which a property was created.
type Channel string

const (
	ChannelSystem     Channel = "SYSTEM"
	ChannelCFCCounter Channel = "CFC_COUNTER"
	ChannelCitizen    Channel = "CITIZEN"
	ChannelDataEntry  Channel = "DATA_ENTRY"
	ChannelMigration  Channel = "MIGRATION"
)

// Source represents the source of construction-detail data.
type Source string

const (
	SourceMunicipalRecords Source = "MUNICIPAL_RECORDS"
	SourceFieldSurvey      Source = "FIELD_SURVEY"
	SourceWebApp           Source = "WEBAPP"
)

// OccupancyType represents owner vs tenant occupancy.
type OccupancyType string

const (
	OccupancyTypeOwner  OccupancyType = "OWNER"
	OccupancyTypeTenant OccupancyType = "TENANT"
)

// Relationship represents father/husband relationship.
type Relationship string

const (
	RelationshipFather  Relationship = "FATHER"
	RelationshipHusband Relationship = "HUSBAND"
)
