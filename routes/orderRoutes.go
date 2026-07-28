package routes

import(
	// "fmt"
	"github.com/gin-gonic/gin"
	controllers "github.com/Lordlance-lanre/Golang-Restaurant-Backend/controllers"
)

func OrderRoutes(incomingRoutes *gin.Engine){
	routesForOrders := incomingRoutes.Group("/orders")

	routesForOrders.GET("/", controllers.GetOrders())
	routesForOrders.GET("/:order_id", controllers.GetOrderById())

	routesForOrders.POST("/create-order", controllers.CreateOrder())
	routesForOrders.PATCH("/:order_id", controllers.UpdateOrder())

	routesForOrders.DELETE("/:order_id", controllers.DeleteOrder())
}