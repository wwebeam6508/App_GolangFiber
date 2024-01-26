package service

import (
	"PBD_backend_go/exception"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func generateJWT(payload jwt.MapClaims) (string, error) {
	//get exp time
	expTime := os.Getenv("JWT_ACCESS_TOKEN_EXPIRE")
	exp, err := time.ParseDuration(expTime)
	if err != nil {
		return "", err
	}
	payload["exp"] = time.Now().Add(exp).Unix()

	secretKey := os.Getenv("JWT_ACCESS_TOKEN_SECRET")
	signOption := jwt.SigningMethodHS256

	token := jwt.NewWithClaims(signOption, payload)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return "Bearer " + tokenString, nil
}

func generateRefreshJWT(payload jwt.MapClaims) (string, error) {
	//get exp time
	expTime := os.Getenv("JWT_REFRESH_TOKEN_EXPIRE")
	exp, err := time.ParseDuration(expTime)
	if err != nil {
		return "", err
	}
	payload["exp"] = time.Now().Add(exp).Unix()
	secretKey := os.Getenv("JWT_REFRESH_TOKEN_SECRET")
	signOption := jwt.SigningMethodHS256

	token := jwt.NewWithClaims(signOption, payload)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return "Bearer " + tokenString, nil
}

func VerifyJWT(tokenString string) (*jwt.Token, error) {
	secretKey := os.Getenv("JWT_ACCESS_TOKEN_SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, exception.ValidationError{Message: "invalid token"}
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func VerifyRefreshJWT(tokenString string) (*jwt.Token, error) {
	secretKey := os.Getenv("JWT_REFRESH_TOKEN_SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, exception.ValidationError{Message: "invalid token"}
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
