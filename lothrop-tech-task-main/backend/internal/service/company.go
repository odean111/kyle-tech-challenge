package service

import (
	"context"
	"fmt"
	"strings"

	"backend/api"
	"backend/internal/repository"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

// CompanyService defines the business logic interface for company operations
type CompanyService interface {
	// ListCompanies retrieves companies with pagination and optional filtering
	ListCompanies(ctx context.Context, params api.GetCompaniesParams) (*api.CompaniesResponse, error)

	// GetCompanyByID retrieves a company by its ID
	GetCompanyByID(ctx context.Context, id openapi_types.UUID) (*api.Company, error)

	// CreateCompany creates a new company with validation
	CreateCompany(ctx context.Context, req api.CreateCompanyRequest) (*api.Company, error)

	// DeleteCompany removes a company by its ID
	DeleteCompany(ctx context.Context, id openapi_types.UUID) error
}

// companyService implements CompanyService
type companyService struct {
	repo repository.CompanyRepository
}

// NewCompanyService creates a new company service
func NewCompanyService(repo repository.CompanyRepository) CompanyService {
	return &companyService{repo: repo}
}

// ListCompanies retrieves companies with pagination and optional filtering
func (s *companyService) ListCompanies(ctx context.Context, params api.GetCompaniesParams) (*api.CompaniesResponse, error) {
	// Set default values
	limit := 20
	offset := 0

	if params.Limit != nil {
		if *params.Limit < 1 || *params.Limit > 100 {
			return nil, fmt.Errorf("limit must be between 1 and 100")
		}
		limit = *params.Limit
	}

	if params.Offset != nil {
		if *params.Offset < 0 {
			return nil, fmt.Errorf("offset must be non-negative")
		}
		offset = *params.Offset
	}

	var jurisdiction *string
	if params.Jurisdiction != nil {
		j := string(*params.Jurisdiction)
		jurisdiction = &j
	}

	companies, total, err := s.repo.GetAll(ctx, limit, offset, jurisdiction)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve companies: %w", err)
	}

	response := &api.CompaniesResponse{
		Companies: companies,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	}

	return response, nil
}

// GetCompanyByID retrieves a company by its ID
func (s *companyService) GetCompanyByID(ctx context.Context, id openapi_types.UUID) (*api.Company, error) {
	company, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve company: %w", err)
	}

	if company == nil {
		return nil, fmt.Errorf("company not found")
	}

	return company, nil
}

// CreateCompany creates a new company with validation
func (s *companyService) CreateCompany(ctx context.Context, req api.CreateCompanyRequest) (*api.Company, error) {
	// Validate required fields
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	company, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create company: %w", err)
	}

	return company, nil
}

// DeleteCompany removes a company by its ID
func (s *companyService) DeleteCompany(ctx context.Context, id openapi_types.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return fmt.Errorf("company not found")
		}
		return fmt.Errorf("failed to delete company: %w", err)
	}

	return nil
}

// validateCreateRequest validates the create company request
func (s *companyService) validateCreateRequest(req api.CreateCompanyRequest) error {
	// Validate company name
	if strings.TrimSpace(req.CompanyName) == "" {
		return fmt.Errorf("company name is required")
	}

	if len(req.CompanyName) > 255 {
		return fmt.Errorf("company name cannot exceed 255 characters")
	}

	// Validate company address
	if strings.TrimSpace(req.CompanyAddress) == "" {
		return fmt.Errorf("company address is required")
	}

	// Validate jurisdiction
	validJurisdictions := []string{"UK", "Singapore", "Caymens"}
	isValidJurisdiction := false
	for _, j := range validJurisdictions {
		if string(req.Jurisdiction) == j {
			isValidJurisdiction = true
			break
		}
	}
	if !isValidJurisdiction {
		return fmt.Errorf("invalid jurisdiction: must be one of %v", validJurisdictions)
	}

	// Validate optional fields
	if req.NumberOfDirectors != nil {
		if *req.NumberOfDirectors < 1 || *req.NumberOfDirectors > 100 {
			return fmt.Errorf("number of directors must be between 1 and 100")
		}
	}

	if req.NumberOfShareholders != nil {
		if *req.NumberOfShareholders < 1 || *req.NumberOfShareholders > 1000 {
			return fmt.Errorf("number of shareholders must be between 1 and 1000")
		}
	}

	return nil
}
