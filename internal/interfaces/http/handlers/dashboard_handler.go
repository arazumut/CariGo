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
		// In production, you might show an error page or log it.
		// For now, we continue with zero values or handle gracefully
		// Since ShowDashboard is a page load, we should probably still render, maybe with an error flash.
		// Let's assume zero stats on error for MVP stability.
		stats = &usecases.DashboardStats{}
	}

	// Helper to format cents to currency string (very basic)
    // In template or here. Let's send raw cents and let template logic or simple division handle it if possible.
	// Or format here.
	formattedTotal := float64(stats.TotalCollected) / 100.0

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"Title":      "Dashboard",
		"ActivePage": "dashboard",
		"Stats": map[string]interface{}{
			"TotalCollected": formattedTotal,
			"OpenInvoices":   stats.OpenInvoices,
		},
	})
}
