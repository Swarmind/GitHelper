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

var Client *github.Client

// Define your API endpoint for handling webhook requests.
func handleWebhook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	client := Client

	// Extract webhook event type from the header
	eventType := github.WebHookType(r)
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading request body: %v", err), http.StatusBadRequest)
		return
	}
	// Parse event based on eventType
	switch eventType {
	case "installation_repositories":
		event := new(github.InstallationRepositoriesEvent)
		err := json.Unmarshal(requestBody, event)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error unmarshalling installation_repositories event: %v", err), http.StatusBadRequest)
			return
		}
		if len(event.RepositoriesAdded) > 0 {
			repoOwner := event.Sender.GetLogin()
			repoName := *event.RepositoriesAdded[0].Name
			fmt.Printf("App installed for repository: %s/%s\n", repoOwner, repoName)
		}

	case "issues":
		event := new(github.IssuesEvent)
		err := json.Unmarshal(requestBody, event)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error unmarshalling issues event: %v", err), http.StatusBadRequest)
			return
		}

		if event.GetAction() == "opened" {
			repoOwner := event.GetRepo().GetOwner().GetLogin()
			repoName := event.GetRepo().GetName()
			issueID := event.GetIssue().GetNumber()
			issueTitle := event.GetIssue().GetTitle()
			issueBody := event.GetIssue().GetBody()

			fmt.Printf("New issue opened: %s/%s Issue: %d Title: %s\n", repoOwner, repoName, issueID, issueTitle)
			fmt.Println("Issue body is: ", issueBody)

			// Respond to the issue
			response := "hardcoded response"
			respond(w, r, client, repoOwner, repoName, int64(issueID), response)
		}

	default:
		http.Error(w, "Unknown event type received", http.StatusBadRequest)
	}
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
	//const gitHost = "https://git.api.com"

	privatePem, err := os.ReadFile(key_path)
	if err != nil {
		log.Fatalf("failed to read pem: %v", err)
	}

	itr, err := ghinstallation.NewAppsTransport(http.DefaultTransport, int64(app_id), privatePem)
	if err != nil {
		log.Fatalf("failed to create app transport: %v\n", err)
	}
	//itr.BaseURL = gitHost

	//create git client with app transport
	client := github.NewClient(
		&http.Client{
			Transport: itr,
			Timeout:   time.Second * 30,
		},
	)
	//)

	if client == nil {
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

	apiClient := github.NewClient(nil).WithAuthToken(
		token.GetToken(),
	)
	if apiClient == nil {
		log.Fatalf("failed to create new git client with token: %v\n", err)
	}

	//log.Println("gh client: ", apiClient)

	return apiClient
}

/* // Example implementation of your 3rd service call.
func callThirdService(content string) (string, error) {
	// Implement the logic to call your 3rd service here.
	// ... (Replace this with your actual API call)
	// For demonstration, return a fixed response.
	return "From 3rd service!", nil
} */

func main() {

	fmt.Println("main process started")

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
	Client = createClient(pk_path, app_id)

	// ... (Set up your webhook endpoint and start the server)
	http.HandleFunc("/webhook", handleWebhook)
	//log.Fatal(http.ListenAndServe(":8086", nil))
	log.Fatal(http.ListenAndServe(":8186", nil))
}
