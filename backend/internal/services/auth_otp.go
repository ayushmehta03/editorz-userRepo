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

		// request body struct for email and otp

		var req struct {
			Email string `json:"email" binding:"required"`
			OTP   string `json:"otp" binding:"required"`
		}


		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or otp"})
			return
		}
		// find the user with the email given in request body

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

		// already verified
		if user.IsEmailVerified {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already verified"})
			return
		}

		// email otp expiry
		if time.Now().After(user.OtpExpiry) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "OTP expired"})
			return
		}

		// verify EMAIL otp (local)
		if err := bcrypt.CompareHashAndPassword(
			[]byte(user.OtpHash),
			[]byte(req.OTP),
		); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid otp"})
			return
		}

		// SEND PHONE OTP VIA MESSAGE CENTRAL

		verificationID, err := utils.MessageCentralSendOTP(user.Phone)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send phone OTP"})
			return
		}

		// update user record: set email verified, save verificationId for phone otp and remove otp hash and expiry
		update := bson.M{
			"$set": bson.M{
				"is_email_verified": true,
				"verification_id":   verificationID,
				"updated_at":        time.Now(),
			},
			"$unset": bson.M{
				"otp_hash":   "",
				"otp_expiry": "",
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

		c.JSON(http.StatusOK, gin.H{
			"message": "Email verified successfully. Phone OTP sent.",
			"user_id": user.ID.Hex(),
		})
	}
}


		// verify phone otp given by user with the otp we generated during email verification and sent via message central

func VerifyPhoneOTP(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		// request body struct for user id and otp

		var req struct {
			UserID string `json:"user_id" binding:"required"`
			OTP    string `json:"otp" binding:"required"`
		}

		// bind the json body to the struct and validate required fields

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		userID, err := primitive.ObjectIDFromHex(req.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		editorCollection := database.OpenCollection("editors", client)

		// find the user with the id given in request body
		var user models.User
		if err := editorCollection.FindOne(ctx, bson.M{
			"_id": userID,
		}).Decode(&user); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// already verified
		if user.IsPhoneVerified {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number already verified"})
			return
		}

		// verificationId must exist (generated when OTP was sent)
		if user.VerificationID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No active phone verification found"})
			return
		}

		//  VERIFY OTP VIA MESSAGE CENTRAL
		if err := utils.MessageCentralVerifyOTP(user.VerificationID, req.OTP); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired OTP"})
			return
		}

		// update user record: set phone verified and remove verificationId
		update := bson.M{
			"$set": bson.M{
				"is_phone_verified": true,
				"updated_at":       time.Now(),
			},
			"$unset": bson.M{
				"verification_id": "",
			},
		}

		// update the user record in database with the above update document

		if _, err := editorCollection.UpdateOne(
			ctx,
			bson.M{"_id": userID},
			update,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Verification failed"})
			return
		}

		// generate auth token for the user after successful phone verification

		token,err:=utils.GenerateToken(user.ID.Hex(),user.Email,user.Role)
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Token generation failed"})
			return 
		}
		// return success response with the token

		c.JSON(http.StatusOK, gin.H{
			"message": "Phone number verified successfully. You can now log in.",
			"token": token,
		})
	}
}
