package controllers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ayushmehta03/editorz-userRepo/backend/internal/database"
	"github.com/ayushmehta03/editorz-userRepo/backend/internal/models"
	"github.com/ayushmehta03/editorz-userRepo/backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUser(client *mongo.Client)gin.HandlerFunc{
	return func(c* gin.Context){

		var user models.User

		if err:=c.ShouldBindJSON(&user);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid input data"})
			return 
		}

		validate:=validator.New()
		if err:=validate.Struct(user);err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": err.Error(),
			})
			return
		}

		hashedPassword,err:=utils.HashPassword(user.PasswordHash)

		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Intrenal server error"})
			return 
		}

		ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)

		defer cancel();


		editorCollection:=database.OpenCollection("editors",client)


		filter := bson.M{
	"$or": []bson.M{
		{"email": user.Email},
		{"username": user.UserName},
	},
}


	count,err:=editorCollection.CountDocuments(ctx,filter)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to check existing user"})
		return 
	}

	if count>0{
		c.JSON(http.StatusConflict,gin.H{"error":"User already exits"})
		return 
	}

	seed:=user.UserName
	if seed==" "{
		seed=user.Email
	} 

	avatarURL := fmt.Sprintf(
			"https://api.dicebear.com/7.x/initials/svg?seed=%s",
			url.QueryEscape(seed),
	)



	otp:=utils.GenerateOtp();

	otpHash,_:=utils.HashPassword(otp)




	otpPurpose:="email"


	if user.Email!=""{
		otpPurpose="email"
	}else if user.Phone!=""{
		otpPurpose="phone"
	}


	user.PasswordHash=hashedPassword
	user.IsEmailVerified=false
	user.IsPhoneVerified=false;
	user.ProfileImage=avatarURL
	user.OtpHash=otpHash
	user.IsHiringListed=false
	user.ShowOnHiringPage=false
	user.EmploymentStatus="Not mentioned"
	user.OtpPurpose=otpPurpose
	user.CreatedAt=time.Now()
	user.UpdatedAt=time.Now()
	user.OtpExpiry= time.Now().Add(10 * time.Minute)

	if _, err := editorCollection.InsertOne(ctx, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
			return
		}




if otpPurpose=="email"{
	go 
}









	}
}