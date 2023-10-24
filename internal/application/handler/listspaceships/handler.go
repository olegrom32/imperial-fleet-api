package listspaceships

import (
	"context"
	"fmt"

	"github.com/olegrom32/imperial-fleet-api/internal/application/domain"
)

type Command struct {
	Name   string
	Class  string
	Status string
}

type spaceshipRepo interface {
	List(name, class, status string) ([]domain.Spaceship, error)
}

type Handler struct {
	repo spaceshipRepo
}

func NewHandler(repo spaceshipRepo) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) ([]domain.Spaceship, error) {
	ships, err := h.repo.List(cmd.Name, cmd.Class, cmd.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to get spaceships from repo: %w", err)
	}

	return ships, nil
}
