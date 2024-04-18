package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Task represents a Todoist task
type Task struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Completed bool   `json:"completed"`
}

var todoistBearerToken string = os.Getenv("TODOIST_BEARER_TOKEN")

func main() {
	// Construct the HTTP request
	req, err := http.NewRequest("GET", "https://api.todoist.com/rest/v2/projects", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", todoistBearerToken))

	// Make the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching projects:", err)
		return
	}
	defer resp.Body.Close()

	// Decode JSON response
	var projects []Task
	err = json.NewDecoder(resp.Body).Decode(&projects)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Print projects
	fmt.Println("Projects:")
	for _, project := range projects {
		fmt.Printf("ID: %s, Content: %s, Completed: %t\n", project.ID, project.Content, project.Completed)
	}
}
