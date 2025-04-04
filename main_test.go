package main

import (
	"os"
	"testing"

	"github.com/JackBekket/GitHelper/pkg/rag/agent"
	"github.com/joho/godotenv"
)

//var AI string
//var API_TOKEN string
//var DB string
var NS string



//TODO: use utils (pkg/github) funcs to create and comment issue
//TODO: add tests for creating, commenting and closing issue



// TODO: refactor it as outdated
func Test_main (t *testing.T)   {


	_ = godotenv.Load()


	//Test getting vectorstore from .env
	// In production name should be replaced by event value
	ai := os.Getenv("AI_ENDPOINT")
	apit := os.Getenv("API_TOKEN")
	db_link := os.Getenv("DB_URL")

	// test data
	var repo_names []string
	var test_prompts []string

	AI = ai
	API_TOKEN = apit
	DB = db_link
	//NS = "gitjob-api"
	model := os.Getenv("MODEL")

	repo_names = []string{"Hellper", "Reflexia"}
	test_prompts = []string{"what is the logic of command package? what is the logic of dialog package?", "where is project config prompt loading happens?" }

	//generateResponse(test_prompts[2],repo_names[2])
	
	
	
	for i := 0; i < 2; i++ {
		//GenerateResponse(test_prompts[i],repo_names[i])
		agent.RunNewAgent(API_TOKEN,model,AI,test_prompts[i],repo_names[i])
	}
		
}




/*
func rag(question string, ai_url string, api_token string, numOfResults int, store vectorstores.VectorStore) (result string, err error) {
	//base_url := os.Getenv("AI_BASEURL")
	base_url := ai_url

	// Create an embeddings client using the specified API and embedding model
	llm, err := openai.New(
		openai.WithBaseURL(base_url),
		openai.WithAPIVersion("v1"),
		openai.WithToken(api_token),
		openai.WithModel("tiger-gemma-9b-v1-i1"),
		openai.WithEmbeddingModel("text-embedding-ada-002"),
	)
	if err != nil {
		return "", err
	}

	//🤕🤕🤕
	searchResults, err := embeddings.SemanticSearch(question, numOfResults, store)
	if err != nil {
		return "", err
	}

	contextBuilder := strings.Builder{}
	for _, doc := range searchResults {
		contextBuilder.WriteString(doc.PageContent)
		contextBuilder.WriteString("\n")
	}
	contexts := contextBuilder.String()

	fullPrompt := fmt.Sprintf("Context: %s\n\nQuestion: %s", contexts, question)

	result, err = chains.Run(
		context.Background(),
		chains.NewRetrievalQAFromLLM(
			llm,
			vectorstores.ToRetriever(store, numOfResults),
		),
		fullPrompt,
		chains.WithMaxTokens(8192),
	)
	if err != nil {
		return "", err
	}

	fmt.Println("====final answer====\n", result)

	return result, nil
}
*/



