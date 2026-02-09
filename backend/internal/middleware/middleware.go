package middleware

import (
	"net/http"
	"strings"

	"github.com/ayushmehta03/editorz-userRepo/backend/internal/utils"
	"github.com/gin-gonic/gin"
)


func AuthMiddleWare()gin.HandlerFunc{
	return func(c *gin.Context){

		// get the header 

		authHeader:=c.GetHeader("Authorization")


		// if it is missing will abort 

		if authHeader==""{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"Authorization header missing"})
			c.Abort()
			return 
		}

		// validity check of the header and token verification

		parts:=strings.SplitN(authHeader," ",2)

		if len(parts)!=2 || parts[0]!="Bearer"{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"Invalid authorzation header"})
			c.Abort()
			return 
		}

		tokenString:=parts[1]

		claims,err:=utils.VerifyToken(tokenString)

		if err!=nil{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"Invalid or expired token"})
			c.Abort()
			return 
		}

		c.Set("user_id",claims.UserID)
		c.Set("email",claims.Email)
		c.Set("role",claims.Role)


		// if all satisfied allow to visit the page 
		
		c.Next()


	}
}