package Models

type JwtRefreshToken struct {
	Id    int    `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"-"`
	Token string `gorm:"column:token;size:32"`
}
