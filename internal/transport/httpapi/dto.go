package httpapi

// CreateRiskRequest represents the payload to create a new IMS risk.
// swagger:model CreateRiskRequest
type CreateRiskRequest struct {
	Title       string `json:"title"`       // Short name of the risk
	Process     string `json:"process"`     // Process where the risk occurs
	Domain      string `json:"domain"`      // Domain: quality|environment|ohs|isms
	Description string `json:"description"` // Detailed risk description
	Likelihood  int    `json:"likelihood"`  // 1 (rare) - 5 (almost certain)
	Impact      int    `json:"impact"`      // 1 (minor) - 5 (catastrophic)
	Owner       string `json:"owner"`       // Responsible person or role
}

// UpdateRiskStatusRequest represents payload to update risk status.
// swagger:model UpdateRiskStatusRequest
type UpdateRiskStatusRequest struct {
	Status string `json:"status"` // Open, Accepted, Mitigated
}

// CreateIncidentRequest represents payload to create an incident.
// swagger:model CreateIncidentRequest
type CreateIncidentRequest struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	Domain        string `json:"domain"`        // quality|environment|ohs|isms
	RelatedRiskID *int   `json:"relatedRiskId"` // Optional link to risk
	Severity      int    `json:"severity"`      // 1-5
	Likelihood    int    `json:"likelihood"`    // 1-5
}

// UpdateIncidentRequest represents payload to update an incident.
// swagger:model UpdateIncidentRequest
type UpdateIncidentRequest struct {
	RootCause *string `json:"rootCause"`
	Status    *string `json:"status"` // Open, Investigation, Closed
}

// CreateAuditRequest represents payload to create an audit.
// swagger:model CreateAuditRequest
type CreateAuditRequest struct {
	Title       string `json:"title"`
	Scope       string `json:"scope"`
	Domain      string `json:"domain"`      // quality|environment|ohs|isms
	PlannedDate string `json:"plannedDate"` // YYYY-MM-DD
	Auditor     string `json:"auditor"`
}

// UpdateAuditRequest represents payload to update an audit.
// swagger:model UpdateAuditRequest
type UpdateAuditRequest struct {
	Status   *string `json:"status"`   // Planned, In Progress, Completed
	Findings *string `json:"findings"` // Summary of audit findings
}

// CreateActionRequest represents payload to create a CAPA action.
// swagger:model CreateActionRequest
type CreateActionRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	SourceType  string `json:"sourceType"` // risk|incident|audit
	SourceID    int    `json:"sourceId"`
	Owner       string `json:"owner"`
	DueDate     string `json:"dueDate"` // YYYY-MM-DD
}

// UpdateActionRequest represents payload to update an action.
// swagger:model UpdateActionRequest
type UpdateActionRequest struct {
	Status  *string `json:"status"`  // Open, In Progress, Done, Overdue
	DueDate *string `json:"dueDate"` // Optional new due date
}
