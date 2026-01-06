package handlers

import (
	"carigo/internal/application/dto"
	"carigo/internal/application/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InvoiceHandler struct {
	createInvoiceUC *usecases.CreateInvoiceUseCase
	listInvoicesUC  *usecases.ListInvoicesUseCase
	listCustomersUC *usecases.ListCustomersUseCase
}

func NewInvoiceHandler(createUC *usecases.CreateInvoiceUseCase, listUC *usecases.ListInvoicesUseCase, listCustUC *usecases.ListCustomersUseCase) *InvoiceHandler {
	return &InvoiceHandler{
		createInvoiceUC: createUC,
		listInvoicesUC:  listUC,
		listCustomersUC: listCustUC,
	}
}

func (h *InvoiceHandler) ShowInvoices(c *gin.Context) {
	invoices, err := h.listInvoicesUC.Execute(c.Request.Context())
	if err != nil {
		invoices = []dto.InvoiceDTO{}
	}
	
	// Fetch customers for the dropdown
	customers, err := h.listCustomersUC.Execute(c.Request.Context())
	if err != nil {
		customers = []dto.CustomerDTO{}
	}

	c.HTML(http.StatusOK, "invoices.html", gin.H{
		"Title":      "Faturalar",
		"ActivePage": "invoices",
		"Invoices":   invoices,
		"Customers":  customers,
	})
}

func (h *InvoiceHandler) CreateInvoice(c *gin.Context) {
	var req dto.CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.createInvoiceUC.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}
