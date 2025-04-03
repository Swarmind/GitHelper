package utils

import (
	"context"
	"fmt"

	"github.com/google/go-github/v62/github"
)

func CreateIssue(client *github.Client, owner string, repo string, id int64,issue_title string ,content string) (*github.Issue,error){
	ctx := context.Background()
	issueRequest := &github.IssueRequest{
		Title: &issue_title,
		Body: &content,
	}
	// issue, response, error
	issue,_, err :=client.Issues.Create(ctx,owner,repo,issueRequest)
	if err != nil {
		return nil, err
	}
	return issue, nil
}



func CommentIssue(client *github.Client, owner string, repo string, id int64, response string) {
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