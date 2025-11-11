package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type WorkOutHandler struct{}

func NewWorkOutHandler() *WorkOutHandler {
	return &WorkOutHandler{}
}

func (wh *WorkOutHandler) HandleGetWorkOutById(w http.ResponseWriter, r *http.Request) {
	paramWorkoutId := chi.URLParam(r, "id")
	if paramWorkoutId == "" {
		http.NotFound(w, r)
		return
	}
	id, err := strconv.ParseInt(paramWorkoutId, 10, 64)
	if err != nil {
		http.NotFound(w, r)
	}

	fmt.Fprintf(w, "The workout ID is: %d", id)
}

func (wh *WorkOutHandler) HandleCreateWorkOut(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Workout created\n")
}
