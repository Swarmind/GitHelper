package rag_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	RAG "github.com/JackBekket/GitHelper/legacy"
	embeddings "github.com/JackBekket/hellper/lib/embeddings"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/vectorstores"
)

var AI string
var API_TOKEN string
var DB string
var NS string

// This test call RAG with optional filteres
func Test_RagWithFilteres(T *testing.T) {
	fmt.Println("This is Rag with Filteres test")

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
	NS = "gitjob-api"

	repo_names = []string{"Hellper", "gitjob_lk", "gitjob-api", "Reflexia"}
	test_prompts = []string{"what is the logic of command package? what is the logic of dialog package?", "Explain how Task API works", "in what file is located activity parser?", "where is project config prompt loading happens?"}

	collection, err := getCollection(AI, API_TOKEN, DB, repo_names[0]) // getting all docs from (whole collection) for namespace (repo_name)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("namespace is: ", repo_names[0])

	filters := map[string]any{
		"type": "doc",
	}
	option := vectorstores.WithFilters(filters)

	response, err := RAG.RagWithOptions(test_prompts[0], AI, API_TOKEN, 2, collection, option)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("(RAG DOC ONLY)Filtered response: ", response)

	/*
		for i := 0; i < 3; i++ {
			generateResponseStuffQA(test_prompts[i],repo_names[i])
		}
	*/

}

// This test function is calling Retrival-Augmented generation with two types of document (doc's and code) and call 'stuffed' method of RAG
func Test_RagReflexia(T *testing.T) {
	fmt.Println("this is Test RAG with 'type: doc + type: code' metadata")

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
	NS = "gitjob-api"

	repo_names = []string{"Hellper", "gitjob_lk", "gitjob-api", "Reflexia"}
	test_prompts = []string{"what is the logic of command package? what is the logic of dialog package?", "Explain how Task API works", "in what file is located activity parser?", "where is project config prompt loading happens?"}

	collection, err := getCollection(AI, API_TOKEN, DB, repo_names[0]) // getting all docs from (whole collection) for namespace (repo_name)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("namespace is: ", repo_names[0])

	response, err := RAG.RagReflexia(test_prompts[0], AI, API_TOKEN, 2, collection)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("(RAG REFLEXIA)Filtered response: ", response)

	/*
		for i := 0; i < 3; i++ {
			generateResponseStuffQA(test_prompts[i],repo_names[i])
		}
	*/
}

// This func is testing stuffed method of RAG
func Test_StuffRag(T *testing.T) {

	println("this is STUFF RAG TEST")

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
	NS = "gitjob-api"

	repo_names = []string{"Hellper", "gitjob_lk", "gitjob-api", "Reflexia"}
	test_prompts = []string{"what is the logic of command package? what is the logic of dialog package?", "Explain how Task API works", "in what file is located activity parser?", "where is project config prompt loading happens?"}

	//generateResponseStuffQA(test_prompts[0],repo_names[0])
	collection, err := getCollection(AI, API_TOKEN, DB, repo_names[0]) // getting all docs from (whole collection) for namespace (repo_name)
	if err != nil {
		log.Println(err)
	}

	response, err := RAG.StuffedQA_Rag(test_prompts[0], AI, API_TOKEN, 2, collection)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("(RAG STUFFED QA)response: ", response)

	/*
		for i := 0; i < 3; i++ {
			generateResponseStuffQA(test_prompts[i],repo_names[i])
		}
	*/

}

// This test is calling Refined QA method of RAG
func Test_RefinedQA_RAG(T *testing.T) {

	println("this is REFINED QA test")

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
	NS = "gitjob-api"

	repo_names = []string{"Hellper", "gitjob_lk", "gitjob-api", "Reflexia"}
	test_prompts = []string{"what is the logic of command package? what is the logic of dialog package?", "Explain how Task API works", "in what file is located activity parser?", "where is project config prompt loading happens?"}

	//generateResponseRefinedQA(test_prompts[0],repo_names[0])

	collection, err := getCollection(AI, API_TOKEN, DB, repo_names[0]) // getting all docs from (whole collection) for namespace (repo_name)
	if err != nil {
		log.Println(err)
	}

	response, err := RAG.StuffedQA_Rag(test_prompts[0], AI, API_TOKEN, 2, collection)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("(RAG REFINED QA)response: ", response)

	/*
		for i := 0; i < 3; i++ {
			generateResponseRefinedQA(test_prompts[i],repo_names[i])
		}
	*/
}

func getCollection(ai_url string, api_token string, db_link string, namespace string) (vectorstores.VectorStore, error) {
	store, err := embeddings.GetVectorStoreWithOptions(ai_url, api_token, db_link, namespace) // ai, api, db, namespace
	if err != nil {
		return nil, err
	}
	return store, nil
}
