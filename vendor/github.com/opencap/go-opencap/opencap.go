package opencap

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
)

// MinUsernameLength is the minimum length allowable for usernames
const MinUsernameLength = 1

// MaxUsernameLength is the maxiumum length allowable for usernames
const MaxUsernameLength = 25

// GetHost returns the highest priority opencap host URL at a given
// domain name.
func GetHost(domain string) (string, error) {
	_, addresses, err := net.LookupSRV("opencap", "tcp", domain)
	if err != nil {
		return "", err
	}
	if len(addresses) < 1 {
		return "", errors.New("No addresses found")
	}

	target := addresses[0].Target

	// strip trailing period for convenience
	if len(target) > 1 {
		if string(target[len(target)-1]) == "." {
			target = target[:len(target)-1]
		}
	}

	return target, nil
}

// ValidateUsername returns true if string is a valid username format
func ValidateUsername(username string) bool {
	Re := regexp.MustCompile(
		fmt.Sprintf(`^[a-z0-9._-]{%v,%v}$`,
			MinUsernameLength,
			MaxUsernameLength,
		),
	)
	return Re.MatchString(username)
}

// ValidateDomain returns true if string is a valid domain format
func ValidateDomain(username string) bool {
	Re := regexp.MustCompile(`^[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(username)
}

// ValidateAlias splits an alias, validates it parts and
// returns err if anything is wrong
func ValidateAlias(alias string) (string, string, error) {
	parts := strings.Split(alias, "$")
	if len(parts) != 2 {
		return "", "", errors.New("Incorrect alias format")
	}

	valid := ValidateUsername(parts[0])
	if !valid {
		return "", "", errors.New("Invalid username format")
	}

	if !ValidateDomain(parts[1]) {
		return "", "", errors.New("Invalid domain format")
	}

	return parts[0], parts[1], nil
}
