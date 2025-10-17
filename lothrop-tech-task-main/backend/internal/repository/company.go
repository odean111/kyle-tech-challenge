package repository

import (
	"context"
	"database/sql"
	"fmt"

	"backend/api"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

// CompanyRepository defines the interface for company data operations
type CompanyRepository interface {
	// GetAll retrieves companies with pagination and optional filtering
	GetAll(ctx context.Context, limit, offset int, jurisdiction *string) ([]api.Company, int, error)

	// GetByID retrieves a company by its ID
	GetByID(ctx context.Context, id openapi_types.UUID) (*api.Company, error)

	// Create creates a new company and returns the created company with generated ID and timestamps
	Create(ctx context.Context, req api.CreateCompanyRequest) (*api.Company, error)

	// Delete removes a company by its ID
	Delete(ctx context.Context, id openapi_types.UUID) error
}

// PostgresCompanyRepository implements CompanyRepository using PostgreSQL
type PostgresCompanyRepository struct {
	db *sql.DB
}

// NewPostgresCompanyRepository creates a new PostgreSQL company repository
func NewPostgresCompanyRepository(db *sql.DB) CompanyRepository {
	return &PostgresCompanyRepository{db: db}
}

// GetAll retrieves companies with pagination and optional filtering
func (r *PostgresCompanyRepository) GetAll(ctx context.Context, limit, offset int, jurisdiction *string) ([]api.Company, int, error) {
	var companies []api.Company
	var total int

	// First, get the total count
	countQuery := "SELECT COUNT(*) FROM companies"
	countArgs := []interface{}{}

	if jurisdiction != nil {
		countQuery += " WHERE jurisdiction = $1"
		countArgs = append(countArgs, *jurisdiction)
	}

	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Then get the companies with pagination
	query := `
		SELECT id, jurisdiction, company_name, company_address, nature_of_business, 
		       number_of_directors, number_of_shareholders, sec_code, date_created, date_updated
		FROM companies`

	args := []interface{}{}
	argIndex := 1

	if jurisdiction != nil {
		query += " WHERE jurisdiction = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *jurisdiction)
		argIndex++
	}

	query += " ORDER BY date_created DESC"
	query += " LIMIT $" + fmt.Sprintf("%d", argIndex) + " OFFSET $" + fmt.Sprintf("%d", argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var company api.Company
		err := rows.Scan(
			&company.Id,
			&company.Jurisdiction,
			&company.CompanyName,
			&company.CompanyAddress,
			&company.NatureOfBusiness,
			&company.NumberOfDirectors,
			&company.NumberOfShareholders,
			&company.SecCode,
			&company.DateCreated,
			&company.DateUpdated,
		)
		if err != nil {
			return nil, 0, err
		}
		companies = append(companies, company)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return companies, total, nil
}

// GetByID retrieves a company by its ID
func (r *PostgresCompanyRepository) GetByID(ctx context.Context, id openapi_types.UUID) (*api.Company, error) {
	query := `
		SELECT id, jurisdiction, company_name, company_address, nature_of_business, 
		       number_of_directors, number_of_shareholders, sec_code, date_created, date_updated
		FROM companies 
		WHERE id = $1`

	var company api.Company
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&company.Id,
		&company.Jurisdiction,
		&company.CompanyName,
		&company.CompanyAddress,
		&company.NatureOfBusiness,
		&company.NumberOfDirectors,
		&company.NumberOfShareholders,
		&company.SecCode,
		&company.DateCreated,
		&company.DateUpdated,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Company not found
		}
		return nil, err
	}

	return &company, nil
}

func (r *PostgresCompanyRepository) Create(ctx context.Context, req api.CreateCompanyRequest) (*api.Company, error) {
	query := `
		INSERT INTO companies (jurisdiction, company_name, company_address, nature_of_business, 
		                      number_of_directors, number_of_shareholders, sec_code)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, jurisdiction, company_name, company_address, nature_of_business, 
		          number_of_directors, number_of_shareholders, sec_code, date_created, date_updated`

	var company api.Company
	err := r.db.QueryRowContext(ctx, query,
		req.Jurisdiction,
		req.CompanyName,
		req.CompanyAddress,
		req.NatureOfBusiness,
		req.NumberOfDirectors,
		req.NumberOfShareholders,
		req.SecCode,
	).Scan(
		&company.Id,
		&company.Jurisdiction,
		&company.CompanyName,
		&company.CompanyAddress,
		&company.NatureOfBusiness,
		&company.NumberOfDirectors,
		&company.NumberOfShareholders,
		&company.SecCode,
		&company.DateCreated,
		&company.DateUpdated,
	)

	if err != nil {
		return nil, err
	}

	return &company, nil
}

// Delete removes a company by its ID
func (r *PostgresCompanyRepository) Delete(ctx context.Context, id openapi_types.UUID) error {
	query := "DELETE FROM companies WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // Company not found
	}

	return nil
}
