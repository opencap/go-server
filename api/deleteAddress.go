package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/opencap/go-server/auth"
)

func validateDeleteAddressParams(req *http.Request) (int, error) {
	params := mux.Vars(req)
	addressType, ok := params["address_type"]
	if !ok {
		return 0, errors.New("Address type bot found")
	}

	addressTypeInt, err := strconv.Atoi(addressType)
	if err != nil {
		return 0, errors.New("Address type must be an ID number")
	}

	return addressTypeInt, nil
}

func (cfg Config) deleteAddressesHandler(w http.ResponseWriter, req *http.Request) {
	addressType, err := validateDeleteAddressParams(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	domain, username, err := auth.Authorize(req.Header, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid authentication")
		return
	}
	user, err := cfg.db.GetUserByDomainUsername(domain, username)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	address, err := cfg.db.GetAddressByAddressType(user, addressType)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Address not found")
		return
	}

	err = cfg.db.DeleteAddress(address)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete address")
		return
	}

	var empty struct{}
	respondWithJSON(w, http.StatusOK, empty)
}
