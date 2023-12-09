package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/cors"
)

type todo struct {
	Id string `json:"id"`
	Todo string `json:"todo"`
	Timestamp int64 `json:"timestamp"`
}
var todos []todo = []todo{}

func todosEndpoint(w http.ResponseWriter, req *http.Request) {
	if(req.Method != "GET") {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	jsonData, err := json.Marshal(todos)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func addTodoEndpoint(w http.ResponseWriter, req *http.Request) {
	if(req.Method != "POST") {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	var newTodo todo

	newTodo.Todo = string(body)
	newTodo.Id = uuid.New().String()
	newTodo.Timestamp = time.Now().Unix()


	todos = append(todos, newTodo)
	
	jsonData, err := json.Marshal(todos)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

func main() {
	const port string = "3001"

	corsHandler := cors.Default()
	http.HandleFunc("/todos", func(w http.ResponseWriter, req *http.Request) {
		var handler http.Handler = corsHandler.Handler(http.HandlerFunc(todosEndpoint))
		handler.ServeHTTP(w, req)
	})
	http.HandleFunc("/addTodo", func(w http.ResponseWriter, req *http.Request) {
		var handler http.Handler = corsHandler.Handler(http.HandlerFunc(addTodoEndpoint))
		handler.ServeHTTP(w, req)
	})
	http.ListenAndServe("localhost:"+port, nil)
}