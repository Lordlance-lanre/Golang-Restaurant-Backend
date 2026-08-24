package main

import (
	"fmt"
	"os"

	"github.com/Lordlance-lanre/Golang-Restaurant-Backend/database"
	"github.com/Lordlance-lanre/Golang-Restaurant-Backend/middleware"
	"github.com/Lordlance-lanre/Golang-Restaurant-Backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
	}
	port := os.Getenv("PORT")
	fmt.Println("Port:", port)

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)

	router.Use(middleware.Authentication())

	routes.FoodRoutes(router)
	routes.OrderRoutes(router)
	// routes.CartRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderItemsRoutes(router)
	routes.InvoiceRoutes(router)

	// 3. Start Gin Server
	router.Run(":" + port)

}
