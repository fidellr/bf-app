package services

import (
	"bf-api/internal/domain/models"
	"bf-api/internal/domain/repositories"
	"context"
	"errors"
	"fmt"
)

type BookService struct {
	repo repositories.BookRepository
}

func NewBookService(repo repositories.BookRepository) *BookService {
	return &BookService{
		repo: repo,
	}
}

func (s *BookService) CreateBook(ctx context.Context, req *models.BookCreateRequest) (*models.Book, error) {

	if err := validateBookCreateRequest(req); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	book := &models.Book{
		Title:     req.Title,
		Author:    req.Author,
		Published: req.Published,
		ISBN:      req.ISBN,
		Pages:     req.Pages,
	}

	if book.Pages < 5 {
		return nil, fmt.Errorf("%w: book must have atleast 5 pages", ErrInvalidInput)
	}

	if err := s.repo.CreateBook(ctx, book); err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}

	return book, nil
}

func (s *BookService) GetByBookID(ctx context.Context, id int) (*models.Book, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: invalid book ID", repositories.ErrInvalidData)
	}

	book, err := s.repo.GetByBookID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrBookNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repository error: %w", err)
	}

	return book, nil
}

func (s *BookService) FetchAllBook(ctx context.Context, page, pageSize int) ([]*models.Book, int, error) {

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	books, total, err := s.repo.FetchAllBook(ctx, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("repository error: %w", err)
	}

	return books, total, nil

}

func (s *BookService) UpdateBook(ctx context.Context, id int, req *models.BookUpdateRequest) (*models.Book, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: invalid book ID", ErrInvalidInput)
	}
	if err := validateBookUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	book, err := s.repo.GetByBookID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrBookNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repository error: %w", err)
	}

	if req.Title != "" {
		book.Title = req.Title
	}
	if req.Author != "" {
		book.Author = req.Author
	}
	if req.ISBN != "" {
		book.ISBN = req.ISBN
	}
	if req.Pages <= 0 {
		book.Pages = req.Pages
	}

	if err := s.repo.UpdateBook(ctx, book); err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}

	return book, nil
}

func (s *BookService) DeleteBook(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("%w: invalid book ID", ErrInvalidInput)
	}

	// biz rule: check if book can be deleted
	_, err := s.repo.GetByBookID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrBookNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("repository error: %w", err)
	}

	if err := s.repo.DeleteBook(ctx, id); err != nil {
		return fmt.Errorf("repository error: %w", err)
	}

	return nil
}

// helper functions
func validateBookCreateRequest(req *models.BookCreateRequest) error {
	if req.Title == "" {
		return errors.New("title is required")
	}
	if len(req.Title) > 200 {
		return errors.New("title too long")
	}
	return nil
}

func validateBookUpdateRequest(req *models.BookUpdateRequest) error {
	if req.Title != "" && len(req.Title) > 200 {
		return errors.New("title too long")
	}
	return nil
}
