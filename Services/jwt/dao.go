package jwt

import (
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"github.com/ruslanfedoseenko/dhtcrawler/Utils"
)

func StoreNewRefreshToken() (token string ,e error) {
	token, e = Utils.GenerateRandomString(32)
	App.Db.Create(&Models.JwtRefreshToken{
		Token: token,
	})
	return
}


func CheckRefreshToken(token string) bool {
	var tokensCount = 0
	App.Db.Model(&Models.JwtRefreshToken{}).Where(&Models.JwtRefreshToken{
		Token: token,
	}).Count(&tokensCount)
	return tokensCount > 0
}


func DeleteRefreshToken(token string) {
	App.Db.Where(&Models.JwtRefreshToken{
		Token: token,
	}).Delete(Models.JwtRefreshToken{})
}

