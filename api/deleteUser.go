package api

import (
	"net/http"

	"github.com/opencap/go-server/auth"
)

func (cfg config) deleteUserHandler(w http.ResponseWriter, req *http.Request) {
	domain, username, err := auth.Authorize(req.Header, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := cfg.db.GetUserByDomainUsername(domain, username)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	err = cfg.db.DeleteUser(user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var empty struct{}
	respondWithJSON(w, http.StatusOK, empty)
}
