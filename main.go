package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	//"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/JackBekket/GitHelper/internal/database"
	reflexia "github.com/JackBekket/GitHelper/internal/reflexia_integration"
	RAG "github.com/JackBekket/GitHelper/pkg/rag"
	"github.com/JackBekket/GitHelper/pkg/rag/agent"
	embd "github.com/JackBekket/hellper/lib/embeddings"
	ghinstallation "github.com/bradleyfalzon/ghinstallation/v2"

	"github.com/rs/zerolog/log"

	"github.com/google/go-github/v65/github"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/vectorstores"
)

var AI string
var API_TOKEN string
var DB string
var whiteList = []string{"GitjobTeam", "JackBekket", "MoonSHRD"} //TODO: change it to load from .env / .yaml and not hardcoded
var APP_ID = -1
var DB_SERVICE database.Service

/*
This is github application, which handle updates from github
*/
func main() {
	fmt.Println("main process started")

	// creating github client from private key
	_ = godotenv.Load()
	_id := os.Getenv("APP_ID")
	app_id, err := strconv.Atoi(_id)
	if err != nil {
		// ... handle error
		panic(err)
	}
	APP_ID = app_id

	// helper url, helper api token, postgres link with embeddings store
	ai := os.Getenv("AI_ENDPOINT")
	apit := os.Getenv("API_TOKEN")
	db_link := os.Getenv("DB_URL")
	AI = ai
	API_TOKEN = apit
	DB = db_link

	dbHandler, err := database.NewHandler(db_link)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create database service")
	}

	if err := dbHandler.DB.Ping(); err != nil {
		log.Fatal().Err(err).Msg("failed to ping database")
	}
	log.Info().Msg("database ping successful")

	db_service, err := database.NewAIService(dbHandler)
	if err != nil {
		log.Fatal().Err(err).Msg("something wrong")
	}
	DB_SERVICE = *db_service

	// ... (Set up your webhook endpoint and start the server)
	http.HandleFunc("/webhook", handleWebhook)
	//log.Fatal(http.ListenAndServe(":8086", nil))
	log.Err(http.ListenAndServe(":8186", nil))
}

/*
Handle events we got from github
if it is installation (meaninig that app is installed to new account or repository) -- print this installation
if it is new issue event -- it tries to call RAG with documents collection associated with this repo
*/
func handleWebhook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Extract webhook event type from the header
	eventType := github.WebHookType(r)
	requestBody, err := io.ReadAll(r.Body)
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

			//
			client, _, err := GetClientByRepoOwner(repoOwner)
			if err != nil {
				log.Print(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Respond to the issue
			response, err := createResponse(issueID, issueBody, repoName)
			if err != nil {
				log.Print("Can't generate response")
				log.Print(err)
				response = "Can't generate response bleep-bloop"
			}
			fmt.Println("response generated")
			respond(client, repoOwner, repoName, int64(issueID), response)
		}
		if event.GetAction() == "closed" {
			repoOwner := event.GetRepo().GetOwner().GetLogin()
			repoName := event.GetRepo().GetName()
			issueID := event.GetIssue().GetNumber()
			issueTitle := event.GetIssue().GetTitle()
			//issueBody := event.GetIssue().GetBody()
			fmt.Printf("Issue closed: %s/%s Issue: %d Title: %s\n", repoOwner, repoName, issueID, issueTitle)

			fmt.Printf("Dropping thread")
			model := os.Getenv("MODEL")
			DB_SERVICE.DropHistory(int64(issueID), repoName, model)

		}
	case "issue_comment":
		event := new(github.IssueCommentEvent)
		err := json.Unmarshal(requestBody, event)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error unmarshalling issues event: %v", err), http.StatusBadRequest)
			return
		}
		if event.GetAction() == "created" {
			repoOwner := event.GetRepo().GetOwner().GetLogin()
			repoName := event.GetRepo().GetName()
			issueID := event.GetIssue().GetNumber()
			issueTitle := event.GetIssue().GetTitle()
			fmt.Printf("New comment on issue: %s/%s Issue: %d Title: %s\n", repoOwner, repoName, issueID, issueTitle)

			comment := event.GetComment()
			comment_body := comment.Body
			commentUser := comment.User
			author := commentUser.Login

			client, installation, err := GetClientByRepoOwner(repoOwner)
			if err != nil {
				log.Print(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if commentUser.GetLogin() == installation.GetAppSlug()+"[bot]" {
				return
			}

			fmt.Printf("Issue Comment:" + *comment_body)
			fmt.Printf("Author: " + *author)
			response, err := genResponse(issueID, *comment_body, repoName)
			if err != nil {
				log.Print("Can't generate response")
				log.Print(err)
				response = "Can't generate response bleep-bloop"
			}
			respond(client, repoOwner, repoName, int64(issueID), response)
		}
	case "push":
		event := new(github.PushEvent)
		err := json.Unmarshal(requestBody, event)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error unmarshalling push event: %v", err), http.StatusBadRequest)
			return
		}

		owner_name := event.Repo.Owner.Login
		log.Print("owner of repo: ", *owner_name)
		check := checkWhitelist(*owner_name)
		if !check {
			http.Error(w, fmt.Sprintf("Error user (owner) is not in the whitelist: %v", err), http.StatusBadRequest)
			return
		}

		ref := event.Ref
		log.Print("ref branch of push: ", *ref)
		repo := event.Repo.Name
		log.Print("push into repo name: ", *repo)
		if *ref == "refs/heads/master" || *ref == "refs/heads/main" {
			repoURL := fmt.Sprintf("https://github.com/%s/%s", *owner_name, *repo)

			pkgRunner, err := reflexia.InitPackageRunner(repoURL)
			if err != nil {
				fmt.Println(err)
			}
			_, _, _, _, err = pkgRunner.RunPackages()
			if err != nil {
				fmt.Println(err)
			}
			log.Print("push to master or main branch of a repo")
		}
	default:
		http.Error(w, "Unknown event type received", http.StatusBadRequest)
	}
}

