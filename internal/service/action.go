package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/xenakil/integraflow-ims/internal/domain"
	"github.com/xenakil/integraflow-ims/internal/repository"
)

type ActionService struct {
	repo      repository.ActionRepository
	riskRepo  repository.RiskRepository
	incRepo   repository.IncidentRepository
	auditRepo repository.AuditRepository
}

func NewActionService(
	repo repository.ActionRepository,
	riskRepo repository.RiskRepository,
	incRepo repository.IncidentRepository,
	auditRepo repository.AuditRepository,
) *ActionService {
	return &ActionService{
		repo:      repo,
		riskRepo:  riskRepo,
		incRepo:   incRepo,
		auditRepo: auditRepo,
	}
}

type CreateActionInput struct {
	Title       string
	Description string
	SourceType  string // risk, incident, audit
	SourceID    int
	Owner       string
	DueDate     string
}

type ActionListFilter struct {
	Status     *string
	SourceType *string
}

func (s *ActionService) CreateAction(in CreateActionInput) (*domain.Action, error) {
	if strings.TrimSpace(in.Title) == "" || in.SourceID == 0 || strings.TrimSpace(in.SourceType) == "" {
		return nil, fmt.Errorf("%w: title, sourceType and sourceId are required", ErrValidation)
	}

	var canonicalType string
	switch strings.ToLower(strings.TrimSpace(in.SourceType)) {
	case "risk":
		canonicalType = "Risk"
		if _, err := s.riskRepo.GetByID(in.SourceID); err != nil {
			if err == repository.ErrNotFound {
				return nil, fmt.Errorf("%w: source risk not found", ErrValidation)
			}
			return nil, err
		}
	case "incident":
		canonicalType = "Incident"
		if _, err := s.incRepo.GetByID(in.SourceID); err != nil {
			if err == repository.ErrNotFound {
				return nil, fmt.Errorf("%w: source incident not found", ErrValidation)
			}
			return nil, err
		}
	case "audit":
		canonicalType = "Audit"
		if _, err := s.auditRepo.GetByID(in.SourceID); err != nil {
			if err == repository.ErrNotFound {
				return nil, fmt.Errorf("%w: source audit not found", ErrValidation)
			}
			return nil, err
		}
	default:
		return nil, fmt.Errorf("%w: sourceType must be risk, incident or audit", ErrValidation)
	}

	now := time.Now().Format(time.RFC3339)

	act := &domain.Action{
		Title:       in.Title,
		Description: in.Description,
		SourceType:  canonicalType,
		SourceID:    in.SourceID,
		Owner:       in.Owner,
		DueDate:     in.DueDate,
		Status:      "Open",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(act); err != nil {
		return nil, err
	}
	return act, nil
}

func (s *ActionService) ListActions(filter ActionListFilter) ([]*domain.Action, error) {
	all, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	out := make([]*domain.Action, 0)
	for _, a := range all {
		if filter.Status != nil && !strings.EqualFold(a.Status, *filter.Status) {
			continue
		}
		if filter.SourceType != nil && !strings.EqualFold(a.SourceType, *filter.SourceType) {
			continue
		}
		out = append(out, a)
	}
	return out, nil
}

type UpdateActionInput struct {
	Status  *string
	DueDate *string
}

func (s *ActionService) UpdateAction(id int, in UpdateActionInput) (*domain.Action, error) {
	a, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if in.Status != nil {
		st := strings.TrimSpace(*in.Status)
		normalized := strings.Title(strings.ToLower(st))
		switch normalized {
		case "Open", "In Progress", "Done", "Overdue":
		default:
			return nil, fmt.Errorf("%w: invalid action status", ErrValidation)
		}
		a.Status = normalized
	}
	if in.DueDate != nil {
		a.DueDate = *in.DueDate
	}
	a.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := s.repo.Update(a); err != nil {
		return nil, err
	}
	return a, nil
}
