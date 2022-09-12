package routes

import (
	"algoAPI/models"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//AlgoResponse defines response object from algo core
type AlgoResponse struct {
	Algorithm     string   `json:"algorithm"`
	Input         string   `json:"input"`
	Parameters    []string `json:"parameters,omitempty"`
	Output        string   `json:"output"`
	ExecutionTime string   `json:"executiontime"`
}

//AlgoRequest defines request object for algo core
type AlgoRequest struct {
	Input      string   `json:"input"`
	Parameters []string `json:"parameters,omitempty"`
}

//Requests defines Requests controller struct
type Requests struct {
	log *log.Logger
	db  *mongo.Database
}

//NewRequests initializes a new Requests struct with given logger and database client
func NewRequests(log *log.Logger, db *mongo.Database) *Requests {
	return &Requests{log: log, db: db}
}

//NewRequest parses the request, calls algo core component, inserts finished request into the database and returns it
//If the request was made with a JWT token, the user is authorized and the request is assigned to the user
func (req *Requests) NewRequest(c echo.Context) error {
	tempRequest := models.Request{}
	userID, _ := primitive.ObjectIDFromHex("")
	var err error
	if c.Get("user") != nil {
		jwtToken := c.Get("user").(*jwt.Token)
		jwtClaims := jwtToken.Claims.(jwt.MapClaims)
		userID, err = primitive.ObjectIDFromHex(jwtClaims["user_id"].(string))
		if err != nil {
			log.Println(err.Error())
			return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "Bad request"}`))
		}
	}
	err = json.NewDecoder(c.Request().Body).Decode(&tempRequest)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "Error parsing input"}`))
	}
	tmpID, _ := primitive.ObjectIDFromHex("")
	tempRequest.ID = tmpID
	tempRequest.Algorithm = c.Param("id")
	tempRequest.Requested = time.Now().Format(time.RFC3339)
	tempRequest.User = userID
	if tempRequest.Parameters == nil {
		tempRequest.Parameters = []string{}
	}
	algoResult, err := callAlgo(tempRequest)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Error running algorithm"}`))
	}
	tempRequest.Completed = time.Now().Format(time.RFC3339)
	tempRequest.ExecutionTime = algoResult.ExecutionTime
	tempRequest.Output = algoResult.Output
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := req.db.Collection("requests").InsertOne(ctx, tempRequest)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Database error"}`))
	}
	id := res.InsertedID.(primitive.ObjectID)
	tempRequest.ID = id
	return c.JSON(http.StatusOK, tempRequest)
}

//NewRequestWithData parses the request, authorizes user, fetches requested data from the database, calls algo core component, inserts finished request into the database and returns it
func (req *Requests) NewRequestWithData(c echo.Context) error {
	tempRequest := models.Request{}
	tempData := models.Data{}
	objectID, err := primitive.ObjectIDFromHex(c.Param("data"))
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "Bad request"}`))
	}
	err = json.NewDecoder(c.Request().Body).Decode(&tempRequest)
	if err != nil {
		if err.Error() != "EOF" {
			log.Println(err.Error())
			return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "Error parsing input"}`))
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = req.db.Collection("datas").FindOne(ctx, bson.M{"_id": objectID}).Decode(&tempData)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Database error fetching data"}`))
	}
	tempRequest.Input = tempData.Data
	tempRequest.Algorithm = c.Param("id")
	tempRequest.Requested = time.Now().Format(time.RFC3339)
	jwtToken := c.Get("user").(*jwt.Token)
	jwtClaims := jwtToken.Claims.(jwt.MapClaims)
	userID, err := primitive.ObjectIDFromHex(jwtClaims["user_id"].(string))
	if err != nil {
		log.Println("Error converting user id despite passing auth: " + err.Error())
		return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Internal server error"}`))
	}
	tempRequest.User = userID
	if tempRequest.User != tempData.User && !checkAllowed(tempData.Allowed, jwtClaims["user_id"].(string)) {
		return c.JSON(http.StatusUnauthorized, json.RawMessage(`{"message": "Not authorized to access data"}`))
	}
	algoResult, err := callAlgo(tempRequest)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Error running algorithm"}`))
	}
	tempRequest.Completed = time.Now().Format(time.RFC3339)
	tempRequest.ExecutionTime = algoResult.ExecutionTime
	tempRequest.Output = algoResult.Output
	if tempRequest.Parameters == nil {
		tempRequest.Parameters = []string{}
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := req.db.Collection("requests").InsertOne(ctx, tempRequest)
	id := res.InsertedID.(primitive.ObjectID)
	tempRequest.ID = id
	return c.JSON(http.StatusOK, tempRequest)
}

