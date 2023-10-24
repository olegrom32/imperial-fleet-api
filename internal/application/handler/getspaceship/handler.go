package getspaceship

import (
	"context"
	"fmt"

	"github.com/olegrom32/imperial-fleet-api/internal/application/domain"
)

type Command struct {
	ID int64
}

type spaceshipRepo interface {
	Get(id int64) (domain.Spaceship, error)
}

type Handler struct {
	repo spaceshipRepo
}

func NewHandler(repo spaceshipRepo) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) (domain.Spaceship, error) {
	ship, err := h.repo.Get(cmd.ID)
	if err != nil {
		return domain.Spaceship{}, fmt.Errorf("failed to get spaceship from repo: %w", err)
	}

	return ship, nil
}
