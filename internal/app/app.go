package app

import (
	"fmt"
	"log/slog"
	"music/internal/config"
	"music/internal/controller"
	"music/internal/handler"
	"music/internal/repository"
	"music/internal/service"
	"music/migration"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	_ "music/docs"

	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
)

type HTTPServer struct {
	Port       string
	httpServer *http.Server
	handler    http.Handler
	db         repository.DB
}

//	@title		Online music
//	@version	1.0

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8000
//	@BasePath	/

//	@schemes	http https

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func NewHTTPServer(cfg config.Config, logger *slog.Logger) *HTTPServer {
	if cfg.Mode == config.ModeProd {
		gin.SetMode(gin.ReleaseMode)
	}

	db := repository.NewPostgresDB(cfg)

	err := migration.StartMigrate(cfg, repository.GetDBUrl(cfg), logger)
	if err != nil {
		logger.Error("Failed to apply migrations", slog.String("error", err.Error()))
	}

	repos := repository.NewRepository(db.GetDB())
	services := service.NewService(repos, logger)
	controllers := controller.NewController(services, logger)
	handlers := handler.NewHandler(controllers)

	return &HTTPServer{
		Port:    cfg.Port,
		handler: handlers.InitRoutes(),
		db:      db,
	}
}

func (s *HTTPServer) Run() error {
	s.httpServer = &http.Server{
		Addr:           ":" + s.Port,
		Handler:        s.handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *HTTPServer) Shutdown() error {
	var shutdownErr error

	// Close the database connection
	if err := s.db.Close(); err != nil {
		shutdownErr = fmt.Errorf("failed to close database: %w", err)
	}

	// Close the HTTP server
	if err := s.httpServer.Close(); err != nil {
		if shutdownErr != nil {
			shutdownErr = fmt.Errorf("%v; failed to close HTTP server: %w", shutdownErr, err)
		} else {
			shutdownErr = fmt.Errorf("failed to close HTTP server: %w", err)
		}
	}

	return shutdownErr
}
