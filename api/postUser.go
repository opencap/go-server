package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	opencap "github.com/opencap/go-opencap"
	"github.com/opencap/go-server/auth"
	"github.com/opencap/go-server/database"
)

type postUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func validatePostUserParams(req *http.Request) (postUserRequest, error) {
	params := postUserRequest{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return params, errors.New("Error parsing request")
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		return params, errors.New("Error parsing request")
	}

	if !auth.ValidatePassword(params.Password) {
		return params, errors.New("Invalid password format, should")
	}

	username, valid := opencap.ValidateUsername(params.Username)
	if !valid {
		return params, fmt.Errorf("Invalid password format. Passwords require at least one upper case letter, at least one special character, at least one number, and must have at least %v characters total", auth.MinPasswordLength)
	}
	params.Username = username

	return params, nil
}

func reqToUser(req postUserRequest, domain string) database.User {
	user := database.User{}
	user.Username = req.Username
	user.Password = req.Password
	user.Domain = domain
	return user
}

func (cfg config) postUserHandler(w http.ResponseWriter, req *http.Request) {
	reqModel, err := validatePostUserParams(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user := reqToUser(reqModel, cfg.domain)

	user.Password, err = auth.HashPassword(user.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = cfg.db.CreateUser(&user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Username already taken")
		return
	}

	var empty struct{}
	respondWithJSON(w, http.StatusOK, empty)
}
