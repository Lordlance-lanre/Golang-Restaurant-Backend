package routes

import (
	// "fmt"
	controllers "github.com/Lordlance-lanre/Golang-Restaurant-Backend/controllers"
	"github.com/gin-gonic/gin"
)

func FoodRoutes(incomingRoutes *gin.Engine) {
	routeForFood := incomingRoutes.Group("api/food")

	routeForFood.GET("/", controllers.GetFoods())
	routeForFood.GET("/:food_id", controllers.GetFoodById())

	routeForFood.POST("/add_food", controllers.CreateFood())
	routeForFood.PATCH("/:food_id", controllers.UpdateFood())

	routeForFood.DELETE("/:food_id", controllers.DeleteFood())
}
