package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/olegrom32/imperial-fleet-api/internal/application"
	"github.com/olegrom32/imperial-fleet-api/internal/application/domain"
)

type spaceshipDB struct {
	ID     int64
	Name   string
	Class  string
	Crew   int64
	Image  string
	Value  float64
	Status string
}

type Spaceship struct {
	db *sql.DB
}

func NewSpaceship(db *sql.DB) *Spaceship {
	return &Spaceship{
		db: db,
	}
}

// Get returns a single record
func (r *Spaceship) Get(id int64) (domain.Spaceship, error) {
	row := r.db.QueryRow(
		"SELECT id, name, class, crew, image, value, status FROM api.spaceship s WHERE id = ?",
		id,
	)

	var dbItem spaceshipDB

	err := row.Scan(&dbItem.ID, &dbItem.Name, &dbItem.Class, &dbItem.Crew, &dbItem.Image, &dbItem.Value, &dbItem.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Spaceship{}, application.ErrNotFound
		}

		return domain.Spaceship{}, err
	}

	return spaceshipDBToDomain(dbItem), nil
}

// List returns list of records with optional filtering
func (r *Spaceship) List(name, class, status string) ([]domain.Spaceship, error) {
	filter := []string{"1 = 1"}

	var args []any

	if name != "" {
		filter = append(filter, "name = ?")
		args = append(args, name)
	}

	if class != "" {
		filter = append(filter, "class = ?")
		args = append(args, class)
	}

	if status != "" {
		filter = append(filter, "status = ?")
		args = append(args, status)
	}

	q := fmt.Sprintf("SELECT id, name, class, crew, image, value, status FROM api.spaceship s WHERE %s", strings.Join(filter, " AND "))

	rows, err := r.db.Query(q, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	var res []domain.Spaceship

	for rows.Next() {
		var dbItem spaceshipDB

		err = rows.Scan(&dbItem.ID, &dbItem.Name, &dbItem.Class, &dbItem.Crew, &dbItem.Image, &dbItem.Value, &dbItem.Status)
		if err != nil {
			return nil, err
		}

		// TODO get spaceship armament (I decided to skip this part and finish unit tests instead, cause I am short on time)

		res = append(res, spaceshipDBToDomain(dbItem))
	}

	return res, nil
}

// Create inserts records into db
func (r *Spaceship) Create(spaceship *domain.Spaceship) error {
	res, err := r.db.Exec(
		"INSERT INTO api.spaceship (name, class, crew, image, value, status) VALUES (?, ?, ?, ?, ?, ?)",
		spaceship.Name,
		spaceship.Class,
		spaceship.Crew,
		spaceship.Image,
		spaceship.Value,
		spaceship.Status,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	spaceship.ID = id

	return nil
}

// Update updates existing records
func (r *Spaceship) Update(spaceship *domain.Spaceship) error {
	_, err := r.db.Exec(
		"UPDATE api.spaceship SET name = ?, class = ?, crew = ?, image = ?, value = ?, status = ? WHERE id = ?",
		spaceship.Name,
		spaceship.Class,
		spaceship.Crew,
		spaceship.Image,
		spaceship.Value,
		spaceship.Status,
		spaceship.ID,
	)

	return err
}

func spaceshipDBToDomain(dbItem spaceshipDB) domain.Spaceship {
	return domain.Spaceship{
		ID:     dbItem.ID,
		Name:   dbItem.Name,
		Class:  dbItem.Class,
		Crew:   dbItem.Crew,
		Image:  dbItem.Image,
		Value:  dbItem.Value,
		Status: dbItem.Status,
	}
}

func spaceshipDomainToDB(m domain.Spaceship) spaceshipDB {
	return spaceshipDB{
		ID:     m.ID,
		Name:   m.Name,
		Class:  m.Class,
		Crew:   m.Crew,
		Image:  m.Image,
		Value:  m.Value,
		Status: m.Status,
	}
}
