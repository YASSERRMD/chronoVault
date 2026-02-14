-- Seed data for ChronoVault
INSERT INTO organizations (id, name, created_at, updated_at) VALUES 
('org-001', 'Demo Corporation', datetime('now'), datetime('now'));

INSERT INTO users (id, organization_id, email, password_hash, full_name, role, created_at, updated_at) VALUES 
('user-001', 'org-001', 'admin@demo.com', '$2a$10$EqKcp1WFKVQISheBBHb0JOj6Wj6HJxW6YJ7.HX6hGyYDq6q5b5r3i', 'Admin User', 'admin', datetime('now'), datetime('now'));

INSERT INTO contracts (id, organization_id, title, counterparty, start_date, end_date, value, status, created_by, created_at, updated_at) VALUES 
('contract-001', 'org-001', 'SaaS Service Agreement', 'TechCorp Inc.', '2024-01-01', '2025-12-31', 150000.00, 'active', 'user-001', datetime('now'), datetime('now'));

INSERT INTO clauses (id, contract_id, title, description, order_index, created_at, updated_at) VALUES 
('clause-001', 'contract-001', 'Payment Terms', 'Payment within 30 days', 1, datetime('now'), datetime('now'));

INSERT INTO obligations (id, clause_id, description, due_date, penalty_amount, penalty_type, responsible_party, status, created_at, updated_at) VALUES 
('ob-001', 'clause-001', 'Pay monthly fee', '2024-01-30', 500, 'fixed', 'Finance', 'fulfilled', datetime('now'), datetime('now')),
('ob-002', 'clause-001', 'Pay monthly fee', '2024-02-28', 500, 'fixed', 'Finance', 'active', datetime('now'), datetime('now'));