//GetRequest authorizes the user, fetches the request and returns it
//If the request is not assigned to a user, it can be requested without authorization
func (req *Requests) GetRequest(c echo.Context) error {
	tempRequest := models.Request{}
	requestedID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(requestedID)
	withUser := false
	userID, _ := primitive.ObjectIDFromHex("")
	publicID, _ := primitive.ObjectIDFromHex("")
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "Invalid id"}`))
	}
	if c.Get("user") != nil {
		jwtToken := c.Get("user").(*jwt.Token)
		jwtClaims := jwtToken.Claims.(jwt.MapClaims)
		userID, err = primitive.ObjectIDFromHex(jwtClaims["user_id"].(string))
		if err != nil {
			log.Println("Error converting user id despite passing auth: " + err.Error())
			return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Internal server error"}`))
		}
		withUser = true
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = req.db.Collection("requests").FindOne(ctx, bson.M{"_id": objectID}).Decode(&tempRequest)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "Request not found"}`))
	}
	if (tempRequest.User != publicID && !withUser) || (withUser && userID != tempRequest.User) {
		return c.JSON(http.StatusUnauthorized, json.RawMessage(`{"message": "Not authorized to access request"}`))
	}
	tempRequest.ID = objectID
	return c.JSON(http.StatusOK, tempRequest)
}

//GetRequests returns all requests in collection, for internal use
func (req *Requests) GetRequests(c echo.Context) error {
	tempRequest := models.Request{}
	tempRequests := []models.Request{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := req.db.Collection("requests").Find(ctx, bson.D{})
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Database error"}`))
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		err := cursor.Decode(&tempRequest)
		if err != nil {
			log.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Database error"}`))
		}
		tempRequests = append(tempRequests, tempRequest)
	}
	return c.JSON(http.StatusOK, tempRequests)
}

//GetRequestsByUser authorizes the user, fetches all of users requests from the database and returns it
func (req *Requests) GetRequestsByUser(c echo.Context) error {
	userID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "Invalid id"}`))
	}
	jwtToken := c.Get("user").(*jwt.Token)
	jwtClaims := jwtToken.Claims.(jwt.MapClaims)
	userID2, err := primitive.ObjectIDFromHex(jwtClaims["user_id"].(string))
	if err != nil {
		log.Println("Error converting user id despite passing auth: " + err.Error())
		return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Something went very wrong, please submit a ticket"}`))
	}
	//check auth
	if userID != userID2 {
		return c.JSON(http.StatusUnauthorized, json.RawMessage(`{"message": "Not authorized to access requests"}`))
	}
	tempRequest := models.Request{}
	tempRequests := []models.Request{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := req.db.Collection("requests").Find(ctx, bson.M{"user": userID})
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "User not found"}`))
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		tempRequest = models.Request{}
		err := cursor.Decode(&tempRequest)
		if err != nil {
			log.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, json.RawMessage(`{"message": "Error fetching requests"}`))
		}
		tempRequests = append(tempRequests, tempRequest)
	}
	return c.JSON(http.StatusOK, tempRequests)
}

//DeleteRequest authorizes the user and deletes the request from the database
func (req *Requests) DeleteRequest(c echo.Context) error {
	tempRequest := models.Request{}
	objectID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "Invalid id"}`))
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
	err = req.db.Collection("requests").FindOne(ctx, bson.M{"_id": objectID}).Decode(&tempRequest)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "Request not found"}`))
	}
	if userID != tempRequest.User {
		return c.JSON(http.StatusUnauthorized, json.RawMessage(`{"message": "Not authorized to access data"}`))
	}
	//delete request
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"_id": objectID}
	res := req.db.Collection("requests").FindOneAndDelete(ctx, filter)
	if res.Err() == mongo.ErrNoDocuments {
		return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "Request not found"}`))
	}
	return c.JSON(http.StatusNoContent, nil)
}

//checkAllowed authorizes user for data usage
func checkAllowed(allowed []string, id string) bool {
	for _, s := range allowed {
		if s == id {
			return true
		}
		i, err := strconv.Atoi(s)
		if err != nil {
			continue
		}
		if i == 0 {
			return true
		}
	}
	return false
}

//callAlgo makes an HTTP Post request to algo core
func callAlgo(request models.Request) (AlgoResponse, error) {
	requestBody, err := json.Marshal(AlgoRequest{
		Input:      request.Input,
		Parameters: request.Parameters,
	})
	if err != nil {
		return AlgoResponse{}, err
	}
	setALGOuri := "http://localhost:8080/algorithm/"
	//setALGOuri := "https://algocore.herokuapp.com/algorithm/"
	response, err := http.Post(setALGOuri+request.Algorithm, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return AlgoResponse{}, err
	}
	defer response.Body.Close()
	tempResponse := AlgoResponse{}
	err = json.NewDecoder(response.Body).Decode(&tempResponse)
	if err != nil {
		return AlgoResponse{}, err
	}
	return tempResponse, nil
}
