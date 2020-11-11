package service

import (
	"log"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Login(req LoginReq) (string, error)
	Register(req RegisterReq) (string, error)
	Logout(req LogoutReq) error
	Auth(req AuthReq) (Payload, error)
}

func NewService(account UserClient) Service {
	return &service{
		account: account,
	}
}

type service struct {
	account UserClient
}

func (s *service) Login(req LoginReq) (string, error) {
	account, err := s.account.GetByEmail(req.Email)
	if err != nil {
		log.Println("User Client Error:", err)
		return "", err
	}
	if account == nil {
		return "Email doesn't registered", nil
	}
	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(req.Password)); err != nil {
		return "Wrong Password", nil
	}

	if err := s.account.AddPushToken(account.ID, req.PushToken); err != nil {
		log.Println("User Client Error:", err)
		return "", err
	}

	payload := Payload{
		ID: account.ID,
	}
	return generateJWT(payload), nil
}

func (s *service) Register(req RegisterReq) (string, error) {
	if req.Email == "" || req.Password == "" || req.PushToken == "" {
		return "Invalid Request", nil
	}
	account, err := s.account.GetByEmail(req.Email)
	if err != nil {
		log.Println("User Client Error:", err)
		return "", err
	}
	if account != nil {
		return "Email already registered", nil
	}
	account, err = s.account.Create(req.Email, hashAndSalt(req.Password), req.PushToken)
	if err != nil {
		log.Println("User Client Error:", err)
		return "", err
	}
	payload := Payload{
		ID: account.ID,
	}
	return generateJWT(payload), nil
}

func (s *service) Logout(req LogoutReq) error {
	if err := s.account.RemovePushToken(req.ID, req.PushToken); err != nil {
		log.Println("User Client Error:", err)
		return err
	}
	return nil
}

func (s *service) Auth(req AuthReq) (Payload, error) {
	payload, err := decodeJWT(req.Token)
	if err != nil {
		log.Println("Decode JWT Error:", err)
		return Payload{}, err
	}
	return payload, nil
}

func hashAndSalt(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Hash failed")
		return ""
	}
	return string(hash)
}

var jwtSecretKey = []byte("Menfessui132jd98132dm&*6sajb23")

type jwtClaims struct {
	jwt.StandardClaims
	Payload Payload
}

func generateJWT(payload Payload) string {
	jwtClaims := &jwtClaims{
		Payload: payload,
	}
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims).SignedString(jwtSecretKey)
	return "Bearer " + token
}

func decodeJWT(token string) (Payload, error) {
	token = token[7:]
	claims := new(jwtClaims)
	jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return token, nil
	})
	return claims.Payload, nil
}
