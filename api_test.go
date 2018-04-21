package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func setup() {
	// Initialize the tasks map
	storage = make(map[string]Task)

	// Add a fake task
	now := time.Now().Format(time.RFC3339)
	storage["6ba7b810-9dad-11d1-80b4-00c04fd430c8"] = Task{
		ID:          "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		Description: "Buy milk",
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func TestTasksHandler(t *testing.T) {
	setup()

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rr := httptest.NewRecorder()

	TasksHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %d", status, http.StatusOK)
	}
	// START code to display in slide OMIT
	tasks := []Task{}

	err := json.NewDecoder(rr.Body).Decode(&tasks)

	if err != nil {
		t.Errorf("unable to decode body: got %s", rr.Body.String())
	}

	if len(tasks) != 1 {
		t.Errorf("unexpected len: got %d want %d", len(tasks), 1)
	}

	task := tasks[0]

	if task.ID != "6ba7b810-9dad-11d1-80b4-00c04fd430c8" {
		t.Errorf("unexpected ID: got %v want %v", task.ID, "6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	}
	// END code to display in slide OMIT

}
func TestCreateTaskHandler(t *testing.T) {

	setup()

	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`
		{
			"description": "Buy bread",
			"completed": false
		}
`))
	rr := httptest.NewRecorder()

	CreateTaskHandler(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("wrong status code: got %v want %d", status, http.StatusOK)
	}

	task := Task{}

	err := json.NewDecoder(rr.Body).Decode(&task)

	if err != nil {
		t.Errorf("unable to decode response body: got %s", rr.Body.String())
	}

	if task.Description != "Buy bread" {
		t.Errorf("unexpected ID: got %v want %v", task.Description, "Buy bread")
	}

	if task.Completed != false {
		t.Errorf("unexpected ID: got %t want %t", task.Completed, false)
	}
}

func TestGetTaskHandler(t *testing.T) {
	setup()

	t.Run("get valid task", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/tasks/{taskId}", nil)
		// sets the URL variables for the given request, to be accessed via mux.Vars for testing route behaviour.
		req = mux.SetURLVars(req, map[string]string{"taskID": "6ba7b810-9dad-11d1-80b4-00c04fd430c8"})

		rr := httptest.NewRecorder()

		GetTaskHandler(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("wrong status code: got %v want %d", status, http.StatusOK)
		}
	})
	t.Run("task not found", func(t *testing.T) {

		req, _ := http.NewRequest(http.MethodGet, "/tasks/{taskId}", nil)
		// sets the URL variables for the given request, to be accessed via mux.Vars for testing route behaviour.
		req = mux.SetURLVars(req, map[string]string{"taskID": "00c04fd430c8-6ba7b810-9dad-11d1-80b4"})

		rr := httptest.NewRecorder()

		GetTaskHandler(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("wrong status code: got %v want %d", status, http.StatusOK)
		}
	})
}

func TestEditTaskHandler(t *testing.T) {
	setup()

	req := httptest.NewRequest(http.MethodPut, "/tasks/{taskId}", strings.NewReader(`
			{
				"description": "Buy bread",
				"completed": true
			}
	`))

	// sets the URL variables for the given request, to be accessed via mux.Vars for testing route behaviour.
	req = mux.SetURLVars(req, map[string]string{"taskID": "6ba7b810-9dad-11d1-80b4-00c04fd430c8"})

	rr := httptest.NewRecorder()

	EditTaskHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %d", status, http.StatusOK)
	}

	task := Task{}

	err := json.NewDecoder(rr.Body).Decode(&task)

	if err != nil {
		t.Errorf("unable to decode response body: got %s", rr.Body.String())
	}

	if task.ID != "6ba7b810-9dad-11d1-80b4-00c04fd430c8" {
		t.Errorf("unexpected ID: got %v want %v", task.ID, "6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	}

	if task.Description != "Buy bread" {
		t.Errorf("unexpected ID: got %v want %v", task.Description, "Buy bread")
	}

	if task.Completed != true {
		t.Errorf("unexpected ID: got %t want %t", task.Completed, true)
	}

}

func TestDeleteTaskHandler(t *testing.T) {
	setup()

	req := httptest.NewRequest(http.MethodDelete, "/tasks/{taskId}", nil)

	// sets the URL variables for the given request, to be accessed via mux.Vars for testing route behaviour.
	req = mux.SetURLVars(req, map[string]string{"taskID": "6ba7b810-9dad-11d1-80b4-00c04fd430c8"})

	rr := httptest.NewRecorder()

	DeleteTaskHandler(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("wrong status code: got %v want %d", status, http.StatusNoContent)
	}
}
