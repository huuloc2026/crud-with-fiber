package utils

import "gorm.io/gorm"

// PaginationResponse contains pagination data
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int64       `json:"totalPages"`
}

// Paginate performs a common pagination query
func Paginate(db *gorm.DB, model interface{}, page int, limit int) PaginationResponse {
	var total int64
	offset := (page - 1) * limit

	// Count total records
	db.Model(model).Count(&total)

	// Get data with limit & offset
	db.Limit(limit).Offset(offset).Find(model)

	// Calculate total pages
	totalPages := (total + int64(limit) - 1) / int64(limit)

	return PaginationResponse{
		Data:       model,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}
