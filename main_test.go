package main

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/JackBekket/GitHelper/pkg/github"
	githubAPI "github.com/JackBekket/GitHelper/pkg/github"
	"github.com/JackBekket/GitHelper/pkg/rag/agent"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func TestAgent(t *testing.T) {
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

func TestGithubAPI(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatal(err)
	}

	appIdStr := os.Getenv("APP_ID")
	appId, err := strconv.ParseInt(appIdStr, 10, 64)
	if err != nil {
		log.Fatal().Err(err).Msg("parsing APP_ID env variable")
	}
	pkPath := os.Getenv("PRIVATE_KEY_NAME")
	GHService := githubAPI.NewGHService(
		appId, pkPath,
		strings.Split(os.Getenv("OWNER_WHITELIST"), ","),
	)

	repoOwner := "JackBekket"

	client, _, err := GHService.GetClientByRepoOwner(repoOwner)
	if err != nil {
		log.Print(err)
		t.Fatal(err)
	}

	repo := "GitHelper"

	issue_title := "test"
	content := "Hey, this is test issue. Explain me how main package works?"
	issue, err := github.CreateIssue(t.Context(), client, repoOwner, repo, issue_title, content)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Issue: %+v", issue)

	content = "this is a test comment"
	comment, err := github.CommentIssue(t.Context(), client, repoOwner, repo, issue.GetNumber(), content)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Comment: %+v", comment)

	issue, err = github.CloseIssue(t.Context(), client, repoOwner, repo, issue.GetNumber())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Closed issue: %+v", issue)
}
