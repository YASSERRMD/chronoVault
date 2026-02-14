package handlers

import (
	"net/http"
	"strconv"

	"chronovault/internal/models"
	"chronovault/internal/services"
	"chronovault/internal/websocket"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	authService       *services.AuthService
	contractService   *services.ContractService
	obligationService *services.ObligationService
	reportService     *services.ReportService
	auditService      *services.AuditService
	wsHub             *websocket.Hub
}

func New(auth *services.AuthService, contract *services.ContractService, obligation *services.ObligationService, report *services.ReportService, audit *services.AuditService, ws *websocket.Hub) *Handler {
	return &Handler{
		authService:       auth,
		contractService:   contract,
		obligationService: obligation,
		reportService:     report,
		auditService:      audit,
		wsHub:             ws,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	OrgID    string `json:"organization_id" binding:"required"`
	Role     string `json:"role"`
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":              user.ID,
			"email":           user.Email,
			"full_name":       user.FullName,
			"role":            user.Role,
			"organization_id": user.OrganizationID,
		},
	})
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role := req.Role
	if role == "" {
		role = "viewer"
	}

	user, token, err := h.authService.Register(req.Email, req.Password, req.FullName, req.OrgID, role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user": gin.H{
			"id":              user.ID,
			"email":           user.Email,
			"full_name":       user.FullName,
			"role":            user.Role,
			"organization_id": user.OrganizationID,
		},
	})
}

func (h *Handler) GetCurrentUser(c *gin.Context) {
	userID := c.GetString("user_id")
	orgID := c.GetString("organization_id")

	c.JSON(http.StatusOK, gin.H{
		"user_id":         userID,
		"organization_id": orgID,
		"role":            c.GetString("role"),
		"email":           c.GetString("email"),
	})
}

func (h *Handler) ListOrganizations(c *gin.Context) {
	orgs, err := h.contractService.ListOrganizations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orgs)
}

func (h *Handler) GetOrganization(c *gin.Context) {
	id := c.Param("id")
	org, err := h.contractService.GetOrganization(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}
	c.JSON(http.StatusOK, org)
}

type OrgRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *Handler) CreateOrganization(c *gin.Context) {
	var req OrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org := &models.Organization{Name: req.Name}
	if err := h.contractService.CreateOrganization(org); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, org)
}

func (h *Handler) UpdateOrganization(c *gin.Context) {
	id := c.Param("id")
	var req OrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org := &models.Organization{ID: id, Name: req.Name}
	if err := h.contractService.UpdateOrganization(org); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, org)
}

