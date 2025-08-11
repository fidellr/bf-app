package handlers

import (
	"bf-api/internal/app/middleware"
	"bf-api/internal/domain/models"
	"bf-api/internal/domain/services"

	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type (
	BookHandler struct {
		service   *services.BookService
		validator *validator.Validate
		logger    *zap.Logger
	}

	ValidationError struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}
)

func NewBookHandler(s *services.BookService, logger *zap.Logger) *BookHandler {

	return &BookHandler{
		service:   s,
		validator: validator.New(),
		logger:    logger,
	}
}

func (h *BookHandler) LogRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		h.logger.Info("incoming request",
			zap.String("method", c.Request().Method),
			zap.String("path", c.Path()),
			zap.String("ip", c.RealIP()),
			zap.String("user_agent", c.Request().UserAgent()),
		)

		err := next(c)

		h.logger.Info("request completed",
			zap.String("method", c.Request().Method),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().Status),
			zap.Duration("duration", time.Since(start)),
		)

		return err
	}
}

// CreateBook godoc
// @Summary Create a new book
// @Description Add a new book to the store
// @Tags books
// @Accept json
// @Produce json
// @Param book body models.BookCreateRequest true "Book data"
// @Success 201 {object} models.Book
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 409 {object} handlers.ErrorResponse
// @Failure 429 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /books [post]
func (h *BookHandler) CreateBook(c echo.Context) error {
	var req models.BookCreateRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Warn("failed to bind request",
			zap.Error(err),
			zap.Any("request_body", c.Request().Body),
		)

		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Details: []ValidationError{{
				Field:   "body",
				Message: "Invalid JSON format",
			}},
		})
	}

	// invalidate request using validator
	if err := h.validator.Struct(req); err != nil {
		var validationErrors []ValidationError
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, ValidationError{
				Field:   err.Field(),
				Message: validationMessage(err),
			})
		}

		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Code:    http.StatusUnprocessableEntity,
			Message: "Validation failed",
			Details: validationErrors,
		})
	}

	h.logger.Debug("create book request validated",
		zap.Any("request", req),
	)

	book, err := h.service.CreateBook(c.Request().Context(), &req)
	if err != nil {
		return handleServiceError(c, h.logger, err)
	}

	c.Response().Header().Set("Cache-Control", "no-store")

	h.logger.Info("book created successfully",
		zap.Int("book_id", book.ID),
		zap.String("isbn", book.ISBN),
	)
	return c.JSON(http.StatusCreated, book)
}

// GetBook godoc
// @Summary Get a book by ID
// @Description Get a single book by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} models.Book
// @Header 200 {string} Cache-Control "max-age=3600, public"
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Failure 429 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /books/{id} [get]
func (h *BookHandler) GetBook(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Code:    http.StatusNotFound,
			Message: "Invalid book ID",
			Details: []ValidationError{{
				Field:   "id",
				Message: "Must be a positive integer",
			}},
		})
	}

	book, err := h.service.GetByBookID(c.Request().Context(), id)
	if err != nil {
		return handleServiceError(c, h.logger, err)
	}

	c.Response().Header().Set("Cache-Control", "max-age=3600, public")
	c.Response().Header().Set("ETag", generateETag(book))

	return c.JSON(http.StatusOK, book)
}

// FetchAllBook godoc
// @Summary List all books
// @Description Get a paginated list of books
// @Tags books
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} models.BookListResponse
// @Header 200 {string} Cache-Control "max-age=60, public"
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 429 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /books [get]
func (h *BookHandler) ListBooks(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	type QueryParams struct {
		Page  int `validate:"min=1"`
		Limit int `validate:"min=1,max=100"`
	}

	params := QueryParams{Page: page, Limit: limit}
	if err := h.validator.Struct(params); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_pagination",
			Code:    http.StatusBadRequest,
			Message: "Invalid pagination parameters",
		})
	}

	books, total, err := h.service.FetchAllBook(c.Request().Context(), page, limit)
	if err != nil {
		return handleServiceError(c, h.logger, err)
	}

	c.Response().Header().Set("Cache-Control", "max-age=60, public")
	return c.JSON(http.StatusOK, models.BookListResponse{
		Data:       books,
		TotalPages: total,
		Page:       page,
		Limit:      limit,
	})
}

