package routes

import(
	"github.com/gin-gonic/gin"
	controllers "github.com/Lordlance-lanre/Golang-Restaurant-Backend/controllers"
	"github.com/Lordlance-lanre/Golang-Restaurant-Backend/middleware"
)

func TableRoutes(incomingRoutes *gin.Engine){
	routesForTable := incomingRoutes.Group("api/table")

	routesForTable.GET("/", middleware.Authentication(), controllers.GetTables())
	routesForTable.GET("/:table_id", middleware.Authentication(), controllers.GetTableById())

	routesForTable.POST("/create_table", middleware.Authentication(), controllers.CreateTable())
	routesForTable.PATCH("/:table_id", middleware.Authentication(), controllers.UpdateTable())

	routesForTable.DELETE("/:table_id", middleware.Authentication(), controllers.DeleteTable())
}