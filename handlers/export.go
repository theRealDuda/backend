package handlers

import (
	"brosok20/models"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

func ExportCharacterHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем ID персонажа из URL
		charID := c.Param("id")
		if charID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Не указан ID персонажа"})
			return
		}

		fmt.Println("[DEBUG] ExportCharacterHandler called for id:", charID)
		// Получаем данные персонажа из БД
		var char models.Character
		err := db.QueryRow(`
			SELECT id, user_id, name, race, class, level, 
				strength, dexterity, constitution, 
				intelligence, wisdom, charisma,
				background, alignment, inventory, spells, features, experience, notes
			FROM characters WHERE id = ?`, charID).Scan(
			&char.ID, &char.UserID, &char.Name, &char.Race, &char.Class, &char.Level,
			&char.Strength, &char.Dexterity, &char.Constitution,
			&char.Intelligence, &char.Wisdom, &char.Charisma,
			&char.Background, &char.Alignment, &char.Inventory, &char.Spells, &char.Features, &char.Experience, &char.Notes,
		)
		if err != nil {
			fmt.Println("[DEBUG] Character not found or DB error:", err)
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Персонаж не найден"})
			return
		}

		fmt.Println("[DEBUG] Character loaded:", char)

		// Создаем PDF
		pdf := CreateCharacterPDF(char)

		// Устанавливаем заголовки для скачивания файла
		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", "attachment; filename=character_"+char.Name+".pdf")

		// Отправляем PDF клиенту
		err = pdf.Output(c.Writer)
		if err != nil {
			fmt.Println("[DEBUG] PDF generation error:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации PDF"})
			return
		}
		fmt.Println("[DEBUG] PDF successfully sent for character:", char.Name)
	}
}

func CreateCharacterPDF(char models.Character) *gofpdf.Fpdf {
	// Создаем новый PDF документ
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8Font("DejaVu", "B", "static/fonts/DejaVuSans-Bold.ttf")
	pdf.AddUTF8Font("DejaVu", "", "static/fonts/DejaVuSans.ttf")
	pdf.SetFont("DejaVu", "", 16)
	pdf.AddPage()

	// Устанавливаем шрифт и добавляем заголовок
	pdf.SetFont("DejaVu", "B", 16)
	pdf.Cell(40, 10, "Лист персонажа D&D 5e")
	pdf.Ln(12)

	// Основная информация о персонаже
	pdf.SetFont("DejaVu", "B", 12)
	pdf.Cell(40, 10, "Основная информация")
	pdf.Ln(8)
	pdf.SetFont("DejaVu", "", 10)

	info := []struct {
		label string
		value string
	}{
		{"Имя:", char.Name},
		{"Раса:", char.Race},
		{"Класс:", char.Class},
		{"Уровень:", strconv.Itoa(char.Level)},
		{"Предыстория:", char.Background},
		{"Мировоззрение:", char.Alignment},
	}

	for _, i := range info {
		pdf.Cell(40, 7, i.label)
		pdf.Cell(60, 7, i.value)
		pdf.Ln(7)
	}
	pdf.Ln(5)

	// Характеристики персонажа
	pdf.SetFont("DejaVu", "B", 12)
	pdf.Cell(40, 10, "Характеристики")
	pdf.Ln(8)
	pdf.SetFont("DejaVu", "", 10)

	abilities := []struct {
		name  string
		value int
		abbr  string
	}{
		{"Сила", char.Strength, "STR"},
		{"Ловкость", char.Dexterity, "DEX"},
		{"Телосложение", char.Constitution, "CON"},
		{"Интеллект", char.Intelligence, "INT"},
		{"Мудрость", char.Wisdom, "WIS"},
		{"Харизма", char.Charisma, "CHA"},
	}

	// Таблица характеристик
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(40, 8, "Характеристика", "1", 0, "C", true, 0, "")
	pdf.CellFormat(20, 8, "Значение", "1", 0, "C", true, 0, "")
	pdf.CellFormat(20, 8, "Модиф.", "1", 1, "C", true, 0, "")
	pdf.SetFillColor(255, 255, 255)

	for _, a := range abilities {
		modifier := (a.value - 10) / 2 // Расчет модификатора
		modStr := fmt.Sprintf("%+d", modifier)

		pdf.CellFormat(40, 8, fmt.Sprintf("%s (%s)", a.name, a.abbr), "1", 0, "", false, 0, "")
		pdf.CellFormat(20, 8, strconv.Itoa(a.value), "1", 0, "C", false, 0, "")
		pdf.CellFormat(20, 8, modStr, "1", 1, "C", false, 0, "")
	}
	pdf.Ln(10)

	// Дополнительные разделы
	addSection(pdf, "Инвентарь", char.Inventory)
	addSection(pdf, "Заклинания", char.Spells)
	addSection(pdf, "Черты и особенности", char.Features)

	return pdf
}

func addSection(pdf *gofpdf.Fpdf, title, content string) {
	if content == "" {
		return
	}
	pdf.SetFont("DejaVu", "B", 12)
	pdf.Cell(40, 10, title)
	pdf.Ln(8)
	pdf.SetFont("DejaVu", "", 10)
	pdf.MultiCell(0, 7, content, "", "", false)
	pdf.Ln(5)
}
