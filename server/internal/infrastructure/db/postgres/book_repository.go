package postgres

import (
	"bf-api/internal/domain/models"
	"bf-api/internal/domain/repositories"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRepository struct {
	pool *pgxpool.Pool
}

func NewBookRepository(pool *pgxpool.Pool) repositories.BookRepository {
	return &BookRepository{pool: pool}
}

func (r *BookRepository) CreateBook(ctx context.Context, book *models.Book) error {
	query := `
		INSERT INTO books (
			title,
			author,
			published,
			isbn,
			pages,
			created_at,
			updated_at
		) VALUES (
			$1, $2, $3, $4, $5, NOW(), NOW()
		)
		RETURNING id, created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		book.Title,
		book.Author,
		book.Published,
		book.ISBN,
		book.Pages,
	).Scan(
		&book.ID,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				return repositories.ErrDuplicateISBN
			case "23503": // foreign_key_violation
				return fmt.Errorf("%w: %s", repositories.ErrInvalidReference, pgErr.Message)
			}
		}
		return fmt.Errorf("failed to create book: %w", err)
	}

	return nil
}

func (r *BookRepository) GetByBookID(ctx context.Context, id int) (*models.Book, error) {
	query := `
	SELECT
		id, title, author, published, isbn, pages, created_at, updated_at
	FROM books
	WHERE id = $1
	`
	var book models.Book
	var pubDate time.Time
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&pubDate,
		&book.ISBN,
		&book.Pages,
		&book.CreatedAt,
		&book.UpdatedAt,
	)
	book.Published = pubDate.Format("2006-01-02")

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repositories.ErrBookNotFound
		}
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	return &book, nil
}

func (r *BookRepository) FetchAllBook(ctx context.Context, page, pageSize int) ([]*models.Book, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM books`
	err := r.pool.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count books: %s", err.Error())
	}

	query := `
		SELECT
			id, title, author, published, isbn, pages, created_at, updated_at
		FROM books
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	offset := (page - 1) * pageSize
	rows, err := r.pool.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch books: %w", err)
	}
	defer rows.Close()

	var books []*models.Book
	for rows.Next() {
		var book models.Book
		var pubDate time.Time
		if err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&pubDate,
			&book.ISBN,
			&book.Pages,
			&book.CreatedAt,
			&book.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan book: %w", err)
		}
		book.Published = pubDate.Format("2006-01-02")
		books = append(books, &book)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return books, total, nil
}

func (r *BookRepository) UpdateBook(ctx context.Context, book *models.Book) error {
	query := `
		UPDATE books
		SET
			title = $1,
			author = $2,
			published = $3,
			isbn = $4,
			pages = $5,
			updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		book.Title,
		book.Author,
		book.Published,
		book.ISBN,
		book.Pages,
		book.ID,
	).Scan(&book.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repositories.ErrBookNotFound
		}
		return fmt.Errorf("failed to update book: %w", err)
	}

	return nil
}
func (r *BookRepository) DeleteBook(ctx context.Context, id int) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	result, err := tx.Exec(ctx, "DELETE FROM books WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}

	if result.RowsAffected() == 0 {
		return repositories.ErrBookNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
