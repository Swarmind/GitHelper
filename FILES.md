# main.go  
package: main  
imports: context, encoding/json, fmt, io, net/http, os, slices, strconv, strings, github.com/JackBekket/GitHelper/internal/database, github.com/JackBekket/GitHelper/internal/reflexia_integration, github.com/JackBekket/GitHelper/pkg/github, github.com/JackBekket/GitHelper/pkg/rag/agent, github.com/JackBekket/hellper/lib/embeddings, github.com/rs/zerolog/log, github.com/google/go-github/v65/github, github.com/joho/godotenv, github.com/tmc/langchaingo/vectorstores  
  
func main():  
	// creating github client from private key  
	_ = godotenv.Load()  
  
	// helper url, helper api token, postgres link with embeddings store  
	AIBaseURL = os.Getenv("AI_URL")  
	AIToken = os.Getenv("API_TOKEN")  
	DBURL = os.Getenv("DB_URL")  
  
	dbHandler, err := database.NewHandler(DBURL)  
	if err != nil {  
		log.Fatal().Err(err).Msg("failed to create database service")  
	}  
  
	if err := dbHandler.DB.Ping(); err != nil {  
		log.Fatal().Err(err).Msg("failed to ping database")  
	}  
	log.Info().Msg("database ping successful")  
  
	DBService, err = database.NewAIService(dbHandler)  
	if err != nil {  
		log.Fatal().Err(err).Msg("initializing AIService")  
	}  
  
	appIdStr := os.Getenv("APP_ID")  
	appId, err := strconv.ParseInt(appIdStr, 10, 64)  
	if err != nil {  
		log.Fatal().Err(err).Msg("parsing APP_ID env variable")  
	}  
	pkPath := os.Getenv("PRIVATE_KEY_NAME")  
	GHService = githubAPI.NewGHService(  
		appId, pkPath,  
		strings.Split(os.Getenv("OWNER_WHITELIST"), ","),  
	)  
  
	// ... (Set up your webhook endpoint and start the server)  
	http.HandleFunc("/webhook", handleWebhook)  
	log.Err(http.ListenAndServe(":8186", nil))  
  
func handleWebhook(w http.ResponseWriter, r *http.Request):  
	ctx := context.Background()  
	defer ctx.Done()  
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
			log.Info().Msgf("App installed for repository: %s/%s", repoOwner, repoName)  
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
			issueId := event.GetIssue().GetNumber()  
			issueTitle := event.GetIssue().GetTitle()  
			issueBody := event.GetIssue().GetBody()  
			author := event.GetIssue().GetUser().GetLogin()  
  
			log.Info().Msgf("New issue opened: #%d %s/%s %s",  
				issueId, repoOwner, repoName, issueTitle,  
			)  
			log.Info().Msgf(  
				"%s: %s", author, issueBody,  
			)  
  
			//  
			client, _, err := GHService.GetClientByRepoOwner(repoOwner)  
			if err != nil {  
				log.Warn().Err(err).Msg("failed to get suitable github client from repo owner installation")  
				http.Error(w, err.Error(), http.StatusBadRequest)  
				return  
			}  
  
			// Respond to the issue  
			response, err := CreateResponse(issueId, issueBody, repoName)  
			if err != nil {  
				log.Warn().Err(err).Msg("failed to generate response")  
				response = "Can't generate response bleep-bloop"  
			}  
			if _, err := githubAPI.CommentIssue(ctx, client, repoOwner, repoName, issueId, response); err != nil {  
				log.Warn().Err(err).Msg("failed to comment issue")  
			}  
  
		}  
		if event.GetAction() == "closed" {  
			repoOwner := event.GetRepo().GetOwner().GetLogin()  
			repoName := event.GetRepo().GetName()  
			issueId := event.GetIssue().GetNumber()  
			issueTitle := event.GetIssue().GetTitle()  
  
			log.Info().Msgf("Issue closed: %s/%s Issue: %d Title: %s",  
				repoOwner, repoName, issueId, issueTitle,  
			)  
  
			model := os.Getenv("MODEL")  
			if err := DBService.DropHistory(int64(issueId), repoName, model); err != nil {  
				log.Warn().Err(err).Msg("failed to drop history")  
			}  
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
			issueId := event.GetIssue().GetNumber()  
			issueTitle := event.GetIssue().GetTitle()  
  
			comment := event.GetComment()  
			commentBody := comment.GetBody()  
			commentUser := comment.GetUser()  
			author := commentUser.GetLogin()  
  
			client, installation, err := GHService.GetClientByRepoOwner(repoOwner)  
			if err != nil {  
				log.Print(err)  
				http.Error(w, err.Error(), http.StatusBadRequest)  
				return  
			}  
  
			if commentUser.GetLogin() == installation.GetAppSlug()+"[bot]" {  
				log.Debug().Msgf("Skipping comment from myself on issue: #%d %s/%s %s",  
					issueId, repoOwner, repoName, issueTitle,  
				)  
				return  
			}  
  
			log.Info().Msgf("New comment on issue: #%d %s/%s %s",  
				issueId, repoOwner, repoName, issueTitle,  
			)  
			log.Info().Msgf(  
				"%s: %s", author, commentBody,  
			)  
  
			response, err := GenerateResponse(issueId, commentBody, repoName)  
			if err != nil {  
				log.Warn().Err(err).Msg("failed to generate response")  
				response = "Can't generate response bleep-bloop"  
			}  
			if _, err := githubAPI.CommentIssue(ctx, client, repoOwner, repoName, issueId, response); err != nil {  
				log.Warn().Err(err).Msg("failed to comment issue")  
			}  
		}  
	case "push":  
		event := new(github.PushEvent)  
		err := json.Unmarshal(requestBody, event)  
		if err != nil {  
			http.Error(w, fmt.Sprintf("Error unmarshalling push event: %v", err), http.StatusBadRequest)  
			return  
		}  
  
		ownerName := event.GetRepo().GetOwner().GetLogin()  
		check := slices.Contains(GHService.WhiteList, ownerName)  
		if !check {  
			log.Warn().Msgf("repo owner %s is not in the whitelist", ownerName)  
			http.Error(w, fmt.Sprintf("Error user (owner) is not in the whitelist: %v", err), http.StatusBadRequest)  
			return  
		}  
  
		ref := event.GetRef()  
		repo := event.GetRepo().GetName()  
		if ref == "refs/heads/master" || ref == "refs/heads/main" {  
			repoURL := fmt.Sprintf("https://github.com/%s/%s", ownerName, repo)  
  
			pkgRunner, err := reflexia.InitPackageRunner(repoURL)  
			if err != nil {  
				log.Warn().Err(err).Msg("initializing reflexia package runner")  
				return  
			}  
			_, _, _, _, err = pkgRunner.RunPackages()  
			if err != nil {  
				log.Warn().Err(err).Msg("package runner RunPackages()")  
				return  
			}  
			log.Info().Msg("push to master or main branch of a repo")  
		} else {  
			log.Warn().Msgf("ref %s is not matched against refs/heads/master || refs/heads/main", ref)  
		}  
	default:  
		http.Error(w, "Unknown event type received", http.StatusBadRequest)  
	}  
  
