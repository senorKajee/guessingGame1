package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type GuessRequestBody struct {
	Guess int `json:"guess"`
}

type GuessResponseBody struct {
	Status       string `json:"status"`
	NumGuesses   int    `json:"numGuesses,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

type LoginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponseBody struct {
	Token string `json:"token"`
}

var numToGuess int
var numGuesses int
var authenticated bool

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/guess", guessHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var requestBody LoginRequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if requestBody.Username == "admin" && requestBody.Password == "password" {
		token := "my-random-token"
		responseBody := LoginResponseBody{Token: token}
		json.NewEncoder(w).Encode(responseBody)
		authenticated = true
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
}

func guessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !authenticated {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var requestBody GuessRequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if numToGuess == 0 {
		numToGuess = rand.Intn(100) + 1
	}

	numGuesses++

	if requestBody.Guess < numToGuess {
		responseBody := GuessResponseBody{Status: "tooLow"}
		json.NewEncoder(w).Encode(responseBody)
	} else if requestBody.Guess > numToGuess {
		responseBody := GuessResponseBody{Status: "tooHigh"}
		json.NewEncoder(w).Encode(responseBody)
	} else {
		responseBody := GuessResponseBody{Status: "success", NumGuesses: numGuesses}
		json.NewEncoder(w).Encode(responseBody)
		numToGuess = 0
		numGuesses = 0
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
