package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/xenakil/integraflow-ims/internal/domain"
	"github.com/xenakil/integraflow-ims/internal/repository"
)

type AuditService struct {
	repo repository.AuditRepository
}

func NewAuditService(repo repository.AuditRepository) *AuditService {
	return &AuditService{repo: repo}
}

type CreateAuditInput struct {
	Title       string
	Scope       string
	Domain      string
	PlannedDate string
	Auditor     string
}

func (s *AuditService) CreateAudit(in CreateAuditInput) (*domain.Audit, error) {
	if strings.TrimSpace(in.Title) == "" || strings.TrimSpace(in.Scope) == "" || strings.TrimSpace(in.Domain) == "" {
		return nil, fmt.Errorf("%w: title, scope and domain are required", ErrValidation)
	}

	dom, err := domain.ParseDomain(in.Domain)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidation, err)
	}

	audit := &domain.Audit{
		Title:       in.Title,
		Scope:       in.Scope,
		Domain:      dom,
		PlannedDate: in.PlannedDate,
		Auditor:     in.Auditor,
		Status:      "Planned",
		Findings:    "",
		CreatedAt:   time.Now().Format(time.RFC3339),
	}

	if err := s.repo.Create(audit); err != nil {
		return nil, err
	}
	return audit, nil
}

func (s *AuditService) ListAudits(statusFilter *string) ([]*domain.Audit, error) {
	all, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	out := make([]*domain.Audit, 0)
	for _, a := range all {
		if statusFilter != nil && !strings.EqualFold(a.Status, *statusFilter) {
			continue
		}
		out = append(out, a)
	}
	return out, nil
}

type UpdateAuditInput struct {
	Status   *string
	Findings *string
}

func (s *AuditService) UpdateAudit(id int, in UpdateAuditInput) (*domain.Audit, error) {
	audit, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if in.Status != nil {
		st := strings.TrimSpace(*in.Status)
		normalized := strings.Title(strings.ToLower(st))
		switch normalized {
		case "Planned", "In Progress", "Completed":
		default:
			return nil, fmt.Errorf("%w: invalid audit status", ErrValidation)
		}
		audit.Status = normalized
	}
	if in.Findings != nil {
		audit.Findings = *in.Findings
	}

	if err := s.repo.Update(audit); err != nil {
		return nil, err
	}
	return audit, nil
}
