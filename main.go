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
	router.Use(loggingMiddleware)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("index.html")
		http.ServeFile(w, r, "static/index.html")
	})

	router.HandleFunc("/guessing", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("guessing.html")
		http.ServeFile(w, r, "static/guessing.html")
	})

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("login part")
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
		fmt.Println("guess part")
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

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Extract the path of the request URL
	path := r.URL.Path

	// Use the path to route the request to different handlers
	switch path {
	case "/":
		homeHandler(w, r)
	case "/login":
		loginHandler(w, r)
	case "/guess":
		guessHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}
