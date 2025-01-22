package models

type Music struct {
	Group string `json:"group" binding:"required"`
	Song  string `json:"song" binding:"required"`
}

type MusicInfo struct {
	ID         int    `json:"id"`
	Group      string `json:"group"`
	Song       string `json:"song"`
	RelaseDate string `json:"release_date"`
	Text       string `json:"text"`
	Link       string `json:"link"`
}

type MusicUpdate struct {
	Group      string `json:"group" db:"music_group"`
	Song       string `json:"song" db:"song"`
	RelaseDate string `json:"release_date" db:"release_date"`
	Text       string `json:"text" db:"text"`
	Link       string `json:"link" db:"link"`
}
