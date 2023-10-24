package repository

import (
	"context"

	"github.com/olegrom32/imperial-fleet-api/internal/application/domain"
)

// FleetMember is an in-memory, very simple implementation of the fleet member repository.
// The prod-ready version should store credentials in the db / dedicated auth microservice.
type FleetMember struct {
}

func NewFleetMember() *FleetMember {
	return &FleetMember{}
}

func (r *FleetMember) List(ctx context.Context) ([]domain.FleetMember, error) {
	return []domain.FleetMember{
		{Login: "user1", Password: "test1"},
		{Login: "user2", Password: "test2"},
	}, nil
}
