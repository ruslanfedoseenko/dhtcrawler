package Models

type User struct {
	Id           int    `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"-"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Mail         string `json:"mail"`
	Role         string `json:"-"`
}
