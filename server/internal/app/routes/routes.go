package routes

import (
	"bf-api/internal/app/handlers"
	bfMiddleware "bf-api/internal/app/middleware"

	"bf-api/internal/domain/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func APIRouter(e *echo.Echo, bookService services.BookService, logger *zap.Logger) {
	bookHandler := handlers.NewBookHandler(bookService, logger)

	e.Use(
		middleware.Recover(),
		middleware.RequestID(),
		bfMiddleware.Tracing(),
		middleware.RequestLoggerWithConfig(
			middleware.RequestLoggerConfig{
				LogURI:    true,
				LogStatus: true,
				LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
					logger.Info("request",
						zap.String("method", c.Request().Method),
						zap.String("uri", v.URI),
						zap.Int("status", v.Status),
						zap.Duration("latency", v.Latency),
						zap.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
					)
					return nil
				},
			},
		),
	)

	v1 := e.Group("/api/v1")

	v1.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	bookRoutes := v1.Group("/books")
	bookRoutes.Use(
		middleware.Gzip(),
		middleware.Secure(),
		middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5)),
	)

	bookRoutes.POST("", bookHandler.CreateBook)
	bookRoutes.GET("", bookHandler.ListBooks)
	bookRoutes.GET("/:id", bookHandler.GetBook)
	bookRoutes.PUT("/:id", bookHandler.UpdateBook)
	bookRoutes.DELETE("/:id", bookHandler.DeleteBook)

}
