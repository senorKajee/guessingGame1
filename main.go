package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Result struct {
	Result string `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var correctNumber = 42

func main() {
	router := mux.NewRouter()
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "login.html")
	})

	router.HandleFunc("/guessing", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "guessing.html")
	})

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request payload"})
			return
		}

		// TODO: Add actual authentication logic here
		if creds.Username != "username" || creds.Password != "password" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid username or password"})
			return
		}

		token := "dummy-token"
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	})

	router.HandleFunc("/guess", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" || token != "Bearer dummy-token" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid or missing token"})
			return
		}

		var guess struct {
			Guess int `json:"guess"`
		}
		err := json.NewDecoder(r.Body).Decode(&guess)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request payload"})
			return
		}

		result := "too low"
		if guess.Guess > correctNumber {
			result = "too high"
		} else if guess.Guess == correctNumber {
			result = "correct"
		}

		json.NewEncoder(w).Encode(Result{Result: result})
	})

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
