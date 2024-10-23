package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	embd "github.com/JackBekket/hellper/lib/embeddings"
	ghinstallation "github.com/bradleyfalzon/ghinstallation/v2"

	RAG "github.com/JackBekket/GitHelper/pkg/rag"

	"github.com/google/go-github/v65/github"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/vectorstores"
)


var AI string
var API_TOKEN string
var DB string
var ClientMap = make(map[string]*github.Client)




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
			client, err := getClientByRepoOwner(repoOwner)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Respond to the issue
			response, err := generateResponse(issueBody, repoName)
			if err != nil {
				log.Println("Can't generate response")
				log.Println(err)
				response = "Can't generate response bleep-bloop"
			}
			fmt.Println("response generated")
			respond(client, repoOwner, repoName, int64(issueID), response)
		}

	default:
		http.Error(w, "Unknown event type received", http.StatusBadRequest)
	}
}



// Generates retrival-augmented generation taking issue body as prompt, generating response and post it as a comment to github issue
func generateResponse(prompt string, namespace string) (string, error) {
	collection, err := getCollection(AI, API_TOKEN, DB, namespace) // getting all docs from (whole collection) for namespace (repo_name)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("namespace is: ", namespace)

	// join doc and code docs, two of each
	response, err := RAG.RagReflexia(prompt, AI, API_TOKEN, 2, collection)	// call retrival-augmented generation with vectorstore of documents (with type:code and type:doc metadata of it). RAG package DO NOT handle any git operation, such as cloning and so on
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



func createClients(app_id int) (error)  {

	//var white_list []string
	white_list := []string{"GitjobTeam","JackBekket"}	//TODO: change it to load from .env / .yaml and not hardcoded

	var result []*github.Client	// not necessary
	tr := http.DefaultTransport
	pk_name := os.Getenv("PRIVATE_KEY_NAME")

	itr, err := ghinstallation.NewAppsTransportKeyFromFile(tr, int64(app_id), pk_name)
    if err != nil {
        log.Fatal(err)
		return err
    }

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
		return err
	}

	installations, _, err := client.Apps.ListInstallations(context.Background(), &github.ListOptions{})
	if err != nil {
		log.Fatalf("failed to list installations: %v\n", err)
		return err
	}

	// log installations
	log.Println("installations: ", installations)


	//capture our installationId for our app
	//we need this for the access token
	var installID int64
	for _, val := range installations {
		installID = val.GetID()

		user := val.GetAccount()
		user_name := user.GetLogin()
		log.Println("installed by entity_name:", user_name)  // repo owner (?)
		target_type := val.GetTargetType()
		log.Println("target tyoe: ", target_type)


		//If whitelist does not contain our names throw error
		if !contains(white_list, user_name) {
			//log.Println("User %s is not in the whitelist. Skipping creating client for it.", user_name)
			return fmt.Errorf("User %s is not in the whitelist. Skipping creating client for it.", user_name)
		  }


		token, _, err := client.Apps.CreateInstallationToken(
			context.Background(),
			installID,
			&github.InstallationTokenOptions{})
		if err != nil {
			log.Fatalf("failed to create installation token: %v\n", err)
			return err
		}

		apiClient := github.NewClient(nil).WithAuthToken(
			token.GetToken(),
		)
		if apiClient == nil {
			log.Fatalf("failed to create new git client with token: %v\n", err)
			return err
		}

		//append
		result = append(result, apiClient)

		// add to global map user => client
		ClientMap[user_name] = apiClient
	}

	// print clients map (not necessary and can be removed)
	log.Println("clients created: ", result)
	return nil //if not error then we 
}


func getClientByRepoOwner(owner string) (*github.Client,error) {
	client,ok := ClientMap[owner]
	if ok {
		return client, nil
	}
	return nil, fmt.Errorf("client not found for key: %s", owner)


}


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
	// creating clients for each installation of the app
	err = createClients(app_id)
	if err != nil {
		log.Println("error creating clients: ", err)	// might be not in whitelist / unauthorized
	}

	// helper url, helper api token, postgres link with embeddings store
	ai := os.Getenv("AI_ENDPOINT")
	apit := os.Getenv("API_TOKEN")
	db_link := os.Getenv("DB_URL")
	AI = ai
	API_TOKEN = apit
	DB = db_link

	// ... (Set up your webhook endpoint and start the server)
	http.HandleFunc("/webhook", handleWebhook)
	//log.Fatal(http.ListenAndServe(":8086", nil))
	log.Fatal(http.ListenAndServe(":8186", nil))
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
  