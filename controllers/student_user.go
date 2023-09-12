package controllers

import (
	"context"
	"fmt"
	"github/Sahil-4555/ratham-backend/configs"
	"github/Sahil-4555/ratham-backend/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var studentCollection *mongo.Collection = configs.GetCollection(configs.DB, "student_user")
var validate_student = validator.New()

// Find User By Email Search
func FindStudentby_UniversityId(ctx context.Context, universityid *string) (models.User, error) {
	var foundStudent models.User
	err := studentCollection.FindOne(ctx, bson.M{"universityid": universityid}).Decode(&foundStudent)
	if err != nil {
		return foundStudent, err
	}
	return foundStudent, err
}

// HashPassword is used to encrypt the password before it is stored in the DB
func HashPassword_Student(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		fmt.Printf(err.Error())
	}

	return string(bytes)
}

// VerifyPassword checks the input password while verifying it with the passward in the DB.
func VerifyPassword_Student(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or passowrd is incorrect")
		check = false
	}

	return check, msg
}

func Register_Student(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	// now := time.Now().In(time.FixedZone("IST", 5*60*60+30*60))
	// user.Createdat = now.Format("02-01-2006 15:04:05") // INDIAN STANDARD TIME
	defer cancel()

	//validate_student the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	if validationErr := validate_student.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": validationErr.Error()})
	}

	err := studentCollection.FindOne(ctx, bson.M{"universityid": user.UniversityId}).Decode(&user)

	if err == nil {
		return c.Status(http.StatusConflict).JSON(fiber.Map{"message": "Student already exists with this Unversity Id"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	user.Password = string(hashedPassword)

	user.Id = primitive.NewObjectID()

	result, err := studentCollection.InsertOne(ctx, user)
	
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "success",
		"user":    result,
	})
}

func Login_Student(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	var user models.User
	defer cancel()

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"data": "Invalid JSON Provided"})
	}

	if validationErr := validate_student.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"data": validationErr.Error()})
	}
	
	foundStudent, err := FindStudentby_UniversityId(ctx, &user.UniversityId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Error Occured While Checking For The University Id"})
		} else {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundStudent.Password), []byte(user.Password))
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid UniversityID Or Password"})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    foundStudent.Id.Hex(),
		ExpiresAt: jwt.At(time.Now().Add(time.Hour * 24)), // 1 day
	})

	token, err := claims.SignedString([]byte(configs.SecretKey()))

	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Could Not Login"})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    foundStudent,
		"token":   token,
	})
}

func Student(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.SecretKey()), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "UnAuthenticated"})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var foundStudent models.User
	_id, err := primitive.ObjectIDFromHex(claims.Issuer)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
	}
	err = studentCollection.FindOne(ctx, bson.M{"id": _id}).Decode(&foundStudent)
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"user":    foundStudent,
	})
}

func Logout_Student(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}