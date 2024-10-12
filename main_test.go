package main_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	embd "github.com/JackBekket/hellper/lib/embeddings"
	embeddings "github.com/JackBekket/hellper/lib/embeddings"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/vectorstores"
)


var AI string
var API_TOKEN string
var DB string
var NS string



func Test_main (t *testing.T)   {


	_ = godotenv.Load()


	//Test getting vectorstore from .env
	// In production name should be replaced by event value
	ai := os.Getenv("AI_ENDPOINT")
	apit := os.Getenv("API_TOKEN")
	db_link := os.Getenv("DB_URL")
	//namesp := os.Getenv("REPO_NAME")

	AI = ai
	API_TOKEN = apit
	DB = db_link
	//NS = "Hellper"
	NS = "gitjob_lk"
	//NS = "gitjob_api"

	//generateResponse("what is the logic of command package? what is the logic of dialog package?",NS)
	generateResponse("where is located login alhorithm and how does it work?",NS)
	//generateResponse("where is located subscription alhorithm and how does it work?",NS)
}


func getCollection(ai_url string, api_token string, db_link string, namespace string) (vectorstores.VectorStore, error) {
	store, err := embd.GetVectorStoreWithOptions(ai_url, api_token, db_link, namespace) // ai, api, db, namespace
	if err != nil {
		return nil, err
	}
	return store, nil
}



func generateResponse(prompt string, namespace string) (string, error) {
	collection, err := getCollection(AI, API_TOKEN, DB, namespace) // getting all docs from (whole collection) for namespace (repo_name)
	if err != nil {
		log.Println(err)
	}
	/* opts := vectorstores.WithFilters(map[string]string{
		"type": "doc",
	}) */

		fmt.Println("namespace is: ", namespace)

		//🤕🤕🤕
		searchResults, err := embeddings.SemanticSearch(prompt, 2, collection)
		if err != nil {
			return "", err
		}
	
		contextBuilder := strings.Builder{}
		for _, doc := range searchResults {
			contextBuilder.WriteString(doc.PageContent)
			contextBuilder.WriteString("\n")
		}
		contexts := contextBuilder.String()
	
		fmt.Sprintf("Context: %s\n\nQuestion: %s", contexts, prompt)


	
	response, err := rag(prompt, AI, API_TOKEN, 1, collection)
	if err != nil {
		return "", err
	}
	return response, nil
	
}


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




