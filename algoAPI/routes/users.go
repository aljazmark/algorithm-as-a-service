package routes

import (
	"algoAPI/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//Users defines Users controller struct
type Users struct {
	log *log.Logger
	db  *mongo.Database
}

//tokenS defines successful login response object
type tokenS struct {
	UserID   string `json:"userID"`
	Token    string `json:"token"`
	Username string `json:"username"`
}

//NewUsers initializes a new Users struct with given logger and database client
func NewUsers(log *log.Logger, db *mongo.Database) *Users {
	return &Users{log: log, db: db}
}

//NewUser parses the request, creates a new user, generates a JWT token and returns a tokenS object
func (usr *Users) NewUser(c echo.Context) error {
	tempUser := models.User{}
	err := json.NewDecoder(c.Request().Body).Decode(&tempUser)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "Invalid input"}`))
	}
	if checkUsername(tempUser.Username, usr.db.Collection("users")) {
		return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "Username not available"}`))
	}
	if checkEmail(tempUser.Email, usr.db.Collection("users")) {
		return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "Email not available"}`))
	}
	pwHash, err := passwordHash(tempUser.Password)
	tempUser.Password = pwHash
	tempUser.Created = time.Now().Format(time.RFC3339)
	tempUser.Updated = time.Now().Format(time.RFC3339)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := usr.db.Collection("users").InsertOne(ctx, tempUser)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Database error"}`))
	}
	fmt.Println(res)
	tempUser.ID = res.InsertedID.(primitive.ObjectID)
	token, err := tempUser.ToJWT()
	if err != nil {
		log.Println(err.Error())
	}
	tokenS := tokenS{}
	tokenS.Token = token
	tokenS.UserID = res.InsertedID.(primitive.ObjectID).Hex()
	tokenS.Username = tempUser.Username
	return c.JSON(http.StatusOK, tokenS)
}

//LoginUser authenticates user login request, generates a JWT token and returns a tokenS object
func (usr *Users) LoginUser(c echo.Context) error {
	tempUser := models.User{}
	dbUser := models.User{}
	token := ""
	err := json.NewDecoder(c.Request().Body).Decode(&tempUser)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "Bad request"}`))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"username": tempUser.Username}
	res := usr.db.Collection("users").FindOne(ctx, filter)
	if res.Err() == mongo.ErrNoDocuments {
		return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "User/Password combination not found"}`))
	}
	err = res.Decode(&dbUser)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Internal server error"}`))
	}
	if verifyPassword(tempUser.Password, dbUser.Password) {
		token, err = dbUser.ToJWT()
		if err != nil {
			log.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Internal server error"}`))
		}
	} else {
		return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "User/Password combination not found"}`))
	}
	tokenS := tokenS{}
	tokenS.Token = token
	tokenS.UserID = dbUser.ID.Hex()
	tokenS.Username = dbUser.Username
	return c.JSON(http.StatusOK, tokenS)
}

//GetUser authorizes the user, fetches the user from database and returns it
func (usr *Users) GetUser(c echo.Context) error {
	objectID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "Invalid id"}`))
	}
	jwtToken := c.Get("user").(*jwt.Token)
	jwtClaims := jwtToken.Claims.(jwt.MapClaims)
	userIDtoken, err := primitive.ObjectIDFromHex(jwtClaims["user_id"].(string))
	if err != nil {
		log.Println("Error converting user id despite passing auth: " + err.Error())
		return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Something went very wrong, please submit a ticket"}`))
	}
	if objectID != userIDtoken {
		return c.JSON(http.StatusUnauthorized, json.RawMessage(`{"message": "Not authorized to access user"}`))
	}
	tempUser := models.User{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = usr.db.Collection("users").FindOne(ctx, bson.M{"_id": objectID}).Decode(&tempUser)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Database error"}`))
	}
	tempUser.Password = ""
	return c.JSON(http.StatusOK, tempUser)
}

