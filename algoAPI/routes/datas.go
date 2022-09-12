package routes

import (
	"algoAPI/models"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo/v4"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//Datas defines Datas controller struct
type Datas struct {
	log *log.Logger
	db  *mongo.Database
}

//NewDatas initializes a new Datas struct with given logger and database client
func NewDatas(log *log.Logger, db *mongo.Database) *Datas {
	return &Datas{log: log, db: db}
}

//NewData parses user input, authorizes the user, creates a new Data object and inserts it into the database
func (req *Datas) NewData(c echo.Context) error {
	tempData := models.Data{}
	err := json.NewDecoder(c.Request().Body).Decode(&tempData)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusBadRequest,
			json.RawMessage(`{"message": "Error decoding request body"}`))
	}
	jwtToken := c.Get("user").(*jwt.Token)
	jwtClaims := jwtToken.Claims.(jwt.MapClaims)
	userID, err := primitive.ObjectIDFromHex(jwtClaims["user_id"].(string))
	if err != nil {
		log.Println("Error converting user id" + err.Error())
		return c.JSON(http.StatusInternalServerError,
			json.RawMessage(`{"message": "Internal server error"}`))
	}
	tempData.User = userID
	tempData.Allowed = []string{}
	tempData.Created = time.Now().Format(time.RFC3339)
	tempData.Updated = time.Now().Format(time.RFC3339)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := req.db.Collection("datas").InsertOne(ctx, tempData)
	if err != nil {
		log.Println("Error inserting data" + err.Error())
		return c.JSON(http.StatusInternalServerError,
			json.RawMessage(`{"message": "Database error inserting data"}`))
	}
	tempData.ID = res.InsertedID.(primitive.ObjectID)
	return c.JSON(http.StatusOK, tempData)
}

//GetData authorizes the user, fetches the requested data from the database and returns it
func (req *Datas) GetData(c echo.Context) error {
	tempData := models.Data{}
	requestedID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(requestedID)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "Bad request"}`))
	}
	jwtToken := c.Get("user").(*jwt.Token)
	jwtClaims := jwtToken.Claims.(jwt.MapClaims)
	userID, err := primitive.ObjectIDFromHex(jwtClaims["user_id"].(string))
	if err != nil {
		log.Println("Error converting user id despite passing auth: " + err.Error())
		return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Something went very wrong, please submit a ticket"}`))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = req.db.Collection("datas").FindOne(ctx, bson.M{"_id": objectID}).Decode(&tempData)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Database error"}`))
	}
	if userID != tempData.User && !checkAllowed(tempData.Allowed, jwtClaims["user_id"].(string)) {
		return c.JSON(http.StatusUnauthorized, json.RawMessage(`{"message": "Not authorized to access data"}`))
	}
	tempData.ID = objectID
	return c.JSON(http.StatusOK, tempData)
}

//GetDatasByUser authorizes the user, fetches all of users data from the database and returns it
func (req *Datas) GetDatasByUser(c echo.Context) error {
	userID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "ID not valid"}`))
	}
	jwtToken := c.Get("user").(*jwt.Token)
	jwtClaims := jwtToken.Claims.(jwt.MapClaims)
	userIDtoken, err := primitive.ObjectIDFromHex(jwtClaims["user_id"].(string))
	if err != nil {
		log.Println("Error converting user id: " + err.Error())
		return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Something went very wrong, please submit a ticket"}`))
	}
	if userID != userIDtoken {
		return c.JSON(http.StatusUnauthorized, json.RawMessage(`{"message": "Not authorized to access data"}`))
	}
	tempData := models.Data{}
	tempDatas := []models.Data{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := req.db.Collection("datas").Find(ctx, bson.M{"user": userID})
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Internal server error"}`))
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		tempData = models.Data{}
		err := cursor.Decode(&tempData)
		if err != nil {
			log.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Internal server error"}`))
		}
		tempDatas = append(tempDatas, tempData)
	}
	return c.JSON(http.StatusOK, tempDatas)
}

//UpdateData parses users request, authorizes the user and updates the database
func (req *Datas) UpdateData(c echo.Context) error {
	tempData := models.Data{}
	err := json.NewDecoder(c.Request().Body).Decode(&tempData)
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
		return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Internal server error"}`))
	}
	tempDataAuthCheck := models.Data{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = req.db.Collection("datas").FindOne(ctx, bson.M{"_id": objectID}).Decode(&tempDataAuthCheck)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Database error"}`))
	}
	if userIDtoken != tempDataAuthCheck.User {
		return c.JSON(http.StatusUnauthorized, json.RawMessage(`{"message": "Not authorized to access data"}`))
	}
	allowedUpdate := (len(tempData.Allowed) > 0)
	dataUpdate := tempData.Data != ""
	update := bson.M{}
	if allowedUpdate && dataUpdate {
		update = bson.M{
			"$set": bson.M{
				"data":    tempData.Data,
				"allowed": tempData.Allowed,
				"updated": time.Now().Format(time.RFC3339),
			},
		}
	} else if allowedUpdate {
		update = bson.M{
			"$set": bson.M{
				"allowed": tempData.Allowed,
				"updated": time.Now().Format(time.RFC3339),
			},
		}
	} else if dataUpdate {
		update = bson.M{
			"$set": bson.M{
				"data":    tempData.Data,
				"updated": time.Now().Format(time.RFC3339),
			},
		}
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = req.db.Collection("datas").UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		if err.Error() == "update document must have at least one element" {
			return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "Empty update"}`))
		}
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Database error"}`))
	}
	return c.JSON(http.StatusOK, nil)
}

//DeleteData authorizes the user and deletes data from the database
func (req *Datas) DeleteData(c echo.Context) error {
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
	tempDataAuthCheck := models.Data{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = req.db.Collection("datas").FindOne(ctx, bson.M{"_id": objectID}).Decode(&tempDataAuthCheck)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Database error"}`))
	}
	if userIDtoken != tempDataAuthCheck.User {
		return c.JSON(http.StatusUnauthorized, json.RawMessage(`{"message": "Not authorized to access data"}`))
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"_id": objectID}
	res := req.db.Collection("datas").FindOneAndDelete(ctx, filter)
	if res.Err() == mongo.ErrNoDocuments {
		return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "Dara not found"}`))
	}
	return c.JSON(http.StatusNoContent, nil)
}
