package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"

	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

const minSecretKeySize = 32

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	secretKey *rsa.PrivateKey
}

type JWTVerifier struct {
	publicKey *rsa.PublicKey
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(tokenPath string) (*JWTMaker, error) {

	// = domain.GoDotEnvVariable("SECRET_KEY_PATH" + "samp.pem")
	// path := common.GoDotEnvVariable("SECRET_KEY_PATH") + "samp.pem"
	content, err := ioutil.ReadFile(tokenPath + "/samp.pem")
	if err != nil {
		logrus.Fatal(err)
	}
	rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM(content)
	if err != nil {
		logrus.Fatal(err)
	}

	if rsaKey.Size() < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey: rsaKey}, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *JWTMaker) CreateJWTToken(userID string, duration time.Duration) (string, error) {
	payload := Payload{
		UserID:    userID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(time.Duration(time.Second) * 300),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, &payload)
	token, err := jwtToken.SignedString(maker.secretKey)
	return token, err
}

func NewJWTVerifier(tokenPath string) (*JWTVerifier, error) {

	content, err := ioutil.ReadFile(tokenPath + "/samp.pub")
	if err != nil {
		logrus.Fatal(err)
	}
	rsaKey, err := jwt.ParseRSAPublicKeyFromPEM(content)
	if err != nil {
		logrus.Fatal(err)
	}

	if rsaKey.Size() < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTVerifier{publicKey: rsaKey}, nil
}

func (verifier *JWTVerifier) GetMetaData(token string) (*Payload, error) {

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		return verifier.publicKey, nil
	})

	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil

}

// VerifyToken checks if the token is valid or not
func (verifier *JWTVerifier) IsValidToken(token string) (bool, error) {
	payload, err := verifier.GetMetaData(token)
	if err != nil {
		return false, err
	}

	return payload != nil, nil
}
