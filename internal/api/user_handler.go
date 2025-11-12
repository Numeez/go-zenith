package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/Numeez/go-zenith/internal/store"
	"github.com/Numeez/go-zenith/internal/utils"
)

type registerUserStruct struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type UserHandler struct {
	store  store.UserStore
	logger *log.Logger
}

func NewUserHandler(store store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		store:  store,
		logger: logger,
	}
}

func (h *UserHandler) validateRegisterUser(req *registerUserStruct) error {
	if req.Username == "" {
		return errors.New("username cannot be empty")
	}
	if len(req.Username) > 50 {
		return errors.New("username is too long")
	}
	if req.Email == "" {
		return errors.New("email cannot be empty")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("Invalid email")
	}

	if req.Password == "" {
		return errors.New("password cannot be empty")
	}
	return nil

}

func (h *UserHandler) HandlerRegisterUser(w http.ResponseWriter, r *http.Request) {
	var request registerUserStruct
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Printf("ERROR: decoding register request: %v", err)
		_ = utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	err := h.validateRegisterUser(&request)
	if err != nil {
		_ = utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	user := &store.User{
		Username: request.Username,
		Email:    request.Email,
	}
	if request.Bio != "" {
		user.Bio = request.Bio
	}
	err = user.PasswordHash.Set(request.Password)
	if err != nil {
		h.logger.Printf("ERROR: hashing password failed: %v", err)
		_ = utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	createdUser, err := h.store.CreateUser(user)
	if err != nil {
		h.logger.Printf("ERROR: failed to create user: %v", err)
		_ = utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": err.Error()})
		return
	}
	_ = utils.WriteJson(w, http.StatusCreated, utils.Envelope{"user": createdUser})

}