func CreateResponse(issueId int, prompt string, namespace string):  
	_, err := getCollection(AIBaseURL, AIToken, DBURL, namespace) // getting all docs from (whole collection) for namespace (repo_name)  
	if err != nil {  
		log.Print(err)  
		return "error", err  
	}  
  
	fmt.Println("namespace is: ", namespace)  
  
	model := os.Getenv("MODEL")  
  
	//response, err := RAG.RagReflexia(prompt, AI, API_TOKEN, 2, collection) // call retrival-augmented generation with vectorstore of documents (with type:code and type:doc metadata of it). RAG package DO NOT handle any git operation, such as cloning and so on  
	dialogGraph, response, err := agent.RunNewAgent(AIToken, model, AIBaseURL, prompt, namespace)  
	if err != nil {  
		return "", err  
	}  
	for _, msg := range dialogGraph.ConversationBuffer {  
		err := DBService.UpdateHistory(int64(issueId), namespace, model, msg)  
		if err != nil {  
			return "", err  
		}  
	}  
	return response, nil  
  
func GenerateResponse(issueId int, prompt string, namespace string):  
	_, err := getCollection(AIBaseURL, AIToken, DBURL, namespace) // getting all docs from (whole collection) for namespace (repo_name)  
	if err != nil {  
		log.Print(err)  
		return "error", err  
	}  
  
	fmt.Println("namespace is: ", namespace)  
  
	model := os.Getenv("MODEL")  
  
	buffer, err := DBService.GetHistory(int64(issueId), namespace, model)  
	if err != nil {  
		log.Err(err)  
		return "", err  
	}  
  
	dialogState := database.NewChatSessionGraph(buffer)  
  
	//response, err := RAG.RagReflexia(prompt, AI, API_TOKEN, 2, collection) // call retrival-augmented generation with vectorstore of documents (with type:code and type:doc metadata of it). RAG package DO NOT handle any git operation, such as cloning and so on  
	dialogGraph, response, err := agent.ContinueAgent(AIToken, model, AIBaseURL, prompt, dialogState)  
	if err != nil {  
		return "", err  
	}  
  
	for _, msg := range dialogGraph.ConversationBuffer[len(dialogGraph.ConversationBuffer)-2:] {  
		err := DBService.UpdateHistory(int64(issueId), namespace, model, msg)  
		if err != nil {  
			return "", err  
		}  
	}  
	return response, nil  
  
func getCollection(ai_url string, api_token string, db_link string, namespace string):  
	store, err := embd.GetVectorStoreWithOptions(ai_url, api_token, db_link, namespace) // ai, api, db, namespace  
	if err != nil {  
		return nil, err  
	}  
	return store, nil  
  
  
# main_test.go  
package: main  
imports: os, strconv, strings, testing, github.com/JackBekket/GitHelper/pkg/github, github.com/JackBekket/GitHelper/pkg/github, github.com/JackBekket/GitHelper/pkg/rag/agent, github.com/joho/godotenv, github.com/rs/zerolog/log  
  
func TestAgent(t *testing.T):  
	- loads .env file  
	- iterates through a map of repositories and prompts  
	- for each prompt, it runs a new agent with the given token, model, base URL, prompt, and repository  
	- if the response is empty, it fails the test  
	- logs the request and response  
func TestGithubAPI(t *testing.T):  
	- loads .env file  
	- parses the APP_ID environment variable  
	- creates a new GHService instance with the given app ID, private key path, and owner whitelist  
	- gets a client by repository owner  
	- creates a new issue with the given title and content  
	- comments on the issue with the given content  
	- closes the issue  
	- logs the issue, comment, and closed issue  
  
  
