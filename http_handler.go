package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type error struct {
	Error string `json:"error"`
}

type body struct {
	Body string `json:"body"`
}

type valid struct {
	Valid bool `json:"valid"`
}

func handlerHealth(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}

func handlerValidatChirp(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	params := body{}

	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		marshalError(w, 500, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		marshalError(w, 400, "Chirp is too long")
	} else {
		marshalOk(w)
	}

}

func marshalError(w http.ResponseWriter, code int, msg string) {
	error := error{
		Error: msg,
	}

	dat, err := json.Marshal(error)
	if err != nil {
		log.Printf("Error marshal JSON: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func marshalOk(w http.ResponseWriter) {
	valid := valid{
		Valid: true,
	}

	dat, err := json.Marshal(valid)
	if err != nil {
		log.Printf("Error marshal JSON: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}
