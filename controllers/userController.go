package controller

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	// "fmt"
	"github.com/Lordlance-lanre/Golang-Restaurant-Backend/database"
	helpers "github.com/Lordlance-lanre/Golang-Restaurant-Backend/helpers"
	"github.com/Lordlance-lanre/Golang-Restaurant-Backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	// "github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GetUsers() gin.HandlerFunc {
	// fmt.Println("Get User Controller")
	return func(c *gin.Context) {
		// c.Writer.Write([]byte("Get User Controller"))
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordsPerPage, err := strconv.Atoi(c.Query("recordsPerPage"))
		if err != nil || recordsPerPage < 1 {
			recordsPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordsPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))
		if err != nil {
			startIndex = (page - 1) * recordsPerPage
		}

		matchStage := bson.D{{Key: "$match", Value: bson.D{}}}
		groupStage := bson.D{
			{
				Key: "$group",
				Value: bson.D{
					{Key: "_id", Value: "null"},
					{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
					{
						Key: "user_data",
						Value: bson.D{
							{Key: "$push", Value: "$$ROOT"},
						},
					},
				},
			},
		}
		projectStage := bson.D{
			{
				Key: "$project",
				Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "total_count", Value: 1},
					{
						Key: "user_items",
						Value: bson.D{
							{
								Key:   "$slice",
								Value: []interface{}{"$user_data", startIndex, recordsPerPage},
							},
						},
					},
				},
			},
		}

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage,
			groupStage,
			projectStage,
		})
		if err != nil {
			msg := "Error occurred while listing the user items"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		var allUsers []bson.M
		if err = result.All(ctx, &allUsers); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allUsers[0])
	}
}

func GetUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Writer.Write([]byte("Get User By Id Controller"))

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Param("user_id")

		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		if err != nil {
			msg := "Error occurred while fetching user: " + err.Error()
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Writer.Write([]byte("Signup Controller"))
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			msg := "Please provide valid data"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return

		}

		validateErr := validate.Struct(user)
		if validateErr != nil {
			msg := validateErr.Error()
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		userCount, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			msg := err.Error()
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		if userCount > 0 {
			msg := "Email already exists"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		userCount, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			msg := err.Error()
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		if userCount > 0 {
			msg := "Phone number already exists"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		// token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, user.User_id, *user.First_name, *user.Last_name)
		// user.Token = &token
		// user.Refresh_token = &refreshToken

		_, err = userCollection.InsertOne(ctx, user)
		if err != nil {
			msg := err.Error()
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Writer.Write([]byte("Login Controller"))
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			msg := err.Error()
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)

		if err != nil {
			msg := "Invalid email or password provided"
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}

		passwordMatch, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if !passwordMatch {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, foundUser.User_id, *foundUser.First_name, *foundUser.Last_name)

		helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		c.JSON(http.StatusOK, gin.H{
			"access_token":  token,
			"refresh_token": refreshToken,
			"user_id": foundUser.User_id,
		})
	}
}

func HashPassword(password string) string {
	// TODO: Implement password hashing
	// return password
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		msg := "Error occurred while hashing password"
		log.Println(msg)
		return ""
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	// TODO: Implement password verification
	// return true, "Password Verified"
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	if err != nil {
		msg := "invalid password"
		log.Println(msg)
		return false, msg
	}
	return true, "Password Verified"
}
