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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// verify email otp given by user with the otp we generated during registration

// verify email otp and auto-send phone otp
func VerifyOtpEmail(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req struct {
			Email string `json:"email"`
			OTP   string `json:"otp"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or otp"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		editorCollection := database.OpenCollection("editors", client)

		var user models.User

		if err := editorCollection.FindOne(ctx, bson.M{
			"email": req.Email,
		}).Decode(&user); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		if user.IsEmailVerified {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already verified"})
			return
		}

		if time.Now().After(user.OtpExpiry) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "OTP expired"})
			return
		}

		if err := bcrypt.CompareHashAndPassword(
			[]byte(user.OtpHash),
			[]byte(req.OTP),
		); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid otp"})
			return
		}

		// generate PHONE otp
		phoneOtp := utils.GenerateOtp()
		phoneOtpHash, _ := utils.HashPassword(phoneOtp)

		// single atomic update
		update := bson.M{
			"$set": bson.M{
				"is_email_verified": true,
				"otp_hash":          phoneOtpHash,
				"otp_expiry":        time.Now().Add(5 * time.Minute),
				"updated_at":        time.Now(),
			},
		}

		if _, err := editorCollection.UpdateOne(
			ctx,
			bson.M{"_id": user.ID},
			update,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Verification failed"})
			return
		}

		// send PHONE otp
		go utils.SendSMSOTP(user.Phone, phoneOtp)

		c.JSON(http.StatusOK, gin.H{
			"message": "Email verified successfully. Phone OTP sent.",
			"user_id": user.ID.Hex(),
		})
	}
}



func VerifyPhoneOTP(client *mongo.Client)gin.HandlerFunc{
	return func(c *gin.Context){

		var body struct{
			UserID string `json:"user_id" validate:"required"`
			OTP string `json:"otp" validate:"required"`

		}

		if err:=c.ShouldBindJSON(&body);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid input"})
			return 
		}

		userId,_:=primitive.ObjectIDFromHex(body.UserID)

		ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
		defer cancel()


		editorCollection:=database.OpenCollection("editors",client)


		var user models.User

		if err:=editorCollection.FindOne(ctx,bson.M{
			"_id":userId,
		}).Decode(&user);err!=nil{
			c.JSON(http.StatusNotFound,gin.H{"error":"User not found"})
			return 
		}
		if user.IsPhoneVerified{
			c.JSON(http.StatusBadRequest,gin.H{"error":"User already verified with phone"})
			return 
		}


		if time.Now().After(user.OtpExpiry){
			c.JSON(http.StatusBadRequest,gin.H{"error":"Otp expired"})
			return 
		}

		if err := bcrypt.CompareHashAndPassword(
			[]byte(user.OtpHash),
			[]byte(body.OTP),
		); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid otp"})
			return
		}

		update := bson.M{
			"$set": bson.M{
				"is_phone_verified": true,
				"updated_at":       time.Now(),
			},
			"$unset": bson.M{
				"otp_hash":   "",
				"otp_expiry": "",
			},
		}
		if _, err := editorCollection.UpdateOne(
			ctx,
			bson.M{"_id": userId},
			update,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Verification failed"})
			return
		}

		token, err := utils.GenerateToken(user.ID.Hex(), user.Email, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Phone number verified",
			"token":   token,
		})
	}
}