package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

const (
	version = "1.0.0"
	docRoot = "public"
)

// storage is a map used to simulate our storage system
var storage map[string]Task

// Task represents a Task
type Task struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// TaskRequestOptions represents the options used to create or edit Task
type TaskRequestOptions struct {
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func main() {

	// Initialize the storage map
	storage = make(map[string]Task)

	// Create the router
	router := mux.NewRouter()

	// Task endpoints
	router.Handle("/tasks", http.HandlerFunc(TasksHandler)).Methods(http.MethodGet)
	router.Handle("/tasks", http.HandlerFunc(CreateTaskHandler)).Methods(http.MethodPost)
	router.Handle("/task/{taskID}", http.HandlerFunc(GetTaskHandler)).Methods(http.MethodGet)
	router.Handle("/task/{taskID}", http.HandlerFunc(EditTaskHandler)).Methods(http.MethodPut)
	router.Handle("/task/{taskID}", http.HandlerFunc(DeleteTaskHandler)).Methods(http.MethodDelete)

	// Serve static files: this will serve our docs under http://localhost:9999/static/docs/
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir(docRoot))),
	)

	log.Printf("Starting Todo App Server version='%s' on port :9999\n", version)
	log.Fatal(http.ListenAndServe(":9999", router))
}

// TasksHandler returns the Task collection
func TasksHandler(w http.ResponseWriter, r *http.Request) {

	tasks := []Task{}

	for _, task := range storage {
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(tasks)
	w.Write(body)
}

// CreateTaskHandler creates an new Task
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {

	// Parse and validate request
	options := TaskRequestOptions{}
	err := json.NewDecoder(r.Body).Decode(&options)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		body, _ := json.Marshal(ErrorResponse{
			Message: "unable to parse the body request",
			Error:   err.Error(),
		})
		w.Write(body)
		return
	}

	// Create the Task
	now := time.Now().Format(time.RFC3339)
	task := Task{
		ID:          uuid.NewV4().String(),
		Description: options.Description,
		Completed:   options.Completed,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Persist the Task
	storage[task.ID] = task

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	body, _ := json.Marshal(task)
	w.Write(body)
}

// GetTaskHandler returns a Task by its ID
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	taskID := vars["taskID"]

	var task Task
	var ok bool
	if task, ok = storage[taskID]; !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		body, _ := json.Marshal(ErrorResponse{
			Message: "Not found",
		})
		w.Write(body)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(task)
	w.Write(body)
}

// EditTaskHandler edits a Task by its ID
func EditTaskHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	taskID := vars["taskID"]

	var task Task
	var ok bool
	if task, ok = storage[taskID]; !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		body, _ := json.Marshal(ErrorResponse{
			Message: "Not found",
		})
		w.Write(body)
		return
	}

	options := TaskRequestOptions{}
	err := json.NewDecoder(r.Body).Decode(&options)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		body, _ := json.Marshal(ErrorResponse{
			Message: "unable to parse the body request",
			Error:   err.Error(),
		})
		w.Write(body)
		return
	}

	task.Description = options.Description
	task.Completed = options.Completed
	task.UpdatedAt = time.Now().Format(time.RFC3339)

	storage[task.ID] = task

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(task)
	w.Write(body)
}

// DeleteTaskHandler deletes a Task by its ID
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["taskID"]

	delete(storage, taskID)

	w.WriteHeader(http.StatusNoContent)
}
