package models

import "github.com/golang-jwt/jwt/v5"

type JwtClaims struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}
