package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend/api"
	"backend/internal/service"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// CompanyHandlers contains the HTTP handlers for company operations
type CompanyHandlers struct {
	service service.CompanyService
	logger  *zap.Logger
}

// NewCompanyHandlers creates a new company handlers instance
func NewCompanyHandlers(service service.CompanyService, logger *zap.Logger) *CompanyHandlers {
	return &CompanyHandlers{
		service: service,
		logger:  logger,
	}
}

// GetCompanies handles GET /api/v1/companies
func (h *CompanyHandlers) GetCompanies(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Getting companies list")

	// Parse query parameters
	params := api.GetCompaniesParams{}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			params.Limit = &limit
		} else {
			h.sendErrorResponse(w, http.StatusBadRequest, "Invalid limit parameter")
			return
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			params.Offset = &offset
		} else {
			h.sendErrorResponse(w, http.StatusBadRequest, "Invalid offset parameter")
			return
		}
	}

	if jurisdictionStr := r.URL.Query().Get("jurisdiction"); jurisdictionStr != "" {
		jurisdiction := api.GetCompaniesParamsJurisdiction(jurisdictionStr)
		params.Jurisdiction = &jurisdiction
	}

	// Call service
	response, err := h.service.ListCompanies(r.Context(), params)
	if err != nil {
		h.logger.Error("Failed to get companies", zap.Error(err))
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve companies")
		return
	}

	h.sendJSONResponse(w, http.StatusOK, response)
}

// GetCompanyByID handles GET /api/v1/companies/{id}
func (h *CompanyHandlers) GetCompanyByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	h.logger.Info("Getting company by ID", zap.String("id", idStr))

	// Parse UUID
	parsedID, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Error("Invalid UUID format", zap.String("id", idStr), zap.Error(err))
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid company ID format")
		return
	}
	id := openapi_types.UUID(parsedID)

	// Call service
	company, err := h.service.GetCompanyByID(r.Context(), id)
	if err != nil {
		if err.Error() == "company not found" {
			h.sendErrorResponse(w, http.StatusNotFound, "Company not found")
			return
		}
		h.logger.Error("Failed to get company", zap.Error(err))
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve company")
		return
	}

	h.sendJSONResponse(w, http.StatusOK, company)
}

// CreateCompany handles POST /api/v1/companies
func (h *CompanyHandlers) CreateCompany(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Creating new company")

	// Parse request body
	var req api.CreateCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Call service
	company, err := h.service.CreateCompany(r.Context(), req)
	if err != nil {
		h.logger.Error("Failed to create company", zap.Error(err))
		h.sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.sendJSONResponse(w, http.StatusCreated, company)
}

// DeleteCompany handles DELETE /api/v1/companies/{id}
func (h *CompanyHandlers) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	h.logger.Info("Deleting company", zap.String("id", idStr))

	// Parse UUID
	parsedID, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Error("Invalid UUID format", zap.String("id", idStr), zap.Error(err))
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid company ID format")
		return
	}
	id := openapi_types.UUID(parsedID)

	// Call service
	err = h.service.DeleteCompany(r.Context(), id)
	if err != nil {
		if err.Error() == "company not found" {
			h.sendErrorResponse(w, http.StatusNotFound, "Company not found")
			return
		}
		h.logger.Error("Failed to delete company", zap.Error(err))
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to delete company")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// sendJSONResponse sends a JSON response
func (h *CompanyHandlers) sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response", zap.Error(err))
	}
}

// sendErrorResponse sends an error response
func (h *CompanyHandlers) sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	response := api.ErrorResponse{
		Error: true,
		Msg:   message,
	}
	h.sendJSONResponse(w, statusCode, response)
}
