package repository

import (
	"database/sql"
	"encoding/json"
	"time"

	"chronovault/internal/models"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) DB() *sql.DB {
	return r.db
}

func (r *Repository) CreateOrganization(org *models.Organization) error {
	org.ID = uuid.New().String()
	org.CreatedAt = time.Now()
	org.UpdatedAt = time.Now()

	_, err := r.db.Exec(`
		INSERT INTO organizations (id, name, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`, org.ID, org.Name, org.CreatedAt, org.UpdatedAt)
	return err
}

func (r *Repository) GetOrganization(id string) (*models.Organization, error) {
	org := &models.Organization{}
	err := r.db.QueryRow(`
		SELECT id, name, created_at, updated_at, deleted_at
		FROM organizations WHERE id = ? AND deleted_at IS NULL
	`, id).Scan(&org.ID, &org.Name, &org.CreatedAt, &org.UpdatedAt, &org.DeletedAt)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (r *Repository) ListOrganizations() ([]models.Organization, error) {
	rows, err := r.db.Query(`
		SELECT id, name, created_at, updated_at, deleted_at
		FROM organizations WHERE deleted_at IS NULL
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []models.Organization
	for rows.Next() {
		var org models.Organization
		if err := rows.Scan(&org.ID, &org.Name, &org.CreatedAt, &org.UpdatedAt, &org.DeletedAt); err != nil {
			return nil, err
		}
		orgs = append(orgs, org)
	}
	return orgs, nil
}

func (r *Repository) UpdateOrganization(org *models.Organization) error {
	org.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		UPDATE organizations SET name = ?, updated_at = ? WHERE id = ?
	`, org.Name, org.UpdatedAt, org.ID)
	return err
}

func (r *Repository) CreateUser(user *models.User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := r.db.Exec(`
		INSERT INTO users (id, organization_id, email, password_hash, full_name, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, user.ID, user.OrganizationID, user.Email, user.PasswordHash, user.FullName, user.Role, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *Repository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(`
		SELECT id, organization_id, email, password_hash, full_name, role, created_at, updated_at, deleted_at
		FROM users WHERE email = ? AND deleted_at IS NULL
	`, email).Scan(&user.ID, &user.OrganizationID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) GetUserByID(id string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(`
		SELECT id, organization_id, email, password_hash, full_name, role, created_at, updated_at, deleted_at
		FROM users WHERE id = ? AND deleted_at IS NULL
	`, id).Scan(&user.ID, &user.OrganizationID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) CreateContract(contract *models.Contract, userID string) error {
	contract.ID = uuid.New().String()
	contract.CreatedAt = time.Now()
	contract.UpdatedAt = time.Now()
	contract.CreatedBy = userID

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO contracts (id, organization_id, title, counterparty, start_date, end_date, value, status, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, contract.ID, contract.OrganizationID, contract.Title, contract.Counterparty, contract.StartDate, contract.EndDate, contract.Value, contract.Status, contract.CreatedBy, contract.CreatedAt, contract.UpdatedAt)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO contract_versions (id, contract_id, version, title, counterparty, start_date, end_date, value, created_by, created_at)
		VALUES (?, ?, 1, ?, ?, ?, ?, ?, ?, ?)
	`, uuid.New().String(), contract.ID, contract.Title, contract.Counterparty, contract.StartDate, contract.EndDate, contract.Value, userID, contract.CreatedAt)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) GetContract(id, orgID string) (*models.Contract, error) {
	contract := &models.Contract{}
	err := r.db.QueryRow(`
		SELECT id, organization_id, title, counterparty, start_date, end_date, value, status, created_by, created_at, updated_at, deleted_at
		FROM contracts WHERE id = ? AND organization_id = ? AND deleted_at IS NULL
	`, id, orgID).Scan(&contract.ID, &contract.OrganizationID, &contract.Title, &contract.Counterparty, &contract.StartDate, &contract.EndDate, &contract.Value, &contract.Status, &contract.CreatedBy, &contract.CreatedAt, &contract.UpdatedAt, &contract.DeletedAt)
	if err != nil {
		return nil, err
	}
	return contract, nil
}

func (r *Repository) ListContracts(orgID string, status string, page, limit int) ([]models.Contract, int, error) {
	offset := (page - 1) * limit

	var total int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM contracts WHERE organization_id = ? AND deleted_at IS NULL`, orgID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, organization_id, title, counterparty, start_date, end_date, value, status, created_by, created_at, updated_at, deleted_at
		FROM contracts WHERE organization_id = ? AND deleted_at IS NULL
	`
	args := []interface{}{orgID}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var contracts []models.Contract
	for rows.Next() {
		var c models.Contract
		if err := rows.Scan(&c.ID, &c.OrganizationID, &c.Title, &c.Counterparty, &c.StartDate, &c.EndDate, &c.Value, &c.Status, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt); err != nil {
			return nil, 0, err
		}
		contracts = append(contracts, c)
	}
	return contracts, total, nil
}

func (r *Repository) UpdateContract(contract *models.Contract, userID string) error {
	contract.UpdatedAt = time.Now()

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE contracts SET title = ?, counterparty = ?, start_date = ?, end_date = ?, value = ?, status = ?, updated_at = ?
		WHERE id = ? AND organization_id = ?
	`, contract.Title, contract.Counterparty, contract.StartDate, contract.EndDate, contract.Value, contract.Status, contract.UpdatedAt, contract.ID, contract.OrganizationID)
	if err != nil {
		return err
	}

	var version int
	err = tx.QueryRow(`SELECT MAX(version) FROM contract_versions WHERE contract_id = ?`, contract.ID).Scan(&version)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO contract_versions (id, contract_id, version, title, counterparty, start_date, end_date, value, created_by, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, uuid.New().String(), contract.ID, version+1, contract.Title, contract.Counterparty, contract.StartDate, contract.EndDate, contract.Value, userID, contract.UpdatedAt)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) DeleteContract(id, orgID string) error {
	_, err := r.db.Exec(`
		UPDATE contracts SET deleted_at = datetime('now') WHERE id = ? AND organization_id = ?
	`, id, orgID)
	return err
}

func (r *Repository) GetContractVersions(contractID, orgID string) ([]models.ContractVersion, error) {
	rows, err := r.db.Query(`
		SELECT cv.id, cv.contract_id, cv.version, cv.title, cv.counterparty, cv.start_date, cv.end_date, cv.value, cv.created_by, cv.created_at
		FROM contract_versions cv
		JOIN contracts c ON cv.contract_id = c.id
		WHERE cv.contract_id = ? AND c.organization_id = ?
		ORDER BY cv.version DESC
	`, contractID, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []models.ContractVersion
	for rows.Next() {
		var v models.ContractVersion
		if err := rows.Scan(&v.ID, &v.ContractID, &v.Version, &v.Title, &v.Counterparty, &v.StartDate, &v.EndDate, &v.Value, &v.CreatedBy, &v.CreatedAt); err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}
	return versions, nil
}

func (r *Repository) CreateClause(clause *models.Clause) error {
	clause.ID = uuid.New().String()
	clause.CreatedAt = time.Now()
	clause.UpdatedAt = time.Now()

	_, err := r.db.Exec(`
		INSERT INTO clauses (id, contract_id, title, description, order_index, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, clause.ID, clause.ContractID, clause.Title, clause.Description, clause.OrderIndex, clause.CreatedAt, clause.UpdatedAt)
	return err
}

func (r *Repository) GetClause(id, orgID string) (*models.Clause, error) {
	clause := &models.Clause{}
	err := r.db.QueryRow(`
		SELECT c.id, c.contract_id, c.title, c.description, c.order_index, c.created_at, c.updated_at, c.deleted_at
		FROM clauses c
		JOIN contracts co ON c.contract_id = co.id
		WHERE c.id = ? AND co.organization_id = ? AND c.deleted_at IS NULL
	`, id, orgID).Scan(&clause.ID, &clause.ContractID, &clause.Title, &clause.Description, &clause.OrderIndex, &clause.CreatedAt, &clause.UpdatedAt, &clause.DeletedAt)
	if err != nil {
		return nil, err
	}
	return clause, nil
}

func (r *Repository) ListClauses(contractID, orgID string) ([]models.Clause, error) {
	rows, err := r.db.Query(`
		SELECT c.id, c.contract_id, c.title, c.description, c.order_index, c.created_at, c.updated_at, c.deleted_at
		FROM clauses c
		JOIN contracts co ON c.contract_id = co.id
		WHERE c.contract_id = ? AND co.organization_id = ? AND c.deleted_at IS NULL
		ORDER BY c.order_index
	`, contractID, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clauses []models.Clause
	for rows.Next() {
		var c models.Clause
		if err := rows.Scan(&c.ID, &c.ContractID, &c.Title, &c.Description, &c.OrderIndex, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt); err != nil {
			return nil, err
		}
		clauses = append(clauses, c)
	}
	return clauses, nil
}

func (r *Repository) UpdateClause(clause *models.Clause) error {
	clause.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		UPDATE clauses SET title = ?, description = ?, order_index = ?, updated_at = ?
		WHERE id = ?
	`, clause.Title, clause.Description, clause.OrderIndex, clause.UpdatedAt, clause.ID)
	return err
}

func (r *Repository) DeleteClause(id, orgID string) error {
	_, err := r.db.Exec(`
		UPDATE clauses SET deleted_at = datetime('now')
		WHERE id = ? AND contract_id IN (SELECT id FROM contracts WHERE organization_id = ?)
	`, id, orgID)
	return err
}

func (r *Repository) GetContractClauses(contractID, orgID string) ([]models.Clause, error) {
	rows, err := r.db.Query(`
		SELECT c.id, c.contract_id, c.title, c.description, c.order_index, c.created_at, c.updated_at, c.deleted_at
		FROM clauses c
		JOIN contracts co ON c.contract_id = co.id
		WHERE c.contract_id = ? AND co.organization_id = ? AND c.deleted_at IS NULL
		ORDER BY c.order_index
	`, contractID, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clauses []models.Clause
	for rows.Next() {
		var c models.Clause
		if err := rows.Scan(&c.ID, &c.ContractID, &c.Title, &c.Description, &c.OrderIndex, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt); err != nil {
			return nil, err
		}
		clauses = append(clauses, c)
	}
	return clauses, nil
}

func (r *Repository) CreateObligation(obligation *models.Obligation) error {
	obligation.ID = uuid.New().String()
	obligation.CreatedAt = time.Now()
	obligation.UpdatedAt = time.Now()
	obligation.Status = "pending"

	_, err := r.db.Exec(`
		INSERT INTO obligations (id, clause_id, description, activation_condition, due_date_rule, due_date, penalty_amount, penalty_type, responsible_party, depends_on_id, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, obligation.ID, obligation.ClauseID, obligation.Description, obligation.ActivationCondition, obligation.DueDateRule, obligation.DueDate, obligation.PenaltyAmount, obligation.PenaltyType, obligation.ResponsibleParty, obligation.DependsOnID, obligation.Status, obligation.CreatedAt, obligation.UpdatedAt)
	return err
}

func (r *Repository) GetObligation(id, orgID string) (*models.Obligation, error) {
	ob := &models.Obligation{}
	err := r.db.QueryRow(`
		SELECT o.id, o.clause_id, o.description, o.activation_condition, o.due_date_rule, o.due_date, o.penalty_amount, o.penalty_type, o.responsible_party, o.depends_on_id, o.status, o.fulfilled_at, o.breached_at, o.created_at, o.updated_at, o.deleted_at
		FROM obligations o
		JOIN clauses c ON o.clause_id = c.id
		JOIN contracts co ON c.contract_id = co.id
		WHERE o.id = ? AND co.organization_id = ? AND o.deleted_at IS NULL
	`, id, orgID).Scan(&ob.ID, &ob.ClauseID, &ob.Description, &ob.ActivationCondition, &ob.DueDateRule, &ob.DueDate, &ob.PenaltyAmount, &ob.PenaltyType, &ob.ResponsibleParty, &ob.DependsOnID, &ob.Status, &ob.FulfilledAt, &ob.BreachedAt, &ob.CreatedAt, &ob.UpdatedAt, &ob.DeletedAt)
	if err != nil {
		return nil, err
	}
	return ob, nil
}

func (r *Repository) ListObligations(orgID string, status string, page, limit int) ([]models.Obligation, int, error) {
	offset := (page - 1) * limit

	var total int
	err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM obligations o
		JOIN clauses c ON o.clause_id = c.id
		JOIN contracts co ON c.contract_id = co.id
		WHERE co.organization_id = ? AND o.deleted_at IS NULL
	`, orgID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT o.id, o.clause_id, o.description, o.activation_condition, o.due_date_rule, o.due_date, o.penalty_amount, o.penalty_type, o.responsible_party, o.depends_on_id, o.status, o.fulfilled_at, o.breached_at, o.created_at, o.updated_at, o.deleted_at
		FROM obligations o
		JOIN clauses c ON o.clause_id = c.id
		JOIN contracts co ON c.contract_id = co.id
		WHERE co.organization_id = ? AND o.deleted_at IS NULL
	`
	args := []interface{}{orgID}

	if status != "" {
		query += " AND o.status = ?"
		args = append(args, status)
	}

	query += " ORDER BY o.due_date ASC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var obligations []models.Obligation
	for rows.Next() {
		var o models.Obligation
		if err := rows.Scan(&o.ID, &o.ClauseID, &o.Description, &o.ActivationCondition, &o.DueDateRule, &o.DueDate, &o.PenaltyAmount, &o.PenaltyType, &o.ResponsibleParty, &o.DependsOnID, &o.Status, &o.FulfilledAt, &o.BreachedAt, &o.CreatedAt, &o.UpdatedAt, &o.DeletedAt); err != nil {
			return nil, 0, err
		}
		obligations = append(obligations, o)
	}
	return obligations, total, nil
}

func (r *Repository) UpdateObligation(obligation *models.Obligation) error {
	obligation.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		UPDATE obligations SET description = ?, activation_condition = ?, due_date_rule = ?, due_date = ?, penalty_amount = ?, penalty_type = ?, responsible_party = ?, depends_on_id = ?, updated_at = ?
		WHERE id = ?
	`, obligation.Description, obligation.ActivationCondition, obligation.DueDateRule, obligation.DueDate, obligation.PenaltyAmount, obligation.PenaltyType, obligation.ResponsibleParty, obligation.DependsOnID, obligation.UpdatedAt, obligation.ID)
	return err
}

func (r *Repository) DeleteObligation(id, orgID string) error {
	_, err := r.db.Exec(`
		UPDATE obligations SET deleted_at = datetime('now')
		WHERE id = ? AND clause_id IN (SELECT id FROM clauses WHERE contract_id IN (SELECT id FROM contracts WHERE organization_id = ?))
	`, id, orgID)
	return err
}

func (r *Repository) FulfillObligation(id, orgID string) error {
	now := time.Now()
	_, err := r.db.Exec(`
		UPDATE obligations SET status = 'fulfilled', fulfilled_at = ?, updated_at = ?
		WHERE id = ? AND clause_id IN (SELECT id FROM clauses WHERE contract_id IN (SELECT id FROM contracts WHERE organization_id = ?))
	`, now, now, id, orgID)
	return err
}

func (r *Repository) CreateObligationEvaluation(eval *models.ObligationEvaluation) error {
	eval.ID = uuid.New().String()
	eval.EvaluatedAt = time.Now()

	_, err := r.db.Exec(`
		INSERT INTO obligation_evaluations (id, obligation_id, evaluated_at, status_before, status_after, notes)
		VALUES (?, ?, ?, ?, ?, ?)
	`, eval.ID, eval.ObligationID, eval.EvaluatedAt, eval.StatusBefore, eval.StatusAfter, eval.Notes)
	return err
}

func (r *Repository) GetObligationHistory(id, orgID string) ([]models.ObligationEvaluation, error) {
	rows, err := r.db.Query(`
		SELECT oe.id, oe.obligation_id, oe.evaluated_at, oe.status_before, oe.status_after, oe.notes
		FROM obligation_evaluations oe
		JOIN obligations o ON oe.obligation_id = o.id
		JOIN clauses c ON o.clause_id = c.id
		JOIN contracts co ON c.contract_id = co.id
		WHERE oe.obligation_id = ? AND co.organization_id = ?
		ORDER BY oe.evaluated_at DESC
	`, id, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var evaluations []models.ObligationEvaluation
	for rows.Next() {
		var e models.ObligationEvaluation
		if err := rows.Scan(&e.ID, &e.ObligationID, &e.EvaluatedAt, &e.StatusBefore, &e.StatusAfter, &e.Notes); err != nil {
			return nil, err
		}
		evaluations = append(evaluations, e)
	}
	return evaluations, nil
}

func (r *Repository) GetActiveObligationsForEvaluation(orgID string) ([]models.Obligation, error) {
	rows, err := r.db.Query(`
		SELECT o.id, o.clause_id, o.description, o.activation_condition, o.due_date_rule, o.due_date, o.penalty_amount, o.penalty_type, o.responsible_party, o.depends_on_id, o.status, o.fulfilled_at, o.breached_at, o.created_at, o.updated_at, o.deleted_at
		FROM obligations o
		JOIN clauses c ON o.clause_id = c.id
		JOIN contracts co ON c.contract_id = co.id
		WHERE co.organization_id = ?
		AND o.deleted_at IS NULL
		AND o.status IN ('pending', 'active')
		AND co.status = 'active'
		AND co.end_date >= date('now')
		ORDER BY o.due_date ASC
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var obligations []models.Obligation
	for rows.Next() {
		var o models.Obligation
		if err := rows.Scan(&o.ID, &o.ClauseID, &o.Description, &o.ActivationCondition, &o.DueDateRule, &o.DueDate, &o.PenaltyAmount, &o.PenaltyType, &o.ResponsibleParty, &o.DependsOnID, &o.Status, &o.FulfilledAt, &o.BreachedAt, &o.CreatedAt, &o.UpdatedAt, &o.DeletedAt); err != nil {
			return nil, err
		}
		obligations = append(obligations, o)
	}
	return obligations, nil
}

func (r *Repository) GetObligationDependencies(obligationID string) (*models.Obligation, error) {
	o := &models.Obligation{}
	err := r.db.QueryRow(`
		SELECT id, clause_id, description, status, depends_on_id
		FROM obligations WHERE id = ?
	`, obligationID).Scan(&o.ID, &o.ClauseID, &o.Description, &o.Status, &o.DependsOnID)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (r *Repository) UpdateObligationStatus(id, status string) error {
	now := time.Now()
	if status == "breached" {
		_, err := r.db.Exec(`
			UPDATE obligations SET status = ?, breached_at = ?, updated_at = ? WHERE id = ?
		`, status, now, now, id)
		return err
	}
	_, err := r.db.Exec(`
		UPDATE obligations SET status = ?, updated_at = ? WHERE id = ?
	`, status, now, id)
	return err
}

func (r *Repository) GetFinancialSummary(orgID string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	var totalContractValue float64
	var totalContracts int64

	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(value), 0), COALESCE(COUNT(*), 0)
		FROM contracts WHERE organization_id = ? AND deleted_at IS NULL
	`, orgID).Scan(&totalContractValue, &totalContracts)
	if err != nil {
		return nil, err
	}
	result["total_contract_value"] = totalContractValue
	result["total_contracts"] = totalContracts

	var totalPenalties float64
	err = r.db.QueryRow(`
		SELECT COALESCE(SUM(o.penalty_amount), 0)
		FROM obligations o
		JOIN clauses c ON o.clause_id = c.id
		JOIN contracts co ON c.contract_id = co.id
		WHERE co.organization_id = ? AND o.status = 'breached' AND o.deleted_at IS NULL
	`, orgID).Scan(&totalPenalties)
	if err != nil {
		return nil, err
	}
	result["total_penalties"] = totalPenalties

	var activeObligations int64
	err = r.db.QueryRow(`
		SELECT COALESCE(COUNT(*), 0)
		FROM obligations o
		JOIN clauses c ON o.clause_id = c.id
		JOIN contracts co ON c.contract_id = co.id
		WHERE co.organization_id = ? AND o.status = 'active' AND o.deleted_at IS NULL
	`, orgID).Scan(&activeObligations)
	if err != nil {
		return nil, err
	}
	result["active_obligations"] = activeObligations

	return result, nil
}

func (r *Repository) GetPenaltyTracking(orgID string) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(`
		SELECT 
			c.id as contract_id,
			c.title as contract_title,
			COALESCE(SUM(o.penalty_amount), 0) as total_penalty,
			COUNT(o.id) as obligation_count,
			COUNT(CASE WHEN o.status = 'breached' THEN 1 END) as breached_count
		FROM contracts c
		LEFT JOIN clauses cl ON c.id = cl.contract_id
		LEFT JOIN obligations o ON cl.id = o.clause_id AND o.deleted_at IS NULL
		WHERE c.organization_id = ? AND c.deleted_at IS NULL
		GROUP BY c.id
		ORDER BY total_penalty DESC
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var contractID, title string
		var totalPenalty float64
		var obligationCount, breachedCount int
		if err := rows.Scan(&contractID, &title, &totalPenalty, &obligationCount, &breachedCount); err != nil {
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"contract_id":      contractID,
			"contract_title":   title,
			"total_penalty":    totalPenalty,
			"obligation_count": obligationCount,
			"breached_count":   breachedCount,
		})
	}
	return results, nil
}

