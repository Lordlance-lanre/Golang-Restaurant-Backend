package routes

import(
	// "fmt"
	"github.com/gin-gonic/gin"
	controllers "github.com/Lordlance-lanre/Golang-Restaurant-Backend/controllers"
)

func InvoiceRoutes(incomingRoutes *gin.Engine){
	routesForInvoices := incomingRoutes.Group("api/invoices")

	routesForInvoices.GET("/", controllers.GetInvoices())
	routesForInvoices.GET("/:invoice_id", controllers.GetInvoiceById())

	routesForInvoices.POST("/create-invoice", controllers.CreateInvoice())
	routesForInvoices.PATCH("/:invoice_id", controllers.UpdateInvoice())

	routesForInvoices.DELETE("/:invoice_id", controllers.DeleteInvoice())
}