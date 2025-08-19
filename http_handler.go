package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type error struct {
	Error string `json:"error"`
}

type body struct {
	Body string `json:"body"`
}

type cleaned_body struct {
	Cleaned_body string `json:"cleaned_body"`
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
		// not allowed words
		badWords := map[string]bool{"kerfuffle": true, "sharbert": true, "fornax": true}
		uncensuredText := strings.Split(params.Body, " ")

		msg := strings.Builder{}

		if params.Body != "" {
			for _, word := range uncensuredText {
				if _, exist := badWords[word]; exist {
					msg.WriteString("****" + " ")
				} else {
					msg.WriteString(word + " ")
				}
			}
		} else {
			msg.WriteString(params.Body)
		}

		marshalOk(w, strings.TrimSpace(msg.String()))
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

func marshalOk(w http.ResponseWriter, msg string) {
	valid := cleaned_body{
		Cleaned_body: msg,
	}

	dat, err := json.Marshal(valid)
	if err != nil {
		log.Printf("Error marshal JSON: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}
