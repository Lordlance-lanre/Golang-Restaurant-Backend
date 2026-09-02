package routes

import (
	// "fmt"
	controllers "github.com/Lordlance-lanre/Golang-Restaurant-Backend/controllers"
	"github.com/Lordlance-lanre/Golang-Restaurant-Backend/middleware"
	"github.com/gin-gonic/gin"
)

func OrderRoutes(incomingRoutes *gin.Engine) {
	routesForOrders := incomingRoutes.Group("api/orders")

	routesForOrders.GET("/", middleware.Authentication(), controllers.GetOrders())
	routesForOrders.GET("/:order_id", middleware.Authentication(), controllers.GetOrderById())

	routesForOrders.POST("/create_order", middleware.Authentication(), controllers.CreateOrder())
	routesForOrders.PATCH("/:order_id", middleware.Authentication(), controllers.UpdateOrder())

	routesForOrders.DELETE("/:order_id", middleware.Authentication(), controllers.DeleteOrder())
}
