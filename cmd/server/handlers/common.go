package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

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
			errorResponse(response, "Something went wrong preparing the response. Check the server logs for more info.", http.StatusInternalServerError)
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
	w.Write(jsonResp)
}
