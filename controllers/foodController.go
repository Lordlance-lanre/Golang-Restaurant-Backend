package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Lordlance-lanre/Golang-Restaurant-Backend/database"
	"github.com/Lordlance-lanre/Golang-Restaurant-Backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"github.com/go-playground/validator/v10"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
var validate = validator.New()

func GetFoods() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Write([]byte("Food Controller"))
	}
}

func GetFoodById() gin.HandlerFunc{
	return func(c *gin.Context){
		// c.Writer.Write([]byte("Food By Id Controller"))
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		foodId := c.Param("food_id")

		var food models.Food

		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while fetching food item: " + err.Error()})
		}
		c.JSON(http.StatusOK, food)
	}
}

func CreateFood() gin.HandlerFunc{
	return func(c *gin.Context){
		// c.Writer.Write([]byte("Create Food Controller"))
		
	var ctx,cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var menu models.Menu
	var food models.Food

	if err := c.BindJSON(&food); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	validationError := validate.Struct(food)
	if validationError != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
		return
	}
	err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
	defer cancel()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Menu ID"})
		return
	}

	food.Created_at, err = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.Updated_at, err = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	food.ID = primitive.NewObjectID()
	food.Food_id = food.ID.Hex()
	var num = toFixed(*food.Price, 2)
	food.Price = &num


	result, err := foodCollection.InsertOne(ctx,food)
	if err != nil{
		msg := fmt.Sprintf("food item was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	defer cancel()

	c.JSON(http.StatusOK, result)
}

func roundNums(num float64) int{
	fmt.Println("Rounding Number")
	return int(num)
}

func toFixed(num float64, precision int) float64{
	fmt.Println("To Fixed")
	return num
}

func UpdateFood() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Write([]byte("Update Food Controller"))
	}
}

func DeleteFood() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Write([]byte("Delete Food Controller"))
	}
}
