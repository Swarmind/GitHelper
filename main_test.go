package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/JackBekket/GitHelper/pkg/github"
	"github.com/JackBekket/GitHelper/pkg/rag/agent"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

var NS string


//TODO: add tests for creating, commenting and closing issue


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


func Test_createIssue(t *testing.T) {

	err := godotenv.Load()
	if err != nil {
		t.Fatal(err)
	}
	repoOwner := "JackBekket"
	repo := "GitHelper"

	issue_title := "test"
	content := "Hey, this is test issue. Explain me how main package works?"

	client, _, err := GetClientByRepoOwner(repoOwner)
	if err != nil {
		log.Print(err)
		t.Fatal(err)
	}
	
	//lastIssueId := ?		// we need to get last id? is it autoincrement?

	issue,err :=github.CreateIssue(*client,repoOwner,repo,lastIssueId,issue_title,content)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf(issue)

}


func Test_commentIssue(t *testing.T,id int64) {

	err := godotenv.Load()
	if err != nil {
		t.Fatal(err)
	}
	repoOwner := "JackBekket"
	repo := "GitHelper"

	content := "Hey, this is test comment."

	client, _, err := GetClientByRepoOwner(repoOwner)
	if err != nil {
		log.Print(err)
		t.Fatal(err)
	}

	issue,err :=github.CommentIssue(*client,repoOwner,repo,id,content)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf(issue)


}

