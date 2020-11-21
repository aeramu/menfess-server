package jwt

import (
	"github.com/aeramu/menfess-server/internal/auth/constants"
	"github.com/aeramu/menfess-server/internal/auth/service"
	"github.com/dgrijalva/jwt-go"
)

var jwtSecretKey = []byte("Menfessui132jd98132dm&*6sajb23")

type jwtClaims struct {
	jwt.StandardClaims
	Payload service.Payload
}

func GenerateJWT(payload service.Payload) string {
	jwtClaims := &jwtClaims{
		Payload: payload,
	}
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims).SignedString(jwtSecretKey)
	return "Bearer " + token
}

func DecodeJWT(token string) (*service.Payload, error) {
	if len(token) < 8 {
		return nil, constants.ErrInvalidToken
	}
	token = token[7:]
	claims := new(jwtClaims)

	// TODO: handle parsing error
	jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return token, nil
	})
	return &claims.Payload, nil
}
