package controller

import(
	// "fmt"
	"github.com/gin-gonic/gin"
)

func GetOrders() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Write([]byte("Get Orders"))
	}
}

func GetOrderById() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Write([]byte("Get Order By Id"))
	}
}

func CreateOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Write([]byte("Create Order"))
	}
}

func UpdateOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Write([]byte("Update Order"))
	}
}

func DeleteOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Write([]byte("Delete Order"))
	}
}