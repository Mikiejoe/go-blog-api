package middlewares

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Mikiejoe/go-blog-api/config"
	"github.com/Mikiejoe/go-blog-api/types"
	"github.com/Mikiejoe/go-blog-api/utils"
	"github.com/golang-jwt/jwt/v5"
)
type contextKey string
const UserKey contextKey = "userId"

func AuthMiddleWare(handlerfunc http.HandlerFunc,u types.UserInTerface) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr:=getToken(r)
		token,err:=validateToken(tokenStr)
		if err!=nil{
			log.Println("error validation is ",err)
			permissionDenied(w)
			return
		}
		if !token.Valid{
			log.Println("invalid token")
			permissionDenied(w)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		uIdstr := claims["userId"].(string)
		user,err := u.GetUserByID(uIdstr)
		if err!=nil{
			log.Println("error user not found")
			permissionDenied(w)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx,UserKey,user.ID.Hex())
		r = r.WithContext(ctx)

		handlerfunc(w,r)
	}
}

func getToken(r *http.Request) string{
	token:= r.Header.Get("Authorization")
	if token !=""{
		return token
	}
	return ""
}

func validateToken(t string)(*jwt.Token,error){
	return jwt.Parse(t,func(t *jwt.Token)(interface{},error){
		if _,ok :=t.Method.(*jwt.SigningMethodHMAC);!ok{
			return nil,fmt.Errorf("unexpected signing method %v",t.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret),nil
	})

}

func permissionDenied(w http.ResponseWriter){
	utils.WriteError(w,http.StatusUnauthorized,fmt.Errorf("invalid token"))
}

func GetUseridFromCtx(ctx context.Context) (string){
	userID:=ctx.Value(UserKey).(string)
	
	return userID
}