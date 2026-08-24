package routes

import(
	// "fmt"
	"github.com/gin-gonic/gin"
	controllers "github.com/Lordlance-lanre/Golang-Restaurant-Backend/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine){
	routesForUsers := incomingRoutes.Group("api/users")
	
	routesForUsers.GET("/all_users", controllers.GetUsers())
	routesForUsers.GET("/:user_id", controllers.GetUserById())
	
	routesForUsers.POST("/signup", controllers.Signup())
	routesForUsers.POST("/login", controllers.Login())

}