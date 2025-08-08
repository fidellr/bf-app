package models

import "time"

type Book struct {
	ID        int       `json:"id"`
	Title     string    `json:"title" validate:"required,min=1,max=200"`
	Author    string    `json:"author" validate:"required,min=1,max=100"`
	Published string    `json:"published" validate:"required,datetime=2006-01-02"`
	ISBN      string    `json:"isbn" validate:"required,isbn"`
	Pages     int       `json:"pages" validate:"required,min=5"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type (
	BookCreateRequest struct {
		Title     string `json:"title" validate:"required,min=1,max=200"`
		Author    string `json:"author" validate:"required,min=1,max=100"`
		Published string `json:"published" validate:"required,datetime=2006-01-02"`
		ISBN      string `json:"isbn" validate:"required,isbn"`
		Pages     int    `json:"pages" validate:"required,min=5,gt=0"`
	}

	BookUpdateRequest struct {
		Title     string `json:"title" validate:"omitempty,min=1,max=200"`
		Author    string `json:"author" validate:"omitempty,min=1,max=100"`
		Published string `json:"published" validate:"omitempty,datetime=2006-01-02"`
		ISBN      string `json:"isbn" validate:"omitempty,isbn"`
		Pages     int    `json:"pages" validate:"omitempty,min=5"`
	}
	BookGetByIDRequest struct {
		ID        int    `json:"id" validate:"required"`
		Title     string `json:"title" validate:"omitempty,min=1,max=200"`
		Author    string `json:"author" validate:"omitempty,min=1,max=100"`
		Published string `json:"published" validate:"omitempty,datetime=2006-01-02"`
		ISBN      string `json:"isbn" validate:"omitempty,isbn"`
		Pages     int    `json:"pages" validate:"omitempty,min=5"`
	}
	BookFetchAllRequest struct {
		ID        int    `json:"id" validate:"omitempty"`
		Title     string `json:"title" validate:"omitempty,min=1,max=200"`
		Author    string `json:"author" validate:"omitempty,min=1,max=100"`
		Published string `json:"published" validate:"omitempty,datetime=2006-01-02"`
		ISBN      string `json:"isbn" validate:"omitempty,isbn"`
		Pages     int    `json:"pages" validate:"omitempty,min=5"`
		PageSize  int    `json:"page_size" validate:"omitempty"`
	}
	BookDeleteRequest struct {
		ID int `json:"id" validate:"required"`
	}
)

type (
	BookListResponse struct {
		Data       []*Book `json:"data"`
		Page       int     `json:"page" example:"10"`
		PageSize   int     `json:"page_size" example:"10"`
		TotalPages int     `json:"total_pages" example:"10"`
		Limit      int     `json:"limit"`
	}
)
