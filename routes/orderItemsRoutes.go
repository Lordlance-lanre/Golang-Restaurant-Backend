package routes

import(
	// "fmt"
	"github.com/gin-gonic/gin"
	controllers "github.com/Lordlance-lanre/Golang-Restaurant-Backend/controllers"
)

func OrderItemsRoutes(incomingRoutes *gin.Engine){
	routesForOrderItems := incomingRoutes.Group("/orderItems")

	routesForOrderItems.GET("/", controllers.GetOrderItems())
	routesForOrderItems.GET("/:orderItem_id", controllers.GetOrderItemById())
	routesForOrderItems.GET("/order/:order_id", controllers.GetOrderItemsByOrder())
	routesForOrderItems.POST("/create-OrderItem", controllers.CreateOrderItem())
	routesForOrderItems.PATCH("/:orderItem_id", controllers.UpdateOrderItem())

	routesForOrderItems.DELETE("/:orderItem_id", controllers.DeleteOrderItem())
}