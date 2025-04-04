package main

import (
	"os"
	"strings"
	"testing"

	"github.com/JackBekket/GitHelper/pkg/rag/agent"
	"github.com/joho/godotenv"
)

var NS string

//TODO: use utils (pkg/github) funcs to create and comment issue
//TODO: add tests for creating, commenting and closing issue

// TODO: refactor it as outdated
func Test_main(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatal(err)
	}

	//Test getting vectorstore from .env
	// In production name should be replaced by event value
	baseURL := os.Getenv("AI_URL")
	token := os.Getenv("API_TOKEN")
	model := os.Getenv("MODEL")

	for repo, prompts := range map[string][]string{
		"Hellper": {
			"what is the logic of command package?\nwhat is the logic of dialog package?",
		},
		"Reflexia": {
			"where is project config prompt loading happens?",
		},
	} {
		for _, prompt := range prompts {
			_, resp, err := agent.RunNewAgent(token, model, baseURL, prompt, repo)
			if err != nil {
				t.Fatal(err)
			}
			if strings.TrimSpace(resp) == "" {
				t.Fail()
				continue
			}
			t.Logf("Request: %s\nResponse: %s", prompt, resp)
		}
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
