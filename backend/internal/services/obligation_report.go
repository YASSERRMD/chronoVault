package services

import (
	"chronovault/internal/models"
	"chronovault/internal/repository"
	"chronovault/internal/websocket"
	"encoding/json"
	"time"
)

type ObligationService struct {
	repo  *repository.Repository
	wsHub *websocket.Hub
}

func NewObligationService(repo *repository.Repository, wsHub *websocket.Hub) *ObligationService {
	return &ObligationService{repo: repo, wsHub: wsHub}
}

func (s *ObligationService) CreateObligation(obligation *models.Obligation) error {
	if obligation.Description == "" || obligation.ClauseID == "" {
		return &ValidationError{Message: "description and clause_id are required"}
	}
	return s.repo.CreateObligation(obligation)
}

func (s *ObligationService) GetObligation(id, orgID string) (*models.Obligation, error) {
	return s.repo.GetObligation(id, orgID)
}

func (s *ObligationService) ListObligations(orgID, status string, page, limit int) ([]models.Obligation, int, error) {
	return s.repo.ListObligations(orgID, status, page, limit)
}

func (s *ObligationService) UpdateObligation(obligation *models.Obligation) error {
	return s.repo.UpdateObligation(obligation)
}

func (s *ObligationService) DeleteObligation(id, orgID string) error {
	return s.repo.DeleteObligation(id, orgID)
}

func (s *ObligationService) FulfillObligation(id, orgID string) error {
	err := s.repo.FulfillObligation(id, orgID)
	if err != nil {
		return err
	}

	obligation, _ := s.repo.GetObligation(id, orgID)
	if obligation != nil {
		s.notifyWebsocket("obligation_fulfilled", obligation)
	}

	return nil
}

func (s *ObligationService) GetObligationHistory(id, orgID string) ([]models.ObligationEvaluation, error) {
	return s.repo.GetObligationHistory(id, orgID)
}

func (s *ObligationService) EvaluateObligations(orgID string) error {
	obligations, err := s.repo.GetActiveObligationsForEvaluation(orgID)
	if err != nil {
		return err
	}

	for _, obligation := range obligations {
		s.evaluateObligation(&obligation)
	}

	return nil
}

func (s *ObligationService) evaluateObligation(obligation *models.Obligation) {
	now := time.Now()
	oldStatus := obligation.Status
	newStatus := oldStatus

	if obligation.DependsOnID != nil && *obligation.DependsOnID != "" {
		dep, err := s.repo.GetObligationDependencies(*obligation.DependsOnID)
		if err == nil && dep.Status != "fulfilled" {
			if dep.Status == "breached" {
				newStatus = "expired"
				s.repo.UpdateObligationStatus(obligation.ID, newStatus)
				s.createEvaluation(obligation.ID, oldStatus, newStatus, "Auto-expired due to dependency breach")
				s.notifyWebsocket("obligation_expired", obligation)
				return
			}
			newStatus = "pending"
			if oldStatus != newStatus {
				s.repo.UpdateObligationStatus(obligation.ID, newStatus)
				s.createEvaluation(obligation.ID, oldStatus, newStatus, "Waiting for dependency")
			}
			return
		}
	}

	if obligation.DueDate != "" {
		dueDate, err := time.Parse("2006-01-02", obligation.DueDate)
		if err == nil {
			if now.After(dueDate) && oldStatus == "active" {
				newStatus = "breached"
				s.repo.UpdateObligationStatus(obligation.ID, newStatus)
				s.createEvaluation(obligation.ID, oldStatus, newStatus, "Exceeded due date")
				s.notifyWebsocket("obligation_breached", obligation)
				s.notifyWebsocket("penalty_applied", map[string]interface{}{
					"obligation_id":  obligation.ID,
					"penalty_amount": obligation.PenaltyAmount,
				})
				return
			}

			if now.Before(dueDate) && oldStatus == "pending" {
				contractEndDate, _ := time.Parse("2006-01-02", "2025-12-31")
				if now.After(contractEndDate) {
					newStatus = "expired"
				} else {
					newStatus = "active"
				}
				if oldStatus != newStatus {
					s.repo.UpdateObligationStatus(obligation.ID, newStatus)
					s.createEvaluation(obligation.ID, oldStatus, newStatus, "Activated based on due date")
					s.notifyWebsocket("obligation_activated", obligation)
				}
			}
		}
	}
}

func (s *ObligationService) createEvaluation(obligationID, statusBefore, statusAfter, notes string) {
	eval := &models.ObligationEvaluation{
		ObligationID: obligationID,
		StatusBefore: statusBefore,
		StatusAfter:  statusAfter,
		Notes:        notes,
	}
	s.repo.CreateObligationEvaluation(eval)
}

func (s *ObligationService) notifyWebsocket(event string, data interface{}) {
	if s.wsHub == nil {
		return
	}
	message := websocket.Message{
		Type: event,
		Data: data,
	}
	jsonMsg, _ := json.Marshal(message)
	s.wsHub.BroadcastMessage(jsonMsg)
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type ReportService struct {
	repo *repository.Repository
}

func NewReportService(repo *repository.Repository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetFinancialSummary(orgID string) (map[string]interface{}, error) {
	return s.repo.GetFinancialSummary(orgID)
}

func (s *ReportService) GetPenaltyTracking(orgID string) ([]map[string]interface{}, error) {
	return s.repo.GetPenaltyTracking(orgID)
}

func (s *ReportService) GetRiskExposure(orgID string) ([]map[string]interface{}, error) {
	return s.repo.GetRiskExposure(orgID)
}

func (s *ReportService) GetYearlyImpact(orgID, year string) ([]map[string]interface{}, error) {
	return s.repo.GetYearlyImpact(orgID, year)
}

type AuditService struct {
	repo *repository.Repository
}

func NewAuditService(repo *repository.Repository) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) ListAuditLogs(orgID string, page, limit int) ([]models.AuditLog, int, error) {
	return s.repo.ListAuditLogs(orgID, page, limit)
}

func (s *AuditService) GetEntityAuditLogs(entityType, entityID, orgID string) ([]models.AuditLog, error) {
	return s.repo.GetEntityAuditLogs(entityType, entityID, orgID)
}

func (s *AuditService) LogAction(orgID, userID, entityType, entityID, action string, oldVals, newVals interface{}) error {
	return s.repo.LogAudit(orgID, userID, entityType, entityID, action, oldVals, newVals)
}
