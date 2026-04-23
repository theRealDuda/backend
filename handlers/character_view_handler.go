package handlers

import (
	"brosok20/models"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CharacterViewHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, ok := getUserID(c)
		if !ok {
			if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Необходима авторизация"})
				return
			}
			c.Redirect(http.StatusFound, "/login")
			return
		}

		var character models.Character
		err := db.QueryRow(`
			SELECT id, user_id, name, race, class, level, strength, dexterity, constitution,
			       intelligence, wisdom, charisma, background, alignment, experience, inventory,
			       spells, features, notes
			FROM characters
			WHERE id = ? AND user_id = ?
			LIMIT 1`, id, userID).
			Scan(&character.ID, &character.UserID, &character.Name, &character.Race, &character.Class,
				&character.Level, &character.Strength, &character.Dexterity, &character.Constitution,
				&character.Intelligence, &character.Wisdom, &character.Charisma, &character.Background,
				&character.Alignment, &character.Experience, &character.Inventory, &character.Spells,
				&character.Features, &character.Notes)
		if err != nil {
			if err == sql.ErrNoRows {
				if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
					c.JSON(http.StatusNotFound, gin.H{"error": "Персонаж не найден"})
					return
				}
				c.String(http.StatusNotFound, "Персонаж не найден")
				return
			}

			if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
				return
			}
			c.String(http.StatusInternalServerError, "Ошибка сервера")
			return
		}

		if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
			c.JSON(http.StatusOK, character)
			return
		}

		c.Redirect(http.StatusFound, "/character/edit/"+id)
	}
}
