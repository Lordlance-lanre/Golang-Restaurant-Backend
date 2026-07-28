package controller

import(
	// "fmt"
	"github.com/gin-gonic/gin"
)

func GetUsers() gin.HandlerFunc {
	// fmt.Println("Get User Controller")
	return func(c *gin.Context){
		c.Writer.Write([]byte("Get User Controller"))
	}
}

func GetUserById() gin.HandlerFunc {
	return func(c *gin.Context){
		c.Writer.Write([]byte("Get User By Id Controller"))
	}
}


func Signup() gin.HandlerFunc {
	return func(c *gin.Context){
		c.Writer.Write([]byte("Signup Controller"))
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context){
		c.Writer.Write([]byte("Login Controller"))
	}
}

func HashPassword(password string) string {
	// TODO: Implement password hashing
	return password
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	// TODO: Implement password verification
	return true, "Password Verified"
}