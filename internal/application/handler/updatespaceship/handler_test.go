package updatespaceship

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/olegrom32/imperial-fleet-api/internal/application/domain"
	"github.com/olegrom32/imperial-fleet-api/internal/application/handler/updatespaceship/mocks"
)

//go:generate ../../../../bin/mockgen -source=handler.go -destination=mocks/handler.go -package mocks -typed

func TestHandler_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	sr := mocks.NewMockspaceshipRepo(ctrl)
	ar := mocks.NewMockarmamentRepo(ctrl)

	h := NewHandler(sr, ar)

	expectedErr := errors.New("error")

	t.Run("when spaceship exists and request is valid, should return no errors", func(t *testing.T) {
		sr.EXPECT().Get(int64(1)).Return(domain.Spaceship{
			ID:       1,
			Name:     "",
			Class:    "",
			Armament: nil,
			Crew:     0,
			Image:    "",
			Value:    0,
			Status:   "",
		}, nil)

		sr.EXPECT().Update(&domain.Spaceship{
			ID:       1,
			Name:     "name1",
			Class:    "class1",
			Armament: make([]domain.SpaceshipArmament, 0),
			Crew:     123,
			Image:    "img1",
			Value:    345.67,
			Status:   "new",
		}).Return(nil)

		err := h.Handle(context.Background(), Command{
			SpaceshipID: 1,
			Name:        "name1",
			Class:       "class1",
			Crew:        123,
			Image:       "img1",
			Value:       345.67,
			Status:      "new",
		})
		require.NoError(t, err)
	})

	t.Run("when invalid request is given, should return an error", func(t *testing.T) {
		err := h.Handle(context.Background(), Command{
			SpaceshipID: 1,
			Name:        "",
			Class:       "class1",
			Crew:        123,
			Image:       "img1",
			Value:       345.67,
			Status:      "new",
		})

		assert.ErrorIs(t, err, ErrInvalidField)
	})

	t.Run("when spaceship is readonly, should return an error", func(t *testing.T) {
		sr.EXPECT().Get(gomock.Any()).Return(domain.Spaceship{Status: "deleted"}, nil)

		err := h.Handle(context.Background(), Command{
			SpaceshipID: 1,
			Name:        "name1",
			Class:       "class1",
			Crew:        123,
			Image:       "img1",
			Value:       345.67,
			Status:      "new",
		})

		assert.ErrorIs(t, err, ErrReadonly)
	})

	t.Run("when spaceship fetch fails, should return an error", func(t *testing.T) {
		sr.EXPECT().Get(gomock.Any()).Return(domain.Spaceship{}, expectedErr)

		err := h.Handle(context.Background(), Command{
			SpaceshipID: 1,
			Name:        "name1",
			Class:       "class1",
			Crew:        123,
			Image:       "img1",
			Value:       345.67,
			Status:      "new",
		})

		assert.ErrorIs(t, err, expectedErr)
	})

	t.Run("when spaceship save fails, should return an error", func(t *testing.T) {
		sr.EXPECT().Get(gomock.Any()).Return(domain.Spaceship{}, nil)

		sr.EXPECT().Update(gomock.Any()).Return(expectedErr)

		err := h.Handle(context.Background(), Command{
			SpaceshipID: 1,
			Name:        "name1",
			Class:       "class1",
			Crew:        123,
			Image:       "img1",
			Value:       345.67,
			Status:      "new",
		})

		assert.ErrorIs(t, err, expectedErr)
	})
}
