package controllers

import (
	"djong-reader-engine/rest/models"
	"djong-reader-engine/rest/services"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type SalesPipelineController struct {
	Service *services.SalesPipelineService
}

// NewSalesPipelineController creates a new instance of SalesPipelineController
func NewSalesPipelineController(service *services.SalesPipelineService) *SalesPipelineController {
	return &SalesPipelineController{Service: service}
}

// HandleSalesPipeline routes requests based on HTTP method
func (c *SalesPipelineController) HandleSalesPipeline(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		c.CreateSalesPipeline(w, r)
	case http.MethodPut:
		c.UpdateSalesPipeline(w, r)
	case http.MethodDelete:
		c.DeleteSalesPipeline(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.APIResponse{
			ResponseMessage: "internal error",
			ResponseCode:    "99",
		})
	}
}

// CreateSalesPipeline handles POST requests to create a new sales pipeline
func (c *SalesPipelineController) CreateSalesPipeline(w http.ResponseWriter, r *http.Request) {
	var req models.SalesPipelineRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.APIResponse{
			ResponseMessage: "invalid request body format",
			ResponseCode:    "99",
		})
		return
	}

	// Call service layer
	response, err := c.Service.CreateSalesPipeline(r.Context(), req)
	if err != nil {
		log.Printf("Error creating sales pipeline: %v", err)
		
		// Check if error is validation related (date format, etc.)
		if strings.Contains(err.Error(), "invalid date format") {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.APIResponse{
				ResponseMessage: err.Error(),
				ResponseCode:    "99",
			})
			return
		}
		
		// Internal server error
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.APIResponse{
			ResponseMessage: "failed to create sales pipeline",
			ResponseCode:    "99",
		})
		return
	}

	// Return success response with 201 status
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.NewSuccessResponse(response))
}

// UpdateSalesPipeline handles PUT requests to update an existing sales pipeline
func (c *SalesPipelineController) UpdateSalesPipeline(w http.ResponseWriter, r *http.Request) {
	var req models.SalesPipelineRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.APIResponse{
			ResponseMessage: "invalid request body format",
			ResponseCode:    "99",
		})
		return
	}

	// Validate ID is present
	if req.ID == nil || *req.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.APIResponse{
			ResponseMessage: "id is required for update",
			ResponseCode:    "99",
		})
		return
	}

	// Call service layer
	_, err := c.Service.UpdateSalesPipeline(r.Context(), req)
	if err != nil {
		log.Printf("Error updating sales pipeline: %v", err)
		
		// Check if error is validation related
		if strings.Contains(err.Error(), "invalid date format") || strings.Contains(err.Error(), "required") {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.APIResponse{
				ResponseMessage: err.Error(),
				ResponseCode:    "99",
			})
			return
		}
		
		if strings.Contains(err.Error(), "not found") {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.APIResponse{
				ResponseMessage: "sales pipeline not found",
				ResponseCode:    "99",
			})
			return
		}
		
		// Internal server error
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.APIResponse{
			ResponseMessage: "failed to update sales pipeline",
			ResponseCode:    "99",
		})
		return
	}

	// Return success response with request body as per requirement
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.NewSuccessResponse(req))
}

// DeleteSalesPipeline handles DELETE requests to soft delete a sales pipeline
func (c *SalesPipelineController) DeleteSalesPipeline(w http.ResponseWriter, r *http.Request) {
	// Get ID from query parameter
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.APIResponse{
			ResponseMessage: "id parameter is required",
			ResponseCode:    "99",
		})
		return
	}

	// Call service layer
	err := c.Service.DeleteSalesPipeline(r.Context(), id)
	if err != nil {
		log.Printf("Error deleting sales pipeline: %v", err)
		
		if strings.Contains(err.Error(), "not found") {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.APIResponse{
				ResponseMessage: "sales pipeline not found or already deleted",
				ResponseCode:    "99",
			})
			return
		}
		
		// Internal server error
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.APIResponse{
			ResponseMessage: "failed to delete sales pipeline",
			ResponseCode:    "99",
		})
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.NewSuccessResponse(map[string]string{"id": id}))
}
