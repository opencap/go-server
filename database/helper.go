package database

import (
	"errors"
	"strconv"
)

func getAddressByType(addressType int, slice []Address) (Address, error) {
	for _, v := range slice {
		if v.AddressType == addressType {
			return v, nil
		}
	}
	return Address{}, errors.New("Address" + strconv.Itoa(addressType) + "not found in slice")
}
