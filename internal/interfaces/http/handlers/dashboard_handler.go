package handlers

import (
	"carigo/internal/application/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	statsUC *usecases.GetDashboardStatsUseCase
}

func NewDashboardHandler(uc *usecases.GetDashboardStatsUseCase) *DashboardHandler {
	return &DashboardHandler{statsUC: uc}
}

func (h *DashboardHandler) ShowDashboard(c *gin.Context) {
	stats, err := h.statsUC.Execute(c.Request.Context())
	if err != nil {
		stats = &usecases.DashboardStats{}
	}

	formattedTotal := float64(stats.TotalCollected) / 100.0
	formattedRevenue := float64(stats.TotalRevenue) / 100.0
	formattedPending := float64(stats.PendingBalance) / 100.0

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"Title":      "Dashboard",
		"ActivePage": "dashboard",
		"Stats": map[string]interface{}{
			"TotalCollected": formattedTotal,
			"OpenInvoices":   stats.OpenInvoices,
			"TotalRevenue":   formattedRevenue,
			"TotalCustomers": stats.TotalCustomers,
			"PendingBalance": formattedPending,
		},
	})
}
