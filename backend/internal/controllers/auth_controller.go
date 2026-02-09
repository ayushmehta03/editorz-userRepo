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
func RegisterUser(client *mongo.Client)gin.HandlerFunc{


	return func(c* gin.Context){




		// copy the exact struct to take user input for editors info 

		var user models.User


		// should bind json -> it fills the data in the strcut coming from frontend response on the baisis of json tags

		if err:=c.ShouldBindJSON(&user);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid input data"})
			return 
		}

		// added the validator to veirfy the input for the struct coming from frontend

		validate:=validator.New()
		if err:=validate.Struct(user);err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": err.Error(),
			})
			return
		}

		// hashing password 

		hashedPassword,err:=utils.HashPassword(user.PasswordHash)

		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Intrenal server error"})
			return 
		}


		// setting context time of 10 seconds 
		ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)

		defer cancel();

		// accessing the editor collection inside the database 

		editorCollection:=database.OpenCollection("editors",client)

		// we will cehck the old user with two params username and email both must be unique 

		filter := bson.M{
	"$or": []bson.M{
		{"email": user.Email},
		{"username": user.UserName},
	},
}


	// count document will return if there is any exitsing user earlier 
	count,err:=editorCollection.CountDocuments(ctx,filter)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to check existing user"})
		return 
	}

	// if count is greater than 0 we will return conflict and close the program

	if count>0{
		c.JSON(http.StatusConflict,gin.H{"error":"User already exits"})
		return 
	}


	// required for default profile picture api on the basis of username it will generate random avtar for profile pic later on usercan update it

	seed:=user.UserName
	if seed==" "{
		seed=user.Email
	} 

	avatarURL := fmt.Sprintf(
			"https://api.dicebear.com/7.x/initials/svg?seed=%s",
			url.QueryEscape(seed),
	)



	// generating the otp and hashing it for verification purpose 

	otp:=utils.GenerateOtp();

	otpHash,_:=utils.HashPassword(otp)



	// otp purpose will start with email and then verify phone

	otpPurpose:="email"


	// putting up all the required values for insertion in the document 


	user.PasswordHash=hashedPassword
	user.IsEmailVerified=false
	user.IsPhoneVerified=false;
	user.ProfileImage=avatarURL
	user.OtpHash=otpHash
	user.IsHiringListed=false
	user.ShowOnHiringPage=false
	user.EmploymentStatus="Not mentioned"
	user.OtpPurpose=otpPurpose
	user.Role="editor"
	user.CreatedAt=time.Now()
	user.UpdatedAt=time.Now()

	// 10 minutes timer for otp expiry

	user.OtpExpiry= time.Now().Add(10 * time.Minute)


	// inserting the document into the user collection
	if _, err := editorCollection.InsertOne(ctx, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
			return
		}

		if otpPurpose=="email"{

		}















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


