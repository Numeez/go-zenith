package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Numeez/go-zenith/internal/store"
	"github.com/Numeez/go-zenith/internal/tokens"
	"github.com/Numeez/go-zenith/internal/utils"
)

type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (h *TokenHandler) HandlerCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("ERROR: decoding request: %v", err)
		_ = utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	user, err := h.userStore.GetUserByName(req.Username)
	if err != nil || user == nil {
		h.logger.Printf("ERROR: GetUserByName: %v", err)
		_ = utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	passwordMatch, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		h.logger.Printf("ERROR : %v", err)
		_ = utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if !passwordMatch {
		h.logger.Printf("ERROR : %v", err)
		_ = utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credential"})
		return
	}
	token, err := h.tokenStore.CreateNewToken(user.Id, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		h.logger.Printf("ERROR : %v", err)
		_ = utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return

	}
	_ = utils.WriteJson(w, http.StatusCreated, utils.Envelope{"authToken": token})
}
