package handlers

import (
	"carigo/internal/application/dto"
	"carigo/internal/application/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	createCustomerUC *usecases.CreateCustomerUseCase
	listCustomersUC  *usecases.ListCustomersUseCase
}

func NewCustomerHandler(create *usecases.CreateCustomerUseCase, list *usecases.ListCustomersUseCase) *CustomerHandler {
	return &CustomerHandler{
		createCustomerUC: create,
		listCustomersUC:  list,
	}
}

func (h *CustomerHandler) ShowCustomers(c *gin.Context) {
	customers, err := h.listCustomersUC.Execute(c.Request.Context())
	if err != nil {
		customers = []dto.CustomerDTO{}
	}

	c.HTML(http.StatusOK, "customers.html", gin.H{
		"Title":      "Müşteriler",
		"ActivePage": "customers",
		"Customers":  customers,
	})
}

func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var req dto.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.createCustomerUC.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}
