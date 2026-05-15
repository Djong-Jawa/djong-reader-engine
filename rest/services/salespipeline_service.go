package services

import (
	"context"
	"djong-reader-engine/rest/models"
	"fmt"
	"log"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SalesPipelineService struct {
	DB *pgxpool.Pool
}

// NewSalesPipelineService creates a new instance of SalesPipelineService
func NewSalesPipelineService(db *pgxpool.Pool) *SalesPipelineService {
	return &SalesPipelineService{DB: db}
}

// CreateSalesPipeline inserts a new sales pipeline record
func (s *SalesPipelineService) CreateSalesPipeline(ctx context.Context, req models.SalesPipelineRequest) (*models.SalesPipelineResponse, error) {
	// Generate new UUID for ID
	createdAt := time.Now()
	isActive := true

	// Parse estimated_close_date if provided
	var estimatedCloseDate *time.Time
	if req.EstimatedCloseDate != nil && *req.EstimatedCloseDate != "" {
		// Try multiple date formats
		formats := []string{
			"2006-01-02",
			"2006-01-02T15:04:05Z07:00",
			time.RFC3339,
		}
		
		var parsedTime time.Time
		var parseErr error
		for _, format := range formats {
			parsedTime, parseErr = time.Parse(format, *req.EstimatedCloseDate)
			if parseErr == nil {
				estimatedCloseDate = &parsedTime
				break
			}
		}
		
		if parseErr != nil {
			return nil, fmt.Errorf("invalid date format for estimated_close_date, expected YYYY-MM-DD")
		}
	}

	query := `
		INSERT INTO mst_sales_pipeline 
		(value, estimated_close_date, sp_stage_id, pax_name, pic_name, created_at, created_by, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING value, estimated_close_date, sp_stage_id, pax_name, pic_name, created_at, created_by, is_active
	`

	var response models.SalesPipelineResponse
	err := s.DB.QueryRow(ctx, query,
		req.Value,
		estimatedCloseDate,
		req.SpStageID,
		req.PaxName,
		req.PicName,
		createdAt,
		req.CreatedBy,
		isActive,
	).Scan(
		&response.Value,
		&response.EstimatedCloseDate,
		&response.SpStageID,
		&response.PaxName,
		&response.PicName,
		&response.CreatedAt,
		&response.CreatedBy,
		&response.IsActive,
	)

	if err != nil {
		log.Printf("Error creating sales pipeline: %v", err)
		// Return user-friendly error message
		return nil, fmt.Errorf("failed to create sales pipeline")
	}

	return &response, nil
}

// UpdateSalesPipeline updates an existing sales pipeline record
func (s *SalesPipelineService) UpdateSalesPipeline(ctx context.Context, req models.SalesPipelineRequest) (*models.SalesPipelineResponse, error) {
	if req.ID == nil || *req.ID == "" {
		return nil, fmt.Errorf("id is required for update")
	}

	updatedAt := time.Now()

	// Parse estimated_close_date if provided
	var estimatedCloseDate *time.Time
	if req.EstimatedCloseDate != nil && *req.EstimatedCloseDate != "" {
		// Try multiple date formats
		formats := []string{
			"2006-01-02",
			"2006-01-02T15:04:05Z07:00",
			time.RFC3339,
		}
		
		var parsedTime time.Time
		var parseErr error
		for _, format := range formats {
			parsedTime, parseErr = time.Parse(format, *req.EstimatedCloseDate)
			if parseErr == nil {
				estimatedCloseDate = &parsedTime
				break
			}
		}
		
		if parseErr != nil {
			return nil, fmt.Errorf("invalid date format for estimated_close_date, expected YYYY-MM-DD")
		}
	}

	query := `
		UPDATE mst_sales_pipeline 
		SET value = $1, 
		    estimated_close_date = $2, 
		    sp_stage_id = $3,
		    pax_name = $4, 
		    pic_name = $5, 
		    updated_at = $6, 
		    updated_by = $7
		WHERE id = $8 AND is_active = true
		RETURNING id, value, estimated_close_date, sp_stage_id, pax_name, pic_name, created_at, created_by, updated_at, updated_by, is_active
	`

	var response models.SalesPipelineResponse
	err := s.DB.QueryRow(ctx, query,
		req.Value,
		estimatedCloseDate,
		req.SpStageID,
		req.PaxName,
		req.PicName,
		updatedAt,
		req.UpdatedBy,
		*req.ID,
	).Scan(
		&response.ID,
		&response.Value,
		&response.EstimatedCloseDate,
		&response.SpStageID,
		&response.PaxName,
		&response.PicName,
		&response.CreatedAt,
		&response.CreatedBy,
		&response.UpdatedAt,
		&response.UpdatedBy,
		&response.IsActive,
	)

	if err != nil {
		log.Printf("Error updating sales pipeline: %v", err)
		// Return user-friendly error message
		return nil, fmt.Errorf("failed to update sales pipeline")
	}

	return &response, nil
}

// DeleteSalesPipeline soft deletes a sales pipeline record by setting is_active to false
func (s *SalesPipelineService) DeleteSalesPipeline(ctx context.Context, id string) error {
	updatedAt := time.Now()

	query := `
		UPDATE mst_sales_pipeline 
		SET is_active = false, updated_at = $1
		WHERE id = $2 AND is_active = true
	`

	commandTag, err := s.DB.Exec(ctx, query, updatedAt, id)
	if err != nil {
		log.Printf("Error deleting sales pipeline: %v", err)
		return fmt.Errorf("failed to delete sales pipeline")
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("sales pipeline not found or already deleted")
	}

	return nil
}
