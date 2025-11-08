package domain

import (
	"errors"
	"strings"
)

// --------- Errors & helpers ---------

var ErrInvalidDomain = errors.New("invalid domain")

func normalize(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

// Domain represents the IMS domain a record belongs to.
type Domain string

const (
	DomainQuality Domain = "Quality"              // ISO 9001
	DomainEnv     Domain = "Environment"          // ISO 14001
	DomainOHS     Domain = "OHS"                  // ISO 45001
	DomainISMS    Domain = "Information Security" // ISMS / ISO 27001
)

// ParseDomain maps string codes to Domain constants.
func ParseDomain(s string) (Domain, error) {
	switch normalize(s) {
	case "quality", "qms":
		return DomainQuality, nil
	case "environment", "ems":
		return DomainEnv, nil
	case "ohs", "safety", "ohsms":
		return DomainOHS, nil
	case "isms", "information security", "security":
		return DomainISMS, nil
	default:
		return "", ErrInvalidDomain
	}
}

// RiskLevelFromScore returns Low/Medium/High from numeric score.
func RiskLevelFromScore(score int) string {
	switch {
	case score >= 16:
		return "High"
	case score >= 8:
		return "Medium"
	default:
		return "Low"
	}
}

// --------- Core IMS models ---------

// Risk represents a risk in the integrated management system.
// swagger:model Risk
type Risk struct {
	ID          int    `json:"id"`          // Auto-generated risk ID
	Title       string `json:"title"`       // Short risk title
	Process     string `json:"process"`     // Process where risk occurs
	Domain      Domain `json:"domain"`      // IMS Domain (Quality/Environment/OHS/Information Security)
	Description string `json:"description"` // Detailed description
	Likelihood  int    `json:"likelihood"`  // 1-5
	Impact      int    `json:"impact"`      // 1-5
	Score       int    `json:"score"`       // Likelihood * Impact
	Level       string `json:"level"`       // Low/Medium/High
	Owner       string `json:"owner"`       // Responsible person / role
	Status      string `json:"status"`      // Open, Accepted, Mitigated
	CreatedAt   string `json:"createdAt"`   // RFC3339 timestamp
}

// Incident represents an incident / nonconformity.
// swagger:model Incident
type Incident struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Domain        Domain `json:"domain"`
	RelatedRiskID *int   `json:"relatedRiskId,omitempty"`
	Severity      int    `json:"severity"`   // 1-5
	Likelihood    int    `json:"likelihood"` // 1-5
	RiskScore     int    `json:"riskScore"`
	RiskLevel     string `json:"riskLevel"` // Low/Medium/High
	RootCause     string `json:"rootCause"`
	Status        string `json:"status"`    // Open, Investigation, Closed
	CreatedAt     string `json:"createdAt"` // RFC3339
	UpdatedAt     string `json:"updatedAt"` // RFC3339
}

// Audit represents an internal IMS audit.
// swagger:model Audit
type Audit struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Scope       string `json:"scope"`
	Domain      Domain `json:"domain"`      // Main focus area
	PlannedDate string `json:"plannedDate"` // YYYY-MM-DD
	Auditor     string `json:"auditor"`
	Status      string `json:"status"`   // Planned, In Progress, Completed
	Findings    string `json:"findings"` // Text field
	CreatedAt   string `json:"createdAt"`
}

// Action represents a corrective / preventive action (CAPA).
// swagger:model Action
type Action struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	SourceType  string `json:"sourceType"` // Risk, Incident, Audit
	SourceID    int    `json:"sourceId"`
	Owner       string `json:"owner"`
	DueDate     string `json:"dueDate"` // YYYY-MM-DD
	Status      string `json:"status"`  // Open, In Progress, Done, Overdue
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// Dashboard aggregates KPIs for IMS.
// swagger:model Dashboard
type Dashboard struct {
	TotalRisks        int            `json:"totalRisks"`
	HighRisks         int            `json:"highRisks"`
	TotalIncidents    int            `json:"totalIncidents"`
	OpenIncidents     int            `json:"openIncidents"`
	ActionsByStatus   map[string]int `json:"actionsByStatus"`
	IncidentsByDomain map[Domain]int `json:"incidentsByDomain"`
}
