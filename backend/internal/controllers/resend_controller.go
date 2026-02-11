package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ayushmehta03/editorz-userRepo/backend/internal/database"
	"github.com/ayushmehta03/editorz-userRepo/backend/internal/models"
	"github.com/ayushmehta03/editorz-userRepo/backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// resend the email otp with cooldown

func ResendEmailOtp(client *mongo.Client, redisClient *redis.Client)gin.HandlerFunc{
	return func (c *gin.Context){

		var req struct{
			Email string `json:"email" binding:"required"`

		}


		if err:=c.ShouldBindJSON(&req);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Email required"})
			return 
		}

		err:=utils.HandleOtpResendBackoff(
			redisClient,
			fmt.Sprintf("otp:email:%s", req.Email),
		)
		if err != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
			return
		}


		ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)

		defer cancel()

		editorCollection:=database.OpenCollection("editors",client)

		var user models.User

		if err:=editorCollection.FindOne(ctx,bson.M{
			"email":req.Email,
		}).Decode(&user);err!=nil{
			c.JSON(http.StatusNotFound,gin.H{"error":"User not found"})
			return 
		}

		if user.IsEmailVerified{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Email already verified"})
			return 
		}

		otp:=utils.GenerateOtp()

		otpHash,_:=utils.HashPassword(otp)


		_,err=editorCollection.UpdateOne(
			ctx,
			bson.M{"_id":user.ID},
			bson.M{
				"$set":bson.M{
					"otp_hash":otpHash,
					"otp_expiry":time.Now().Add(10*time.Minute),
					"updated_at":time.Now(),
				},
			},
		)

		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to resend otp"})
			return 
		}

		go utils.SendOTPEmail(user.Email,otp)

		c.JSON(http.StatusOK,gin.H{"message":"OTP resent successfully"})


	}
}

 // resend phone otp with cooldown

 func ResendPhoneOtp(client *mongo.Client,redisClient *redis.Client)gin.HandlerFunc{
	return func(c *gin.Context){

		var req struct{
			UserID string `json:"user_id" binding:"required"`
		}


		if err:=c.ShouldBindJSON(&req);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"User Id required"})
			return 
		}

		userID,err:=primitive.ObjectIDFromHex(req.UserID)

		if err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid user id"})
			return 
		}

		ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)

		defer cancel()

		editorCollection:=database.OpenCollection("editors",client)

		var user models.User

		if err:=editorCollection.FindOne(ctx,bson.M{"_id":userID}).Decode(&user);err!=nil{
			c.JSON(http.StatusNotFound,gin.H{"error":"User not found"})
			return 
		}
		if user.IsPhoneVerified {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone already verified"})
			return
		}

		countryCode := os.Getenv("MESSAGE_CENTRAL_COUNTRY_CODE")

		// Rate limit
		err = utils.HandleOtpResendBackoff(
			redisClient,
			fmt.Sprintf("otp:phone:%s:%s", countryCode, user.Phone),
		)
		if err != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
			return
		}

		verificationID, err := utils.MessageCentralSendOTP(user.Phone)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
			return
		}

		_, err = editorCollection.UpdateOne(
			ctx,
			bson.M{"_id": user.ID},
			bson.M{
				"$set": bson.M{
					"verification_id": verificationID,
					"updated_at":      time.Now(),
				},
			},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update verification"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Phone OTP resent successfully",
		})




	}
 }