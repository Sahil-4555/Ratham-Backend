package controllers

import (
	"context"
	"net/http"
	"time"
	"github/Sahil-4555/ratham-backend/models"
	"github/Sahil-4555/ratham-backend/configs"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/dgrijalva/jwt-go/v4"
	"go.mongodb.org/mongo-driver/bson"
)
	
var sessionCollection *mongo.Collection = configs.GetCollection(configs.DB, "Sessions")
var dean_Collection *mongo.Collection = configs.GetCollection(configs.DB, "dean_user")
var student_Collection *mongo.Collection = configs.GetCollection(configs.DB, "student_user")

func AuthenticateDean(c *fiber.Ctx) (*models.User, error) {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.SecretKey()), nil
	})

	if err != nil {
		return nil, err
	}

	claims := token.Claims.(*jwt.StandardClaims)
	_id, err := primitive.ObjectIDFromHex(claims.Issuer)

	if err != nil {
		return nil, err
	}

	var foundDean models.User
	err = dean_Collection.FindOne(c.Context(), bson.M{"id": _id}).Decode(&foundDean)
	
	if err != nil {
		return nil, err
	}

	return &foundDean, nil
}

func AuthenticateStudent(c *fiber.Ctx) (*models.User, error) {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.SecretKey()), nil
	})

	if err != nil {
		return nil, err
	}

	claims := token.Claims.(*jwt.StandardClaims)
	_id, err := primitive.ObjectIDFromHex(claims.Issuer)

	if err != nil {
		return nil, err
	}

	var foundStudent models.User
	err = student_Collection.FindOne(c.Context(), bson.M{"id": _id}).Decode(&foundStudent)
	
	if err != nil {
		return nil, err
	}

	return &foundStudent, nil
}

func AddNewSession(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	foundDean, err := AuthenticateDean(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{"message": "UnAuthenticated or Wrong Credentials",
		})
	}

	var req struct {
		IsFree 		bool				`json:"isfree,omitempty"`
		StartTime 	string 				`json:"starttime,omitempty"`
		EndTime		string	 			`json:"endtime,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Request or Failed to Fetch Data",
		})
	}

	newSession := models.Session{
		Id: primitive.NewObjectID(),
		Dean: foundDean.UniversityId,
		Students: []models.Student{},
		IsFree: req.IsFree,
		StartTime: req.StartTime,
		EndTime: req.EndTime,
	}

	result, err := sessionCollection.InsertOne(ctx, newSession)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error Inserting User Into MongoDB",
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "success",
		"id":      result.InsertedID,
	})
}

type SessionResponse struct {
	ID        primitive.ObjectID `json:"id,omitempty"`
	Dean      string             `json:"dean,omitempty"`
	IsFree    bool               `json:"isfree,omitempty"`
	StartTime string             `json:"starttime,omitempty"`
	EndTime   string             `json:"endtime,omitempty"`
}

func GetAllFreeSessions(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := AuthenticateStudent(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "UnAuthenticated or Wrong Credentials",
		})
	}

	now := time.Now().In(time.FixedZone("IST", 5*60*60+30*60))
	currentTimeStr := now.Format("02-01-2006 15:04:05")

	filter := bson.M{
		"isfree":    true,
		"starttime": bson.M{"$gt": currentTimeStr},
	}

	cursor, err := sessionCollection.Find(ctx, filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching free sessions",
		})
	}
	defer cursor.Close(ctx)

	var sessions []SessionResponse
	for cursor.Next(ctx) {
		var session models.Session
		if err := cursor.Decode(&session); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error decoding sessions",
			})
		}

		// Create a SessionResponse without the Students field
		sessionResponse := SessionResponse{
			ID:        session.Id,
			Dean:      session.Dean,
			StartTime: session.StartTime,
			EndTime:   session.EndTime,
		}

		sessions = append(sessions, sessionResponse)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
		"sessions": sessions,
	})
}

func BookASession(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	foundStudent, err := AuthenticateStudent(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "UnAuthenticated or Wrong Credentials",
		})
	}

	sessionIDParam := c.Params("id")
	sessionID, err := primitive.ObjectIDFromHex(sessionIDParam)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid session ID format",
		})
	}

	filter := bson.M{"id": sessionID}
	var existingSession models.Session
	err = sessionCollection.FindOne(ctx, filter).Decode(&existingSession)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"message": "Session not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching session",
		})
	}

	existingSession.Students = append(existingSession.Students, models.Student{UniversityId: foundStudent.UniversityId})

	update := bson.M{"$set": bson.M{"students": existingSession.Students}}
	_, err = sessionCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating session",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Session Booked Successfully",
	})
}

type SessionInfo struct {
	StartTime string         `json:"starttime,omitempty"`
	EndTime   string         `json:"endtime,omitempty"`
	Students  []models.Student `json:"students,omitempty"`
}

func GetUpcomingFreeSessions(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := AuthenticateDean(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{"message": "UnAuthenticated or Wrong Credentials",
		})
	}

	now := time.Now().In(time.FixedZone("IST", 5*60*60+30*60))
	currentTimeStr := now.Format("02-01-2006 15:04:05")

	filter := bson.M{
		"isfree":    true,
		"starttime": bson.M{"$gt": currentTimeStr},
	}

	cursor, err := sessionCollection.Find(ctx, filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching upcoming free sessions",
		})
	}
	defer cursor.Close(ctx)

	var sessionInfos []SessionInfo
	for cursor.Next(ctx) {
		var session models.Session
		if err := cursor.Decode(&session); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error decoding sessions",
			})
		}

		sessionInfo := SessionInfo{
			StartTime: session.StartTime,
			EndTime:   session.EndTime,
			Students:  session.Students,
		}

		sessionInfos = append(sessionInfos, sessionInfo)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Successfully Fetched UpComing Free Sessions",
		"sessions": sessionInfos,
	})
}




