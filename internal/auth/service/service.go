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

func NewService(repo Repository, user UserClient, notif NotificationClient) Service {
	return &service{
		repo:  repo,
		user:  user,
		notif: notif,
	}
}

type service struct {
	repo  Repository
	user  UserClient
	notif NotificationClient
}

func (s *service) Login(req LoginReq) (string, error) {
	// Check registered email
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		log.WithFields(log.Fields{
			"err":   err,
			"email": req.Email,
		}).Errorln("[Login] Failed get user by email")
		return "", constants.ErrInternalServer
	}
	if user == nil {
		return "", constants.ErrEmailNotRegistered
	}

	// Check password
	if !comparePassword(user.Password, req.Password) {
		return "", constants.ErrPasswordWrong
	}

	// Add push token for notification
	if err := s.notif.AddPushToken(user.ID, req.PushToken); err != nil {
		log.WithFields(log.Fields{
			"err":       err,
			"id":        user.ID,
			"pushToken": req.PushToken,
		}).Errorln("[Login] Failed add push token to user service")
		return "", constants.ErrInternalServer
	}

	// Generate token token
	payload := Payload{
		ID: user.ID,
	}
	token := generateJWT(payload)

	return token, nil
}

func (s *service) Register(req RegisterReq) (string, error) {
	// validate request
	if err := req.Validate(); err != nil {
		return "", err
	}

	// Check registered email
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		log.WithFields(log.Fields{
			"err":   err,
			"email": req.Email,
		}).Errorln("[Register] Failed get user by email from user service")
		return "", constants.ErrInternalServer
	}
	if user != nil {
		return "", constants.ErrEmailRegistered
	}

	// Create user
	user, err = s.user.Create()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Errorln("[Register] Failed create user")
		return "", constants.ErrInternalServer
	}
	hash, err := hashAndSalt(req.Password)
	if err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"password": req.Password,
		}).Errorln("[Register] Failed hash password")
		return "", constants.ErrInternalServer
	}
	user.Email = req.Email
	user.Password = hash
	if err := s.repo.Save(*user); err != nil {
		log.WithFields(log.Fields{
			"err":   err,
			"email": req.Email,
		}).Errorln("[Register] Failed save user to repository")
		return "", constants.ErrInternalServer
	}

	// Add push token for notification
	if err := s.notif.AddPushToken(user.ID, req.PushToken); err != nil {
		log.WithFields(log.Fields{
			"err":       err,
			"id":        user.ID,
			"pushToken": req.PushToken,
		}).Errorln("[Register] Failed add push token")
		return "", constants.ErrInternalServer
	}

	// Generate JWT Token
	payload := Payload{
		ID:        user.ID,
		PushToken: req.PushToken,
	}
	token := generateJWT(payload)

	return token, nil
}

func (s *service) Logout(req LogoutReq) error {
	payload, err := s.Auth(AuthReq{Token: req.Token})
	if err != nil {
		log.WithFields(log.Fields{
			"err":   err,
			"token": req.Token,
		}).Errorln("[Logout] Failed get jwt payload")
		return constants.ErrInternalServer
	}

	// Remove push token from notification
	if err := s.notif.RemovePushToken(payload.ID, payload.PushToken); err != nil {
		log.WithFields(log.Fields{
			"err":       err,
			"id":        payload.ID,
			"pushToken": payload.PushToken,
		}).Errorln("[Logout] Failed remove push token")
		return constants.ErrInternalServer
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
		return nil, constants.ErrInternalServer
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
	return token
}

func decodeJWT(token string) (*Payload, error) {
	claims := new(jwtClaims)

	// TODO: handle parsing error
	jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return token, nil
	})
	return &claims.Payload, nil
}
