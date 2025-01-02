package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(secret []byte, userId string) (string, error) {
	duration:= time.Second * time.Duration(60*60*12)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"userId":userId,
		"expireAt": time.Now().Add(duration).Unix(),
	})
	tokenStr,err:=token.SignedString(secret)
	if err!=nil{
		return "",err
	}
	return tokenStr,nil
}