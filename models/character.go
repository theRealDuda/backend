package models

type Character struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	Name         string `json:"name"`
	Race         string `json:"race"`
	Class        string `json:"class"`
	Level        int    `json:"level"`
	Strength     int    `json:"strength"`
	Dexterity    int    `json:"dexterity"`
	Constitution int    `json:"constitution"`
	Intelligence int    `json:"intelligence"`
	Wisdom       int    `json:"wisdom"`
	Charisma     int    `json:"charisma"`
	Background   string `json:"background"`
	Alignment    string `json:"alignment"`
	Inventory    string `json:"inventory"`
	Spells       string `json:"spells"`
	Features     string `json:"features"`
	Experience   int    `json:"experience"`
	Notes        string `json:"notes"`
}
