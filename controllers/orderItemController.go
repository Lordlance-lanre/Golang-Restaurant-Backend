package controller

import(
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetOrderItems() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Write([]byte("Get Order Items"))
	}
}

func GetOrderItemsByOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Write([]byte("Get Order Item By Order"))
	}
}

func ItemsByOrder(id string) (OrderItems []primitive.M, err error){
	fmt.Println("Items By Order")
	return
}

func GetOrderItemById() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Write([]byte("Get Order Item By Id"))
	}
}

func CreateOrderItem() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Write([]byte("Create Order Item"))
	}
}

func UpdateOrderItem() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Write([]byte("Update Order Item"))
	}
}

func DeleteOrderItem() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Write([]byte("Delete Order Item"))
	}
}