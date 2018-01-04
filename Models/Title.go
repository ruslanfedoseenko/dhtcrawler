package Models

import "database/sql/driver"

type TitleType string

const (
	TV    TitleType = "tv"
	Movie TitleType = "movie"
	Game  TitleType = "game"
)

func (u *TitleType) Scan(value interface{}) error { *u = TitleType(value.([]byte)); return nil }
func (u TitleType) Value() (driver.Value, error)  { return string(u), nil }

type Title struct {
	Id          uint32    `gorm:"column:id"`
	Year        uint32    `gorm:"column:year"`
	Title       string    `gorm:"column:title"`
	TitleType   TitleType `gorm:"column:title_type"`
	Description string    `gorm:"column:description"`
	Ganres      []string  `gorm:"-"`
	PosterUrl   string    `gorm:"column:poster"`
}

func (Title) TableName() string {
	return "titles"
}
