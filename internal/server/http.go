package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/olegrom32/imperial-fleet-api/internal/application"
	"github.com/olegrom32/imperial-fleet-api/internal/application/handler/getspaceship"

	"github.com/olegrom32/imperial-fleet-api/internal/application/handler/createspaceship"
	"github.com/olegrom32/imperial-fleet-api/internal/application/handler/deletespaceship"
	"github.com/olegrom32/imperial-fleet-api/internal/application/handler/listspaceships"
	"github.com/olegrom32/imperial-fleet-api/internal/application/handler/updatespaceship"
)

var successResponse = `{"success":true}`

func RegisterRoutes(r *chi.Mux, di *DIContainer) {
	// List spaceships
	r.Get("/spaceship", func(w http.ResponseWriter, r *http.Request) {
		res, err := di.ListSpaceshipsHandler.Handle(r.Context(), listspaceships.Command{
			Name:   r.URL.Query().Get("name"),
			Class:  r.URL.Query().Get("class"),
			Status: r.URL.Query().Get("status"),
		})
		if err != nil {
			// TODO replace with proper logger
			log.Print(err)

			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(errorResponse(err))

			return
		}

		resp := struct {
			Data []any `json:"data"`
		}{
			Data: make([]any, len(res)),
		}

		for i := range res {
			resp.Data[i] = struct {
				ID     int64  `json:"id"`
				Name   string `json:"name"`
				Status string `json:"status"`
			}{
				ID:     res[i].ID,
				Name:   res[i].Name,
				Status: res[i].Status,
			}
		}

		respBytes, err := json.Marshal(resp)
		if err != nil {
			log.Print(err)

			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(errorResponse(err))

			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respBytes)
	})

	// Create spaceship
	r.Post("/spaceship", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name     string `json:"name"`
			Class    string `json:"class"`
			Armament []struct {
				ID  int64 `json:"id"`
				Qty int64 `json:"qty"`
			} `json:"armament"`
			Crew   int64   `json:"crew"`
			Image  string  `json:"image"`
			Value  float64 `json:"value"`
			Status string  `json:"status"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(errorResponse(err))

			return
		}

		cmd := createspaceship.Command{
			Name:     req.Name,
			Class:    req.Class,
			Armament: make([]createspaceship.CommandArmament, len(req.Armament)),
			Crew:     req.Crew,
			Image:    req.Image,
			Value:    req.Value,
			Status:   req.Status,
		}

		for i := range req.Armament {
			cmd.Armament[i] = createspaceship.CommandArmament{
				ArmamentID: req.Armament[i].ID,
				Quantity:   req.Armament[i].Qty,
			}
		}

		_, err := di.CreateSpaceshipsHandler.Handle(r.Context(), cmd)
		if err != nil {
			// TODO replace with proper logger
			log.Print(err)

			switch {
			case errors.Is(err, createspaceship.ErrInvalidField), errors.Is(err, createspaceship.ErrInvalidArmament):
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write(errorResponse(err))
			default:
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write(errorResponse(err))
			}

			return
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(successResponse))
	})

	// Get spaceship info
	r.Get("/spaceship/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(errorResponse(err))

			return
		}

		res, err := di.GetSpaceshipHandler.Handle(r.Context(), getspaceship.Command{
			ID: int64(idInt),
		})
		if err != nil {
			// TODO replace with proper logger
			log.Print(err)

			switch {
			case errors.Is(err, application.ErrNotFound):
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(""))
			default:
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write(errorResponse(err))
			}

			return
		}

		respBytes, err := json.Marshal(struct {
			ID     int64   `json:"id"`
			Name   string  `json:"name"`
			Class  string  `json:"class"`
			Crew   int64   `json:"crew"`
			Image  string  `json:"image"`
			Value  float64 `json:"value"`
			Status string  `json:"status"`
		}{
			ID:     res.ID,
			Name:   res.Name,
			Class:  res.Class,
			Crew:   res.Crew,
			Image:  res.Image,
			Value:  res.Value,
			Status: res.Status,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(errorResponse(err))

			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respBytes)
	})

	// Update spaceship
	r.Put("/spaceship/{id}", func(w http.ResponseWriter, r *http.Request) {
		// req happens to be the same as for creation but from a business logic perspective
		// they may be different, e.g. limited amount if fields may be allowed to update
		var req struct {
			Name     string `json:"name"`
			Class    string `json:"class"`
			Armament []struct {
				ID  int64 `json:"id"`
				Qty int64 `json:"qty"`
			} `json:"armament"`
			Crew   int64   `json:"crew"`
			Image  string  `json:"image"`
			Value  float64 `json:"value"`
			Status string  `json:"status"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(errorResponse(err))

			return
		}

		id := chi.URLParam(r, "id")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(errorResponse(err))

			return
		}

		cmd := updatespaceship.Command{
			SpaceshipID: int64(idInt),
			Name:        req.Name,
			Class:       req.Class,
			Armament:    make([]updatespaceship.CommandArmament, len(req.Armament)),
			Crew:        req.Crew,
			Image:       req.Image,
			Value:       req.Value,
			Status:      req.Status,
		}

		for i := range req.Armament {
			cmd.Armament[i] = updatespaceship.CommandArmament{
				ArmamentID: req.Armament[i].ID,
				Quantity:   req.Armament[i].Qty,
			}
		}

		if err := di.UpdateSpaceshipsHandler.Handle(r.Context(), cmd); err != nil {
			// TODO replace with proper logger
			log.Print(err)

			switch {
			case errors.Is(err, updatespaceship.ErrInvalidField),
				errors.Is(err, updatespaceship.ErrInvalidArmament),
				errors.Is(err, updatespaceship.ErrReadonly):
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write(errorResponse(err))
			case errors.Is(err, application.ErrNotFound):
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(""))
			default:
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write(errorResponse(err))
			}

			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(successResponse))
	})

	// Delete spaceship
	r.Delete("/spaceship/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(errorResponse(err))

			return
		}

		if err := di.DeleteSpaceshipsHandler.Handle(r.Context(), deletespaceship.Command{
			SpaceshipID: int64(idInt),
		}); err != nil {
			// TODO replace with proper logger
			log.Print(err)

			switch {
			case errors.Is(err, updatespaceship.ErrInvalidField),
				errors.Is(err, updatespaceship.ErrInvalidArmament),
				errors.Is(err, updatespaceship.ErrReadonly):
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write(errorResponse(err))
			case errors.Is(err, application.ErrNotFound):
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(""))
			default:
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write(errorResponse(err))
			}

			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(successResponse))
	})
}

func errorResponse(err error) []byte {
	return []byte(fmt.Sprintf(`{"success":false","message":%q}`, err.Error()))
}
