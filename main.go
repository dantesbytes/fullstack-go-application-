package mian

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"


)

type User struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

// main functions

func main() {
    db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
    if err!= nil {
        log.Fatal(err)
    }
    defer db.Close()

	// create table if not exists 

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT, email TEXT)")
	if err!= nil {
        log.Fatal(err)
    }


	// routers 

    router := mux.NewRouter()
	router.HandleFunc("/api/go/users", getUsers(db)).Methods("GET")
	router.HandleFunc("/api/go/users", createUser(db)).Methods("POST")
	router.HandleFunc("/api/go/users/{id}", getUser(db)).Methods("GET")
	router.HandleFunc("/api/go/users/{id}", updateUser(db)).Methods("PUT")
	router.HandleFunc("/api/go/users/{id}", deleteUser(db)).Methods("DELETE")


	//wrap router with cors and json content type middleware
	
	enhancedRouter := enableCors(jsonContentTypeMiddleware(router))

	//start server 
	log.Fatal(http.ListenAndServe(":8080", enhancedRouter))


	func enableCors(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

            w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")

			//check if the request is for cors prefligh
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
                return
            }

			next.ServeHTTP(w, r)
            
        })

    }


	func jsonContentTypeMiddleware(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//set json content-type
            w.Header().Set("Content-Type", "application/json")
            next.ServeHTTP(w, r)
        })
    }


	//get all users 
	



}

