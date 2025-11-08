package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/xenakil/integraflow-ims/internal/domain"
	"github.com/xenakil/integraflow-ims/internal/repository"
)

type IncidentService struct {
	incRepo  repository.IncidentRepository
	riskRepo repository.RiskRepository
}

func NewIncidentService(incRepo repository.IncidentRepository, riskRepo repository.RiskRepository) *IncidentService {
	return &IncidentService{incRepo: incRepo, riskRepo: riskRepo}
}

type CreateIncidentInput struct {
	Title         string
	Description   string
	Domain        string
	RelatedRiskID *int
	Severity      int
	Likelihood    int
}

type IncidentListFilter struct {
	Domain *domain.Domain
	Status *string
}

func (s *IncidentService) CreateIncident(in CreateIncidentInput) (*domain.Incident, error) {
	if strings.TrimSpace(in.Title) == "" || strings.TrimSpace(in.Description) == "" {
		return nil, fmt.Errorf("%w: title and description are required", ErrValidation)
	}

	dom, err := domain.ParseDomain(in.Domain)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidation, err)
	}

	if in.Severity < 1 || in.Severity > 5 || in.Likelihood < 1 || in.Likelihood > 5 {
		return nil, fmt.Errorf("%w: severity and likelihood must be between 1 and 5", ErrValidation)
	}

	if in.RelatedRiskID != nil {
		if _, err := s.riskRepo.GetByID(*in.RelatedRiskID); err != nil {
			if err == repository.ErrNotFound {
				return nil, fmt.Errorf("%w: related risk ID does not exist", ErrValidation)
			}
			return nil, err
		}
	}

	score := in.Severity * in.Likelihood
	level := domain.RiskLevelFromScore(score)
	now := time.Now().Format(time.RFC3339)

	inc := &domain.Incident{
		Title:         in.Title,
		Description:   in.Description,
		Domain:        dom,
		RelatedRiskID: in.RelatedRiskID,
		Severity:      in.Severity,
		Likelihood:    in.Likelihood,
		RiskScore:     score,
		RiskLevel:     level,
		RootCause:     "",
		Status:        "Open",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.incRepo.Create(inc); err != nil {
		return nil, err
	}
	return inc, nil
}

func (s *IncidentService) ListIncidents(filter IncidentListFilter) ([]*domain.Incident, error) {
	all, err := s.incRepo.GetAll()
	if err != nil {
		return nil, err
	}

	out := make([]*domain.Incident, 0)
	for _, inc := range all {
		if filter.Domain != nil && inc.Domain != *filter.Domain {
			continue
		}
		if filter.Status != nil && !strings.EqualFold(inc.Status, *filter.Status) {
			continue
		}
		out = append(out, inc)
	}
	return out, nil
}

func (s *IncidentService) GetIncident(id int) (*domain.Incident, error) {
	return s.incRepo.GetByID(id)
}

type UpdateIncidentInput struct {
	RootCause *string
	Status    *string
}

func (s *IncidentService) UpdateIncident(id int, in UpdateIncidentInput) (*domain.Incident, error) {
	inc, err := s.incRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if in.RootCause != nil {
		inc.RootCause = strings.TrimSpace(*in.RootCause)
	}
	if in.Status != nil {
		status := strings.TrimSpace(*in.Status)
		normalized := strings.Title(strings.ToLower(status))
		switch normalized {
		case "Open", "Investigation", "Closed":
		default:
			return nil, fmt.Errorf("%w: invalid incident status", ErrValidation)
		}
		inc.Status = normalized
	}
	inc.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := s.incRepo.Update(inc); err != nil {
		return nil, err
	}
	return inc, nil
}
