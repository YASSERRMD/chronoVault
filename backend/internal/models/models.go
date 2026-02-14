package models

import "time"

type Organization struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type User struct {
	ID             string     `json:"id"`
	OrganizationID string     `json:"organization_id"`
	Email          string     `json:"email"`
	PasswordHash   string     `json:"-"`
	FullName       string     `json:"full_name"`
	Role           string     `json:"role"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type Contract struct {
	ID             string     `json:"id"`
	OrganizationID string     `json:"organization_id"`
	Title          string     `json:"title"`
	Counterparty   string     `json:"counterparty"`
	StartDate      string     `json:"start_date"`
	EndDate        string     `json:"end_date"`
	Value          float64    `json:"value"`
	Status         string     `json:"status"`
	CreatedBy      string     `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type ContractVersion struct {
	ID           string    `json:"id"`
	ContractID   string    `json:"contract_id"`
	Version      int       `json:"version"`
	Title        string    `json:"title"`
	Counterparty string    `json:"counterparty"`
	StartDate    string    `json:"start_date"`
	EndDate      string    `json:"end_date"`
	Value        float64   `json:"value"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
}

type Clause struct {
	ID          string     `json:"id"`
	ContractID  string     `json:"contract_id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	OrderIndex  int        `json:"order_index"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type Obligation struct {
	ID                  string     `json:"id"`
	ClauseID            string     `json:"clause_id"`
	Description         string     `json:"description"`
	ActivationCondition string     `json:"activation_condition,omitempty"`
	DueDateRule         string     `json:"due_date_rule,omitempty"`
	DueDate             string     `json:"due_date,omitempty"`
	PenaltyAmount       float64    `json:"penalty_amount"`
	PenaltyType         string     `json:"penalty_type"`
	ResponsibleParty    string     `json:"responsible_party,omitempty"`
	DependsOnID         *string    `json:"depends_on_id,omitempty"`
	Status              string     `json:"status"`
	FulfilledAt         *time.Time `json:"fulfilled_at,omitempty"`
	BreachedAt          *time.Time `json:"breached_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty"`
}

type ObligationEvaluation struct {
	ID           string    `json:"id"`
	ObligationID string    `json:"obligation_id"`
	EvaluatedAt  time.Time `json:"evaluated_at"`
	StatusBefore string    `json:"status_before"`
	StatusAfter  string    `json:"status_after"`
	Notes        string    `json:"notes,omitempty"`
}

type AuditLog struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	UserID         string    `json:"user_id,omitempty"`
	EntityType     string    `json:"entity_type"`
	EntityID       string    `json:"entity_id"`
	Action         string    `json:"action"`
	OldValues      string    `json:"old_values,omitempty"`
	NewValues      string    `json:"new_values,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}
