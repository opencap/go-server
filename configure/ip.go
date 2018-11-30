package configure

import (
	"errors"

	externalip "github.com/glendc/go-external-ip"
)

// GetPublicIP preferred outbound ip of this machine
func GetPublicIP() (string, error) {
	consensus := externalip.DefaultConsensus(nil, nil)
	ip, err := consensus.ExternalIP()
	if err != nil {
		return "", errors.New("Couldn't get ip address")
	}
	return ip.String(), nil
}
