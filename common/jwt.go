package common

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tshivanshu9/budget-be/internal/models"
)

type CustomJWTClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJWT(user models.UserModel) (*string, *string, error) {
	userClaims := CustomJWTClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	signedAccessToken, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomJWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})

	signedRefreshToken, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, nil, err
	}
	return &signedAccessToken, &signedRefreshToken, nil
}

func ParseJWT(token string) (*CustomJWTClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &CustomJWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		fmt.Println("Error in parsing jwt", err)
		return nil, err
	} else if claims, ok := parsedToken.Claims.(CustomJWTClaims); ok {
		return &claims, nil
	} else {
		return nil, errors.New("unknown jwt claims")
	}

}