// UpdateBook godoc
// @Summary Update a book
// @Description Update an existing book by ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body models.BookUpdateRequest true "Book data"
// @Success 200 {object} models.Book
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Failure 403 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /books/{id} [put]
func (h *BookHandler) UpdateBook(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Code:    http.StatusNotFound,
			Message: "Invalid book ID",
		})
	}

	var req models.BookUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
		})
	}

	book, err := h.service.UpdateBook(c.Request().Context(), id, &req)
	if err != nil {
		return handleServiceError(c, h.logger, err)
	}

	return c.JSON(http.StatusOK, book)
}

// DeleteBook godoc
// @Summary Delete a book
// @Description Delete a book by ID (soft delete)
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 204
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Failure 403 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBook(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Code:    http.StatusNotFound,
			Message: "Invalid book ID",
		})
	}

	if err := h.service.DeleteBook(c.Request().Context(), id); err != nil {
		return handleServiceError(c, h.logger, err)
	}

	return c.NoContent(http.StatusNoContent)
}

type ErrorResponse struct {
	Error   string            `json:"error" example:"not found"`
	Message string            `json:"message" example:"book not found"`
	Code    int               `json:"code" example:"404"`
	Details []ValidationError `json:"details"`
}

// Helper functions
func validationMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "Value too small"
	case "max":
		return "Value too large"
	case "email":
		return "Invalid email format"
	default:
		return fieldError.Error()
	}
}

func handleServiceError(c echo.Context, logger *zap.Logger, err error) error {
	ctx := c.Request().Context()

	if valErr, ok := err.(validator.ValidationErrors); ok {
		validationErrors := make([]ValidationError, len(valErr))
		for i, fieldErr := range valErr {
			validationErrors[i] = ValidationError{
				Field:   fieldErr.Field(),
				Message: validationMessage(fieldErr),
			}
		}
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: "Validation failed",
			Details: validationErrors,
		})
	}

	switch {
	case errors.Is(err, services.ErrNotFound):
		logger.Warn("resource not found",
			zap.Error(err),
			zap.String("path", c.Path()),
			zap.String("trace_id", getTraceID(ctx)),
		)

		return c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "not_found",
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
	case errors.Is(err, services.ErrInvalidInput):
		logger.Warn("invalid input",
			zap.Error(err),
			zap.Any("request_body", c.Request().Body),
			zap.String("trace_id", getTraceID(ctx)),
		)

		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_input",
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	case errors.Is(err, services.ErrConflict):
		logger.Warn("conflict",
			zap.Error(err),
			zap.Any("conflict", c.Request().Body),
			zap.String("trace_id", getTraceID(ctx)),
		)

		return c.JSON(http.StatusConflict, ErrorResponse{
			Error:   "conflict",
			Code:    http.StatusConflict,
			Message: err.Error(),
		})
	case errors.Is(err, context.DeadlineExceeded):
		logger.Warn("timeout",
			zap.Error(err),
			zap.Any("timeout", c.Request().Body),
			zap.String("trace_id", getTraceID(ctx)),
		)

		return c.JSON(http.StatusGatewayTimeout, ErrorResponse{
			Error:   "timeout",
			Code:    http.StatusGatewayTimeout,
			Message: "Request timed out",
		})
	default:
		logger.Error("unexpected error",
			zap.Error(err),
			zap.String("path", c.Path()),
			zap.String("method", c.Request().Method),
			zap.String("trace_id", getTraceID(ctx)),
			zap.Stack("stack"),
		)

		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}
}

func generateETag(book *models.Book) string {
	return strconv.Itoa(book.ID) + "-" + strconv.FormatInt(book.UpdatedAt.Unix(), 10)
}

func getTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(middleware.TraceIDKey).(string); ok {
		return traceID
	}
	return "not_available"
}
