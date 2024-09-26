package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/go-github/v45/github"
)

// Define your API endpoint for handling webhook requests.
func handleWebhook(w http.ResponseWriter, r *http.Request) {

	var repoOwner *string
	var repoName *string

	// Extract the issue event details from the webhook payload.
	// ... (Logic to handle webhook payload and extract issue content)
	// read the request body
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(w, "Error reading request body: ", err)
		return
	}

	// unmarshal the JSON payload
	var payload interface{}
	err = json.Unmarshal(requestBody, &payload)
	if err != nil {
		fmt.Fprintln(w, "Error unmarshaling JSON payload: ", err)
		return
	}

	switch event := payload.(type) {
	case github.InstallationRepositoriesEvent:
		// Handle app installation events:
		repoOwner = event.Sender.Login
		repoName = event.RepositoriesAdded[0].Name
		fmt.Println("App installed for repository:", repoOwner, repoName)
		/*
				case github.Installation:
			        // Handle app installation events:
			        repoOwner := event.Repository.Owner
			        repoName := event.Repository.Name
			        fmt.Println("App installed for repository:", repoOwner, repoName)
		*/
	case github.Issue:
		// Handle new issue events:
		issueTitle := event.Title
		issueBody := event.Body
		issueID := event.ID
		fmt.Println("New issue with title:", issueTitle, "and body:", issueBody)

		// ... (Rest of your logic to extract issue content and call 3rd service)
		response := "hardcoded response"

		respond(w, r, *repoOwner, *repoName, *issueID, response)

	default:
		fmt.Fprintln(w, "Unknown event type received")
	}

	//response := "hardcoded response"

	/*
	   // Call your 3rd service with the extracted issue content.
	   response, err := callThirdService(issueContent)
	   if err != nil {
	       fmt.Fprintln(w, "Error calling 3rd service: ", err)
	       return
	   }
	*/

}

func respond(w http.ResponseWriter, r *http.Request, owner string, repo string, id int64, response string) {
	ctx := r.Context()

	// Craft a reply message from the response from the 3rd service.
	replyMessage := fmt.Sprintf("Here's the response from our 3rd service:\n%s", response)

	// Create a new comment on the issue using the GitHub API.
	client := github.NewClient(nil)
	_, _, err := client.Issues.CreateComment(ctx, owner, repo, int(id), &github.IssueComment{
		// Configure the comment with the issue's ID and other necessary details.
		Body: &replyMessage,
	})

	if err != nil {
		fmt.Fprintln(w, "Error creating comment on issue: ", err)
	} else {
		fmt.Fprintln(w, "Comment successfully created!")
	}

}

// Example implementation of your 3rd service call.
func callThirdService(content string) (string, error) {
	// Implement the logic to call your 3rd service here.
	// ... (Replace this with your actual API call)
	// For demonstration, return a fixed response.
	return "From 3rd service!", nil
}

func main() {
	// ... (Set up your webhook endpoint and start the server)
	http.HandleFunc("/webhook", handleWebhook)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
