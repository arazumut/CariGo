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
	// 1. Config
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "carigo.db" // Local dev
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 2. Infrastrucure
	baseRepo, custRepo, invRepo, payRepo, allocRepo, err := sqlite.NewRepositories(dbPath)
	if err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}
	
	realClock := ports.RealClock{}

	// 3. Application (UseCases)
	registerPaymentUC := usecases.NewRegisterPaymentUseCase(payRepo, invRepo, allocRepo, baseRepo, realClock)
	createInvoiceUC := usecases.NewCreateInvoiceUseCase(invRepo, realClock)
	listInvoicesUC := usecases.NewListInvoicesUseCase(invRepo)
	listPaymentsUC := usecases.NewListPaymentsUseCase(payRepo)
	dashboardStatsUC := usecases.NewGetDashboardStatsUseCase(payRepo, invRepo)
	
	// Customer UseCases
	createCustomerUC := usecases.NewCreateCustomerUseCase(custRepo)
	listCustomersUC := usecases.NewListCustomersUseCase(custRepo)
	getCustomerStatementUC := usecases.NewGetCustomerStatementUseCase(custRepo, invRepo, payRepo)

	// 4. Interfaces (HTTP)
	paymentHandler := handlers.NewPaymentHandler(registerPaymentUC, listPaymentsUC, listCustomersUC)
	invoiceHandler := handlers.NewInvoiceHandler(createInvoiceUC, listInvoicesUC, listCustomersUC)
	dashboardHandler := handlers.NewDashboardHandler(dashboardStatsUC)
	customerHandler := handlers.NewCustomerHandler(createCustomerUC, listCustomersUC, getCustomerStatementUC)

	r := gin.Default()
	// Fix [GIN-debug] You trusted all proxies.
	// For MVP/Local, we trust no one or localhost.
	r.SetTrustedProxies(nil)

	// Static Files & Templates
	r.Static("/assets", "./web/assets")
	r.LoadHTMLGlob("web/templates/**/*")

	// Helper for checking active page in template (if needed, mostly passed via gin.H)

	// Initial Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP", "version": "MVP+"})
	})

	// UI Routes
	r.GET("/", dashboardHandler.ShowDashboard)
	r.GET("/invoices", invoiceHandler.ShowInvoices)
	r.GET("/payments", paymentHandler.ShowPayments)
	r.GET("/customers", customerHandler.ShowCustomers)
	r.GET("/customers/:id", customerHandler.ShowCustomerStatement)

	// API Routes
	api := r.Group("/api/v1")
	{
		api.POST("/invoices", invoiceHandler.CreateInvoice)
		api.POST("/payments", paymentHandler.RegisterPayment)
		api.POST("/customers", customerHandler.CreateCustomer)
	}

	// --------------------------------------------------------------------------------
	// SEED DATA (FOR DEMO/MVP ONLY - REMOVE IN PROD)
	// --------------------------------------------------------------------------------
	// Normally run via migration script, but here to allow USER to test immediately
	// seedData(custRepo, invRepo) 
	// (Skipped to avoid "SafeToAutoRun" confusing, but you can request it)
	
	log.Printf("Starting server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
