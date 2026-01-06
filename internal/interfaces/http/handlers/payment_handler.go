package handlers

import (
	"carigo/internal/application/dto"
	"carigo/internal/application/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	registerPaymentUC *usecases.RegisterPaymentUseCase
	listPaymentsUC    *usecases.ListPaymentsUseCase
}

func NewPaymentHandler(registerUC *usecases.RegisterPaymentUseCase, listUC *usecases.ListPaymentsUseCase) *PaymentHandler {
	return &PaymentHandler{
		registerPaymentUC: registerUC,
		listPaymentsUC:    listUC,
	}
}

func (h *PaymentHandler) ShowPayments(c *gin.Context) {
	payments, err := h.listPaymentsUC.Execute(c.Request.Context())
	if err != nil {
		payments = []dto.PaymentDTO{}
	}

	c.HTML(http.StatusOK, "payments.html", gin.H{
		"Title":      "Ödemeler",
		"ActivePage": "payments",
		"Payments":   payments,
	})
}

func (h *PaymentHandler) RegisterPayment(c *gin.Context) {
	var req dto.RegisterPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.registerPaymentUC.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
