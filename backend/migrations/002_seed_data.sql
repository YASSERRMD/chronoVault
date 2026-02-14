-- Seed data for ChronoVault
-- Run this after migrations

-- Insert demo organization
INSERT INTO organizations (id, name, created_at, updated_at) VALUES 
('org-001', 'Demo Corporation', datetime('now'), datetime('now'));

-- Insert demo users (password: password123 for all - hash is bcrypt for "password")
INSERT INTO users (id, organization_id, email, password_hash, full_name, role, created_at, updated_at) VALUES 
('user-001', 'org-001', 'admin@demo.com', '$2a$10$EqKcp1WFKVQISheBBHb0JOj6Wj6HJxW6YJ7.HX6hGyYDq6q5b5r3i', 'Admin User', 'admin', datetime('now'), datetime('now')),
('user-002', 'org-001', 'legal@demo.com', '$2a$10$EqKcp1WFKVQISheBBHb0JOj6Wj6HJxW6YJ7.HX6hGyYDq6q5b5r3i', 'Legal Manager', 'legal_manager', datetime('now'), datetime('now')),
('user-003', 'org-001', 'analyst@demo.com', '$2a$10$EqKcp1WFKVQISheBBHb0JOj6Wj6HJxW6YJ7.HX6hGyYDq6q5b5r3i', 'Analyst', 'analyst', datetime('now'), datetime('now')),
('user-004', 'org-001', 'viewer@demo.com', '$2a$10$EqKcp1WFKVQISheBBHb0JOj6Wj6HJxW6YJ7.HX6hGyYDq6q5b5r3i', 'Viewer', 'viewer', datetime('now'), datetime('now'));

-- Insert demo contracts
INSERT INTO contracts (id, organization_id, title, counterparty, start_date, end_date, value, status, created_by, created_at, updated_at) VALUES 
('contract-001', 'org-001', 'SaaS Service Agreement', 'TechCorp Inc.', '2024-01-01', '2025-12-31', 150000.00, 'active', 'user-001', datetime('now'), datetime('now')),
('contract-002', 'org-001', 'Office Lease Agreement', 'Property Holdings LLC', '2024-06-01', '2026-05-31', 250000.00, 'active', 'user-001', datetime('now'), datetime('now')),
('contract-003', 'org-001', 'Consulting Services', 'Business Consultants Group', '2024-03-01', '2024-12-31', 75000.00, 'expired', 'user-001', datetime('now'), datetime('now')),
('contract-004', 'org-001', 'Software License', 'Software Vendor Co.', '2024-09-01', '2025-08-31', 50000.00, 'active', 'user-002', datetime('now'), datetime('now'));

-- Insert contract versions
INSERT INTO contract_versions (id, contract_id, version, title, counterparty, start_date, end_date, value, created_by, created_at) VALUES 
('ver-001', 'contract-001', 1, 'SaaS Service Agreement', 'TechCorp Inc.', '2024-01-01', '2025-12-31', 150000.00, 'user-001', datetime('now')),
('ver-002', 'contract-002', 1, 'Office Lease Agreement', 'Property Holdings LLC', '2024-06-01', '2026-05-31', 250000.00, 'user-001', datetime('now')),
('ver-003', 'contract-003', 1, 'Consulting Services', 'Business Consultants Group', '2024-03-01', '2024-12-31', 75000.00, 'user-001', datetime('now')),
('ver-004', 'contract-004', 1, 'Software License', 'Software Vendor Co.', '2024-09-01', '2025-08-31', 50000.00, 'user-002', datetime('now'));

-- Insert clauses for Contract 001
INSERT INTO clauses (id, contract_id, title, description, order_index, created_at, updated_at) VALUES 
('clause-001', 'contract-001', 'Payment Terms', 'Payment shall be made within 30 days of invoice date', 1, datetime('now'), datetime('now')),
('clause-002', 'contract-001', 'Service Level Agreement', 'System uptime of 99.9% guaranteed', 2, datetime('now'), datetime('now')),
('clause-003', 'contract-001', 'Data Protection', 'All customer data shall be encrypted and protected', 3, datetime('now'), datetime('now'));

-- Insert clauses for Contract 002
INSERT INTO clauses (id, contract_id, title, description, order_index, created_at, updated_at) VALUES 
('clause-004', 'contract-002', 'Rent Payment', 'Monthly rent due on 1st of each month', 1, datetime('now'), datetime('now')),
('clause-005', 'contract-002', 'Maintenance', 'Landlord responsible for building maintenance', 2, datetime('now'), datetime('now'));

-- Insert clauses for Contract 003
INSERT INTO clauses (id, contract_id, title, description, order_index, created_at, updated_at) VALUES 
('clause-006', 'contract-003', 'Deliverables', 'Monthly reports and analysis', 1, datetime('now'), datetime('now')),
('clause-007', 'contract-003', 'Payment Schedule', 'Quarterly payments', 2, datetime('now'), datetime('now'));

-- Insert clauses for Contract 004
INSERT INTO clauses (id, contract_id, title, description, order_index, created_at, updated_at) VALUES 
('clause-008', 'contract-004', 'License Grant', 'Non-exclusive license to use software', 1, datetime('now'), datetime('now')),
('clause-009', 'contract-004', 'Support', '24/7 technical support included', 2, datetime('now'), datetime('now'));

