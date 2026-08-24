package routes

import(
	"github.com/gin-gonic/gin"
	controllers "github.com/Lordlance-lanre/Golang-Restaurant-Backend/controllers"
)

func TableRoutes(incomingRoutes *gin.Engine){
	routesForTable := incomingRoutes.Group("api/table")

	routesForTable.GET("/", controllers.GetTables())
	routesForTable.GET("/:table_id", controllers.GetTableById())

	routesForTable.POST("/create-table", controllers.CreateTable())
	routesForTable.PATCH("/:table_id", controllers.UpdateTable())

	routesForTable.DELETE("/:table_id", controllers.DeleteTable())
}