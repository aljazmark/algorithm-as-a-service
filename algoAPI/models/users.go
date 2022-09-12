package models

import (
	//"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//User defines user structure
type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string             `json:"username"`
	Email    string             `json:"email"`
	Password string             `json:"password,omitempty"`
	Created  string             `json:"created"`
	Updated  string             `json:"updated"`
}

//var secret = os.Getenv("JWTSecret")
var secret = "ExampleSecret42"

//ToJWT returns a JWT token from a User
func (u *User) ToJWT() (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = u.ID
	claims["user_username"] = u.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

//FromJWT returns a User id from a JWT token
func (u *User) FromJWT(token string) (string, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	return claims["user_id"].(string), nil
}
