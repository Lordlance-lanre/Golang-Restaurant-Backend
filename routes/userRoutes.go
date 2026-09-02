package routes

import(
	// "fmt"
	"github.com/gin-gonic/gin"
	controllers "github.com/Lordlance-lanre/Golang-Restaurant-Backend/controllers"
	"github.com/Lordlance-lanre/Golang-Restaurant-Backend/middleware"
)

func UserRoutes(incomingRoutes *gin.Engine){
	routesForUsers := incomingRoutes.Group("api/users")
	
	routesForUsers.GET("/all_users", middleware.Authentication(),controllers.GetUsers())
	routesForUsers.GET("/:user_id", middleware.Authentication(), controllers.GetUserById())
	
	routesForUsers.POST("/signup", controllers.Signup())
	routesForUsers.POST("/login", controllers.Login())

}