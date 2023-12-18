package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

var jwtSecret = []byte(viper.GetString("server.jwtSecret"))

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

func GenerateToken(userID int64) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour)
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "douyin",
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
	return token, err
}

// 解析token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
