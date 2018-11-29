package portforward

import (
	"fmt"

	"github.com/NebulousLabs/go-upnp"
)

// Open opens port 443 and forwards it to us
func Open(port uint16) (string, error) {
	d, err := upnp.Discover()
	if err != nil {
		return "", fmt.Errorf("Error discovering router: %v", err.Error())
	}
	// discover external IP
	ip, err := d.ExternalIP()
	if err != nil {
		return "", fmt.Errorf("Error fetching external IP address: %v", err.Error())
	}

	err = d.Forward(port, "OpenCAP server")
	if err != nil {
		return "", fmt.Errorf("Error closing port forward: %v", err.Error())
	}
	return ip, nil
}

// Close closes port 443 forwarding
func Close(port uint16) error {
	d, err := upnp.Discover()
	if err != nil {
		return fmt.Errorf("Error discovering router: %v", err.Error())
	}

	err = d.Clear(port)
	if err != nil {
		return fmt.Errorf("Error closing port forward: %v", err.Error())
	}
	return nil
}
