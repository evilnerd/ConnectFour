package handlers

import (
	"connectfour/internal/db"
	"connectfour/internal/service"
	"context"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

var (
	userService  *service.UserService
	gamesService *service.GamesService
)

func init() {
	userService = service.NewUserService(db.NewMariaDbUserRepository(), time.Minute*2)
	gamesService = service.NewGamesService(
		userService,
		db.NewMariaDbGameRepository())
}

func marshal(obj interface{}, response http.ResponseWriter) bool {
	response.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(response)
	err := encoder.Encode(obj)
	return handleError(err, response)
}

func unmarshal[T interface{}](response http.ResponseWriter, request *http.Request) (T, bool) {
	var req T
	log.Debugf("Unmarshalling request to %s", request.RequestURI)
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	return req, handleError(err, response)
}

func handleError(err error, response http.ResponseWriter) bool {
	var unmarshalErr *json.UnmarshalTypeError
	var marshalErr *json.MarshalerError

	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(response, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else if errors.As(err, &marshalErr) {
			errorResponse(response, "Something went wrong preparing the response. Check the api logs for more info.", http.StatusInternalServerError)
		} else {
			errorResponse(response, "Bad Request: "+err.Error(), http.StatusBadRequest)
		}
		return false
	}
	return true
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	log.Warnf("Returning error response: %s (%d)", message, httpStatusCode)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	_, _ = w.Write(jsonResp)
}

func JwtValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		log.Debugf("Auth header: %s", auth)

		tokenString, found := strings.CutPrefix(auth, "Bearer ")
		if !found || tokenString == "" {
			log.Warnf("No bearer token found in request to: %s", r.URL.Path)
			errorResponse(w, "No bearer token found", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil {
			log.Warnf("Invalid bearer token for request to %s: %v", r.URL.Path, err)
			errorResponse(w, "Invalid bearer token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			exp, ok := claims["exp"].(float64)
			if !ok {
				log.Warnf("Invalid expiration time in token for request to %s", r.URL.Path)
				errorResponse(w, "Invalid expiration time", http.StatusUnauthorized)
				return
			}
			if expired(exp) {
				log.Warnf("Token expired for request to %s", r.URL.Path)
				errorResponse(w, "Token is expired", http.StatusUnauthorized)
				return
			}

			// Valid token, proceed
			email := claims["email"].(string)
			log.Debugf("Valid JWT for user %s, accessing %s", email, r.URL.Path)
			ctx := context.WithValue(r.Context(), "email", email)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			log.Warnf("Invalid token claims for request to %s", r.URL.Path)
			errorResponse(w, "Invalid token", http.StatusUnauthorized)
		}
	})
}

func expired(exp float64) bool {
	return time.Unix(int64(exp), 0).Before(time.Now())
}
