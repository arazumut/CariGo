package main

import (
	"carigo/internal/application/ports"
	"carigo/internal/application/usecases"
	"carigo/internal/infrastructure/persistence/sqlite"
	"carigo/internal/interfaces/http/handlers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "carigo.db"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	baseRepo, custRepo, invRepo, payRepo, allocRepo, err := sqlite.NewRepositories(dbPath)
	if err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}
	
	realClock := ports.RealClock{}

	registerPaymentUC := usecases.NewRegisterPaymentUseCase(payRepo, invRepo, allocRepo, baseRepo, realClock)
	createInvoiceUC := usecases.NewCreateInvoiceUseCase(invRepo, realClock)
	listInvoicesUC := usecases.NewListInvoicesUseCase(invRepo)
	listPaymentsUC := usecases.NewListPaymentsUseCase(payRepo)
	dashboardStatsUC := usecases.NewGetDashboardStatsUseCase(payRepo, invRepo, custRepo)
	
	createCustomerUC := usecases.NewCreateCustomerUseCase(custRepo)
	listCustomersUC := usecases.NewListCustomersUseCase(custRepo)
	getCustomerStatementUC := usecases.NewGetCustomerStatementUseCase(custRepo, invRepo, payRepo)

	paymentHandler := handlers.NewPaymentHandler(registerPaymentUC, listPaymentsUC, listCustomersUC)
	invoiceHandler := handlers.NewInvoiceHandler(createInvoiceUC, listInvoicesUC, listCustomersUC)
	dashboardHandler := handlers.NewDashboardHandler(dashboardStatsUC)
	customerHandler := handlers.NewCustomerHandler(createCustomerUC, listCustomersUC, getCustomerStatementUC)

	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.Static("/assets", "./web/assets")
	r.LoadHTMLGlob("web/templates/**/*")

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP", "version": "MVP+"})
	})
	r.GET("/", dashboardHandler.ShowDashboard)
	r.GET("/invoices", invoiceHandler.ShowInvoices)
	r.GET("/payments", paymentHandler.ShowPayments)
	r.GET("/customers", customerHandler.ShowCustomers)
	r.GET("/customers/:id", customerHandler.ShowCustomerStatement)

	api := r.Group("/api/v1")
	{
		api.POST("/invoices", invoiceHandler.CreateInvoice)
		api.POST("/payments", paymentHandler.RegisterPayment)
		api.POST("/customers", customerHandler.CreateCustomer)
	}

	log.Printf("Starting server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
