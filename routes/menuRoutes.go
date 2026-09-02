package routes

import(
	// "fmt"
	"github.com/gin-gonic/gin"
	controllers "github.com/Lordlance-lanre/Golang-Restaurant-Backend/controllers"
	"github.com/Lordlance-lanre/Golang-Restaurant-Backend/middleware"
)

func MenuRoutes(incomingRoutes *gin.Engine){
	routesForMenu := incomingRoutes.Group("api/menu")

	routesForMenu.GET("/", middleware.Authentication(), controllers.GetMenus())
	routesForMenu.GET("/:menu_id", middleware.Authentication(), controllers.GetMenuById())

	routesForMenu.POST("/create_menu", middleware.Authentication(),controllers.CreateMenu())
	routesForMenu.PATCH("/:menu_id", middleware.Authentication(),controllers.UpdateMenu())

	routesForMenu.DELETE("/:menu_id", middleware.Authentication(),controllers.DeleteMenu())
}