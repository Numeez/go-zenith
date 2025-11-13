package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Numeez/go-zenith/internal/middleware"
	"github.com/Numeez/go-zenith/internal/store"
	"github.com/Numeez/go-zenith/internal/utils"
)

type WorkOutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkOutHandler(store store.WorkoutStore, logger *log.Logger) *WorkOutHandler {
	return &WorkOutHandler{
		workoutStore: store,
		logger:       logger,
	}
}

func (wh *WorkOutHandler) HandleGetWorkOutById(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIdParam(r)
	if err != nil {
		http.NotFound(w, r)
	}
	workout, err := wh.workoutStore.GetWorkOutById(id)
	if err != nil {
		wh.logger.Print(err.Error())
		_ = utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": err})
		return
	}
	_ = utils.WriteJson(w, http.StatusOK, utils.Envelope{"workout": workout})

}

func (wh *WorkOutHandler) HandleCreateWorkOut(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	if err := json.NewDecoder(r.Body).Decode(&workout); err != nil {
		wh.logger.Print(err.Error())
		_ = utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create workout"})
		return
	}
	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		_ = utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "user should be logged in"})
		return
	}
	workout.UserId = currentUser.Id
	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		wh.logger.Print(err.Error())
		_ = utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create workout"})
		return
	}
	_ = utils.WriteJson(w, http.StatusOK, utils.Envelope{"workout": createdWorkout})
}

func (wh *WorkOutHandler) HandlerUpdateWorkoutById(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIdParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: ReadIdParam: %v", err)
		if err := utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workou id"}); err != nil {
			wh.logger.Printf("ERROR: writing json error: %v", err)
		}
		return
	}
	existingWorkout, err := wh.workoutStore.GetWorkOutById(id)
	if err != nil {
		wh.logger.Printf("ERROR: GetWorkoutById: %v", err)
		if err := utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "unable to fetch the workout"}); err != nil {
			wh.logger.Printf("ERROR: writing json error: %v", err)
		}
		return
	}
	if existingWorkout == nil {
		wh.logger.Printf("ERROR: GetWorkoutById: %v", err)
		if err := utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "workout does not exists"}); err != nil {
			wh.logger.Printf("ERROR: writing json error: %v", err)
		}
		return
	}
	type updateRequest struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if request.Title != nil {
		existingWorkout.Title = *request.Title
	}
	if request.Description != nil {
		existingWorkout.Description = *request.Description
	}
	if request.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *request.DurationMinutes
	}
	if request.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *request.CaloriesBurned
	}
	if request.Entries != nil {
		existingWorkout.Entries = request.Entries
	}
	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		_ = utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "user should be logged in"})
		return
	}
	workOutOwner, err := wh.workoutStore.GetWorkoutOwner(int64(existingWorkout.Id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = utils.WriteJson(w, http.StatusNotFound, utils.Envelope{"error": "work does not found"})
			return
		}
		_ = utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if workOutOwner != currentUser.Id {
		_ = utils.WriteJson(w, http.StatusForbidden, utils.Envelope{"error": "you are not authorized to update this workout"})
		return

	}
	err = wh.workoutStore.UpdateWorkout(existingWorkout)
	if err != nil {
		wh.logger.Printf("Update workout failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = utils.WriteJson(w, http.StatusOK, utils.Envelope{"workout": existingWorkout})

}

func (wh *WorkOutHandler) HandlerDeleteWorkout(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIdParam(r)
	if err != nil {
		wh.logger.Print(err.Error())
		_ = utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "param id is not given"})
	}
	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		_ = utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "user should be logged in"})
		return
	}
	workOutOwner, err := wh.workoutStore.GetWorkoutOwner(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = utils.WriteJson(w, http.StatusNotFound, utils.Envelope{"error": "work does not found"})
			return
		}
		_ = utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if workOutOwner != currentUser.Id {
		_ = utils.WriteJson(w, http.StatusForbidden, utils.Envelope{"error": "you are not authorized to update this workout"})
		return

	}
	if err := wh.workoutStore.DeleteWorkout(id); err != nil {
		wh.logger.Print(err.Error())
		_ = utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": err})
		return
	}
	_ = utils.WriteJson(w, http.StatusNoContent, utils.Envelope{"message": "Workout Deleted"})
}
