package models

import (
	"math"
)

// PaginatedResponse represents a paginated list of posts
type PaginatedResponse struct {
	Posts       []Post `json:"posts"`
	CurrentPage int    `json:"currentPage"`
	TotalPages  int    `json:"totalPages"`
	TotalPosts  int    `json:"totalPosts"`
	PageSize    int    `json:"pageSize"`
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(posts []Post, currentPage, totalPosts, pageSize int) PaginatedResponse {
	totalPages := int(math.Ceil(float64(totalPosts) / float64(pageSize)))

	return PaginatedResponse{
		Posts:       posts,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		TotalPosts:  totalPosts,
		PageSize:    pageSize,
	}
}
