package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)



// required data fields for jwt token generation

type JWTClaims struct{
	UserID string `json:"user_id"`
	Email string `json:"email"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// generate the jwt token for authorization header 

func GenerateToken(userId,email,role string)(string,error){



	secret:=os.Getenv("JWT_SECRET")


	claims:=JWTClaims{
		UserID: userId,
		Email: email,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24*time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			
		},
	}

	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	return token.SignedString([]byte(secret))

}


func VerifyToken(tokenStr string)(*JWTClaims,error){
	
		secret:=os.Getenv("JWT_SECRET")

		token,err:=jwt.ParseWithClaims(
			tokenStr,
			&JWTClaims{},
			func (token *jwt.Token)(interface{},error){
				return []byte(secret),nil
			},
		)
		if err!=nil{
			return nil,err
		}

		claims,ok:=token.Claims.(*JWTClaims)

		if !ok || !token.Valid{
			return nil,jwt.ErrTokenInvalidClaims
		}

		return claims,nil


}


