package pkg

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	configs "github.com/nycholasmarques/socialdev/config"
)

var cfg = configs.LoadConfig()

func CreateToken(userId string, KeepLogin bool) (string, error) {
	var expiresIn int64
	if KeepLogin {
			expiresIn = time.Now().Add((time.Hour * 24) * 30).Unix()
	} else {
			expiresIn = time.Now().Add((time.Hour * 24) * 7).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.MapClaims{
					"user_id": userId,
					"exp":     expiresIn,
					"iat":     time.Now().Unix(),
			})

	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
			return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("token inválido")
	}

	return nil
}
