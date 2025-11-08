package sqlite

import (
	"database/sql"
	"errors"

	_ "github.com/glebarez/sqlite"

	"github.com/xenakil/integraflow-ims/internal/domain"
	"github.com/xenakil/integraflow-ims/internal/repository"
)

func NewDB(path string) (*sql.DB, error) {
	// Note: driver name is "sqlite" (glebarez/sqlite), not "sqlite3"
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if err := initSchema(db); err != nil {
		return nil, err
	}
	return db, nil
}

func initSchema(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS risks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			process TEXT NOT NULL,
			domain TEXT NOT NULL,
			description TEXT,
			likelihood INTEGER NOT NULL,
			impact INTEGER NOT NULL,
			score INTEGER NOT NULL,
			level TEXT NOT NULL,
			owner TEXT,
			status TEXT NOT NULL,
			created_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS incidents (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			domain TEXT NOT NULL,
			related_risk_id INTEGER,
			severity INTEGER NOT NULL,
			likelihood INTEGER NOT NULL,
			risk_score INTEGER NOT NULL,
			risk_level TEXT NOT NULL,
			root_cause TEXT,
			status TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS audits (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			scope TEXT NOT NULL,
			domain TEXT NOT NULL,
			planned_date TEXT NOT NULL,
			auditor TEXT,
			status TEXT NOT NULL,
			findings TEXT,
			created_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS actions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			source_type TEXT NOT NULL,
			source_id INTEGER NOT NULL,
			owner TEXT,
			due_date TEXT,
			status TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return err
		}
	}
	return nil
}

// ---------- Risk repository ----------

type RiskRepository struct {
	db *sql.DB
}

func NewRiskRepository(db *sql.DB) *RiskRepository {
	return &RiskRepository{db: db}
}

