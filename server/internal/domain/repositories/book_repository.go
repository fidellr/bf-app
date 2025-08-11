package repositories

import (
	"bf-api/internal/domain/models"
	"context"
)

type BookRepository interface {
	CreateBook(ctx context.Context, book *models.Book) error
	GetByBookID(ctx context.Context, id int) (*models.Book, error)
	FetchAllBook(ctx context.Context, page, pageSize int) ([]*models.Book, int, error)
	UpdateBook(ctx context.Context, book *models.Book) error
	DeleteBook(ctx context.Context, id int) error
}
