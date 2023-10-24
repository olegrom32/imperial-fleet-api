package deletespaceship

import (
	"context"
	"fmt"

	"github.com/olegrom32/imperial-fleet-api/internal/application/domain"
)

type Command struct {
	SpaceshipID int64
}

type spaceshipRepo interface {
	Get(id int64) (domain.Spaceship, error)
	Update(m *domain.Spaceship) error
}

type Handler struct {
	repo spaceshipRepo
}

func NewHandler(repo spaceshipRepo) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) error {
	ship, err := h.repo.Get(cmd.SpaceshipID)
	if err != nil {
		return fmt.Errorf("failed to find the spaceship in the repo: %w", err)
	}

	if ship.Deleted() {
		// Already deleted - no action needed (idempotent API)
		return nil
	}

	ship.Status = domain.SpaceshipStatusDeleted

	if err := h.repo.Update(&ship); err != nil {
		return fmt.Errorf("failed to update the spaceship in repo: %w", err)
	}

	return nil
}
