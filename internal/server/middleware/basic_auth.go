package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/olegrom32/imperial-fleet-api/internal/application/domain"
)

type fleetMemberRepo interface {
	List(ctx context.Context) ([]domain.FleetMember, error)
}

func BasicAuth(realm string, repo fleetMemberRepo) (func(next http.Handler) http.Handler, error) {
	members, err := repo.List(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get list of fleet members: %w", err)
	}

	creds := make(map[string]string)

	for _, m := range members {
		creds[m.Login] = m.Password
	}

	return middleware.BasicAuth(realm, creds), nil
}
