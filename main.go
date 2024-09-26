package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	ghinstallation "github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v65/github"
	"github.com/joho/godotenv"
)

// Define your API endpoint for handling webhook requests.
func handleWebhook(w http.ResponseWriter, r *http.Request) {

	var repoOwner *string
	var repoName *string

	// creating github client from private key
	_ = godotenv.Load()
	_id := os.Getenv("APP_ID")
	//wh_secret := os.Getenv("WEBHOOK_SECRET")
	pk_path := os.Getenv("PRIVATE_KEY_PATH")
	app_id, err := strconv.Atoi(_id)
	if err != nil {
		// ... handle error
		panic(err)
	}
	client := createClient(pk_path, app_id)

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

		respond(w, r, client, *repoOwner, *repoName, *issueID, response)

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

func respond(w http.ResponseWriter, r *http.Request, client *github.Client, owner string, repo string, id int64, response string) {
	ctx := r.Context()

	// Craft a reply message from the response from the 3rd service.
	replyMessage := fmt.Sprintf("Here's the response from our 3rd service:\n%s", response)

	// Create a new comment on the issue using the GitHub API.

	//client := github.NewClient(nil)
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

func createClient(key_path string, app_id int) *github.Client {
	const gitHost = "https://git.api.com"

	privatePem, err := os.ReadFile(key_path)
	if err != nil {
		log.Fatalf("failed to read pem: %v", err)
	}

	itr, err := ghinstallation.NewAppsTransport(http.DefaultTransport, int64(app_id), privatePem)
	if err != nil {
		log.Fatalf("failed to create app transport: %v\n", err)
	}
	itr.BaseURL = gitHost

	//create git client with app transport
	client, err := github.NewClient(
		&http.Client{
			Transport: itr,
			Timeout:   time.Second * 30,
		},
	).WithEnterpriseURLs(gitHost, gitHost)
	//)

	if err != nil {
		log.Fatalf("failed to create git client for app: %v\n", err)
	}

	installations, _, err := client.Apps.ListInstallations(context.Background(), &github.ListOptions{})
	if err != nil {
		log.Fatalf("failed to list installations: %v\n", err)
	}

	//capture our installationId for our app
	//we need this for the access token
	var installID int64
	for _, val := range installations {
		installID = val.GetID()
	}

	token, _, err := client.Apps.CreateInstallationToken(
		context.Background(),
		installID,
		&github.InstallationTokenOptions{})
	if err != nil {
		log.Fatalf("failed to create installation token: %v\n", err)
	}

	apiClient, err := github.NewClient(nil).WithAuthToken(
		token.GetToken(),
	).WithEnterpriseURLs(gitHost, gitHost)
	if err != nil {
		log.Fatalf("failed to create new git client with token: %v\n", err)
	}

	return apiClient
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
