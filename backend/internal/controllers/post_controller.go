package controllers

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
)


func CreatePost(client *mongo.Client)gin.HandlerFunc{
	return func(c *gin.Context){

		userId,exists:=c.Get("user_id")

		if !exists{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"Unauthorized"})
			return 
		}

		var post models.Post

		if err:=c.ShouldBindJSON(&post);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid input"})
			return 
		}



		authorObjId,err:=primitive.ObjectIDFromHex(userId.(string))

		if err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid user id "})
			return 
		}

		ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)

		defer cancel()

		var user models.User


		editorCol:=database.OpenCollection("editors",client)


		if err:=editorCol.FindOne(ctx,bson.M{"_id":authorObjId}).Decode(&user);err!=nil{
			c.JSON(http.StatusNotFound,gin.H{"error":"User not found"})
			return 
		}

		post.ID=primitive.NewObjectID()
		post.AuthorID=authorObjId
		post.Published=true
		post.Slug=utils.GenerateUniqueSlug(post.Title)
		post.CreatedAt=time.Now()
		post.UpdatedAt=time.Now()

		post.User=models.PostUser{
			ID: user.ID,
			Name: user.UserName,
			ProfileImage: user.ProfileImage,
		}

		postCol:=database.OpenCollection("posts",client)

		_,err=postCol.InsertOne(ctx,post)

		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to create post"})
			return 
		}
		

c.JSON(http.StatusCreated, gin.H{
	"message": "Post created successfully",
	"post":    post,
})

	}
}