package configure

import (
	"fmt"
	"strconv"

	"github.com/NebulousLabs/go-upnp"
)

// OpenPort opens port and forwards it to us
func OpenPort(portString string) error {
	port, err := strconv.Atoi(portString)
	if err != nil {
		return fmt.Errorf("Error parsing PORT: %v", err.Error())
	}

	d, err := upnp.Discover()
	if err != nil {
		return fmt.Errorf("Error discovering router: %v", err.Error())
	}

	err = d.Forward(uint16(port), "OpenCAP server")
	if err != nil {
		return fmt.Errorf("Error closing port forward: %v", err.Error())
	}
	return nil
}

// ClosePort closes port forwarding
func ClosePort(portString string) error {
	port, err := strconv.Atoi(portString)
	if err != nil {
		return fmt.Errorf("Error parsing PORT: %v", err.Error())
	}

	d, err := upnp.Discover()
	if err != nil {
		return fmt.Errorf("Error discovering router: %v", err.Error())
	}

	err = d.Clear(uint16(port))
	if err != nil {
		return fmt.Errorf("Error closing port forward: %v", err.Error())
	}
	return nil
}
