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
}

func NewInvoiceHandler(createUC *usecases.CreateInvoiceUseCase, listUC *usecases.ListInvoicesUseCase) *InvoiceHandler {
	return &InvoiceHandler{
		createInvoiceUC: createUC,
		listInvoicesUC:  listUC,
	}
}

func (h *InvoiceHandler) ShowInvoices(c *gin.Context) {
	invoices, err := h.listInvoicesUC.Execute(c.Request.Context())
	if err != nil {
		// MVP: Log error, return empty list or error page.
		// For now, let's return error page or just empty list with error flash check?
		// We'll just return HTML with empty list and maybe an alert if we had one.
		invoices = []dto.InvoiceDTO{}
	}

	c.HTML(http.StatusOK, "invoices.html", gin.H{
		"Title":      "Faturalar",
		"ActivePage": "invoices",
		"Invoices":   invoices,
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
