-- migrations/001_initial_schema.sql
-- Organizations table
CREATE TABLE IF NOT EXISTS organizations (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at DATETIME DEFAULT NULL
);

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    organization_id TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    full_name TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'viewer',
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at DATETIME DEFAULT NULL,
    FOREIGN KEY (organization_id) REFERENCES organizations(id)
);

-- Contracts table
CREATE TABLE IF NOT EXISTS contracts (
    id TEXT PRIMARY KEY,
    organization_id TEXT NOT NULL,
    title TEXT NOT NULL,
    counterparty TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    value REAL NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'draft',
    created_by TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at DATETIME DEFAULT NULL,
    FOREIGN KEY (organization_id) REFERENCES organizations(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
);

-- Contract versions table
CREATE TABLE IF NOT EXISTS contract_versions (
    id TEXT PRIMARY KEY,
    contract_id TEXT NOT NULL,
    version INTEGER NOT NULL,
    title TEXT NOT NULL,
    counterparty TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    value REAL NOT NULL DEFAULT 0,
    created_by TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (contract_id) REFERENCES contracts(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
);

-- Clauses table
CREATE TABLE IF NOT EXISTS clauses (
    id TEXT PRIMARY KEY,
    contract_id TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    order_index INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at DATETIME DEFAULT NULL,
    FOREIGN KEY (contract_id) REFERENCES contracts(id)
);

-- Obligations table
CREATE TABLE IF NOT EXISTS obligations (
    id TEXT PRIMARY KEY,
    clause_id TEXT NOT NULL,
    description TEXT NOT NULL,
    activation_condition TEXT,
    due_date_rule TEXT,
    due_date DATE,
    penalty_amount REAL DEFAULT 0,
    penalty_type TEXT DEFAULT 'fixed',
    responsible_party TEXT,
    depends_on_id TEXT DEFAULT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    fulfilled_at DATETIME DEFAULT NULL,
    breached_at DATETIME DEFAULT NULL,
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at DATETIME DEFAULT NULL,
    FOREIGN KEY (clause_id) REFERENCES clauses(id),
    FOREIGN KEY (depends_on_id) REFERENCES obligations(id)
);

-- Obligation evaluations table (to prevent duplicates)
CREATE TABLE IF NOT EXISTS obligation_evaluations (
    id TEXT PRIMARY KEY,
    obligation_id TEXT NOT NULL,
    evaluated_at DATETIME NOT NULL DEFAULT (datetime('now')),
    status_before TEXT NOT NULL,
    status_after TEXT NOT NULL,
    notes TEXT,
    FOREIGN KEY (obligation_id) REFERENCES obligations(id)
);

-- Audit logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id TEXT PRIMARY KEY,
    organization_id TEXT NOT NULL,
    user_id TEXT,
    entity_type TEXT NOT NULL,
    entity_id TEXT NOT NULL,
    action TEXT NOT NULL,
    old_values TEXT,
    new_values TEXT,
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (organization_id) REFERENCES organizations(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_organization ON users(organization_id);
CREATE INDEX IF NOT EXISTS idx_contracts_organization ON contracts(organization_id);
CREATE INDEX IF NOT EXISTS idx_contracts_status ON contracts(status);
CREATE INDEX IF NOT EXISTS idx_clauses_contract ON clauses(contract_id);
CREATE INDEX IF NOT EXISTS idx_obligations_clause ON obligations(clause_id);
CREATE INDEX IF NOT EXISTS idx_obligations_status ON obligations(status);
CREATE INDEX IF NOT EXISTS idx_obligations_due_date ON obligations(due_date);
CREATE INDEX IF NOT EXISTS idx_obligations_depends_on ON obligations(depends_on_id);
CREATE INDEX IF NOT EXISTS idx_audit_organization ON audit_logs(organization_id);
CREATE INDEX IF NOT EXISTS idx_audit_entity ON audit_logs(entity_type, entity_id);
