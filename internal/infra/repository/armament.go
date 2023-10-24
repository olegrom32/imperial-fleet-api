package repository

import (
	"database/sql"
	"errors"

	"github.com/olegrom32/imperial-fleet-api/internal/application"
	"github.com/olegrom32/imperial-fleet-api/internal/application/domain"
)

type armamentDB struct {
	ID   int64
	Name string
}

type Armament struct {
	db *sql.DB
}

func NewArmament(db *sql.DB) *Armament {
	return &Armament{
		db: db,
	}
}

// Get returns a single record
func (r *Armament) Get(id int64) (domain.Armament, error) {
	row := r.db.QueryRow(
		"SELECT id, name FROM api.armament s WHERE id = ?",
		id,
	)

	var dbItem armamentDB

	err := row.Scan(&dbItem.ID, &dbItem.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Armament{}, application.ErrNotFound
		}

		return domain.Armament{}, err
	}

	return domain.Armament{
		ID:   dbItem.ID,
		Name: dbItem.Name,
	}, nil
}
