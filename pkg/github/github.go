package github

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v65/github"
	"github.com/rs/zerolog/log"
)

var (
	IssueStateClosed            = "closed"
	IssueClosedReasonCompleted  = "completed"
	IssueClosedReasonNotPlanned = "not_planned"
)

type Service struct {
	WhiteList []string

	appId  int64
	pkPath string
}

func NewGHService(appId int64, pkPath string, whitelist []string) *Service {
	return &Service{
		WhiteList: whitelist,

		appId:  appId,
		pkPath: pkPath,
	}
}

func CloseIssue(ctx context.Context, client *github.Client, owner, repo string, id int) (*github.Issue, error) {
	issue, _, err := client.Issues.Edit(ctx, owner, repo, id, &github.IssueRequest{
		State:       &IssueStateClosed,
		StateReason: &IssueClosedReasonCompleted,
	})

	return issue, err
}

func CreateIssue(ctx context.Context, client *github.Client, owner, repo string, issueTitle string, content string) (*github.Issue, error) {
	issue, _, err := client.Issues.Create(ctx, owner, repo, &github.IssueRequest{
		Title: &issueTitle,
		Body:  &content,
	})
	if err != nil {
		return nil, err
	}
	return issue, nil
}

func CommentIssue(ctx context.Context, client *github.Client, owner, repo string, id int, response string) (*github.IssueComment, error) {
	// Craft a reply message from the response from the 3rd service.
	replyMessage := response

	// Create a new comment on the issue using the GitHub API.
	issueComment, _, err := client.Issues.CreateComment(ctx, owner, repo, id, &github.IssueComment{
		// Configure the comment with the issue's ID and other necessary details.
		Body: &replyMessage,
	})

	return issueComment, err
}

func (s Service) GetClientByRepoOwner(owner string) (*github.Client, *github.Installation, error) {
	tr := http.DefaultTransport

	itr, err := ghinstallation.NewAppsTransportKeyFromFile(tr, s.appId, s.pkPath)
	if err != nil {
		return nil, nil, err
	}

	//create git client with app transport
	client := github.NewClient(
		&http.Client{
			Transport: itr,
			Timeout:   time.Second * 30,
		},
	)

	if client == nil {
		return nil, nil, fmt.Errorf("create git client for app")
	}

	installations, _, err := client.Apps.ListInstallations(context.Background(), &github.ListOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("list installations: %+v", err)
	}

	for _, installation := range installations {
		log.Info().Msgf("installation: %+v",
			*installation,
		)
		log.Info().Msgf("repository selection: %s",
			installation.GetRepositorySelection(),
		)
	}

	//capture our installationId for our app
	//we need this for the access token
	var installID int64
	for _, val := range installations {
		installID = val.GetID()

		user := val.GetAccount()
		username := user.GetLogin()
		targetType := val.GetTargetType()

		log.Debug().Msgf("installed by %s, target type: %s", username, targetType)

		if username != owner {
			continue
		}

		//If whitelist does not contain our names throw error
		if !slices.Contains(s.WhiteList, username) {
			return nil, nil, fmt.Errorf("user %s is not in the whitelist. Skipping creating client for it", username)
		}

		token, _, err := client.Apps.CreateInstallationToken(
			context.Background(),
			installID,
			&github.InstallationTokenOptions{})
		if err != nil {
			return nil, nil, fmt.Errorf("creating installation token: %+v", err)
		}

		apiClient := github.NewClient(nil).WithAuthToken(
			token.GetToken(),
		)
		if apiClient == nil {
			return nil, nil, fmt.Errorf("empty client")
		}

		return apiClient, val, nil
	}

	return nil, nil, fmt.Errorf("client not found for key: %s", owner)
}
