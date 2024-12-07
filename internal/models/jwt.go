package models

import "github.com/golang-jwt/jwt/v5"

type JWT struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	GUID  string `json:"guid"`
	IP    string `json:"ip"`
	jwt.RegisteredClaims
}

type User struct {
	GUID string
	IP string
	RefreshTokenHash string
	Email string
}
