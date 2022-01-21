package service

import (
	"joranvest/commons"
	"joranvest/helper"
	"os"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"

	"fmt"
	"time"
)

//-- JWTService is What is JWT can do
type JWTService interface {
	GenerateToken(UserId string, EntityId string) string
	ValidateToken(token string) (*jwt.Token, error)
	GetUserByToken(token string) helper.UserIdentity
}

type jwtCustomClaim struct {
	UserId   string `json:"user_id"`
	EntityId string `json:"entity_id"`
	jwt.StandardClaims
}

type jwtService struct {
	secretKey string
	issuer    string
}

//-- New Instance JWT Service
func NewJWTService() JWTService {
	return &jwtService{
		issuer:    "putraops",
		secretKey: getSecretKey(),
	}
}

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRETKEY")
	if secretKey != "" {
		secretKey = "createdbyputraops"
	}
	return secretKey
}

func (j *jwtService) GenerateToken(UserID string, EntityId string) string {
	claims := &jwtCustomClaim{
		UserID,
		EntityId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Local().UTC().AddDate(1, 0, 0).Unix(),
			Issuer:    j.issuer,
			IssuedAt:  time.Now().Local().UTC().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (j *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", t_.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
}

func (c *jwtService) GetUserByToken(token string) helper.UserIdentity {
	commons.Logger()
	var userIdentity helper.UserIdentity
	aToken, err := c.ValidateToken(token)
	if err != nil {
		log.Error("jwt-service.go:GetUserByToken")
		log.Error(err.Error())
		panic(err.Error())
	}
	claims := aToken.Claims.(jwt.MapClaims)
	userIdentity.Token = token
	userIdentity.UserId = fmt.Sprintf("%v", claims["user_id"])
	userIdentity.EntityId = fmt.Sprintf("%v", claims["entity_id"])
	return userIdentity
}
