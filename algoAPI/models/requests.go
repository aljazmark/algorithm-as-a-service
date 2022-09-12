package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Request defines algorithm request
type Request struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Algorithm     string             `json:"algorithm"`
	Input         string             `json:"input"`
	Parameters    []string           `json:"parameters"`
	Output        string             `json:"output"`
	Requested     string             `json:"requested"`
	Completed     string             `json:"completed"`
	ExecutionTime string             `json:"executionTime"`
	User          primitive.ObjectID `json:"user"`
}
