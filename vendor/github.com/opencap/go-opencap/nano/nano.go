package nano

import (
	"errors"
	"regexp"
	"strings"
)

// ValidateAddress returns true if string is a valid username format
func ValidateAddress(address string) error {
	Re := regexp.MustCompile(`^(xrb_|nano_)[0-9a-z]{60}$`)
	if !Re.MatchString(address) {
		return errors.New("invalid format")
	}

	sections := strings.Split(address, "_")
	if len(sections) < 2 {
		return errors.New("invalid number of seperators")
	}

	if len(sections[1]) < 60 {
		return errors.New("too short")
	}

	if !validateChecksum(sections[1][:52], sections[1][52:]) {
		return errors.New("invalid checksum")
	}

	return nil
}

func validateChecksum(address, checksum string) bool {
	// TODO: Actually check the checksum
	return len(address) == 52 && len(checksum) == 8
}
