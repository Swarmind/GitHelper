## Package: main

This package handles webhook events from GitHub and interacts with a database to store and retrieve information related to GitHub repositories.

### Imports:

- context
- encoding/json
- fmt
- io
- net/http
- os
- slices
- strconv
- strings
- github.com/JackBekket/GitHelper/internal/database
- github.com/JackBekket/GitHelper/internal/reflexia_integration
- github.com/JackBekket/GitHelper/pkg/agent/rag
- github.com/JackBekket/GitHelper/pkg/github
- github.com/JackBekket/hellper/lib/embeddings
- github.com/rs/zerolog/log
- github.com/google/go-github/v65/github
- github.com/joho/godotenv
- github.com/tmc/langchaingo/vectorstores

### External Data and Input Sources:

- Environment variables: AIBaseURL, AIToken, DBURL, APP_ID, PRIVATE_KEY_NAME, MODEL
- GitHub webhook events: installation_repositories, issues, issue_comment, push

### TODOs:

- Delete the obsolete generateResponse function when the previous genResponse function is tested.

### Code Summary:

1. **Initialization:**
   - Loads environment variables and sets up database connection.
   - Creates a GitHub API service instance.

2. **Webhook Handler:**
   - Handles incoming webhook events from GitHub.
   - Processes events based on their type:
     - Installation: Logs the installation of the app in a repository.
     - Issue opened: Creates a response to the issue using the RAG agent.
     - Issue closed: Drops the history of the issue from the database.
     - Issue comment: Generates a response to the comment using the RAG agent.
     - Push: Checks if the push is to the master or main branch and runs the Reflexia package runner if it is.

3. **Response Generation:**
   - Creates a response to a GitHub issue using the RAG agent.
   - Updates the history of the issue in the database.

4. **Continuation of Thread:**
   - Generates a response to a comment on a GitHub issue, continuing the existing conversation.
   - Updates the history of the issue in the database.

5. **Collection Retrieval:**
   - Retrieves a collection of documents from the database for a given namespace.

### Summary:

This package is responsible for handling GitHub webhook events and interacting with a database to store and retrieve information related to GitHub repositories. It uses the RAG agent to generate responses to issues and comments, and it updates the history of the conversation in the database.

main_test.go
Package: main

Imports:
- "os"
- "strconv"
- "strings"
- "testing"
- "github.com/JackBekket/GitHelper/pkg/agent/rag"
- "github.com/JackBekket/GitHelper/pkg/github"
- "github.com/JackBekket/GitHelper/pkg/github"
- "github.com/joho/godotenv"
- "github.com/rs/zerolog/log"

External data, input sources:
- .env file for environment variables

TODOs:
- None

Summary:
- The code defines two test functions, TestAgent and TestGithubAPI, which test the functionality of the agent and GitHub API components, respectively.
- TestAgent tests the agent's ability to retrieve and process prompts from a vectorstore, using environment variables to configure the agent's settings.
- TestGithubAPI tests the GitHub API's ability to create, comment on, and close issues, using environment variables to configure the API's settings.
- Both test functions use the godotenv package to load environment variables from a .env file.
- The code also imports and uses the agent, GitHub, and godotenv packages to perform the tests.

