package controller

import(
	// "fmt"
	"github.com/gin-gonic/gin"
)

func GetTables() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Write([]byte("Get Tables"))
	}
}

func GetTableById() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Write([]byte("Get Table By Id"))
	}
}

func CreateTable() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Write([]byte("Create Table"))
	}
}

func UpdateTable() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Write([]byte("Update Table"))
	}
}

func DeleteTable() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Write([]byte("Delete Table"))
	}
}