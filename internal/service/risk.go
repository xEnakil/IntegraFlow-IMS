package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/xenakil/integraflow-ims/internal/domain"
	"github.com/xenakil/integraflow-ims/internal/repository"
)

var ErrValidation = errors.New("validation error")

type RiskService struct {
	repo repository.RiskRepository
}

func NewRiskService(repo repository.RiskRepository) *RiskService {
	return &RiskService{repo: repo}
}

type CreateRiskInput struct {
	Title       string
	Process     string
	Domain      string
	Description string
	Likelihood  int
	Impact      int
	Owner       string
}

type RiskListFilter struct {
	Domain *domain.Domain
	Status *string
}

func (s *RiskService) CreateRisk(in CreateRiskInput) (*domain.Risk, error) {
	if strings.TrimSpace(in.Title) == "" {
		return nil, fmt.Errorf("%w: title is required", ErrValidation)
	}
	if strings.TrimSpace(in.Process) == "" {
		return nil, fmt.Errorf("%w: process is required", ErrValidation)
	}

	dom, err := domain.ParseDomain(in.Domain)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidation, err)
	}

	if in.Likelihood < 1 || in.Likelihood > 5 || in.Impact < 1 || in.Impact > 5 {
		return nil, fmt.Errorf("%w: likelihood and impact must be between 1 and 5", ErrValidation)
	}

	score := in.Likelihood * in.Impact
	level := domain.RiskLevelFromScore(score)

	r := &domain.Risk{
		Title:       in.Title,
		Process:     in.Process,
		Domain:      dom,
		Description: in.Description,
		Likelihood:  in.Likelihood,
		Impact:      in.Impact,
		Score:       score,
		Level:       level,
		Owner:       in.Owner,
		Status:      "Open",
		CreatedAt:   time.Now().Format(time.RFC3339),
	}

	if err := s.repo.Create(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *RiskService) ListRisks(filter RiskListFilter) ([]*domain.Risk, error) {
	all, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	out := make([]*domain.Risk, 0)
	for _, r := range all {
		if filter.Domain != nil && r.Domain != *filter.Domain {
			continue
		}
		if filter.Status != nil && !strings.EqualFold(r.Status, *filter.Status) {
			continue
		}
		out = append(out, r)
	}
	return out, nil
}

func (s *RiskService) UpdateStatus(id int, status string) (*domain.Risk, error) {
	status = strings.TrimSpace(status)
	if status == "" {
		return nil, fmt.Errorf("%w: status is required", ErrValidation)
	}

	normalized := strings.Title(strings.ToLower(status))
	switch normalized {
	case "Open", "Accepted", "Mitigated":
	default:
		return nil, fmt.Errorf("%w: invalid risk status", ErrValidation)
	}

	r, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	r.Status = normalized
	if err := s.repo.Update(r); err != nil {
		return nil, err
	}
	return r, nil
}