func checkWhitelist(owner_name string) bool {
	//If whitelist does not contain our names return false
	if !contains(whiteList, owner_name) {
		log.Printf("User %s is not in the whitelist. Skipping creating client for it.", owner_name)
		return false
	} else {
		return true
	}
}

// if issue just opened we creating agent and generating response
func createResponse(issue_id int, prompt string, namespace string) (string, error) {
	_, err := getCollection(AI, API_TOKEN, DB, namespace) // getting all docs from (whole collection) for namespace (repo_name)
	if err != nil {
		log.Print(err)
		return "error", err
	}

	fmt.Println("namespace is: ", namespace)

	model := os.Getenv("MODEL")

	//response, err := RAG.RagReflexia(prompt, AI, API_TOKEN, 2, collection) // call retrival-augmented generation with vectorstore of documents (with type:code and type:doc metadata of it). RAG package DO NOT handle any git operation, such as cloning and so on
	dialog_graph, response, err := agent.RunNewAgent(API_TOKEN, model, AI, prompt, namespace)
	if err != nil {
		return "", err
	}
	err = updateHistoryDb(issue_id, namespace, model, dialog_graph)
	if err != nil {
		return "", err
	}
	return response, nil
}

func updateHistoryDb(issue_id int, repo_name string, model_name string, dialog_stack *database.ChatSessionGraph) error {
	buffer := dialog_stack.ConversationBuffer
	last_msg := buffer[len(buffer)-1]
	err := DB_SERVICE.UpdateHistory(int64(issue_id), repo_name, model_name, last_msg)
	if err != nil {
		return err
	}
	return nil
}

// continue thread
func genResponse(issue_id int, prompt string, namespace string) (string, error) {
	_, err := getCollection(AI, API_TOKEN, DB, namespace) // getting all docs from (whole collection) for namespace (repo_name)
	if err != nil {
		log.Print(err)
		return "error", err
	}

	fmt.Println("namespace is: ", namespace)

	model := os.Getenv("MODEL")

	buffer, err := DB_SERVICE.GetHistory(int64(issue_id), namespace, model)
	if err != nil {
		log.Err(err)
		return "", err
	}

	dialog_state := database.NewChatSessionGraph(buffer)

	//response, err := RAG.RagReflexia(prompt, AI, API_TOKEN, 2, collection) // call retrival-augmented generation with vectorstore of documents (with type:code and type:doc metadata of it). RAG package DO NOT handle any git operation, such as cloning and so on
	dialog_graph, response, err := agent.ContinueAgent(API_TOKEN, model, AI, prompt, dialog_state)
	if err != nil {
		return "", err
	}
	err = updateHistoryDb(issue_id, namespace, model, dialog_graph)
	if err != nil {
		return "", err
	}
	return response, nil

}