-- Insert obligations
INSERT INTO obligations (id, clause_id, description, activation_condition, due_date_rule, due_date, penalty_amount, penalty_type, responsible_party, status, fulfilled_at, breached_at, created_at, updated_at) VALUES 
('ob-001', 'clause-001', 'Pay monthly subscription fee', 'Contract start', 'Monthly by 30th', '2024-01-30', 500, 'fixed', 'Finance', 'fulfilled', datetime('now'), NULL, datetime('now'), datetime('now')),
('ob-002', 'clause-001', 'Pay February subscription fee', 'Contract start', 'Monthly by 30th', '2024-02-28', 500, 'fixed', 'Finance', 'fulfilled', datetime('now'), NULL, datetime('now'), datetime('now')),
('ob-003', 'clause-001', 'Pay March subscription fee', 'Contract start', 'Monthly by 30th', '2024-03-30', 500, 'fixed', 'Finance', 'fulfilled', datetime('now'), NULL, datetime('now'), datetime('now')),
('ob-004', 'clause-001', 'Pay April subscription fee', 'Contract start', 'Monthly by 30th', '2024-04-30', 500, 'fixed', 'Finance', 'breached', NULL, datetime('now'), datetime('now'), datetime('now')),
('ob-005', 'clause-001', 'Pay May subscription fee', 'Contract start', 'Monthly by 30th', '2024-05-30', 500, 'fixed', 'Finance', 'active', NULL, NULL, datetime('now'), datetime('now')),
('ob-006', 'clause-002', 'Maintain 99.9% uptime', 'Contract start', 'Continuous', '2024-12-31', 1000, 'daily', 'IT', 'active', NULL, NULL, datetime('now'), datetime('now')),
('ob-007', 'clause-003', 'Implement encryption at rest', 'Contract start', 'Within 60 days', '2024-03-01', 2000, 'fixed', 'Security', 'fulfilled', datetime('now'), NULL, datetime('now'), datetime('now')),
('ob-008', 'clause-003', 'Annual security audit', 'Contract start', 'Annual', '2024-12-31', 5000, 'fixed', 'Security', 'active', NULL, NULL, datetime('now'), datetime('now')),
('ob-009', 'clause-004', 'Pay June rent', 'Month start', 'Monthly by 1st', '2024-06-01', 20000, 'fixed', 'Finance', 'fulfilled', datetime('now'), NULL, datetime('now'), datetime('now')),
('ob-010', 'clause-004', 'Pay July rent', 'Month start', 'Monthly by 1st', '2024-07-01', 20000, 'fixed', 'Finance', 'fulfilled', datetime('now'), NULL, datetime('now'), datetime('now')),
('ob-011', 'clause-004', 'Pay August rent', 'Month start', 'Monthly by 1st', '2024-08-01', 20000, 'fixed', 'Finance', 'active', NULL, NULL, datetime('now'), datetime('now')),
('ob-012', 'clause-008', 'Annual license fee', 'Contract start', 'Annual', '2024-09-01', 50000, 'fixed', 'Finance', 'fulfilled', datetime('now'), NULL, datetime('now'), datetime('now')),
('ob-013', 'clause-009', 'Submit support tickets via portal', 'Contract start', 'As needed', NULL, 0, 'fixed', 'IT', 'active', NULL, NULL, datetime('now'), datetime('now')),
('ob-014', 'clause-006', 'Submit Q1 report', 'Quarter start', 'Q1 end', '2024-03-31', 1000, 'fixed', 'Consultant', 'fulfilled', datetime('now'), NULL, datetime('now'), datetime('now')),
('ob-015', 'clause-006', 'Submit Q2 report', 'Quarter start', 'Q2 end', '2024-06-30', 1000, 'fixed', 'Consultant', 'breached', NULL, datetime('now'), datetime('now'), datetime('now')),
('ob-016', 'clause-006', 'Submit Q3 report', 'Quarter start', 'Q3 end', '2024-09-30', 1000, 'fixed', 'Consultant', 'expired', NULL, NULL, datetime('now'), datetime('now'));

-- Insert some evaluation records
INSERT INTO obligation_evaluations (id, obligation_id, evaluated_at, status_before, status_after, notes) VALUES 
('eval-001', 'ob-001', datetime('now', '-4 months'), 'pending', 'active', 'Activated on contract start'),
('eval-002', 'ob-001', datetime('now', '-3 months'), 'active', 'fulfilled', 'Payment received'),
('eval-003', 'ob-004', datetime('now', '-2 months'), 'active', 'breached', 'Payment overdue by 30 days'),
('eval-004', 'ob-015', datetime('now', '-8 months'), 'active', 'breached', 'Report not submitted');

-- Insert audit logs
INSERT INTO audit_logs (id, organization_id, user_id, entity_type, entity_id, action, old_values, new_values, created_at) VALUES 
('audit-001', 'org-001', 'user-001', 'contract', 'contract-001', 'create', NULL, '{"title":"SaaS Service Agreement"}', datetime('now')),
('audit-002', 'org-001', 'user-001', 'contract', 'contract-002', 'create', NULL, '{"title":"Office Lease Agreement"}', datetime('now')),
('audit-003', 'org-001', 'user-002', 'obligation', 'ob-001', 'update', '{"status":"active"}', '{"status":"fulfilled"}', datetime('now'));
