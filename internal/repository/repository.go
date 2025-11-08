package repository

import (
	"errors"

	"github.com/xenakil/integraflow-ims/internal/domain"
)

var ErrNotFound = errors.New("not found")

type RiskRepository interface {
	Create(r *domain.Risk) error
	Update(r *domain.Risk) error
	GetAll() ([]*domain.Risk, error)
	GetByID(id int) (*domain.Risk, error)
}

type IncidentRepository interface {
	Create(i *domain.Incident) error
	Update(i *domain.Incident) error
	GetAll() ([]*domain.Incident, error)
	GetByID(id int) (*domain.Incident, error)
}

type AuditRepository interface {
	Create(a *domain.Audit) error
	Update(a *domain.Audit) error
	GetAll() ([]*domain.Audit, error)
	GetByID(id int) (*domain.Audit, error)
}

type ActionRepository interface {
	Create(a *domain.Action) error
	Update(a *domain.Action) error
	GetAll() ([]*domain.Action, error)
	GetByID(id int) (*domain.Action, error)
}
