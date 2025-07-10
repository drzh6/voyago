package handler

import (
	"time"
)

type TripAddRequestBody struct {
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	Status        string    `json:"status"`
	IsPublic      bool      `json:"is_public"`
	CoverImageUrl string    `json:"cover_image_url"`
}

type TripRequestBody struct {
	Id            string
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	OwnerId       string    `json:"owner_id"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	Status        string    `json:"status"`
	IsPublic      bool      `json:"is_public"`
	InviteCode    string    `json:"invite_code"`
	CoverImageUrl string    `json:"cover_image_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
