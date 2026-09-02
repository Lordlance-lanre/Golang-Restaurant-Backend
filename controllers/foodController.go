package controller

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/Lordlance-lanre/Golang-Restaurant-Backend/database"
	"github.com/Lordlance-lanre/Golang-Restaurant-Backend/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
var validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Writer.Write([]byte("Food Controller"))

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))

		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}
		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))
		if err != nil {
			startIndex = (page - 1) * recordPerPage
		}

		matchStage := bson.D{{Key: "$match", Value: bson.D{}}}
		groupStage := bson.D{
			{
				Key: "$group",
				Value: bson.D{
					{Key: "_id", Value: "null"},
					{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
					{
						Key: "data",
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
						Key: "food_items",
						Value: bson.D{
							{
								Key: "$slice",
								Value: []interface{}{"$data", startIndex, recordPerPage},
							},
						},
					},
				},
			},
		}
		result, err := foodCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while fetching all food: " + err.Error()})
			return
		}

		var allFoods []bson.M
		if err = result.All(ctx, &allFoods); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allFoods)
	}
}

func GetFoodById() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Writer.Write([]byte("Food By Id Controller"))
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		foodId := c.Param("food_id")

		var food models.Food

		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while fetching food item: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, food)
	}
}

func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Writer.Write([]byte("Create Food Controller"))

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var menu models.Menu
		var food models.Food

		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		validationError := validate.Struct(food)
		if validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
			return
		}
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Menu ID"})
			return
		}

		food.Created_at, err = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.Updated_at, err = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		food.ID = primitive.NewObjectID()
		foodID := food.ID.Hex()
		food.Food_id = &foodID
		var num = toFixed(*food.Price, 2)
		food.Price = &num

		_, err = foodCollection.InsertOne(ctx, food)
		if err != nil {
			msg := "food item was not created"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, food)
	}
}

func roundNums(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(roundNums(num*output)) / output
}

func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Writer.Write([]byte("Update Food Controller"))

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var menu models.Menu
		var food models.Food

		foodId := c.Param("food_id")

		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj bson.D

		if food.Name != nil {
			updateObj = append(updateObj, bson.E{Key: "name", Value: food.Name})
		}
		if food.Price != nil {
			updateObj = append(updateObj, bson.E{Key: "price", Value: food.Price})
		}
		if food.Food_image != nil {
			updateObj = append(updateObj, bson.E{Key: "food_image", Value: food.Food_image})
		}
		if food.Menu_id != nil {
			err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
			if err != nil {
				msg := "Invalid Menu ID"
				c.JSON(http.StatusBadRequest, gin.H{"error": msg})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "menu_id", Value: food.Menu_id})
		}
		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: food.Updated_at})

		upsert := true
		filter := bson.M{"food_id": foodId}
		opts := options.UpdateOne().SetUpsert(upsert)
		_, err := foodCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			opts,
		)
		if err != nil {
			msg := "Food Item failed to fetch"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, updateObj)
	}
}

func DeleteFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		foodId := c.Param("food_id")
		filter := bson.M{"food_id": foodId}

		result, err := foodCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while deleting food item: " + err.Error()})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Food item not found"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
