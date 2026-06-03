package bpa

import (
	"encoding/json"

	"digit-oss/noc-services/internal/domain"
)

// Address represents a physical address.
type Address struct {
	TenantID        string               `json:"tenantId"`
	DoorNo          string               `json:"doorNo"`
	PlotNo          string               `json:"plotNo"`
	ID              string               `json:"id"`
	Landmark        string               `json:"landmark"`
	City            string               `json:"city"`
	District        string               `json:"district"`
	Region          string               `json:"region"`
	State           string               `json:"state"`
	Country         string               `json:"country"`
	Pincode         string               `json:"pincode"`
	AdditionDetails string               `json:"additionDetails"`
	BuildingName    string               `json:"buildingName"`
	Street          string               `json:"street"`
	Locality        *Boundary            `json:"locality"`
	GeoLocation     *GeoLocation         `json:"geoLocation"`
	AuditDetails    *domain.AuditDetails `json:"auditDetails"`
}

// Boundary represents a geographic/administrative boundary.
type Boundary struct {
	Code             string     `json:"code"`
	Name             string     `json:"name"`
	Label            string     `json:"label"`
	Latitude         string     `json:"latitude"`
	Longitude        string     `json:"longitude"`
	Children         []Boundary `json:"children"`
	MaterializedPath string     `json:"materializedPath"`
}

// GeoLocation represents geographic coordinates.
type GeoLocation struct {
	Latitude          *float64        `json:"latitude"`
	Longitude         *float64        `json:"longitude"`
	AdditionalDetails json.RawMessage `json:"additionalDetails"`
}
