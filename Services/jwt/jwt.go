package jwt

import (
	"time"
	"errors"
	"log"
	"github.com/dgrijalva/jwt-go"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"github.com/ruslanfedoseenko/dhtcrawler/Utils"
)


func CreateNewTokens(uuid string, role string) (authTokenString, refreshTokenString, csrfSecret string, err error) {
	// generate the csrf secret
	csrfSecret, err = Utils.GenerateCSRFSecret()
	if err != nil {
		return
	}


	// generate the refresh token
	refreshTokenString, err = createRefreshTokenString(uuid, role, csrfSecret)


	// generate the auth token
	authTokenString, err = createAuthTokenString(uuid, role, csrfSecret)
	if err != nil {
		return
	}
	// don't need to check for err bc we're returning everything anyway
	return
}

// @adam-hanna: check if refreshToken["sub"] == authToken["sub"]?
// I don't think this is necessary bc a valid refresh token will always generate
// a valid auth token of the same "sub"
func CheckAndRefreshTokens(oldAuthTokenString string, oldRefreshTokenString string, oldCsrfSecret string) (newAuthTokenString, newRefreshTokenString, newCsrfSecret string, err error) {
	// first, check that a csrf token was provided
	if oldCsrfSecret == "" {
		log.Println("No CSRF token!")
		err = errors.New("Unauthorized")
		return
	}
	// now, check that it matches what's in the auth token claims
	authToken, err := jwt.ParseWithClaims(oldAuthTokenString, &Models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	authTokenClaims, ok := authToken.Claims.(*Models.TokenClaims)
	if !ok {
		return
	}
	if oldCsrfSecret != authTokenClaims.Csrf {
		log.Println("CSRF token doesn't match jwt!")
		err = errors.New("Unauthorized")
		return
	}


	// next, check the auth token in a stateless manner
	if authToken.Valid {
		log.Println("Auth token is valid")
		// auth token has not expired
		// we need to return the csrf secret bc that's what the function calls for
		newCsrfSecret = authTokenClaims.Csrf

		// update the exp of refresh token string, but don't save to the db
		// we don't need to check if our refresh token is valid here
		// because we aren't renewing the auth token, the auth token is already valid
		newRefreshTokenString, err = updateRefreshTokenExp(oldRefreshTokenString)
		newAuthTokenString = oldAuthTokenString
		return
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		log.Println("Auth token is not valid")
		if ve.Errors&(jwt.ValidationErrorExpired) != 0 {
			log.Println("Auth token is expired")
			// auth token is expired
			newAuthTokenString, newCsrfSecret, err = updateAuthTokenString(oldRefreshTokenString, oldAuthTokenString)
			if err != nil {
				return
			}

			// update the exp of refresh token string
			newRefreshTokenString, err = updateRefreshTokenExp(oldRefreshTokenString)
			if err != nil {
				return
			}

			// update the csrf string of the refresh token
			newRefreshTokenString, err = updateRefreshTokenCsrf(newRefreshTokenString, newCsrfSecret)
			return
		} else {
			log.Println("Error in auth token")
			err = errors.New("Error in auth token")
			return
		}
	} else {
		log.Println("Error in auth token")
		err = errors.New("Error in auth token")
		return
	}

	// if we get here, there was some error validating the token
	err = errors.New("Unauthorized")
	return
}

func createAuthTokenString(uuid string, role string, csrfSecret string) (authTokenString string, err error) {
	authTokenExp 	:= time.Now().Add(Models.AuthTokenValidTime).Unix()
	authClaims 		:= Models.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			Subject: uuid,
			ExpiresAt: authTokenExp,
		},
		Role: role,
		Csrf: csrfSecret,
	}

	// create a signer for rsa 256
	authJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), authClaims)

	// generate the auth token string
	authTokenString, err = authJwt.SignedString(signKey)
	return
}

func createRefreshTokenString(uuid string, role string, csrfString string) (refreshTokenString string, err error) {
	refreshTokenExp := time.Now().Add(Models.RefreshTokenValidTime).Unix()
	refreshJti, err := StoreNewRefreshToken()
	if err != nil {
		return
	}

	refreshClaims := Models.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			Id: refreshJti, // jti
			Subject: uuid,
			ExpiresAt: refreshTokenExp,
		},
		Role: role,
		Csrf: csrfString,
	}

	// create a signer for rsa 256
	refreshJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), refreshClaims)

	// generate the refresh token string
	refreshTokenString, err = refreshJwt.SignedString(signKey)
	return
}


