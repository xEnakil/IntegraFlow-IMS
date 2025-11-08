package service

import (
	"github.com/xenakil/integraflow-ims/internal/domain"
	"github.com/xenakil/integraflow-ims/internal/repository"
)

type DashboardService struct {
	riskRepo   repository.RiskRepository
	incRepo    repository.IncidentRepository
	actionRepo repository.ActionRepository
}

func NewDashboardService(
	riskRepo repository.RiskRepository,
	incRepo repository.IncidentRepository,
	actionRepo repository.ActionRepository,
) *DashboardService {
	return &DashboardService{
		riskRepo:   riskRepo,
		incRepo:    incRepo,
		actionRepo: actionRepo,
	}
}

func (s *DashboardService) GetDashboard() (*domain.Dashboard, error) {
	risks, err := s.riskRepo.GetAll()
	if err != nil {
		return nil, err
	}
	incidents, err := s.incRepo.GetAll()
	if err != nil {
		return nil, err
	}
	actions, err := s.actionRepo.GetAll()
	if err != nil {
		return nil, err
	}

	dash := &domain.Dashboard{
		ActionsByStatus:   make(map[string]int),
		IncidentsByDomain: make(map[domain.Domain]int),
	}

	dash.TotalRisks = len(risks)
	for _, r := range risks {
		if r.Level == "High" {
			dash.HighRisks++
		}
	}

	dash.TotalIncidents = len(incidents)
	for _, inc := range incidents {
		if inc.Status == "Open" {
			dash.OpenIncidents++
		}
		dash.IncidentsByDomain[inc.Domain]++
	}

	for _, a := range actions {
		dash.ActionsByStatus[a.Status]++
	}

	return dash, nil
}
