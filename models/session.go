package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	UniversityId string `json:"universityid,omitempty"`
}

type Session struct {
	Id 			primitive.ObjectID 	`json:"id,omitempty"`
	Dean 		string 				`json:"dean,omitempty"`
	Students 	[]Student 			`json:"students,omitempty"`
	IsFree 		bool				`json:"isfree,omitempty"`
	StartTime 	string 				`json:"starttime,omitempty"`
	EndTime		string	 			`json:"endtime,omitempty"`
}