package updatespaceship

import (
	"context"
	"errors"
	"fmt"

	"github.com/olegrom32/imperial-fleet-api/internal/application/domain"
)

var (
	ErrInvalidArmament = errors.New("invalid armament")
	ErrInvalidField    = errors.New("field is invalid")
	ErrReadonly        = errors.New("spaceship is in read-only state")
)

type Command struct {
	SpaceshipID int64
	Name        string
	Class       string
	Armament    []CommandArmament
	Crew        int64
	Image       string
	Value       float64
	Status      string
}

type CommandArmament struct {
	ArmamentID int64
	Quantity   int64
}

type spaceshipRepo interface {
	Get(id int64) (domain.Spaceship, error)
	Update(m *domain.Spaceship) error
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

func (h *Handler) Handle(ctx context.Context, cmd Command) error {
	// Validate command. Assuming name, class and status are required, the rest is optional
	// Validation happens to be the same as in createspaceship handler. But it could be different,
	// the business constraints for updating could be different so we cannot extract this code to a
	// single service
	if cmd.Name == "" {
		return fmt.Errorf("%w: name cannot be blank", ErrInvalidField)
	}

	if cmd.Class == "" {
		return fmt.Errorf("%w: class cannot be blank", ErrInvalidField)
	}

	if cmd.Status == "" {
		return fmt.Errorf("%w: status cannot be blank", ErrInvalidField)
	}

	if cmd.Value < 0 {
		return fmt.Errorf("%w: value cannot be negative", ErrInvalidField)
	}

	if cmd.Crew < 0 {
		return fmt.Errorf("%w: crew cannot be negative", ErrInvalidField)
	}

	arms := make([]domain.SpaceshipArmament, len(cmd.Armament))

	for i := range cmd.Armament {
		// Validate armament
		if cmd.Armament[i].Quantity < 0 {
			return fmt.Errorf("%w: quantity cannot be negative", ErrInvalidArmament)
		}

		// Check if armament exists
		arm, err := h.armRepo.Get(cmd.Armament[i].ArmamentID)
		if err != nil {
			return fmt.Errorf("failed to find the armament in the repo: %w", err)
		}

		arms[i] = domain.SpaceshipArmament{
			Armament: arm,
			Quantity: cmd.Armament[i].Quantity,
		}
	}

	ship, err := h.shipRepo.Get(cmd.SpaceshipID)
	if err != nil {
		return fmt.Errorf("failed to find the spaceship in the repo: %w", err)
	}

	// Check if ship is in deleted state
	if ship.Deleted() {
		return ErrReadonly
	}

	ship.Name = cmd.Name
	ship.Class = cmd.Class
	ship.Status = cmd.Status
	ship.Crew = cmd.Crew
	ship.Value = cmd.Value
	ship.Image = cmd.Image
	ship.Armament = arms

	if err := h.shipRepo.Update(&ship); err != nil {
		return fmt.Errorf("failed to update the spaceship in repo: %w", err)
	}

	return nil
}
