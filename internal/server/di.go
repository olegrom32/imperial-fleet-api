package server

import (
	"database/sql"

	"github.com/olegrom32/imperial-fleet-api/internal/application/handler/createspaceship"
	"github.com/olegrom32/imperial-fleet-api/internal/application/handler/deletespaceship"
	"github.com/olegrom32/imperial-fleet-api/internal/application/handler/getspaceship"
	"github.com/olegrom32/imperial-fleet-api/internal/application/handler/listspaceships"
	"github.com/olegrom32/imperial-fleet-api/internal/application/handler/updatespaceship"
	"github.com/olegrom32/imperial-fleet-api/internal/infra/repository"
)

type DIContainer struct {
	DB *sql.DB

	FleetMembersRepo *repository.FleetMember

	ListSpaceshipsHandler   *listspaceships.Handler
	GetSpaceshipHandler     *getspaceship.Handler
	CreateSpaceshipsHandler *createspaceship.Handler
	UpdateSpaceshipsHandler *updatespaceship.Handler
	DeleteSpaceshipsHandler *deletespaceship.Handler
}
