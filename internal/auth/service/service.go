package service

import (
	"github.com/aeramu/menfess-server/internal/auth/constants"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockery --all --output=../mock --case=underscore

type Service interface {
	Login(req LoginReq) (string, error)
	Register(req RegisterReq) (string, error)
	Logout(req LogoutReq) error
	Auth(req AuthReq) (*Payload, error)
}

func NewService(user UserClient, notif NotificationClient) Service {
	return &service{
		user:  user,
		notif: notif,
	}
}

type service struct {
	user  UserClient
	notif NotificationClient
}

func (s *service) Login(req LoginReq) (string, error) {
	// Check registered email
	user, err := s.user.GetByEmail(req.Email)
	if err != nil {
		log.WithFields(log.Fields{
			"err":   err,
			"email": req.Email,
		}).Errorln("[Login] Failed get user by email")
		return "", err
	}
	if user == nil {
		return "Email doesn't registered", nil
	}

	// Check password
	if !comparePassword(user.Password, req.Password) {
		return "Wrong Password", nil
	}

	// Add push token for notification
	if err := s.notif.AddPushToken(user.ID, req.PushToken); err != nil {
		log.WithFields(log.Fields{
			"err":       err,
			"id":        user.ID,
			"pushToken": req.PushToken,
		}).Errorln("[Login] Failed add push token to user service")
		return "", err
	}

	// Generate token token
	payload := Payload{
		ID: user.ID,
	}
	token := generateJWT(payload)

	return token, nil
}

func (s *service) Register(req RegisterReq) (string, error) {
	// Validate request
	// TODO: create validate method on register req
	if req.Email == "" || req.Password == "" || req.PushToken == "" {
		return "Invalid Request", nil
	}

	// Check registered email
	user, err := s.user.GetByEmail(req.Email)
	if err != nil {
		log.WithFields(log.Fields{
			"err":   err,
			"email": req.Email,
		}).Errorln("[Register] Failed get user by email from user service")
		return "", err
	}
	if user != nil {
		return "Email already registered", nil
	}

	// Create user
	hash, err := hashAndSalt(req.Password)
	if err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"password": req.Password,
		}).Errorln("[Register] Failed hash password")
		return "", err
	}
	user, err = s.user.Create(req.Email, hash, req.PushToken)
	if err != nil {
		log.WithFields(log.Fields{
			"err":       err,
			"email":     req.Email,
			"pushToken": req.PushToken,
		}).Errorln("[Register] Failed create user")
		return "", err
	}

	// Add push token for notification
	if err := s.notif.AddPushToken(user.ID, req.PushToken); err != nil {
		log.WithFields(log.Fields{
			"err":       err,
			"id":        user.ID,
			"pushToken": req.PushToken,
		}).Errorln("[Register] Failed add push token")
		return "", err
	}

	// Generate JWT Token
	payload := Payload{
		ID: user.ID,
	}
	token := generateJWT(payload)

	return token, nil
}

func (s *service) Logout(req LogoutReq) error {
	// Remove push token from notification
	if err := s.notif.RemovePushToken(req.ID, req.PushToken); err != nil {
		log.WithFields(log.Fields{
			"err":       err,
			"id":        req.ID,
			"pushToken": req.PushToken,
		}).Errorln("[Logout] Failed remove push token")
		return err
	}

	return nil
}

func (s *service) Auth(req AuthReq) (*Payload, error) {
	// Decode jwt
	payload, err := decodeJWT(req.Token)
	if err != nil {
		log.WithFields(log.Fields{
			"err":   err,
			"token": req.Token,
		}).Errorln("[Auth] Failed decode jwt token")
		return nil, err
	}

	return payload, nil
}

// hashAndSalt hash password
func hashAndSalt(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

// comparePassword return true when hash and password is equal
func comparePassword(hash string, pwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd)); err != nil {
		return false
	}
	return true
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

func decodeJWT(token string) (*Payload, error) {
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