func (r *Repository) GetRiskExposure(orgID string) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(`
		SELECT 
			c.id,
			c.title,
			c.value,
			COALESCE(SUM(o.penalty_amount), 0) as potential_penalty,
			COUNT(CASE WHEN o.status = 'active' THEN 1 END) as active_obligations,
			COUNT(CASE WHEN o.status = 'pending' THEN 1 END) as pending_obligations
		FROM contracts c
		LEFT JOIN clauses cl ON c.id = cl.contract_id
		LEFT JOIN obligations o ON cl.id = o.clause_id AND o.deleted_at IS NULL
		WHERE c.organization_id = ? AND c.deleted_at IS NULL AND c.status = 'active'
		GROUP BY c.id
		ORDER BY potential_penalty DESC
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id, title string
		var value, potentialPenalty float64
		var activeObligations, pendingObligations int
		if err := rows.Scan(&id, &title, &value, &potentialPenalty, &activeObligations, &pendingObligations); err != nil {
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"contract_id":         id,
			"contract_title":      title,
			"contract_value":      value,
			"potential_penalty":   potentialPenalty,
			"active_obligations":  activeObligations,
			"pending_obligations": pendingObligations,
			"risk_level":          r.calculateRiskLevel(potentialPenalty, value),
		})
	}
	return results, nil
}

func (r *Repository) calculateRiskLevel(penalty, value float64) string {
	if value == 0 {
		return "low"
	}
	ratio := penalty / value
	if ratio >= 0.5 {
		return "high"
	} else if ratio >= 0.2 {
		return "medium"
	}
	return "low"
}

func (r *Repository) GetYearlyImpact(orgID, year string) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(`
		SELECT 
			strftime('%m', o.breached_at) as month,
			COALESCE(SUM(o.penalty_amount), 0) as penalties
		FROM obligations o
		JOIN clauses c ON o.clause_id = c.id
		JOIN contracts co ON c.contract_id = co.id
		WHERE co.organization_id = ?
		AND o.status = 'breached'
		AND o.breached_at IS NOT NULL
		AND strftime('%Y', o.breached_at) = ?
		GROUP BY strftime('%m', o.breached_at)
		ORDER BY month
	`, orgID, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]map[string]interface{}, 0)
	months := map[string]string{"01": "January", "02": "February", "03": "March", "04": "April", "05": "May", "06": "June", "07": "July", "08": "August", "09": "September", "10": "October", "11": "November", "12": "December"}

	for rows.Next() {
		var month, penalties string
		if err := rows.Scan(&month, &penalties); err != nil {
			return nil, err
		}
		monthName := months[month]
		if monthName == "" {
			monthName = month
		}
		results = append(results, map[string]interface{}{
			"month":     monthName,
			"penalties": penalties,
		})
	}
	return results, nil
}

func (r *Repository) CreateAuditLog(log *models.AuditLog) error {
	log.ID = uuid.New().String()
	log.CreatedAt = time.Now()

	_, err := r.db.Exec(`
		INSERT INTO audit_logs (id, organization_id, user_id, entity_type, entity_id, action, old_values, new_values, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, log.ID, log.OrganizationID, log.UserID, log.EntityType, log.EntityID, log.Action, log.OldValues, log.NewValues, log.CreatedAt)
	return err
}

