package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"
	"unicode"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// GLOBAL, TREAT LIKE CONSTANT
var globalJWTSigningMethod = jwt.SigningMethodHS256

const (
	// MinPassLength is the minimum length allowable for passwords
	MinPassLength = 8
	// MaxPassLength is the maxiumum length allowable for passwords
	MaxPassLength = 50

	// MinUsernameLength is the minimum length allowable for usernames
	MinUsernameLength = 3
	// MaxUsernameLength is the maxiumum length allowable for usernames
	MaxUsernameLength = 15

	// MinPasswordLength is the minimum length allowable for passwords
	MinPasswordLength = 8
)

// Claims for JWT
type Claims struct {
	Domain   string `json:"domain"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// ValidatePassword returns true if string is a valid password format
// At least one upper case English letter
// At least one lower case English letter
// At least one digit
// At least one special character
// Minimum eight in length
func ValidatePassword(password string) bool {
	length := 0
	number := false
	upper := false
	special := false
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsSpace(c):
			return false
		}
		length++
	}
	return (length > MinPasswordLength) && number && upper && special
}

// HashPassword hashes a password securely
func HashPassword(password string) (string, error) {
	const cost = 10 // min=4, max=31, default=10
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

// CheckPasswordHash checks a password hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// MakeToken makes a new jwt token for a user
func MakeToken(domain, username, secret string, nowUTC time.Time, expiresIn time.Duration) (string, error) {
	signingKey := []byte(secret)

	claims := Claims{
		Username: username,
		Domain:   domain,
		StandardClaims: jwt.StandardClaims{
			Issuer:    domain,
			IssuedAt:  nowUTC.Unix(),
			ExpiresAt: nowUTC.Add(expiresIn).Unix(),
		},
	}

	token := jwt.NewWithClaims(globalJWTSigningMethod, claims)
	return token.SignedString(signingKey)
}

// ValidateToken returns the username and domain associate with the token
// if it is valid
func ValidateToken(tokenString, secret string, nowUTC time.Time) (string, string, error) {
	signingKey := []byte(secret)

	claimsStructure := Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStructure,
		func(token *jwt.Token) (interface{}, error) { return signingKey, nil },
	)
	if err != nil {
		return "", "", errors.New("Couldn't parse claims")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", "", errors.New("Couldn't parse claims")
	}

	if claims.ExpiresAt < nowUTC.Unix() {
		return "", "", errors.New("Expired token")
	}

	return claims.Domain, claims.Username, nil
}

// Authorize takes the headers and returns the username of the user if the token is valid
func Authorize(headers http.Header, secret string) (string, string, error) {
	authString := headers.Get("Authorization")
	splitAuth := strings.Split(authString, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", "", errors.New("Malformed authorization header")
	}

	// Verify the token
	domain, username, err := ValidateToken(splitAuth[1], secret, time.Now().UTC())
	if err != nil {
		return "", "", errors.New("Invalid JWT, can't authorize")
	}
	return domain, username, nil
}