// OBSOLETE, todo: delete this when prev. genResponse is tested
// Generates retrival-augmented generation taking issue body as prompt, generating response and post it as a comment to github issue
func generateResponse(prompt string, namespace string) (string, error) {
	collection, err := getCollection(AI, API_TOKEN, DB, namespace) // getting all docs from (whole collection) for namespace (repo_name)
	if err != nil {
		log.Print(err)
		return "error", err
	}

	fmt.Println("namespace is: ", namespace)

	// join doc and code docs, two of each

	response, err := RAG.RagReflexia(prompt, AI, API_TOKEN, 2, collection) // call retrival-augmented generation with vectorstore of documents (with type:code and type:doc metadata of it). RAG package DO NOT handle any git operation, such as cloning and so on
	if err != nil {
		return "", err
	}
	return response, nil
}

func respond(client *github.Client, owner string, repo string, id int64, response string) {
	ctx := context.Background()
	// Craft a reply message from the response from the 3rd service.
	replyMessage := response

	// Create a new comment on the issue using the GitHub API.
	a, b, err := client.Issues.CreateComment(ctx, owner, repo, int(id), &github.IssueComment{
		// Configure the comment with the issue's ID and other necessary details.
		Body: &replyMessage,
	})

	fmt.Println("Var #1: ", a)
	fmt.Println("Var #2: ", b)
	if err != nil {
		fmt.Println("Error creating comment on issue: ", err)
	} else {
		fmt.Println("Comment successfully created!")
	}

}

func GetClientByRepoOwner(owner string) (*github.Client, *github.Installation, error) {
	tr := http.DefaultTransport
	pkName := os.Getenv("PRIVATE_KEY_NAME")

	itr, err := ghinstallation.NewAppsTransportKeyFromFile(tr, int64(APP_ID), pkName)
	if err != nil {
		log.Err(err)
	}

	//create git client with app transport
	client := github.NewClient(
		&http.Client{
			Transport: itr,
			Timeout:   time.Second * 30,
		},
	)

	if client == nil {
		log.Print("failed to create git client for app: %v\n", err)
	}

	installations, _, err := client.Apps.ListInstallations(context.Background(), &github.ListOptions{})
	if err != nil {
		log.Print("failed to list installations: %v\n", err)
	}

	for _, installation := range installations {
		log.Print("installation : ", *installation)
		log.Print("repository selection: ", *installation.RepositorySelection)
	}

	//capture our installationId for our app
	//we need this for the access token
	var installID int64
	for _, val := range installations {
		installID = val.GetID()

		user := val.GetAccount()
		username := user.GetLogin()
		log.Print("installed by entity_name:", username) // repo owner (?)
		targetType := val.GetTargetType()
		log.Print("target type: ", targetType)

		if username != owner {
			continue
		}

		//If whitelist does not contain our names throw error
		if !contains(whiteList, username) {
			//log.Println("User %s is not in the whitelist. Skipping creating client for it.", user_name)
			return nil, nil, fmt.Errorf("user %s is not in the whitelist. Skipping creating client for it", username)
		}

		token, _, err := client.Apps.CreateInstallationToken(
			context.Background(),
			installID,
			&github.InstallationTokenOptions{})
		if err != nil {
			log.Print("failed to create installation token: %v\n", err)
		}

		apiClient := github.NewClient(nil).WithAuthToken(
			token.GetToken(),
		)
		if apiClient == nil {
			log.Print("failed to create new git client with token: %v\n", err)
		}

		return apiClient, val, nil
	}

	return nil, nil, fmt.Errorf("client not found for key: %s", owner)
}

func getCollection(ai_url string, api_token string, db_link string, namespace string) (vectorstores.VectorStore, error) {
	store, err := embd.GetVectorStoreWithOptions(ai_url, api_token, db_link, namespace) // ai, api, db, namespace
	if err != nil {
		return nil, err
	}
	return store, nil
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
