package handlers

import (
	"connectfour/internal/service"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"
)

var secretKey = []byte("connectfour is the ultimate game")

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req service.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	user, err := userService.FindUserByEmail(req.Email)

	if err != nil {
		errorResponse(w, "Could not load user", http.StatusInternalServerError)
		return
	}

	if verifyPassword(req.Password, user.Token) {
		tokenString, err := createToken(user.Email, user.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Internal api error while creating JWT"))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, tokenString)
		return
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Invalid credentials"))
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req service.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// All is good, let's create the user.
	if user, err := userService.CreateUser(req.Email, req.Name, hashPassword(req.Password)); err != nil {
		// User creation failed
		if errors.Is(err, service.UserExistsError{}) {
			errorResponse(w, "User already exists", http.StatusConflict)
			return
		}
		errorResponse(w, "User creation failed", http.StatusInternalServerError)
		return
	} else {
		// User creation succeeded
		w.WriteHeader(http.StatusCreated)
		userService.Cache(&user)
		_ = json.NewEncoder(w).Encode(service.NewCreateUserResponse(user))
		return
	}
}

func verifyPassword(password, hash string) bool {
	return hashPassword(password) == hash
}

func hashPassword(password string) string {
	// create sha256 hash of specified password
	h := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", h)
}

func createToken(email string, name string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			"name":  name,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func emailFromContext(r *http.Request) string {
	ctx := r.Context()
	if (ctx.Value("email")) != nil {
		return ctx.Value("email").(string)
	}
	return ""
}
