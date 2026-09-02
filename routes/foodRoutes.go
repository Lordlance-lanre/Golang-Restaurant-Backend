package routes

import (
	// "fmt"
	controllers "github.com/Lordlance-lanre/Golang-Restaurant-Backend/controllers"
	"github.com/gin-gonic/gin"
	"github.com/Lordlance-lanre/Golang-Restaurant-Backend/middleware"
)

func FoodRoutes(incomingRoutes *gin.Engine) {
	routeForFood := incomingRoutes.Group("api/food")

	routeForFood.GET("/", middleware.Authentication(), controllers.GetFoods())
	routeForFood.GET("/:food_id", middleware.Authentication(), controllers.GetFoodById())

	routeForFood.POST("/add_food", middleware.Authentication(), controllers.CreateFood())
	routeForFood.PATCH("/:food_id", middleware.Authentication(), controllers.UpdateFood())

	routeForFood.DELETE("/:food_id", middleware.Authentication(), controllers.DeleteFood())
}
