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
	"golang.org/x/crypto/bcrypt"
)

// register user api
func RegisterUser(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}

		validate := validator.New()
		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": err.Error(),
			})
			return
		}

		hashedPassword, err := utils.HashPassword(user.PasswordHash)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		editorCollection := database.OpenCollection("editors", client)

		filter := bson.M{
			"$or": []bson.M{
				{"email": user.Email},
				{"username": user.UserName},
			},
		}

		count, err := editorCollection.CountDocuments(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing user"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}

		seed := user.UserName
		if seed == "" {
			seed = user.Email
		}

		avatarURL := fmt.Sprintf(
			"https://api.dicebear.com/7.x/initials/svg?seed=%s",
			url.QueryEscape(seed),
		)

		// generate EMAIL otp
		otp := utils.GenerateOtp()
		otpHash, _ := utils.HashPassword(otp)

		user.PasswordHash = hashedPassword
		user.IsEmailVerified = false
		user.IsPhoneVerified = false
		user.ProfileImage = avatarURL
		user.OtpHash = otpHash
		user.OtpExpiry = time.Now().Add(10 * time.Minute)
		user.Role = "editor"
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		if _, err := editorCollection.InsertOne(ctx, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
			return
		}

		// send EMAIL otp (plain otp, async)
		go utils.SendEmailOtp(user.Email, otp)

		c.JSON(http.StatusCreated, gin.H{
			"message": "Account registered successfully. Please verify your email.",
		})
	}
}

// login via email or username + password

func LoginWithPassword(client *mongo.Client)gin.HandlerFunc{
	return func(c *gin.Context){

		var req struct{
			Identifier string `json:"identifier"`
			Password string `json:"password"`
		}

		if err:=c.ShouldBindJSON(&req);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid input"})
			return 
		}

		ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)

		defer cancel()


		editorCollection:=database.OpenCollection("editors",client)

		var user models.User

		filter:=bson.M{
			"$or":[]bson.M{
				{"email":req.Identifier},
				{"username":req.Identifier},
			},
		}

		if err:=editorCollection.FindOne(ctx,filter).Decode(&user);err!=nil{
			c.JSON(http.StatusNotFound,gin.H{"error":"No such user found"});
			return 
		}

		if err:=bcrypt.CompareHashAndPassword(
			[]byte(user.PasswordHash),
			[]byte(req.Password),
		);err!=nil{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"Invalid credentials"})
			return 
		}

		if !user.IsEmailVerified || !user.IsPhoneVerified{
			c.JSON(http.StatusForbidden,gin.H{"error":"Email or Phone Not verified"})
			return 
		}

		token, err := utils.GenerateToken(user.ID.Hex(), user.Email, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"token":   token,
		})

	}
}


