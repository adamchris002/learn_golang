package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"backend/database"
	"backend/models"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type userResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	//decode request body to user struct
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	//check if user with the same username already exists
	var existingUser models.User

	result := database.DB.
		Where("username = ?", user.Username).
		First(&existingUser)

	if result.Error == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Register Fail",
			"message":      fmt.Sprintf("User with username %s exist", user.Username),
		})
		return
	}

	//check if the error is not a record not found error
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Register Fail",
			"message":      "Database error",
		})
		return
	}

	//encode the password before saving the user to the database
	err = EncodePassword(&user)

	//check if there was an error when encoding the password
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Register Fail",
			"message":      "Fail to encrypt password",
		})
		return
	}

	//save the user to the database
	database.DB.Create(&user)

	//send a success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"messageTitle": "Register Success",
		"message":      "User created",
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	//decode request body to user struct
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)

	//check if user with the same username already exists
	var searchUser models.User

	result := database.DB.Where("username = ?", user.Username).First(&searchUser)

	if result.Error != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Login Fail",
			"message":      fmt.Sprintf("User with username %s does not exist", user.Username),
		})
		return
	} else {
		checkPassword := VerifyPassword([]byte(searchUser.Password), []byte(user.Password))

		if !checkPassword {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]string{
				"messageTitle": "Login Fail",
				"message":      "Password is incorrect",
			})
			return
		}

		token, err := GenerateToken(searchUser.ID)

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)

			json.NewEncoder(w).Encode(map[string]string{
				"messageTitle": "Login Fail",
				"message":      "Fail to generate token",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"user":         userResponse{ID: searchUser.ID, Username: searchUser.Username, Email: searchUser.Email, FirstName: searchUser.FirstName, LastName: searchUser.LastName},
			"token":        token,
			"messageTitle": "Login Success",
			"message":      "Login Successful",
		})
	}

	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
}

func EncodePassword(user *models.User) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.Password = string(hashPassword)
	return nil
}

func VerifyPassword(hashedPassword, password []byte) bool {

	err := bcrypt.CompareHashAndPassword(hashedPassword, password)

	if err != nil {
		return false
	}

	return true
}

func GenerateToken(userID uint) (string, error) {
	// Create a new JWT token with the user ID as a claim
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
