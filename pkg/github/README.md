# github

This package provides a service for interacting with GitHub's API. It allows users to create, comment on, and close issues, as well as retrieve a client for a given repository owner.

```
github/
├── github.go
```

## Service

The `Service` type is responsible for managing the interaction with the GitHub API. It has the following methods:

- `NewGHService`: This function creates a new instance of the `Service` type.
- `CloseIssue`: This method closes an issue on GitHub.
- `CreateIssue`: This method creates a new issue on GitHub.
- `CommentIssue`: This method adds a comment to an existing issue on GitHub.
- `GetClientByRepoOwner`: This method retrieves a GitHub client for a given repository owner.

## Constants

The package defines the following constants:

- `IssueStateClosed`: This constant represents the closed state of an issue.
- `IssueClosedReasonCompleted`: This constant represents the completed reason for closing an issue.
- `IssueClosedReasonNotPlanned`: This constant represents the not planned reason for closing an issue.

## Configuration

The package does not appear to have any configuration files or environment variables.

## Edge Cases

The package does not appear to have any specific edge cases for launching the application.

