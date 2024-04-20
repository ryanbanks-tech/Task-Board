package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Project represents a Todoist project
type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Task represents a Todoist task
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

func main() {

	// Define a map to store ID-user mappings
	userMap := make(map[string]string)

	// Construct the HTTP request for projects
	req, err := http.NewRequest("GET", "https://api.todoist.com/rest/v2/projects", nil)
	if err != nil {
		fmt.Println("Error creating project request:", err)
		return
	}

	// Set the authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", todoistBearerToken))

	// Construct the HTTP request for project tasks
	taskReq, err := http.NewRequest("GET", "https://api.todoist.com/rest/v2/tasks", nil)
	if err != nil {
		fmt.Println("Error creating project tasks request:", err)
		return
	}

	// Set the authorization header for the tasks request
	taskReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", todoistBearerToken))

	// Construct the HTTP request for project users
	userReq, err := http.NewRequest("GET", "https://api.todoist.com/rest/v2/projects/2325372207/collaborators", nil)
	if err != nil {
		fmt.Println("Error creating project users request:", err)
		return
	}

	// Set the authorization header for the usersa request
	userReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", todoistBearerToken))

	// Make the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching projects:", err)
		return
	}
	defer resp.Body.Close()

	// Make the request for project tasks
	taskClient := http.Client{}
	taskResp, err := taskClient.Do(taskReq)
	if err != nil {
		fmt.Println("Error fetching project tasks:", err)
		return
	}
	defer taskResp.Body.Close()

	// Make the request for project users
	userClient := http.Client{}
	userResp, err := userClient.Do(userReq)
	if err != nil {
		fmt.Println("Error fetching project users:", err)
		return
	}
	defer userResp.Body.Close()

	// Decode JSON response
	var projects []Project
	err = json.NewDecoder(resp.Body).Decode(&projects)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	var projectTasks []Task
	err = json.NewDecoder(taskResp.Body).Decode(&projectTasks)
	if err != nil {
		fmt.Println("Error decoding project tasks JSON:", err)
		return
	}

	var projectUsers []User
	err = json.NewDecoder(userResp.Body).Decode(&projectUsers)
	if err != nil {
		fmt.Println("Error decoding project users JSON:", err)
		return
	}

	// User mapping
	for _, projectUser := range projectUsers {
		userMap[projectUser.ID] = projectUser.Name
	}

	// Print projects
	fmt.Println("Projects:")
	for _, project := range projects {
		fmt.Printf("ID: %s, Name: %s\n", project.ID, project.Name)
	}

	// Print project tasks
	fmt.Println("Project Tasks")
	for _, projectTask := range projectTasks {
		fmt.Printf("ID: %s, Assignee: %s:%s, Project: %s, Content: %s, Completed: %t\n", projectTask.ID, projectTask.Assignee, userMap[projectTask.Assignee], projectTask.Project, projectTask.Content, projectTask.Completed)
	}
}