func updateRefreshTokenExp(oldRefreshTokenString string) (newRefreshTokenString string, err error) {
	refreshToken, err := jwt.ParseWithClaims(oldRefreshTokenString, &Models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	oldRefreshTokenClaims, ok := refreshToken.Claims.(*Models.TokenClaims)
	if !ok {
		return
	}

	refreshTokenExp := time.Now().Add(Models.RefreshTokenValidTime).Unix()

	refreshClaims := Models.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			Id: oldRefreshTokenClaims.StandardClaims.Id, // jti
			Subject: oldRefreshTokenClaims.StandardClaims.Subject,
			ExpiresAt: refreshTokenExp,
		},
		Role: oldRefreshTokenClaims.Role,
		Csrf: oldRefreshTokenClaims.Csrf,
	}

	// create a signer for rsa 256
	refreshJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), refreshClaims)

	// generate the refresh token string
	newRefreshTokenString, err = refreshJwt.SignedString(signKey)
	return
}

func updateAuthTokenString(refreshTokenString string, oldAuthTokenString string) (newAuthTokenString, csrfSecret string, err error) {
	refreshToken, err := jwt.ParseWithClaims(refreshTokenString, &Models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	refreshTokenClaims, ok := refreshToken.Claims.(*Models.TokenClaims)
	if !ok {
		err = errors.New("Error reading jwt claims")
		return
	}

	// check if the refresh token has been revoked
	if CheckRefreshToken(refreshTokenClaims.StandardClaims.Id) {
		// the refresh token has not been revoked
		// has it expired?
		if refreshToken.Valid {
			// nope, the refresh token has not expired
			// issue a new auth token
			authToken, _ := jwt.ParseWithClaims(oldAuthTokenString, &Models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
				return verifyKey, nil
			})

			oldAuthTokenClaims, ok := authToken.Claims.(*Models.TokenClaims)
			if !ok {
				err = errors.New("Error reading jwt claims")
				return
			}

			// our policy is to regenerate the csrf secret for each new auth token
			csrfSecret, err = Utils.GenerateCSRFSecret()
			if err != nil {
				return
			}

			newAuthTokenString, err = createAuthTokenString(oldAuthTokenClaims.StandardClaims.Subject, oldAuthTokenClaims.Role, csrfSecret)

			return
		} else {
			log.Println("Refresh token has expired!")
			// the refresh token has expired!
			// Revoke the token in our db and require the user to login again
			DeleteRefreshToken(refreshTokenClaims.StandardClaims.Id)

			err = errors.New("Unauthorized")
			return
		}
	} else {
		log.Println("Refresh token has been revoked!")
		// the refresh token has been revoked!
		err = errors.New("Unauthorized")
		return
	}
}


func RevokeRefreshToken(refreshTokenString string) error {
	refreshToken, err := jwt.ParseWithClaims(refreshTokenString, &Models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	if err != nil {
		return errors.New("Could not parse refresh token with claims")
	}

	refreshTokenClaims, ok := refreshToken.Claims.(*Models.TokenClaims)
	if !ok {
		return errors.New("Could not read refresh token claims")
	}

	DeleteRefreshToken(refreshTokenClaims.StandardClaims.Id)

	return nil
}

func updateRefreshTokenCsrf(oldRefreshTokenString string, newCsrfString string) (newRefreshTokenString string, err error) {
	refreshToken, err := jwt.ParseWithClaims(oldRefreshTokenString, &Models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	oldRefreshTokenClaims, ok := refreshToken.Claims.(*Models.TokenClaims)
	if !ok {
		return
	}

	refreshClaims := Models.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			Id: oldRefreshTokenClaims.StandardClaims.Id, // jti
			Subject: oldRefreshTokenClaims.StandardClaims.Subject,
			ExpiresAt: oldRefreshTokenClaims.StandardClaims.ExpiresAt,
		},
		Role: oldRefreshTokenClaims.Role,
		Csrf: newCsrfString,
	}

	// create a signer for rsa 256
	refreshJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), refreshClaims)

	// generate the refresh token string
	newRefreshTokenString, err = refreshJwt.SignedString(signKey)
	return
}

func GrabUUID(authTokenString string) (string, error) {
	authToken, _ := jwt.ParseWithClaims(authTokenString, &Models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return "", errors.New("Error fetching claims")
	})
	authTokenClaims, ok := authToken.Claims.(*Models.TokenClaims)
	if !ok {
		return "", errors.New("Error fetching claims")
	}

	return authTokenClaims.StandardClaims.Subject, nil
}