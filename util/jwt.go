package util

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const secret = "some secret"

var ttl = time.Duration(1000 * time.Hour)

type Claims struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
	jwt.StandardClaims
}

func GenerateAuthToken(claims *Claims) string {
	token := newJWTToken(claims)
	tokenString, err := token.SignedString([]byte(secret))
	CheckErr(err, "token.SignedString")

	return tokenString
}

func newJWTToken(claims *Claims) *jwt.Token {
	currentTime := time.Now()
	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: currentTime.Add(ttl).Unix(),
		IssuedAt:  currentTime.Unix(),
		Issuer:    "task2trip",
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}

func ParseJWT(tokenStr string, claims jwt.Claims) error {
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		b := []byte(secret)
		return b, nil
	})
	return err
}
