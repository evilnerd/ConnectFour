package handlers

import (
	"connectfour/internal/service"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
	"time"
)

var secretKey = []byte("connectfour is the ultimate model")

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
		tokenString, err := createToken(user.Email)
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

	user, err := userService.FindUserByEmail(req.Email)

	// Could not determine if the user already exists
	if err != nil {
		errorResponse(w, "Could not determine if user already exists or not", http.StatusInternalServerError)
		return
	}

	// The user (or at least the e-mail address) already exists
	if strings.EqualFold(user.Email, req.Email) {
		errorResponse(w, "This e-mail address is already in use", http.StatusConflict)
		return
	}

	// All is good, let's create the user.
	if user, err := userService.CreateUser(req.Email, req.Name, req.Password); err != nil {
		errorResponse(w, "User creation failed", http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(service.NewCreateUserResponse(user))
		return
	}
}

func verifyPassword(password, hash string) bool {
	// TODO: implement proper hashing.
	return password == hash
}

func createToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
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
	return ctx.Value("email").(string)
}
