package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/opencap/go-opencap/bitcoin"
	"github.com/opencap/go-opencap/nano"
	"github.com/opencap/go-server/auth"
	"github.com/opencap/go-server/database"
)

type putAddressRequest struct {
	AddressType int    `json:"address_type"`
	Address     string `json:"address"`
}

func validatePutAddressParams(req *http.Request) (putAddressRequest, error) {
	params := putAddressRequest{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return params, errors.New("Error parsing request")
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		return params, errors.New("Error parsing request")
	}

	valid := false
	switch params.AddressType {
	case 100: // Bitcoin P2PKH
		valid = bitcoin.ValidateP2PKH(params.Address) == nil
	case 101: // Bitcoin P2SH
		valid = bitcoin.ValidateP2SH(params.Address) == nil
	case 102: // Bitcoin Bech32
		valid = bitcoin.ValidateSegwitBech32(params.Address) == nil
	case 103: // Bitcoin Payment Code
		valid = true
	case 200: // Bitcoin Cash P2PKH
		valid = bitcoin.ValidateP2PKH(params.Address) == nil
	case 201: // Bitcoin Cash P2SH
		valid = bitcoin.ValidateP2SH(params.Address) == nil
	case 300: // Nano
		valid = nano.ValidateAddress(params.Address) == nil
	default:
		valid = false
	}
	if !valid {
		return putAddressRequest{}, fmt.Errorf(
			"Invalid address format")
	}

	return params, nil
}

func reqToAddress(req putAddressRequest) database.Address {
	address := database.Address{}
	address.Address = req.Address
	address.AddressType = req.AddressType
	return address
}

func (cfg config) putAddressHandler(w http.ResponseWriter, req *http.Request) {
	reqModel, err := validatePutAddressParams(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	domain, username, err := auth.Authorize(req.Header, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid login credentials")
		return
	}

	user, err := cfg.db.GetUserByDomainUsername(domain, username)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid login credentials")
		return
	}

	address := reqToAddress(reqModel)

	err = cfg.db.CreateOrUpdateAddress(&user, address)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't update address")
		return
	}

	var empty struct{}
	respondWithJSON(w, http.StatusOK, empty)
}
