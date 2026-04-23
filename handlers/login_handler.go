package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LoginFormHandler(c *gin.Context) {
	tab := c.Query("tab")
	if tab != "register" {
		tab = "login"
	}
	renderLogon(c, tab, http.StatusOK, "")
}

func LoginHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		emailOrUsername := strings.TrimSpace(c.PostForm("email"))
		password := c.PostForm("password")

		if emailOrUsername == "" || password == "" {
			message := "Введите email/имя пользователя и пароль"
			if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
				c.String(http.StatusBadRequest, message)
				return
			}
			renderLogon(c, "login", http.StatusBadRequest, message)
			return
		}

		var userID int
		var dbPassword string
		query := `SELECT id, password FROM users WHERE email = ? OR username = ? LIMIT 1`
		err := db.QueryRow(query, emailOrUsername, emailOrUsername).Scan(&userID, &dbPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				message := "Неверный логин или пароль"
				if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
					c.String(http.StatusUnauthorized, message)
					return
				}
				renderLogon(c, "login", http.StatusUnauthorized, message)
				return
			}
			if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
				c.String(http.StatusInternalServerError, "Ошибка сервера")
				return
			}
			renderLogon(c, "login", http.StatusInternalServerError, "Ошибка сервера")
			return
		}

		if password != dbPassword {
			message := "Неверный логин или пароль"
			if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
				c.String(http.StatusUnauthorized, message)
				return
			}
			renderLogon(c, "login", http.StatusUnauthorized, message)
			return
		}

		session := sessions.Default(c)
		session.Set("user_id", userID)
		if err := session.Save(); err != nil {
			if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
				c.String(http.StatusInternalServerError, "Не удалось сохранить сессию")
				return
			}
			renderLogon(c, "login", http.StatusInternalServerError, "Не удалось сохранить сессию")
			return
		}

		if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
			c.JSON(http.StatusOK, gin.H{"success": true, "redirect": "/dashboard"})
		} else {
			c.Redirect(http.StatusFound, "/dashboard")
		}
	}
}
