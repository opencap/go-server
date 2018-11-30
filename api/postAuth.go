package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	opencap "github.com/opencap/go-opencap"
	"github.com/opencap/go-server/auth"
)

type postAuthResponse struct {
	Token string `json:"jwt"`
}

type postAuthRequest struct {
	Alias    string `json:"alias"`
	Password string `json:"password"`
}

func validatePostAuthParams(req *http.Request) (string, string, string, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", "", "", errors.New("Error parsing request")
	}

	params := postAuthRequest{}
	err = json.Unmarshal(body, &params)
	if err != nil {
		return "", "", "", errors.New("Error parsing request")
	}

	username, domain, err := opencap.ValidateAlias(params.Alias)
	return username, domain, params.Password, err
}

// Handler for postAuth
func (cfg Config) postAuthHandler(w http.ResponseWriter, req *http.Request) {
	username, domain, password, err := validatePostAuthParams(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	dbUser, err := cfg.db.GetUserByDomainUsername(domain, username)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "User not found")
		return
	}

	if !auth.CheckPasswordHash(password, dbUser.Password) {
		respondWithError(w, http.StatusBadRequest, "Incorrect username password combination")
		return
	}

	token, err := auth.MakeToken(dbUser.Domain, dbUser.Username, cfg.jwtSecret, time.Now().UTC(), cfg.jwtExpirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := postAuthResponse{
		Token: token,
	}

	respondWithJSON(w, http.StatusOK, resp)
}
