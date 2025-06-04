package main

import (
	"brosok20/handlers"
	"database/sql"
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	store := cookie.NewStore([]byte("secret"))

	// Настройка роутера
	r := gin.Default()
	r.Use(sessions.Sessions("brosok20session", store))
	// Инициализация базы данных
	db, err := sql.Open("sqlite3", "./brosok20.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создание таблиц, если их нет
	initDB(db)

	// Загрузка шаблонов
	r.LoadHTMLGlob("templates/*")

	// Статические файлы
	r.Static("/static", "./static")

	// Обработчики маршрутов
	r.GET("/", handlers.HomeHandler)
	r.GET("/login", handlers.LoginFormHandler)
	r.POST("/login", handlers.LoginHandler(db))
	r.GET("/register", handlers.RegisterFormHandler)
	r.POST("/register", handlers.RegisterHandler(db))
	r.GET("/dashboard", handlers.AuthMiddleware(db), handlers.DashboardHandler(db))
	r.GET("/character/:id", handlers.AuthMiddleware(db), handlers.CharacterViewHandler(db))
	r.GET("/character/new", handlers.AuthMiddleware(db), handlers.NewCharacterFormHandler)
	r.POST("/character/save", handlers.AuthMiddleware(db), handlers.SaveCharacterHandler(db))
	r.GET("/character/export/:id", handlers.AuthMiddleware(db), handlers.ExportCharacterHandler(db))
	r.GET("/api/characters", handlers.AuthMiddleware(db), handlers.ApiCharacterListHandler(db))
	r.GET("/character/edit/:id", handlers.AuthMiddleware(db), handlers.EditCharacterFormHandler(db))
	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

// Инициализация базы данных
func initDB(db *sql.DB) {
	// Создание таблицы пользователей
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Создание таблицы персонажей
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS characters (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			race TEXT,
			class TEXT,
			level INTEGER,
			strength INTEGER,
			dexterity INTEGER,
			constitution INTEGER,
			intelligence INTEGER,
			wisdom INTEGER,
			charisma INTEGER,
			background TEXT,
			alignment TEXT,
			skills TEXT,
			languages TEXT,
			inventory TEXT,
			spells TEXT,
			features TEXT,
			experience INTEGER DEFAULT 0,
			notes TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}
