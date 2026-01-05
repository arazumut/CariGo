package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

func (h *DashboardHandler) ShowDashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "base.html", gin.H{
		"Title":      "Dashboard",
		"ActivePage": "dashboard",
	})
}
