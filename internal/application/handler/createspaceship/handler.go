package createspaceship

import (
	"context"
	"errors"
	"fmt"

	"github.com/olegrom32/imperial-fleet-api/internal/application/domain"
)

var (
	ErrInvalidArmament = errors.New("invalid armament")
	ErrInvalidField    = errors.New("field is invalid")
)

type Command struct {
	Name     string
	Class    string
	Armament []CommandArmament
	Crew     int64
	Image    string
	Value    float64
	Status   string
}

type CommandArmament struct {
	ArmamentID int64
	Quantity   int64
}

type spaceshipRepo interface {
	Create(m *domain.Spaceship) error
}

type armamentRepo interface {
	Get(id int64) (domain.Armament, error)
}

type Handler struct {
	shipRepo spaceshipRepo
	armRepo  armamentRepo
}

func NewHandler(shipRepo spaceshipRepo, armRepo armamentRepo) *Handler {
	return &Handler{
		shipRepo: shipRepo,
		armRepo:  armRepo,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) (domain.Spaceship, error) {
	// Validate command. Assuming name, class and status are required, the rest is optional

	// TODO check for duplicates here or add unique constraint in the db
	if cmd.Name == "" {
		return domain.Spaceship{}, fmt.Errorf("%w: name cannot be blank", ErrInvalidField)
	}

	if cmd.Class == "" {
		return domain.Spaceship{}, fmt.Errorf("%w: class cannot be blank", ErrInvalidField)
	}

	if cmd.Status == "" {
		return domain.Spaceship{}, fmt.Errorf("%w: status cannot be blank", ErrInvalidField)
	}

	if cmd.Value < 0 {
		return domain.Spaceship{}, fmt.Errorf("%w: value cannot be negative", ErrInvalidField)
	}

	if cmd.Crew < 0 {
		return domain.Spaceship{}, fmt.Errorf("%w: crew cannot be negative", ErrInvalidField)
	}

	arms := make([]domain.SpaceshipArmament, len(cmd.Armament))

	for i := range cmd.Armament {
		// Validate armament
		if cmd.Armament[i].Quantity < 0 {
			return domain.Spaceship{}, fmt.Errorf("%w: quantity cannot be negative", ErrInvalidArmament)
		}

		// Check if armament exists
		arm, err := h.armRepo.Get(cmd.Armament[i].ArmamentID)
		if err != nil {
			return domain.Spaceship{}, fmt.Errorf("%w: %s", ErrInvalidArmament, err.Error())
		}

		arms[i] = domain.SpaceshipArmament{
			Armament: arm,
			Quantity: cmd.Armament[i].Quantity,
		}
	}

	ship := domain.Spaceship{
		Name:     cmd.Name,
		Class:    cmd.Class,
		Armament: arms,
		Crew:     cmd.Crew,
		Image:    cmd.Image,
		Value:    cmd.Value,
		Status:   cmd.Status,
	}

	if err := h.shipRepo.Create(&ship); err != nil {
		return domain.Spaceship{}, fmt.Errorf("failed to create the spaceship in repo: %w", err)
	}

	return ship, nil
}
