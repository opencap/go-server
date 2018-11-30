package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	opencap "github.com/opencap/go-opencap"
	"github.com/opencap/go-server/auth"
	"github.com/opencap/go-server/database"
)

type postUserRequest struct {
	Alias              string `json:"alias"`
	Password           string `json:"password"`
	CreateUserPassword string `json:"create_user_password"`
}

func validatePostUserParams(req *http.Request) (string, string, string, string, error) {
	params := postUserRequest{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", "", "", "", errors.New("Error reading request")
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		return "", "", "", "", errors.New("Error parsing request")
	}

	if !auth.ValidatePassword(params.Password) {
		return "", "", "", "", errors.New("Invalid password format, should")
	}

	username, domain, err := opencap.ValidateAlias(params.Alias)
	if err != nil {
		return "", "", "", "", err
	}
	params.Alias = username

	return username, domain, params.Password, params.CreateUserPassword, nil
}

func (cfg Config) postUserHandler(w http.ResponseWriter, req *http.Request) {
	username, domain, password, createUserPassword, err := validatePostUserParams(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if domain != cfg.domainName {
		respondWithError(w, http.StatusBadRequest, "Alias must use $"+cfg.domainName)
		return
	}

	if createUserPassword != cfg.createUserPassword {
		respondWithError(w, http.StatusBadRequest, "invalid create_user_password")
		return
	}

	user := database.User{
		Username: username,
		Domain:   domain,
		Password: password,
	}

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
