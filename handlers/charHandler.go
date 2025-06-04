package handlers

import (
	"brosok20/models"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func DashboardHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получение списка персонажей пользователя
		var characters []models.Character // Character - структура, определяющая персонажа
		rows, err := db.Query("SELECT id, name FROM characters WHERE user_id = ?", c.GetInt("user_id"))
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

func CharacterViewHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получение и отображение персонажа
	}
}

func NewCharacterFormHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "character_edit.html", gin.H{
		"action": "/character/save",
		"method": "POST",
	})
}

func SaveCharacterHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем user_id из сессии
		sessionUserID := 0
		if session, ok := c.Get("user_id"); ok {
			sessionUserID = session.(int)
		} else {
			// Если не нашли user_id в контексте, пробуем из сессии gin-contrib/sessions
			sess := sessions.Default(c)
			uid := sess.Get("user_id")
			if uidInt, ok := uid.(int64); ok {
				sessionUserID = int(uidInt)
			} else if uidInt, ok := uid.(int); ok {
				sessionUserID = uidInt
			}
		}
		if sessionUserID == 0 {
			c.String(http.StatusUnauthorized, "Необходима авторизация")
			return
		}

		// Получаем данные из формы
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

		// Валидация (минимальная)
		if name == "" {
			c.String(http.StatusBadRequest, "Имя персонажа обязательно")
			return
		}

		// Сохраняем персонажа
		_, err := db.Exec(`
            INSERT INTO characters (
                user_id, name, race, class, level, strength, dexterity, constitution, intelligence, wisdom, charisma,
                background, alignment, inventory, spells, features, experience, notes
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			sessionUserID, name, race, class, level, strength, dexterity, constitution, intelligence, wisdom, charisma,
			background, alignment, inventory, spells, features, experience, notes,
		)
		if err != nil {
			c.String(http.StatusInternalServerError, "Ошибка сохранения персонажа: "+err.Error())
			return
		}

		c.Redirect(http.StatusFound, "/dashboard")
	}
}

func ApiCharacterListHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
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

func EditCharacterFormHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		var character models.Character
		err := db.QueryRow(`
            SELECT id, user_id, name, race, class, level, strength, dexterity, constitution, intelligence, wisdom, charisma,
                   background, alignment, experience, inventory, spells, features, notes
            FROM characters WHERE id = ? AND user_id = ? LIMIT 1`, id, userID).
			Scan(&character.ID, &character.UserID, &character.Name, &character.Race, &character.Class, &character.Level,
				&character.Strength, &character.Dexterity, &character.Constitution, &character.Intelligence, &character.Wisdom, &character.Charisma,
				&character.Background, &character.Alignment, &character.Experience, &character.Inventory, &character.Spells, &character.Features, &character.Notes)
		if err != nil {
			c.String(http.StatusNotFound, "Персонаж не найден")
			return
		}

		c.HTML(http.StatusOK, "character_edit.html", gin.H{
			"character": character,
			"action":    "/character/save",
			"method":    "POST",
		})
	}
}
