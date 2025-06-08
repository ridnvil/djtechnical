package utils

import (
	"DeallsJobsTest/config"
	"DeallsJobsTest/models"
	"github.com/caarlos0/env/v11"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}

func GenerateJWTToken(user models.User) (string, error) {
	// generate token using custome claims
	var envCfg config.EnvConfig
	if errcnf := env.Parse(&envCfg); errcnf != nil {
		return "", errcnf
	}
	jwtSecret := []byte(envCfg.SecretKEY)
	claims := models.JwtClaims{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
		Fullname: user.FullName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "DeallsJobsTest",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
