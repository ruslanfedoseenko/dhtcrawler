package jwt

import (
	"crypto/rsa"
	"io/ioutil"
	"github.com/dgrijalva/jwt-go"
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
)

// location of the files used for signing and verification
const (
	privKeyPath = "keys/app.rsa"     // `$ openssl genrsa -out app.rsa 2048`
	pubKeyPath  = "keys/app.rsa.pub" // `$ openssl rsa -in app.rsa -pubout > app.rsa.pub`
)

// keys are held in global variables
// i havn't seen a memory corruption/info leakage in go yet
// but maybe it's a better idea, just to store the public key in ram?
// and load the signKey on every signing request? depends on  your usage i guess
var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

var App *Config.App

// read the key files before starting http handlers
func InitJWT(app *Config.App) error {
	App = app
	signBytes, err := ioutil.ReadFile(App.Config.HttpConfig.JwtAuthPrivateKeyPath)
	if err != nil {
		return err
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return err
	}

	verifyBytes, err := ioutil.ReadFile(App.Config.HttpConfig.JwtAuthPublicKeyPath)
	if err != nil {
		return err
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return err
	}

	return nil
}