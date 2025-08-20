package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type myError struct {
	Error string `json:"error"`
}

func handlerHealth(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}

// Marshal and write on http.ResponseWriter the status code and error message
// w is the http.ResponseWriter
// statusCode should be any status code for http status response (2xx,5xx,4xx etc..)
// msg is the error message to be marshalled
// If msg is empty, it will be set to "Something went wrong"
// If the marshal fails, it will log the error
// and write a generic error message to the response
// The response will have the Content-Type set to application/json
// and the status code set to the provided statusCode
// If the statusCode is 0, it will default to http.StatusInternalServerError
// The response body will contain a JSON object with the error message
// Example response: {"error": "Something went wrong"}
// Example usage: marshalError(w, 400, "Bad Request")
// Example usage: marshalError(w, 0, "Internal Server Error")
func marshalError(w http.ResponseWriter, statusCode int, msg string) {
	error := myError{
		Error: msg,
	}

	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}

	dat, err := json.Marshal(error)
	if err != nil {
		log.Printf("Error marshal JSON: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(dat)
}

// Marshal and write on http.ResponseWriter the status code and struct marshalled
// w is the http.ResponseWriter
// statusCode should be any status code for http status response (2xx,5xx,4xx etc..)
// data is the struct to be marshalled
// If statusCode is 0, it will default to http.StatusOK
// The response will have the Content-Type set to application/json
// The response body will contain the marshalled data
// Example response: {"id": "123e4567-e89b-12d3-a456-426614174000", "body": "Hello, world!"}
// Example usage: marshalOkJson(w, 200, chirp)
// Example usage: marshalOkJson(w, 0, chirp)
// Example usage: marshalOkJson(w, 201, chirp)
func marshalOkJson(w http.ResponseWriter, statusCode int, data any) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	dat, err := json.Marshal(data)
	ifErrCondition(w, err, http.StatusInternalServerError, "")

	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write(dat)
	}
}

// A wrapper for if err nil condition that also write on the http.ResponseWriter.
// w is the http.ResponseWriter
// statusCode should be any status code for http status response (2xx,5xx,4xx etc..)
// msg if valuated as "" will be set to Something went wrong
func ifErrCondition(w http.ResponseWriter, err error, statusCode int, msg string) {

	if err != nil {

		if msg == "" {
			msg = err.Error()
		}

		marshalError(w, statusCode, msg)
		return
	}
}
