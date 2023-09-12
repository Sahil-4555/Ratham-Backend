package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        		primitive.ObjectID 	`json:"id,omitempty"`
	UniversityId 	string 				`json:"universityid,omitempty"`
	Password		string 				`json:"password,omitempty"`
}