package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Envelope map[string]any

func WriteJson(w http.ResponseWriter, status int, data Envelope) error {
	js, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}
	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(js)
	return nil
}

func ReadIdParam(r *http.Request) (int64, error) {
	paramWorkoutId := chi.URLParam(r, "id")
	if paramWorkoutId == "" {
		return 0, errors.New("id param is not availale")
	}
	id, err := strconv.ParseInt(paramWorkoutId, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil 
}