func (r *RiskRepository) Create(risk *domain.Risk) error {
	res, err := r.db.Exec(`
		INSERT INTO risks (title, process, domain, description, likelihood, impact, score, level, owner, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		risk.Title, risk.Process, string(risk.Domain), risk.Description,
		risk.Likelihood, risk.Impact, risk.Score, risk.Level,
		risk.Owner, risk.Status, risk.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		risk.ID = int(id)
	}
	return nil
}

func (r *RiskRepository) Update(risk *domain.Risk) error {
	res, err := r.db.Exec(`
		UPDATE risks SET title=?, process=?, domain=?, description=?, likelihood=?, impact=?, score=?, level=?, owner=?, status=?, created_at=?
		WHERE id=?`,
		risk.Title, risk.Process, string(risk.Domain), risk.Description,
		risk.Likelihood, risk.Impact, risk.Score, risk.Level,
		risk.Owner, risk.Status, risk.CreatedAt, risk.ID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *RiskRepository) GetAll() ([]*domain.Risk, error) {
	rows, err := r.db.Query(`
		SELECT id, title, process, domain, description, likelihood, impact, score, level, owner, status, created_at
		FROM risks`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*domain.Risk
	for rows.Next() {
		var d string
		risk := &domain.Risk{}
		if err := rows.Scan(
			&risk.ID, &risk.Title, &risk.Process, &d, &risk.Description,
			&risk.Likelihood, &risk.Impact, &risk.Score, &risk.Level,
			&risk.Owner, &risk.Status, &risk.CreatedAt,
		); err != nil {
			return nil, err
		}
		risk.Domain = domain.Domain(d)
		out = append(out, risk)
	}
	return out, nil
}

func (r *RiskRepository) GetByID(id int) (*domain.Risk, error) {
	row := r.db.QueryRow(`
		SELECT id, title, process, domain, description, likelihood, impact, score, level, owner, status, created_at
		FROM risks WHERE id = ?`, id)

	var d string
	risk := &domain.Risk{}
	if err := row.Scan(
		&risk.ID, &risk.Title, &risk.Process, &d, &risk.Description,
		&risk.Likelihood, &risk.Impact, &risk.Score, &risk.Level,
		&risk.Owner, &risk.Status, &risk.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	risk.Domain = domain.Domain(d)
	return risk, nil
}

// ---------- Incident repository ----------

type IncidentRepository struct {
	db *sql.DB
}

func NewIncidentRepository(db *sql.DB) *IncidentRepository {
	return &IncidentRepository{db: db}
}

func (r *IncidentRepository) Create(inc *domain.Incident) error {
	var related interface{} = nil
	if inc.RelatedRiskID != nil {
		related = *inc.RelatedRiskID
	}
	res, err := r.db.Exec(`
		INSERT INTO incidents (title, description, domain, related_risk_id, severity, likelihood, risk_score, risk_level, root_cause, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		inc.Title, inc.Description, string(inc.Domain),
		related, inc.Severity, inc.Likelihood, inc.RiskScore,
		inc.RiskLevel, inc.RootCause, inc.Status,
		inc.CreatedAt, inc.UpdatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		inc.ID = int(id)
	}
	return nil
}

func (r *IncidentRepository) Update(inc *domain.Incident) error {
	var related interface{} = nil
	if inc.RelatedRiskID != nil {
		related = *inc.RelatedRiskID
	}
	res, err := r.db.Exec(`
		UPDATE incidents
		SET title=?, description=?, domain=?, related_risk_id=?, severity=?, likelihood=?, risk_score=?, risk_level=?, root_cause=?, status=?, created_at=?, updated_at=?
		WHERE id=?`,
		inc.Title, inc.Description, string(inc.Domain),
		related, inc.Severity, inc.Likelihood, inc.RiskScore, inc.RiskLevel,
		inc.RootCause, inc.Status, inc.CreatedAt, inc.UpdatedAt, inc.ID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *IncidentRepository) GetAll() ([]*domain.Incident, error) {
	rows, err := r.db.Query(`
		SELECT id, title, description, domain, related_risk_id, severity, likelihood, risk_score, risk_level, root_cause, status, created_at, updated_at
		FROM incidents`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*domain.Incident
	for rows.Next() {
		var d string
		var related sqlNullInt
		inc := &domain.Incident{}
		if err := rows.Scan(
			&inc.ID, &inc.Title, &inc.Description, &d, &related,
			&inc.Severity, &inc.Likelihood, &inc.RiskScore,
			&inc.RiskLevel, &inc.RootCause, &inc.Status,
			&inc.CreatedAt, &inc.UpdatedAt,
		); err != nil {
			return nil, err
		}
		inc.Domain = domain.Domain(d)
		if related.Valid {
			id := related.V
			inc.RelatedRiskID = &id
		}
		out = append(out, inc)
	}
	return out, nil
}

func (r *IncidentRepository) GetByID(id int) (*domain.Incident, error) {
	row := r.db.QueryRow(`
		SELECT id, title, description, domain, related_risk_id, severity, likelihood, risk_score, risk_level, root_cause, status, created_at, updated_at
		FROM incidents WHERE id = ?`, id)

	var d string
	var related sqlNullInt
	inc := &domain.Incident{}
	if err := row.Scan(
		&inc.ID, &inc.Title, &inc.Description, &d, &related,
		&inc.Severity, &inc.Likelihood, &inc.RiskScore,
		&inc.RiskLevel, &inc.RootCause, &inc.Status,
		&inc.CreatedAt, &inc.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	inc.Domain = domain.Domain(d)
	if related.Valid {
		id := related.V
		inc.RelatedRiskID = &id
	}
	return inc, nil
}

// ---------- Audit repository ----------

type AuditRepository struct {
	db *sql.DB
}

func NewAuditRepository(db *sql.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Create(a *domain.Audit) error {
	res, err := r.db.Exec(`
		INSERT INTO audits (title, scope, domain, planned_date, auditor, status, findings, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		a.Title, a.Scope, string(a.Domain), a.PlannedDate, a.Auditor,
		a.Status, a.Findings, a.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		a.ID = int(id)
	}
	return nil
}

func (r *AuditRepository) Update(a *domain.Audit) error {
	res, err := r.db.Exec(`
		UPDATE audits
		SET title=?, scope=?, domain=?, planned_date=?, auditor=?, status=?, findings=?, created_at=?
		WHERE id=?`,
		a.Title, a.Scope, string(a.Domain), a.PlannedDate, a.Auditor,
		a.Status, a.Findings, a.CreatedAt, a.ID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *AuditRepository) GetAll() ([]*domain.Audit, error) {
	rows, err := r.db.Query(`
		SELECT id, title, scope, domain, planned_date, auditor, status, findings, created_at
		FROM audits`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*domain.Audit
	for rows.Next() {
		var d string
		a := &domain.Audit{}
		if err := rows.Scan(
			&a.ID, &a.Title, &a.Scope, &d,
			&a.PlannedDate, &a.Auditor, &a.Status,
			&a.Findings, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		a.Domain = domain.Domain(d)
		out = append(out, a)
	}
	return out, nil
}

func (r *AuditRepository) GetByID(id int) (*domain.Audit, error) {
	row := r.db.QueryRow(`
		SELECT id, title, scope, domain, planned_date, auditor, status, findings, created_at
		FROM audits WHERE id = ?`, id)

	var d string
	a := &domain.Audit{}
	if err := row.Scan(
		&a.ID, &a.Title, &a.Scope, &d,
		&a.PlannedDate, &a.Auditor, &a.Status,
		&a.Findings, &a.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	a.Domain = domain.Domain(d)
	return a, nil
}

// ---------- Action repository ----------

type ActionRepository struct {
	db *sql.DB
}

func NewActionRepository(db *sql.DB) *ActionRepository {
	return &ActionRepository{db: db}
}

func (r *ActionRepository) Create(a *domain.Action) error {
	res, err := r.db.Exec(`
		INSERT INTO actions (title, description, source_type, source_id, owner, due_date, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		a.Title, a.Description, a.SourceType, a.SourceID,
		a.Owner, a.DueDate, a.Status, a.CreatedAt, a.UpdatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		a.ID = int(id)
	}
	return nil
}

func (r *ActionRepository) Update(a *domain.Action) error {
	res, err := r.db.Exec(`
		UPDATE actions
		SET title=?, description=?, source_type=?, source_id=?, owner=?, due_date=?, status=?, created_at=?, updated_at=?
		WHERE id=?`,
		a.Title, a.Description, a.SourceType, a.SourceID,
		a.Owner, a.DueDate, a.Status, a.CreatedAt, a.UpdatedAt, a.ID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *ActionRepository) GetAll() ([]*domain.Action, error) {
	rows, err := r.db.Query(`
		SELECT id, title, description, source_type, source_id, owner, due_date, status, created_at, updated_at
		FROM actions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*domain.Action
	for rows.Next() {
		a := &domain.Action{}
		if err := rows.Scan(
			&a.ID, &a.Title, &a.Description, &a.SourceType, &a.SourceID,
			&a.Owner, &a.DueDate, &a.Status, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

func (r *ActionRepository) GetByID(id int) (*domain.Action, error) {
	row := r.db.QueryRow(`
		SELECT id, title, description, source_type, source_id, owner, due_date, status, created_at, updated_at
		FROM actions WHERE id = ?`, id)

	a := &domain.Action{}
	if err := row.Scan(
		&a.ID, &a.Title, &a.Description, &a.SourceType, &a.SourceID,
		&a.Owner, &a.DueDate, &a.Status, &a.CreatedAt, &a.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return a, nil
}

// small helper for nullable integer
type sqlNullInt struct {
	Valid bool
	V     int
}

func (n *sqlNullInt) Scan(value any) error {
	if value == nil {
		n.Valid = false
		n.V = 0
		return nil
	}
	switch v := value.(type) {
	case int64:
		n.V = int(v)
		n.Valid = true
		return nil
	default:
		return errors.New("invalid type for sqlNullInt")
	}
}
