package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SaveCharacterHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := getUserID(c)
		if !ok || userID == 0 {
			c.String(http.StatusUnauthorized, "Необходима авторизация")
			return
		}

		name := c.PostForm("name")
		race := c.PostForm("race")
		class := c.PostForm("class")
		level, _ := strconv.Atoi(c.PostForm("level"))
		strength, _ := strconv.Atoi(c.PostForm("strength"))
		dexterity, _ := strconv.Atoi(c.PostForm("dexterity"))
		constitution, _ := strconv.Atoi(c.PostForm("constitution"))
		intelligence, _ := strconv.Atoi(c.PostForm("intelligence"))
		wisdom, _ := strconv.Atoi(c.PostForm("wisdom"))
		charisma, _ := strconv.Atoi(c.PostForm("charisma"))
		background := c.PostForm("background")
		alignment := c.PostForm("alignment")
		inventory := c.PostForm("inventory")
		spells := c.PostForm("spells")
		features := c.PostForm("features")
		notes := c.PostForm("notes")
		experience, _ := strconv.Atoi(c.PostForm("experience"))

		if name == "" {
			c.String(http.StatusBadRequest, "Имя персонажа обязательно")
			return
		}

		_, err := db.Exec(`
			INSERT INTO characters (
				user_id, name, race, class, level, strength, dexterity, constitution, intelligence, wisdom, charisma,
				background, alignment, inventory, spells, features, experience, notes
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			userID, name, race, class, level, strength, dexterity, constitution, intelligence, wisdom, charisma,
			background, alignment, inventory, spells, features, experience, notes,
		)
		if err != nil {
			c.String(http.StatusInternalServerError, "Ошибка сохранения персонажа: "+err.Error())
			return
		}

		c.Redirect(http.StatusFound, "/dashboard")
	}
}
