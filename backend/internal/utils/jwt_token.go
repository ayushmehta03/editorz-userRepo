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


