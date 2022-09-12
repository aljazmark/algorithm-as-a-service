package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Data struct definition
type Data struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Data    string             `json:"data"` //json or dataID
	User    primitive.ObjectID `json:"user"`
	Allowed []string           `json:"allowed"`
	Created string             `json:"created"`
	Updated string             `json:"updated"`
}
