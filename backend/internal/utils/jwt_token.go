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


	// refrence for generating token 

	secret:=os.Getenv("JWT_SECRET")

	// place all the value from parameters with expiry and issue time 

	// after 24 hrs token will be invalid 

	claims:=JWTClaims{
		UserID: userId,
		Email: email,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24*time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			
		},
	}


	// geerate the toekn with signing method and claims refrence
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)


	// return the token string 

	return token.SignedString([]byte(secret))

}


// veify token if the toekn is still inside the header 

func VerifyToken(tokenStr string)(*JWTClaims,error){
	
	// get the secret
		secret:=os.Getenv("JWT_SECRET")

		// check the signing method algo ,secretkey and expiry as well as issued time

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

		// place the value inside claims struct

		claims,ok:=token.Claims.(*JWTClaims)

		if !ok || !token.Valid{
			return nil,jwt.ErrTokenInvalidClaims
		}

		// return the claims struct 

		// so that email,userid and role can be trusted

		return claims,nil


}


