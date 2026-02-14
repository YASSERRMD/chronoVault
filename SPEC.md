# ChronoVault - Technical Specification

## Project Overview
- **Project Name**: ChronoVault
- **Type**: Full-stack Time-based Contract Obligation Tracking System
- **Core Functionality**: Track contracts and automatically evaluate whether obligations are active, fulfilled, breached, or expired
- **Target Users**: Legal teams, compliance officers, contract managers

## Technology Stack
- **Frontend**: Vue 3 (Vite)
- **Backend**: Go (Gin)
- **Database**: SQLite
- **Data Access**: Raw SQL only (NO ORM)
- **Real-time**: Gorilla WebSocket
- **Containerized**: Docker

---

## Architecture

### Multi-Tenant Architecture
- Organizations table with isolation
- Users belong to organizations
- JWT authentication with organization context
- Strict data isolation at query level

### Database Schema

#### organizations
| Column | Type | Constraints |
|--------|------|-------------|
| id | TEXT | PRIMARY KEY |
| name | TEXT | NOT NULL |
| created_at | DATETIME | NOT NULL |
| updated_at | DATETIME | NOT NULL |
| deleted_at | DATETIME | NULL |

#### users
| Column | Type | Constraints |
|--------|------|-------------|
| id | TEXT | PRIMARY KEY |
| organization_id | TEXT | FK -> organizations.id |
| email | TEXT | NOT NULL UNIQUE |
| password_hash | TEXT | NOT NULL |
| full_name | TEXT | NOT NULL |
| role | TEXT | NOT NULL |
| created_at | DATETIME | NOT NULL |
| updated_at | DATETIME | NOT NULL |
| deleted_at | DATETIME | NULL |

#### contracts
| Column | Type | Constraints |
|--------|------|-------------|
| id | TEXT | PRIMARY KEY |
| organization_id | TEXT | FK -> organizations.id |
| title | TEXT | NOT NULL |
| counterparty | TEXT | NOT NULL |
| start_date | DATE | NOT NULL |
| end_date | DATE | NOT NULL |
| value | REAL | NOT NULL |
| status | TEXT | NOT NULL |
| created_by | TEXT | FK -> users.id |
| created_at | DATETIME | NOT NULL |
| updated_at | DATETIME | NOT NULL |
| deleted_at | DATETIME | NULL |

#### contract_versions
| Column | Type | Constraints |
|--------|------|-------------|
| id | TEXT | PRIMARY KEY |
| contract_id | TEXT | FK -> contracts.id |
| version | INTEGER | NOT NULL |
| title | TEXT | NOT NULL |
| counterparty | TEXT | NOT NULL |
| start_date | DATE | NOT NULL |
| end_date | DATE | NOT NULL |
| value | REAL | NOT NULL |
| created_by | TEXT | FK -> users.id |
| created_at | DATETIME | NOT NULL |

#### clauses
| Column | Type | Constraints |
|--------|------|-------------|
| id | TEXT | PRIMARY KEY |
| contract_id | TEXT | FK -> contracts.id |
| title | TEXT | NOT NULL |
| description | TEXT | |
| order_index | INTEGER | NOT NULL |
| created_at | DATETIME | NOT NULL |
| updated_at | DATETIME | NOT NULL |
| deleted_at | DATETIME | NULL |

#### obligations
| Column | Type | Constraints |
|--------|------|-------------|
| id | TEXT | PRIMARY KEY |
| clause_id | TEXT | FK -> clauses.id |
| description | TEXT | NOT NULL |
| activation_condition | TEXT | |
| due_date_rule | TEXT | |
| due_date | DATE | |
| penalty_amount | REAL | DEFAULT 0 |
| penalty_type | TEXT | |
| responsible_party | TEXT | |
| depends_on_id | TEXT | FK -> obligations.id |
| status | TEXT | NOT NULL |
| fulfilled_at | DATETIME | NULL |
| breached_at | DATETIME | NULL |
| created_at | DATETIME | NOT NULL |
| updated_at | DATETIME | NOT NULL |
| deleted_at | DATETIME | NULL |

#### obligation_evaluations
| Column | Type | Constraints |
|--------|------|-------------|
| id | TEXT | PRIMARY KEY |
| obligation_id | TEXT | FK -> obligations.id |
| evaluated_at | DATETIME | NOT NULL |
| status_before | TEXT | NOT NULL |
| status_after | TEXT | NOT NULL |
| notes | TEXT | |

#### audit_logs
| Column | Type | Constraints |
|--------|------|-------------|
| id | TEXT | PRIMARY KEY |
| organization_id | TEXT | FK -> organizations.id |
| user_id | TEXT | FK -> users.id |
| entity_type | TEXT | NOT NULL |
| entity_id | TEXT | NOT NULL |
| action | TEXT | NOT NULL |
| old_values | TEXT | |
| new_values | TEXT | |
| created_at | DATETIME | NOT NULL |

---

## API Endpoints

### Auth
- POST /api/auth/login
- POST /api/auth/register
- POST /api/auth/refresh

### Organizations
- GET /api/organizations
- GET /api/organizations/:id
- POST /api/organizations
- PUT /api/organizations/:id

### Contracts
- GET /api/contracts
- GET /api/contracts/:id
- POST /api/contracts
- PUT /api/contracts/:id
- DELETE /api/contracts/:id
- GET /api/contracts/:id/versions
- GET /api/contracts/:id/clauses

### Clauses
- GET /api/clauses
- POST /api/clauses
- PUT /api/clauses/:id
- DELETE /api/clauses/:id

### Obligations
- GET /api/obligations
- GET /api/obligations/:id
- POST /api/obligations
- PUT /api/obligations/:id
- DELETE /api/obligations/:id
- POST /api/obligations/:id/fulfill
- GET /api/obligations/:id/history

### Reports
- GET /api/reports/financial-summary
- GET /api/reports/penalty-tracking
- GET /api/reports/risk-exposure
- GET /api/reports/yearly-impact

### Audit
- GET /api/audit
- GET /api/audit/:entity_type/:entity_id

### WebSocket
- WS /ws

---

## Acceptance Criteria

1. Multi-tenant system with proper organization isolation
2. JWT authentication working correctly
3. Contract CRUD with version history
4. Clause tree structure within contracts
5. Obligation status evaluation (Pending/Active/Fulfilled/Breached/Expired)
6. Dependency logic working correctly
7. Financial penalty calculation and aggregation
8. Real-time WebSocket notifications
9. Complete audit trail
10. Vue 3 frontend with all pages
11. Docker containers for backend and frontend
12. Seed data for testing
