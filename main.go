package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	dbPath string = ".sqlite"
	port   string = ":8000"
)

var (
	cat Category
	db  = connection(dbPath)
)

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("new request from:", r.Method, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		next.ServeHTTP(w, r)
	})
}

func getCategory(w http.ResponseWriter, r *http.Request) {
	categories := cat.read(db)
	json.NewEncoder(w).Encode(&categories)
}

func createCategory(w http.ResponseWriter, r *http.Request) {
	json.NewDecoder(r.Body).Decode(&cat)
	id, err := cat.create(db, cat.Name, cat.Description, cat.IsLimit)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(id)
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	json.NewDecoder(r.Body).Decode(&cat)
	result, err := cat.update(db, params["id"], cat.Description)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(result)
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	result, err := cat.delete(db, params["id"])
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(result)

}

func runAPI() {
	r := mux.NewRouter()
	r.Use(logging)
	r.Use(cors)
	r.HandleFunc("/category", getCategory).Methods("GET")
	r.HandleFunc("/category", createCategory).Methods("POST")
	r.HandleFunc("/category/{id}", updateCategory).Methods("PUT")
	r.HandleFunc("/category/{id}", deleteCategory).Methods("DELETE")
	log.Println("server listen on port", port)
	log.Fatal(http.ListenAndServe(port, r))
}

func main() {
	defer db.Close()
	runAPI()
}
