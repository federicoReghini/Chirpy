package main

import (
	"encoding/json"
	"log"
)

type createUserBodyRequest struct {
	Email string `json:"email"`
}

func (c *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	params := createUserBodyRequest{}

	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		marshalError(w, 500, "Something went wrong")
		return
	}

	user, err := c.db.CreateUser(req.Context(), params.Email)
	if err != nil {
		log.Fatal(err)
	}

	dat, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshal JSON: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(dat)
}
