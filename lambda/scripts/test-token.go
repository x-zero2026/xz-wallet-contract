package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	DID      string `json:"did"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func main() {
	secret := "Ia7gdt+1znW6j9I9XXLg+//MbKYIMa3HW5X7Eqd3gho="
	did := "0x3070deb1c17432b094d30509ccbfd598fb2793435efdca9273dfbc558bc040ca"
	username := "admin"

	claims := Claims{
		DID:      did,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(tokenString)
}
