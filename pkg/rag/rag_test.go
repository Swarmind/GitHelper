package rag_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	RAG "github.com/JackBekket/GitHelper/pkg/rag"
	embeddings "github.com/JackBekket/hellper/lib/embeddings"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/vectorstores"
)


var AI string
var API_TOKEN string
var DB string
var NS string



func Test_RagWithFilteres(T *testing.T)  {
	
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

	repo_names = []string{"Hellper","gitjob_lk","gitjob-api", "Reflexia"}
	test_prompts = []string{"what is the logic of command package? what is the logic of dialog package?","Explain how Task API works", "in what file is located activity parser?", "where is project config prompt loading happens?" }

	collection, err := getCollection(AI, API_TOKEN, DB, repo_names[0]) // getting all docs from (whole collection) for namespace (repo_name)
	if err != nil {
		log.Println(err)
	}
	/* opts := vectorstores.WithFilters(map[string]string{
		"type": "doc",
	}) */

	fmt.Println("namespace is: ", repo_names[0])


	filters := map[string]any{
		"type": "doc",
	}

	option := vectorstores.WithFilters(filters)


	
	response, err := RAG.RagWithOptions(test_prompts[0], AI, API_TOKEN, 2, collection,option)
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

func Test_RagReflexia(T *testing.T)  {
	
	fmt.Println("this is Test RAG Reflexia")

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

	repo_names = []string{"Hellper","gitjob_lk","gitjob-api", "Reflexia"}
	test_prompts = []string{"what is the logic of command package? what is the logic of dialog package?","Explain how Task API works", "in what file is located activity parser?", "where is project config prompt loading happens?" }

	collection, err := getCollection(AI, API_TOKEN, DB, repo_names[0]) // getting all docs from (whole collection) for namespace (repo_name)
	if err != nil {
		log.Println(err)
	}
	/* opts := vectorstores.WithFilters(map[string]string{
		"type": "doc",
	}) */

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


/*
func Test_Rag(T *testing.T)  {
	
	fmt.Println("this is test RAG simple")

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

	repo_names = []string{"Hellper","gitjob_lk","gitjob-api", "Reflexia"}
	test_prompts = []string{"what is the logic of command package? what is the logic of dialog package?","Explain how Task API works", "in what file is located activity parser?", "where is project config prompt loading happens?" }

	//generateResponse(test_prompts[0],repo_names[0])

	
	
	for i := 0; i < 3; i++ {
		generateResponseStuffQA(test_prompts[i],repo_names[i])
	}
	
		
}
*/

func Test_StuffRag(T *testing.T)  {
	
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

	repo_names = []string{"Hellper","gitjob_lk","gitjob-api", "Reflexia"}
	test_prompts = []string{"what is the logic of command package? what is the logic of dialog package?","Explain how Task API works", "in what file is located activity parser?", "where is project config prompt loading happens?" }

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


func Test_RefinedQA_RAG(T *testing.T)  {
	
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

	repo_names = []string{"Hellper","gitjob_lk","gitjob-api", "Reflexia"}
	test_prompts = []string{"what is the logic of command package? what is the logic of dialog package?","Explain how Task API works", "in what file is located activity parser?", "where is project config prompt loading happens?" }

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






/*
func generateResponseRefinedQA(prompt string, namespace string) (string, error) {
	collection, err := getCollection(AI, API_TOKEN, DB, namespace) // getting all docs from (whole collection) for namespace (repo_name)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("namespace is: ", namespace)
	
	response, err := RAG.RefinedQA_RAG(prompt, AI, API_TOKEN, 2, collection)
	if err != nil {
		return "", err
	}
	return response, nil
	
}


func generateResponseStuffQA(prompt string, namespace string) (string, error) {
	collection, err := getCollection(AI, API_TOKEN, DB, namespace) // getting all docs from (whole collection) for namespace (repo_name)
	if err != nil {
		log.Println(err)
	}


	fmt.Println("namespace is: ", namespace)
	
	response, err := RAG.StuffedQA_Rag(prompt, AI, API_TOKEN, 2, collection)
	if err != nil {
		return "", err
	}
	return response, nil
	
}

func generateResponse(prompt string, namespace string) (string, error) {
	collection, err := getCollection(AI, API_TOKEN, DB, namespace) // getting all docs from (whole collection) for namespace (repo_name)
	if err != nil {
		log.Println(err)
	}


	fmt.Println("namespace is: ", namespace)
	
	response, err := RAG.Rag(prompt, AI, API_TOKEN, 2, collection)
	if err != nil {
		return "", err
	}
	return response, nil
}
*/	

func getCollection(ai_url string, api_token string, db_link string, namespace string) (vectorstores.VectorStore, error) {
	store, err := embeddings.GetVectorStoreWithOptions(ai_url, api_token, db_link, namespace) // ai, api, db, namespace
	if err != nil {
		return nil, err
	}
	return store, nil
}

