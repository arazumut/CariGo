package handlers

import (
	"carigo/internal/application/dto"
	"carigo/internal/application/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	registerPaymentUC *usecases.RegisterPaymentUseCase
}

func NewPaymentHandler(uc *usecases.RegisterPaymentUseCase) *PaymentHandler {
	return &PaymentHandler{registerPaymentUC: uc}
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
