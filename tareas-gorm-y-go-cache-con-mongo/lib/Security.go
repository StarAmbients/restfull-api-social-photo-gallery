package lib

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/starambients/tareas-gorm-y-go-cache-con-mongo/data"
	"github.com/starambients/tareas-gorm-y-go-cache-con-mongo/models"
	"github.com/jinzhu/gorm"
)

type JSONToken struct {
	Token string `json:"token"`
}

type TokenJWT struct {
	UserID uint
	jwt.StandardClaims
}

// ValidateCredent Buscar si existe un usuario para estos credenciales
func ValidateCredent(cred *models.Credentials, db *gorm.DB) *models.User {
	user := &models.User{}
	db.Where("username = ? AND password = ?", cred.Username, cred.Password).Find(user)
	if user.ID > 0 {
		return user
	} else {
		return nil
	}
}

func CreateJWT(usr *models.User) (string, error) {
	loginToken := &TokenJWT{UserID: usr.ID}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), loginToken) //Libreria
	return jwtToken.SignedString([]byte(os.Getenv("secret_key")))

}

func CreateToken(usr *models.User, cacheClient data.CacheProvider) (token string, err error) {
	h := hmac.New(sha256.New, []byte(os.Getenv("secret_key")))
	_, err = h.Write([]byte(usr.Username))
	if err == nil {
		token = hex.EncodeToString(h.Sum(nil))
		cacheClient.SetExpiration(token, usr, time.Hour)
	}
	return
}

// GetUserJWT : We have a JWT token and we want to extract a user from it
func GetUserJWT(tokenString string, db *gorm.DB) *models.User {
	tokenStruct := new(TokenJWT)
	token, err := jwt.ParseWithClaims(tokenString, tokenStruct, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("secret_key")), nil
	})
	// In case there is no error and the token we got is valid and there is a UserID in it
	if err == nil && token.Valid && tokenStruct.UserID > 0 {
		// We have to look for a user in the model which has an id like the one in the token we are working with
		user := new(models.User)
		db.First(user, tokenStruct.UserID)
		// We have found a user to which belonged the token
		//TODO We should verify if DB action was fine and return ERROR if not
		return user
	}
	// If there is no user, we return null
	// Detail: any other possibility not taken into account before also will return null
	// for example, if there was an error accessing the DB (later we can write specific
	// error message for this.
	return nil

}

func GetUserTokenCache(tokenString string, cacheClient data.CacheProvider) *models.User {
	if validUser, exists := cacheClient.Get(tokenString); exists && validUser != nil {
		return validUser.(*models.User)
	}
	return nil
}


//aqui a funcao de validao de base da dedos com WHERE
