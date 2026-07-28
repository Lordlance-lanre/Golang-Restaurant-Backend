package controller

import(
	// "fmt"
	"github.com/gin-gonic/gin"
)

func GetInvoices() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Write([]byte("Get Invoices"))
	}
}

func GetInvoiceById() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Write([]byte("Get Invoice By Id"))
	}
}

func CreateInvoice() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Write([]byte("Create Invoice"))
	}
}

func UpdateInvoice() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Write([]byte("Update Invoice"))
	}
}

func DeleteInvoice() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Write([]byte("Delete Invoice"))
	}
}