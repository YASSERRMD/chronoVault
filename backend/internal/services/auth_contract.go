package services

import (
	"errors"
	"time"

	"chronovault/internal/models"
	"chronovault/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      *repository.Repository
	jwtSecret string
}

func NewAuthService(repo *repository.Repository, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret}
}

type Claims struct {
	UserID         string `json:"user_id"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	OrganizationID string `json:"organization_id"`
	jwt.RegisteredClaims
}

func (s *AuthService) Register(email, password, fullName, orgID, role string) (*models.User, string, error) {
	existing, _ := s.repo.GetUserByEmail(email)
	if existing != nil {
		return nil, "", errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{
		Email:          email,
		PasswordHash:   string(hashedPassword),
		FullName:       fullName,
		OrganizationID: orgID,
		Role:           role,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, "", err
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := Claims{
		UserID:         user.ID,
		Email:          user.Email,
		Role:           user.Role,
		OrganizationID: user.OrganizationID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

type ContractService struct {
	repo *repository.Repository
}

func NewContractService(repo *repository.Repository) *ContractService {
	return &ContractService{repo: repo}
}

func (s *ContractService) CreateOrganization(org *models.Organization) error {
	if org.Name == "" {
		return errors.New("organization name is required")
	}
	return s.repo.CreateOrganization(org)
}

func (s *ContractService) GetOrganization(id string) (*models.Organization, error) {
	return s.repo.GetOrganization(id)
}

func (s *ContractService) ListOrganizations() ([]models.Organization, error) {
	return s.repo.ListOrganizations()
}

func (s *ContractService) UpdateOrganization(org *models.Organization) error {
	return s.repo.UpdateOrganization(org)
}

func (s *ContractService) CreateContract(contract *models.Contract, userID string) error {
	if contract.Title == "" || contract.Counterparty == "" {
		return errors.New("title and counterparty are required")
	}
	return s.repo.CreateContract(contract, userID)
}

func (s *ContractService) GetContract(id, orgID string) (*models.Contract, error) {
	return s.repo.GetContract(id, orgID)
}

func (s *ContractService) ListContracts(orgID, status string, page, limit int) ([]models.Contract, int, error) {
	return s.repo.ListContracts(orgID, status, page, limit)
}

func (s *ContractService) UpdateContract(contract *models.Contract, userID string) error {
	return s.repo.UpdateContract(contract, userID)
}

func (s *ContractService) DeleteContract(id, orgID string) error {
	return s.repo.DeleteContract(id, orgID)
}

func (s *ContractService) GetContractVersions(contractID, orgID string) ([]models.ContractVersion, error) {
	return s.repo.GetContractVersions(contractID, orgID)
}

func (s *ContractService) GetContractClauses(contractID, orgID string) ([]models.Clause, error) {
	return s.repo.GetContractClauses(contractID, orgID)
}

func (s *ContractService) CreateClause(clause *models.Clause) error {
	if clause.Title == "" || clause.ContractID == "" {
		return errors.New("title and contract_id are required")
	}
	return s.repo.CreateClause(clause)
}

func (s *ContractService) GetClause(id, orgID string) (*models.Clause, error) {
	return s.repo.GetClause(id, orgID)
}

func (s *ContractService) ListClauses(contractID, orgID string) ([]models.Clause, error) {
	return s.repo.ListClauses(contractID, orgID)
}

func (s *ContractService) UpdateClause(clause *models.Clause) error {
	return s.repo.UpdateClause(clause)
}

func (s *ContractService) DeleteClause(id, orgID string) error {
	return s.repo.DeleteClause(id, orgID)
}
