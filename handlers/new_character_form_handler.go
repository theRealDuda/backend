package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewCharacterFormHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "character_edit.html", gin.H{
		"action": "/character/save",
		"method": "POST",
	})
}
