package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	opencap "github.com/opencap/go-opencap"
	"github.com/opencap/go-server/database"
)

type getAddressesResponse struct {
	Address     string `json:"address"`
	AddressType int    `json:"address_type"`
}

func addressesToResponse(addresses []database.Address) (string, error) {
	respBody := make([]getAddressesResponse, 0)
	for _, v := range addresses {
		respBody = append(respBody, getAddressesResponse{
			Address:     v.Address,
			AddressType: v.AddressType,
		})
	}
	respBodyBytes, err := json.Marshal(respBody)
	return string(respBodyBytes), err
}

func addressToResponse(address database.Address) (string, error) {
	respBody := getAddressesResponse{
		Address:     address.Address,
		AddressType: address.AddressType,
	}
	respBodyBytes, err := json.Marshal(respBody)
	return string(respBodyBytes), err
}

func validateGetAddressParams(req *http.Request) (string, string, int, error) {
	params := req.URL.Query()
	aliasSlice, ok := params["alias"]
	if !ok || len(aliasSlice) < 1 {
		return "", "", 0, errors.New("No alias was included in the request")
	}

	addressTypeSlice, ok := params["address_type"]
	if !ok || len(addressTypeSlice) < 1 {
		addressTypeSlice = []string{"-1"} // -1 means no address type was provided
	}

	addressTypeInt, err := strconv.Atoi(addressTypeSlice[0])
	if err != nil {
		return "", "", 0, errors.New("Address type must be an ID number")
	}

	username, domain, err := opencap.ValidateAlias(aliasSlice[0])
	if err != nil {
		return "", "", 0, err
	}

	return username, domain, addressTypeInt, nil
}

func (cfg config) getAddressHandler(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	username, domain, addressType, err := validateGetAddressParams(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := cfg.db.GetUserByDomainUsername(domain, username)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	// Address type was requested
	if addressType >= 0 {
		address, err := cfg.db.GetAddressByAddressType(user, addressType)
		if err != nil || len(address.Address) == 0 {
			respondWithError(w, http.StatusNotFound, "Address not found")
		}

		body, err := addressToResponse(address)
		if err != nil || len(address.Address) == 0 {
			respondWithError(w, http.StatusInternalServerError, "Address not found")
			return
		}

		respondWithJSON(w, http.StatusOK, body)
		return
	}

	// return all addresses
	addresses, err := cfg.db.GetAddresses(user)
	body, err := addressesToResponse(addresses)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, body)
}
