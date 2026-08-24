package routes

import(
	// "fmt"
	"github.com/gin-gonic/gin"
	controllers "github.com/Lordlance-lanre/Golang-Restaurant-Backend/controllers"
)

func MenuRoutes(incomingRoutes *gin.Engine){
	routesForMenu := incomingRoutes.Group("api/menu")

	routesForMenu.GET("/", controllers.GetMenus())
	routesForMenu.GET("/:menu_id", controllers.GetMenuById())

	routesForMenu.POST("/create_menu", controllers.CreateMenu())
	routesForMenu.PATCH("/:menu_id", controllers.UpdateMenu())

	routesForMenu.DELETE("/:menu_id", controllers.DeleteMenu())
}