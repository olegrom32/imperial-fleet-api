package internal

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/olegrom32/imperial-fleet-api/internal/server"
	"github.com/olegrom32/imperial-fleet-api/internal/server/middleware"
)

func Run() error {
	di, err := wire()
	if err != nil {
		return fmt.Errorf("failed to initialize dependencies: %w", err)
	}

	authMiddleware, err := middleware.BasicAuth("imperial-fleet-api", di.FleetMembersRepo)
	if err != nil {
		return fmt.Errorf("failed to initialize auth middleware: %w", err)
	}

	r := chi.NewRouter()
	r.Use(authMiddleware)
	r.Use(middleware.ContentTypeJSON)
	r.Use(chimiddleware.Timeout(30 * time.Second))

	server.RegisterRoutes(r, di)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return http.ListenAndServe(":"+port, r)
}
