package Models

import (
	"time"
	"github.com/dgrijalva/jwt-go"
)

type TokenClaims struct {
	jwt.StandardClaims
	Role string `json:"role"`
	Csrf string `json:"csrf"`
}

const (
	RefreshTokenValidTime = time.Hour * 72
	AuthTokenValidTime = time.Minute * 15
)

