package middleware

import (
	"fmt"
	"log"
	"option-dance/pkg/util"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var hmacSampleSecret = []byte("8f9dsg3fd3gadfbvnuujy1fdsc")
var ttl int64 = 3600 * 24 * 30 //Expiration time one month

type Payload struct {
	UserId    string
	ExpiredAt int64
}

func GenerateJWT(userId string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    userId,
		"expiredAt": time.Now().Unix() + ttl,
	})
	tokenString, e := token.SignedString(hmacSampleSecret)
	if e != nil {
		log.Print("error", e)
	}
	return tokenString
}

func ParseJWT(tokenString string) (payload Payload) {
	if tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		payload.UserId = util.Interface2String(claims["userId"])
		payload.ExpiredAt = int64(util.Interface2Float64(claims["expiredAt"]))
	} else {
		fmt.Println(err)
	}
	return
}
