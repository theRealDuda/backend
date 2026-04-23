package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApiCharacterListHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := getUserID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
			return
		}

		rows, err := db.Query(`
			SELECT id, name, race, class, level, background, alignment, strength, dexterity, constitution, intelligence, wisdom, charisma, experience, inventory, spells, features, notes
			FROM characters WHERE user_id = ?`, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
			return
		}
		defer rows.Close()

		var characters []map[string]interface{}
		for rows.Next() {
			var ch = make(map[string]interface{})
			var (
				id, level, strength, dexterity, constitution, intelligence, wisdom, charisma, experience int
				name, race, className, background, alignment, inventory, spells, features, notes         string
			)
			err := rows.Scan(&id, &name, &race, &className, &level, &background, &alignment, &strength, &dexterity, &constitution, &intelligence, &wisdom, &charisma, &experience, &inventory, &spells, &features, &notes)
			if err != nil {
				continue
			}
			ch["id"] = id
			ch["name"] = name
			ch["race"] = race
			ch["class"] = className
			ch["level"] = level
			ch["background"] = background
			ch["alignment"] = alignment
			ch["strength"] = strength
			ch["dexterity"] = dexterity
			ch["constitution"] = constitution
			ch["intelligence"] = intelligence
			ch["wisdom"] = wisdom
			ch["charisma"] = charisma
			ch["experience"] = experience
			ch["inventory"] = inventory
			ch["spells"] = spells
			ch["features"] = features
			ch["notes"] = notes
			characters = append(characters, ch)
		}
		c.JSON(http.StatusOK, characters)
	}
}
