package controller

import (
	// "fmt"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Lordlance-lanre/Golang-Restaurant-Backend/database"
	"github.com/Lordlance-lanre/Golang-Restaurant-Backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetTables() gin.HandlerFunc{
	return func(c *gin.Context){
		// c.Writer.Write([]byte("Get Tables"))
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := tableCollection.Find(ctx, bson.M{})
		if err != nil {
			msg := "Error occured while listing the Table items"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		var allTables []bson.M
		if err = result.All(ctx, &allTables); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allTables)
	}
}

func GetTableById() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		tableId := c.Param("table_id")

		var table models.Table

		err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while fetching table item: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, table)
	}
}

func CreateTable() gin.HandlerFunc{
	return func(c *gin.Context){
		// c.Writer.Write([]byte("Create Table"))
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var table models.Table

		if err := c.BindJSON(&table); err != nil{
			msg := "Error occured while binding the json"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		validationErr := validate.Struct(table)

		if validationErr != nil{
			msg := "Invalid input"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		table.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.ID = primitive.NewObjectID()
		table.Table_id = table.ID.Hex()
		
		_, err := tableCollection.InsertOne(ctx, table)
		if err != nil {
			msg := "Order Item was not created"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		
		c.JSON(http.StatusOK, table)
		
	}
}

func UpdateTable() gin.HandlerFunc{
	return func(c *gin.Context){
		// c.Writer.Write([]byte("Update Table"))
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var table models.Table

		var updateObj bson.D
		tableId := c.Param("table_id")

		if err := c.BindJSON(&table); err != nil {
			msg := "Error Updating Order Item"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		if table.Number_of_guests != nil{
			updateObj = append(updateObj, bson.E{Key: "number_of_guests", Value: table.Number_of_guests})
		}

		if table.Table_number != nil{
			updateObj = append(updateObj, bson.E{Key: "table_number", Value: table.Table_number})
		}

		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: table.Updated_at})

		upsert := true
		filter := bson.M{"table_id": tableId}
		opts := options.UpdateOne().SetUpsert(upsert)

		result, err := tableCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			opts,
		)
		if err != nil {
			msg := "Tbale item update failed"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func DeleteTable() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		tableId := c.Param("table_id")
		filter := bson.M{"table_id": tableId}

		result, err := tableCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while deleting table: " + err.Error()})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}