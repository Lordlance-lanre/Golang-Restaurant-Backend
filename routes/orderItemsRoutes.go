package routes

import(
	// "fmt"
	"github.com/gin-gonic/gin"
	controllers "github.com/Lordlance-lanre/Golang-Restaurant-Backend/controllers"
	"github.com/Lordlance-lanre/Golang-Restaurant-Backend/middleware"
)

func OrderItemsRoutes(incomingRoutes *gin.Engine){
	routesForOrderItems := incomingRoutes.Group("api/orderItems")

	routesForOrderItems.GET("/", middleware.Authentication(), controllers.GetOrderItems())
	routesForOrderItems.GET("/:orderItem_id", middleware.Authentication(), controllers.GetOrderItemById())
	routesForOrderItems.GET("/order/:order_id", middleware.Authentication(), controllers.GetOrderItemsByOrder())
	routesForOrderItems.POST("/create-OrderItem", middleware.Authentication(), controllers.CreateOrderItem())
	routesForOrderItems.PATCH("/:orderItem_id", middleware.Authentication(), controllers.UpdateOrderItem())

	routesForOrderItems.DELETE("/:orderItem_id", middleware.Authentication(), controllers.DeleteOrderItem())
}