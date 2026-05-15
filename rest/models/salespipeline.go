package models

import "time"

// SalesPipelineRequest represents the request body for creating/updating a sales pipeline
type SalesPipelineRequest struct {
    ID                 *string `json:"id,omitempty"`
    Value              *int32 `json:"value,omitempty"`
    EstimatedCloseDate *string `json:"estimated_close_date,omitempty"` // Changed to string for flexible date parsing
    SpStageID          *int32  `json:"sp_stage_id,omitempty"`          // Added missing field
    PaxName            *string `json:"pax_name,omitempty"`
    PicName            *string `json:"pic_name,omitempty"`
    CreatedBy          *string `json:"created_by,omitempty"`
    UpdatedBy          *string `json:"updated_by,omitempty"`
}

// SalesPipelineResponse represents the response data for a sales pipeline
type SalesPipelineResponse struct {
    ID                 string     `json:"id"`
    Value              *int32    `json:"value,omitempty"`
    EstimatedCloseDate *time.Time `json:"estimated_close_date,omitempty"` // Fixed typo
    SpStageID          *int32     `json:"sp_stage_id,omitempty"`          // Added missing field
    PaxName            *string    `json:"pax_name,omitempty"`
    PicName            *string    `json:"pic_name,omitempty"`             // Fixed typo
    CreatedAt          time.Time  `json:"created_at"`
    CreatedBy          *string    `json:"created_by,omitempty"`
    UpdatedAt          *time.Time `json:"updated_at,omitempty"`
    UpdatedBy          *string    `json:"updated_by,omitempty"`
    IsActive           bool       `json:"is_active"`
}