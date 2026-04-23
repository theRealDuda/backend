package handlers

import (
	"brosok20/models"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DashboardHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := getUserID(c)
		if !ok {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		var characters []models.Character
		rows, err := db.Query("SELECT id, name FROM characters WHERE user_id = ?", userID)
		if err != nil {
			c.String(http.StatusInternalServerError, "Ошибка получения персонажей")
			return
		}
		defer rows.Close()

		for rows.Next() {
			var character models.Character
			if err := rows.Scan(&character.ID, &character.Name); err != nil {
				c.String(http.StatusInternalServerError, "Ошибка чтения персонажа")
				return
			}
			characters = append(characters, character)
		}

		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"characters": characters,
		})
	}
}
