package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func HomeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "brosok20.html", nil)
}

func LoginFormHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "logon.html", nil)
}

func LoginHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get form values
		emailOrUsername := c.PostForm("email")
		password := c.PostForm("password")

		// Find user by email or username
		var userID int
		var dbPassword string
		query := `SELECT id, password FROM users WHERE email = ? OR username = ? LIMIT 1`
		err := db.QueryRow(query, emailOrUsername, emailOrUsername).Scan(&userID, &dbPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				c.String(http.StatusUnauthorized, "Пользователь не найден")
				return
			}
			c.String(http.StatusInternalServerError, "Ошибка сервера")
			return
		}

		// Compare password (plain for demo; use bcrypt in production!)
		if password != dbPassword {
			c.String(http.StatusUnauthorized, "Неверный пароль")
			return
		}

		// Set session
		session := sessions.Default(c)
		session.Set("user_id", userID)
		session.Save()

		// AJAX or normal redirect
		if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
			c.JSON(http.StatusOK, gin.H{"success": true, "redirect": "/dashboard"})
		} else {
			c.Redirect(http.StatusFound, "/dashboard")
		}
	}
}

func RegisterFormHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "logon.html", nil)
}

func RegisterHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		email := c.PostForm("email")
		password := c.PostForm("password")

		// Basic validation
		if username == "" || email == "" || password == "" {
			c.String(http.StatusBadRequest, "Все поля обязательны для заполнения")
			return
		}

		// Check if username or email already exists
		var exists int
		err := db.QueryRow("SELECT 1 FROM users WHERE username = ? OR email = ? LIMIT 1", username, email).Scan(&exists)
		if err != nil && err != sql.ErrNoRows {
			c.String(http.StatusInternalServerError, "Ошибка сервера")
			return
		}
		if err == nil {
			c.String(http.StatusConflict, "Пользователь с таким email или именем уже существует")
			return
		}

		// Insert new user (plain password for demo; use bcrypt in production!)
		res, err := db.Exec(
			"INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
			username, email, password,
		)
		if err != nil {
			c.String(http.StatusInternalServerError, "Ошибка при создании пользователя")
			return
		}

		userID, err := res.LastInsertId()
		if err != nil {
			c.String(http.StatusInternalServerError, "Ошибка при создании пользователя")
			return
		}

		// Set session
		session := sessions.Default(c)
		session.Set("user_id", userID)
		session.Save()

		// AJAX or normal redirect
		if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
			c.JSON(http.StatusOK, gin.H{"success": true, "redirect": "/dashboard"})
		} else {
			c.Redirect(http.StatusFound, "/dashboard")
		}
	}
}

func AuthMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			// If AJAX request, return 401
			if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется вход"})
				c.Abort()
				return
			}
			// Otherwise, redirect to login page
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		// Set user_id in context for handlers
		c.Set("user_id", userID)
		c.Next()
	}
}