func (r *Repository) ListAuditLogs(orgID string, page, limit int) ([]models.AuditLog, int, error) {
	offset := (page - 1) * limit

	var total int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM audit_logs WHERE organization_id = ?`, orgID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(`
		SELECT id, organization_id, user_id, entity_type, entity_id, action, old_values, new_values, created_at
		FROM audit_logs WHERE organization_id = ?
		ORDER BY created_at DESC LIMIT ? OFFSET ?
	`, orgID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var l models.AuditLog
		if err := rows.Scan(&l.ID, &l.OrganizationID, &l.UserID, &l.EntityType, &l.EntityID, &l.Action, &l.OldValues, &l.NewValues, &l.CreatedAt); err != nil {
			return nil, 0, err
		}
		logs = append(logs, l)
	}
	return logs, total, nil
}

func (r *Repository) GetEntityAuditLogs(entityType, entityID, orgID string) ([]models.AuditLog, error) {
	rows, err := r.db.Query(`
		SELECT id, organization_id, user_id, entity_type, entity_id, action, old_values, new_values, created_at
		FROM audit_logs WHERE entity_type = ? AND entity_id = ? AND organization_id = ?
		ORDER BY created_at DESC
	`, entityType, entityID, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var l models.AuditLog
		if err := rows.Scan(&l.ID, &l.OrganizationID, &l.UserID, &l.EntityType, &l.EntityID, &l.Action, &l.OldValues, &l.NewValues, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

func (r *Repository) LogAudit(orgID, userID, entityType, entityID, action string, oldVals, newVals interface{}) error {
	var oldJSON, newJSON string
	if oldVals != nil {
		b, _ := json.Marshal(oldVals)
		oldJSON = string(b)
	}
	if newVals != nil {
		b, _ := json.Marshal(newVals)
		newJSON = string(b)
	}

	log := &models.AuditLog{
		OrganizationID: orgID,
		UserID:         userID,
		EntityType:     entityType,
		EntityID:       entityID,
		Action:         action,
		OldValues:      oldJSON,
		NewValues:      newJSON,
	}
	return r.CreateAuditLog(log)
}
