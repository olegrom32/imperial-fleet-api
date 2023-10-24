package internal

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/olegrom32/imperial-fleet-api/internal/application/handler/createspaceship"
	"github.com/olegrom32/imperial-fleet-api/internal/application/handler/deletespaceship"
	"github.com/olegrom32/imperial-fleet-api/internal/application/handler/getspaceship"
	"github.com/olegrom32/imperial-fleet-api/internal/application/handler/listspaceships"
	"github.com/olegrom32/imperial-fleet-api/internal/application/handler/updatespaceship"
	"github.com/olegrom32/imperial-fleet-api/internal/infra/repository"
	"github.com/olegrom32/imperial-fleet-api/internal/server"
)

func wire() (*server.DIContainer, error) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}

	dbCfg := mysql.Config{
		User:   "root",
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:3306", dbHost),
		DBName: "api",
	}

	db, err := sql.Open("mysql", dbCfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	// A simple retry to wait for DB
	tries := 60
	for i := 1; i <= tries; i++ {
		err := db.Ping()
		if err == nil {
			break
		}

		if i == tries {
			return nil, fmt.Errorf("db is not ready: %w", err)
		}

		time.Sleep(time.Second)
	}

	spaceshipsRepo := repository.NewSpaceship(db)
	armamentRepo := repository.NewArmament(db)
	fleetMembersRepo := repository.NewFleetMember()

	listSpaceshipHandler := listspaceships.NewHandler(spaceshipsRepo)
	getSpaceshipHandler := getspaceship.NewHandler(spaceshipsRepo)
	createSpaceshipHandler := createspaceship.NewHandler(spaceshipsRepo, armamentRepo)
	updateSpaceshipHandler := updatespaceship.NewHandler(spaceshipsRepo, armamentRepo)
	deleteSpaceshipHandler := deletespaceship.NewHandler(spaceshipsRepo)

	return &server.DIContainer{
		DB:                      db,
		FleetMembersRepo:        fleetMembersRepo,
		ListSpaceshipsHandler:   listSpaceshipHandler,
		GetSpaceshipHandler:     getSpaceshipHandler,
		CreateSpaceshipsHandler: createSpaceshipHandler,
		UpdateSpaceshipsHandler: updateSpaceshipHandler,
		DeleteSpaceshipsHandler: deleteSpaceshipHandler,
	}, nil
}
