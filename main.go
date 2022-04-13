package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
)

// Todo Struct (Model)
type Todo struct {
	ID   string `json:"id"`
	Todo string `json:"todo"`
}

// Init Todo var as a slice Todo struct
var todos []Todo

func Welcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Welcome To My TODO List.")
}

// Get All Todo

func getAlltodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// Get A Todo
func getAtodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get Params

	// Loop through todos and find with id
	for _, item := range todos {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
		}
	}
	json.NewEncoder(w).Encode(&Todo{})
}

// Create A Todo
func createTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todo Todo
	_ = json.NewDecoder(r.Body).Decode(&todo)
	todos = append(todos, todo)
	json.NewEncoder(w).Encode(todo)
}

// Update A Todo
func updateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range todos {
		if item.ID == params["id"] {
			todos = append(todos[:index], todos[index+1:]...)
			var todo Todo
			todo.ID = params["id"]
			_ = json.NewDecoder(r.Body).Decode(&todo)
			todos = append(todos, todo)
			json.NewEncoder(w).Encode(todo)
			return
		}
	}
	json.NewEncoder(w).Encode(todos)
}

// Delete A Todo
func deleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range todos {
		if item.ID == params["id"] {
			todos = append(todos[:index], todos[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(todos)
}

// Connection URI
const uri = "mongodb+srv://rootuser:rootpass@mongodb.7w8dm.mongodb.net/mongoDB?retryWrites=true&w=majority"

func main() {
	// Init Router
	app := mux.NewRouter()

	//Mock Data
	todos = append(todos,
		Todo{"1", "Gym"},
		Todo{"2", "Door Repair"},
		Todo{"3", "Shopping"},
		Todo{"4", "Bike Service"},
	)
	//Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	//client.Database().Collection().
	fmt.Println("Successfully Connected and Pinged")

	// Route Handlers / Endpoints
	app.HandleFunc("/", Welcome).Methods("GET")
	app.HandleFunc("/todo", getAlltodo).Methods("GET")
	app.HandleFunc("/todo/{id}", getAtodo).Methods("GET")
	app.HandleFunc("/todo", createTodo).Methods("POST")
	app.HandleFunc("/todo/{id}", updateTodo).Methods("PUT")
	app.HandleFunc("/todo/{id}", deleteTodo).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", app))
}