//UpdateUser authorizes user, parses the request and updates the database
func (usr *Users) UpdateUser(c echo.Context) error {
	tempUser := models.User{}
	err := json.NewDecoder(c.Request().Body).Decode(&tempUser)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "Bad request"}`))
	}
	requestedID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(requestedID)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "Invalid id"}`))
	}
	jwtToken := c.Get("user").(*jwt.Token)
	jwtClaims := jwtToken.Claims.(jwt.MapClaims)
	userIDtoken, err := primitive.ObjectIDFromHex(jwtClaims["user_id"].(string))
	if err != nil {
		log.Println("Error converting user id despite passing auth: " + err.Error())
		return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Something went very wrong, please submit a ticket"}`))
	}
	tempUserAuthCheck := models.User{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = usr.db.Collection("users").FindOne(ctx, bson.M{"_id": objectID}).Decode(&tempUserAuthCheck)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Database error fetching user"}`))
	}
	if userIDtoken != tempUserAuthCheck.ID {
		return c.JSON(http.StatusUnauthorized, json.RawMessage(`{"message": "Not authorized to access user"}`))
	}
	if tempUser.Username != "" && tempUser.Username != tempUserAuthCheck.Username {
		if checkUsername(tempUser.Username, usr.db.Collection("users")) {
			return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "Username not available"}`))
		}

		tempUserAuthCheck.Username = tempUser.Username
	}
	if tempUser.Email != "" && tempUser.Email != tempUserAuthCheck.Email {
		if checkEmail(tempUser.Email, usr.db.Collection("users")) {
			return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "Email not available"}`))
		}
		tempUserAuthCheck.Email = tempUser.Email
	}
	if tempUser.Password != "" {
		tempUserAuthCheck.Password, _ = passwordHash(tempUser.Password)
	}
	update := bson.M{
		"$set": bson.M{
			"username": tempUserAuthCheck.Username,
			"email":    tempUserAuthCheck.Email,
			"password": tempUserAuthCheck.Password,
			"updated":  time.Now().Format(time.RFC3339),
		},
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = usr.db.Collection("users").UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		if err.Error() == "update document must have at least one element" {
			return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "Empty update"}`))
		}
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Internal server error"}`))
	}
	return c.JSON(http.StatusOK, nil)
}

//DeleteUser authorizes the user and deletes user from the database
func (usr *Users) DeleteUser(c echo.Context) error {
	objectID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "Invalid id"}`))
	}
	jwtToken := c.Get("user").(*jwt.Token)
	jwtClaims := jwtToken.Claims.(jwt.MapClaims)
	userIDtoken, err := primitive.ObjectIDFromHex(jwtClaims["user_id"].(string))
	if err != nil {
		log.Println("Error converting user id despite passing auth: " + err.Error())
		return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Something went very wrong, please contact admin"}`))
	}
	tempUser := models.User{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = usr.db.Collection("users").FindOne(ctx, bson.M{"_id": objectID}).Decode(&tempUser)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Database error"}`))
	}
	if userIDtoken != tempUser.ID {
		return c.JSON(http.StatusUnauthorized, json.RawMessage(`{"message": "Not authorized to delete user"}`))
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"_id": objectID}
	res := usr.db.Collection("users").FindOneAndDelete(ctx, filter)
	if res.Err() == mongo.ErrNoDocuments {
		return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "User not found"}`))
	} else if res.Err() != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Internal server error"}`))
	}
	return c.JSON(http.StatusNoContent, nil)
}

//passwordHash generates hash value from password and returns it
func passwordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

//verifyPassowrd hashes the password and compares it with password hash from the database
func verifyPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return (err == nil)
}

//checkUsername checks if username exists in the database
func checkUsername(username string, db *mongo.Collection) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"username": username}
	res := db.FindOne(ctx, filter)
	if res.Err() == mongo.ErrNoDocuments {
		return false
	}
	return true
}

//checkEmail checks if email exists in the database
func checkEmail(email string, db *mongo.Collection) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"email": email}
	res := db.FindOne(ctx, filter)
	if res.Err() == mongo.ErrNoDocuments {
		return false
	}
	return true
}
