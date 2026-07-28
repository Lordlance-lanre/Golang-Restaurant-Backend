package controller

import(
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Lordlance-lanre/Golang-Restaurant-Backend/database"
	"github.com/Lordlance-lanre/Golang-Restaurant-Backend/models"
	"github.com/gin-gonic/gin"
	// "github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func GetMenus() gin.HandlerFunc{
	return func(c *gin.Context){
		// c.Writer.Write([]byte("Get Menus"))
		var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		result , err := menuCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"": err.Error()})
		}
		var allMenus []bson.M
		if err := result.All(ctx, &allMenus); err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occured: " + err.Error()})
		}
		c.JSON(http.StatusOK, allMenus)

	}
}

func GetMenuById() gin.HandlerFunc{
	return func(c *gin.Context){
		// c.Writer.Write([]byte("Get Menu By Id"))

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		menuId := c.Param("menu_id")

		var menu models.Menu

		err := foodCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while fetching menu: " + err.Error()})
		}
		c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc{
	return func(c *gin.Context){
		// c.Writer.Write([]byte("Create Menu"))

	var ctx,cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var menu models.Menu
	// var food models.Food

	if err := c.BindJSON(&menu); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	validationError := validate.Struct(menu)
	if validationError != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
		return
	}

	menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	menu.ID = primitive.NewObjectID()
	menuID := menu.ID.Hex()
	menu.Menu_id = &menuID


	result, err :=menuCollection.InsertOne(ctx,menu)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occured while creating menu: " + err.Error()})
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, result)

}
}


func UpdateMenu() gin.HandlerFunc{
	return func(c *gin.Context){
		// c.Writer.Write([]byte("Update Menu"))

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}
		menuId := c.Param("menu_id")

		filter := bson.M{"menu_id": menuId}
		
		var updateObj bson.D

		if !menu.Start_Date.IsZero() && !menu.End_Date.IsZero(){
			if !inTimeSpan(menu.Start_Date, menu.End_Date, time.Now()){
				msg := fmt.Sprintf("Please update the end date for menu")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "start_date", Value: menu.Start_Date})
			updateObj = append(updateObj, bson.E{Key: "end_date", Value: menu.End_Date})
		}

		if menu.Name != nil{
			updateObj = append(updateObj, bson.E{Key: "name", Value: menu.Name})
		}
		if menu.Category != nil{
			updateObj = append(updateObj, bson.E{Key: "category", Value: menu.Category})
		}
		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: menu.Updated_at})

		upsert := true
		opts := options.UpdateOne().SetUpsert(upsert)

		result, err := menuCollection.UpdateOne(ctx, filter, bson.D{bson.E{Key: "$set", Value: updateObj}}, opts)

		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Menu update failed: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)

	}
}

func DeleteMenu() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Write([]byte("Delete Menu"))
	}
}