package routes

import(
	// "fmt"
	"github.com/gin-gonic/gin"
	controllers "github.com/Lordlance-lanre/Golang-Restaurant-Backend/controllers"
)

func MenuRoutes(incomingRoutes *gin.Engine){
	routesForMenu := incomingRoutes.Group("/menu")

	routesForMenu.GET("/", controllers.GetMenus())
	routesForMenu.GET("/:menu_id", controllers.GetMenuById())

	routesForMenu.POST("/create-menu", controllers.CreateMenu())
	routesForMenu.PATCH("/:menu_id", controllers.UpdateMenu())

	routesForMenu.DELETE("/:menu_id", controllers.DeleteMenu())
}