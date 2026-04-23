package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RegisterFormHandler(c *gin.Context) {
	renderLogon(c, "register", http.StatusOK, "")
}

func RegisterHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := strings.TrimSpace(c.PostForm("username"))
		email := strings.TrimSpace(c.PostForm("email"))
		password := c.PostForm("password")

		if username == "" || email == "" || password == "" {
			message := "Все поля обязательны для заполнения"
			if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
				c.String(http.StatusBadRequest, message)
				return
			}
			renderLogon(c, "register", http.StatusBadRequest, message)
			return
		}

		var exists int
		err := db.QueryRow("SELECT COUNT(1) FROM users WHERE username = ? OR email = ?", username, email).Scan(&exists)
		if err != nil {
			if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
				c.String(http.StatusInternalServerError, "Ошибка сервера")
				return
			}
			renderLogon(c, "register", http.StatusInternalServerError, "Ошибка сервера")
			return
		}
		if exists > 0 {
			message := "Пользователь с таким email или именем уже существует"
			if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
				c.String(http.StatusConflict, message)
				return
			}
			renderLogon(c, "register", http.StatusConflict, message)
			return
		}

		res, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, password)
		if err != nil {
			if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
				c.String(http.StatusInternalServerError, "Ошибка при создании пользователя")
				return
			}
			renderLogon(c, "register", http.StatusInternalServerError, "Ошибка при создании пользователя")
			return
		}

		userID, err := res.LastInsertId()
		if err != nil {
			if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
				c.String(http.StatusInternalServerError, "Ошибка при создании пользователя")
				return
			}
			renderLogon(c, "register", http.StatusInternalServerError, "Ошибка при создании пользователя")
			return
		}

		session := sessions.Default(c)
		session.Set("user_id", int(userID))
		if err := session.Save(); err != nil {
			if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
				c.String(http.StatusInternalServerError, "Не удалось сохранить сессию")
				return
			}
			renderLogon(c, "register", http.StatusInternalServerError, "Не удалось сохранить сессию")
			return
		}

		if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
			c.JSON(http.StatusOK, gin.H{"success": true, "redirect": "/dashboard"})
		} else {
			c.Redirect(http.StatusFound, "/dashboard")
		}
	}
}
