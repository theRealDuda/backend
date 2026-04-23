package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func renderLogon(c *gin.Context, activeTab string, status int, errorMessage string) {
	if activeTab == "" {
		activeTab = "login"
	}

	data := gin.H{
		"activeTab": activeTab,
	}

	if errorMessage != "" {
		data["error"] = errorMessage
	}

	c.HTML(status, "logon.html", data)
}

func HomeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "brosok20.html", nil)
}
