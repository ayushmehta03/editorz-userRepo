package services

import (
	"context"
	"net/http"
	"time"

	"github.com/ayushmehta03/editorz-userRepo/backend/internal/database"
	"github.com/ayushmehta03/editorz-userRepo/backend/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func VerifyOtpEmail(client *mongo.Client)gin.HandlerFunc{

	return func(c*gin.Context){

	var req struct{
		Email string `json:"email"`
		OTP string `json:"otp"`

	}


	if err:=c.ShouldBindJSON(&req);err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid email or otp"})
		return 
	}

	ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	editorCollection:=database.OpenCollection("editors",client)
		
	var user models.User

	if err:=editorCollection.FindOne(ctx,bson.M{"email":req.Email}).Decode(&user);err!=nil{
		c.JSON(http.StatusNotFound,gin.H{"error":"User not found"})
		return 
	}

	if user.IsEmailVerified{
		c.JSON(http.StatusBadRequest,gin.H{"error":"User alredy verified with email"})
		return 
	}

	if time.Now().After(user.OtpExpiry){
		c.JSON(http.StatusBadRequest,gin.H{"error":"Otp has expired"})
		return 
	}


	if err:=bcrypt.CompareHashAndPassword([]byte(user.OtpHash),[]byte(req.OTP));err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid otp"})
		return 
	}


update := bson.M{
	"$set": bson.M{
		"is_email_verified": true,
		"updated_at":  time.Now(),
	},
	"$unset": bson.M{
		"otp_hash":   "",
		"otp_expiry": "",
	},
}

	_,err:=editorCollection.UpdateOne(ctx,bson.M{"email":req.Email},update)

	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Verification failed"})
		return 
	}


	


	}
}