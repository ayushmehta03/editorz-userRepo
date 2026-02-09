package services

import (
	"context"
	"net/http"
	"time"

	"github.com/ayushmehta03/editorz-userRepo/backend/internal/database"
	"github.com/ayushmehta03/editorz-userRepo/backend/internal/models"
	"github.com/ayushmehta03/editorz-userRepo/backend/internal/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// verify email otp given by user with the otp we generated during registration

func VerifyOtpEmail(client *mongo.Client)gin.HandlerFunc{

	return func(c*gin.Context){


		// required fields email and otp

	var req struct{
		Email string `json:"email"`
		OTP string `json:"otp"`

	}


	// take the value from frontend inside the req struct

	if err:=c.ShouldBindJSON(&req);err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid email or otp"})
		return 
	}

	// context time out 
	ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()


	// connect to the collection in database

	editorCollection:=database.OpenCollection("editors",client)
		

	// user model copying for getting details from db search
	var user models.User


	// search in the editor collection with refrence to the email

	if err:=editorCollection.FindOne(ctx,bson.M{"email":req.Email}).Decode(&user);err!=nil{
		c.JSON(http.StatusNotFound,gin.H{"error":"User not found"})
		return 
	}

	// if user email is alreadfy verified no need to go ahead 

	if user.IsEmailVerified{
		c.JSON(http.StatusBadRequest,gin.H{"error":"User alredy verified with email"})
		return 
	}


	// if otp is expired wrt time stop execution

	if time.Now().After(user.OtpExpiry){
		c.JSON(http.StatusBadRequest,gin.H{"error":"Otp has expired"})
		return 
	}

	// comparing the saved hashed otp with the given otp

	if err:=bcrypt.CompareHashAndPassword([]byte(user.OtpHash),[]byte(req.OTP));err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid otp"})
		return 
	}


	// update the data inside the database after verification with this update query

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

	// generate the jwt token after successfull verification

	token,err:=utils.GenerateToken(user.ID.Hex(),user.Email,user.Role)

	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Token generation failed"})
		return 
	}

	// once done return the token and message 


	c.JSON(http.StatusOK,gin.H{
		"message":"Account verified",
		"token":token,
	})





	}
}