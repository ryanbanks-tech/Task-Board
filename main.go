package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Task struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Completed bool   `json:"completed"`
	Assignee  string `json:"assignee_id"`
	Project   string `json:"project_id"`
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var todoistBearerToken string = os.Getenv("TODOIST_BEARER_TOKEN")

func getTasks(w http.ResponseWriter, r *http.Request) {
	userMap := make(map[string]string)

	// Fetch projects
	projects, err := fetchProjects()
	if err != nil {
		http.Error(w, "Error fetching projects", http.StatusInternalServerError)
		return
	}

	// Fetch tasks
	tasks, err := fetchTasks()
	if err != nil {
		http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
		return
	}

	// Fetch users
	users, err := fetchUsers()
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}

	// Map users
	for _, user := range users {
		userMap[user.ID] = user.Name
	}

	// Create response structure
	var response struct {
		Projects []Project `json:"projects"`
		Tasks    []Task    `json:"tasks"`
	}
	response.Projects = projects
	response.Tasks = tasks

	// Encode response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func fetchProjects() ([]Project, error) {
	req, err := http.NewRequest("GET", "https://api.todoist.com/rest/v2/projects", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", todoistBearerToken))

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, err
	}
	return projects, nil
}

func fetchTasks() ([]Task, error) {
	req, err := http.NewRequest("GET", "https://api.todoist.com/rest/v2/tasks", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", todoistBearerToken))

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tasks []Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func fetchUsers() ([]User, error) {
	req, err := http.NewRequest("GET", "https://api.todoist.com/rest/v2/projects/2325372207/collaborators", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", todoistBearerToken))

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}
	return users, nil
}

func main() {
	http.HandleFunc("/tasks", getTasks)
	log.Fatal(http.ListenAndServe(":8990", nil))
}
