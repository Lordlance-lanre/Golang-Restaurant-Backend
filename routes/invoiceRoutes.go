package routes

import(
	// "fmt"
	"github.com/gin-gonic/gin"
	controllers "github.com/Lordlance-lanre/Golang-Restaurant-Backend/controllers"
	"github.com/Lordlance-lanre/Golang-Restaurant-Backend/middleware"
)

func InvoiceRoutes(incomingRoutes *gin.Engine){
	routesForInvoices := incomingRoutes.Group("api/invoices")

	routesForInvoices.GET("/", middleware.Authentication(), controllers.GetInvoices())
	routesForInvoices.GET("/:invoice_id", middleware.Authentication(), controllers.GetInvoiceById())

	routesForInvoices.POST("/create-invoice", middleware.Authentication(), controllers.CreateInvoice())
	routesForInvoices.PATCH("/:invoice_id", middleware.Authentication(), controllers.UpdateInvoice())

	routesForInvoices.DELETE("/:invoice_id", middleware.Authentication(), controllers.DeleteInvoice())
}