func (h *Handler) ListContracts(c *gin.Context) {
	orgID := c.GetString("organization_id")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	contracts, total, err := h.contractService.ListContracts(orgID, status, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  contracts,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *Handler) GetContract(c *gin.Context) {
	orgID := c.GetString("organization_id")
	id := c.Param("id")

	contract, err := h.contractService.GetContract(id, orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
		return
	}
	c.JSON(http.StatusOK, contract)
}

type ContractRequest struct {
	Title        string  `json:"title" binding:"required"`
	Counterparty string  `json:"counterparty" binding:"required"`
	StartDate    string  `json:"start_date" binding:"required"`
	EndDate      string  `json:"end_date" binding:"required"`
	Value        float64 `json:"value"`
	Status       string  `json:"status"`
}

func (h *Handler) CreateContract(c *gin.Context) {
	orgID := c.GetString("organization_id")
	userID := c.GetString("user_id")

	var req ContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := req.Status
	if status == "" {
		status = "draft"
	}

	contract := &models.Contract{
		OrganizationID: orgID,
		Title:          req.Title,
		Counterparty:   req.Counterparty,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		Value:          req.Value,
		Status:         status,
	}

	if err := h.contractService.CreateContract(contract, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, contract)
}

func (h *Handler) UpdateContract(c *gin.Context) {
	orgID := c.GetString("organization_id")
	userID := c.GetString("user_id")
	id := c.Param("id")

	var req ContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contract := &models.Contract{
		ID:             id,
		OrganizationID: orgID,
		Title:          req.Title,
		Counterparty:   req.Counterparty,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		Value:          req.Value,
		Status:         req.Status,
	}

	if err := h.contractService.UpdateContract(contract, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contract)
}

func (h *Handler) DeleteContract(c *gin.Context) {
	orgID := c.GetString("organization_id")
	id := c.Param("id")

	if err := h.contractService.DeleteContract(id, orgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contract deleted"})
}

func (h *Handler) GetContractVersions(c *gin.Context) {
	orgID := c.GetString("organization_id")
	id := c.Param("id")

	versions, err := h.contractService.GetContractVersions(id, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, versions)
}

func (h *Handler) GetContractClauses(c *gin.Context) {
	orgID := c.GetString("organization_id")
	id := c.Param("id")

	clauses, err := h.contractService.GetContractClauses(id, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, clauses)
}

func (h *Handler) ListClauses(c *gin.Context) {
	orgID := c.GetString("organization_id")
	contractID := c.Query("contract_id")

	clauses, err := h.contractService.ListClauses(contractID, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, clauses)
}

type ClauseRequest struct {
	ContractID  string `json:"contract_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	OrderIndex  int    `json:"order_index"`
}

func (h *Handler) CreateClause(c *gin.Context) {
	var req ClauseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clause := &models.Clause{
		ContractID:  req.ContractID,
		Title:       req.Title,
		Description: req.Description,
		OrderIndex:  req.OrderIndex,
	}

	if err := h.contractService.CreateClause(clause); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, clause)
}

func (h *Handler) UpdateClause(c *gin.Context) {
	id := c.Param("id")
	_ = c.GetString("organization_id")

	var req ClauseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clause := &models.Clause{
		ID:          id,
		ContractID:  req.ContractID,
		Title:       req.Title,
		Description: req.Description,
		OrderIndex:  req.OrderIndex,
	}

	if err := h.contractService.UpdateClause(clause); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, clause)
}

func (h *Handler) DeleteClause(c *gin.Context) {
	orgID := c.GetString("organization_id")
	id := c.Param("id")

	if err := h.contractService.DeleteClause(id, orgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Clause deleted"})
}

func (h *Handler) ListObligations(c *gin.Context) {
	orgID := c.GetString("organization_id")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	obligations, total, err := h.obligationService.ListObligations(orgID, status, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  obligations,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *Handler) GetObligation(c *gin.Context) {
	orgID := c.GetString("organization_id")
	id := c.Param("id")

	obligation, err := h.obligationService.GetObligation(id, orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Obligation not found"})
		return
	}
	c.JSON(http.StatusOK, obligation)
}

type ObligationRequest struct {
	ClauseID            string  `json:"clause_id" binding:"required"`
	Description         string  `json:"description" binding:"required"`
	ActivationCondition string  `json:"activation_condition"`
	DueDateRule         string  `json:"due_date_rule"`
	DueDate             string  `json:"due_date"`
	PenaltyAmount       float64 `json:"penalty_amount"`
	PenaltyType         string  `json:"penalty_type"`
	ResponsibleParty    string  `json:"responsible_party"`
	DependsOnID         *string `json:"depends_on_id"`
}

func (h *Handler) CreateObligation(c *gin.Context) {
	var req ObligationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	penaltyType := req.PenaltyType
	if penaltyType == "" {
		penaltyType = "fixed"
	}

	obligation := &models.Obligation{
		ClauseID:            req.ClauseID,
		Description:         req.Description,
		ActivationCondition: req.ActivationCondition,
		DueDateRule:         req.DueDateRule,
		DueDate:             req.DueDate,
		PenaltyAmount:       req.PenaltyAmount,
		PenaltyType:         penaltyType,
		ResponsibleParty:    req.ResponsibleParty,
		DependsOnID:         req.DependsOnID,
	}

	if err := h.obligationService.CreateObligation(obligation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, obligation)
}

func (h *Handler) UpdateObligation(c *gin.Context) {
	id := c.Param("id")
	_ = c.GetString("organization_id")

	var req ObligationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	obligation := &models.Obligation{
		ID:                  id,
		ClauseID:            req.ClauseID,
		Description:         req.Description,
		ActivationCondition: req.ActivationCondition,
		DueDateRule:         req.DueDateRule,
		DueDate:             req.DueDate,
		PenaltyAmount:       req.PenaltyAmount,
		PenaltyType:         req.PenaltyType,
		ResponsibleParty:    req.ResponsibleParty,
		DependsOnID:         req.DependsOnID,
	}

	if err := h.obligationService.UpdateObligation(obligation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, obligation)
}

func (h *Handler) DeleteObligation(c *gin.Context) {
	orgID := c.GetString("organization_id")
	id := c.Param("id")

	if err := h.obligationService.DeleteObligation(id, orgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Obligation deleted"})
}

func (h *Handler) FulfillObligation(c *gin.Context) {
	orgID := c.GetString("organization_id")
	id := c.Param("id")

	if err := h.obligationService.FulfillObligation(id, orgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Obligation fulfilled"})
}

func (h *Handler) GetObligationHistory(c *gin.Context) {
	orgID := c.GetString("organization_id")
	id := c.Param("id")

	history, err := h.obligationService.GetObligationHistory(id, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

func (h *Handler) GetFinancialSummary(c *gin.Context) {
	orgID := c.GetString("organization_id")

	summary, err := h.reportService.GetFinancialSummary(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (h *Handler) GetPenaltyTracking(c *gin.Context) {
	orgID := c.GetString("organization_id")

	penalties, err := h.reportService.GetPenaltyTracking(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, penalties)
}

func (h *Handler) GetRiskExposure(c *gin.Context) {
	orgID := c.GetString("organization_id")

	risk, err := h.reportService.GetRiskExposure(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, risk)
}

func (h *Handler) GetYearlyImpact(c *gin.Context) {
	orgID := c.GetString("organization_id")
	year := c.DefaultQuery("year", "2024")

	impact, err := h.reportService.GetYearlyImpact(orgID, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, impact)
}

func (h *Handler) ListAuditLogs(c *gin.Context) {
	orgID := c.GetString("organization_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	logs, total, err := h.auditService.ListAuditLogs(orgID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  logs,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *Handler) GetEntityAuditLogs(c *gin.Context) {
	orgID := c.GetString("organization_id")
	entityType := c.Param("entity_type")
	entityID := c.Param("entity_id")

	logs, err := h.auditService.GetEntityAuditLogs(entityType, entityID, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}
