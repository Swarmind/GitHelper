# Agent

This package contains tests for the semantic search agent. The tests cover various scenarios, including one-shot search, memory-based search, and long conversations. The tests use a generic LLM model, which is created using the createGenericLLM function. The createGenericLLM function loads environment variables for the AI URL and API token, and uses them to create an openai.LLM instance. The tests use the agent.OneShotRun function to execute the semantic search agent. The agent.OneShotRun function takes a query, a model, and an optional initial state as input. The initial state is used for memory-based search and long conversations. The tests assert or compare the results with the expected output.

## Project package structure:

- agent.go
- semantic_search_test.go
- superagent.go
- tools/tools.go
- pkg/agent/rag/semantic_search_test.go

## Code entities relations:

The package contains tests for the semantic search agent. The tests cover various scenarios, including one-shot search, memory-based search, and long conversations. The tests use a generic LLM model, which is created using the createGenericLLM function. The createGenericLLM function loads environment variables for the AI URL and API token, and uses them to create an openai.LLM instance. The tests use the agent.OneShotRun function to execute the semantic search agent. The agent.OneShotRun function takes a query, a model, and an optional initial state as input. The initial state is used for memory-based search and long conversations. The tests assert or compare the results with the expected output.

## Environment variables, flags, cmdline arguments, files and their paths that can be used for configuration:

- AI_URL
- API_TOKEN

## Edge cases of how application can be launched:

- The application can be launched by running the tests using the "go test" command.

## Unclear places:

- Should the AI_URL and API_TOKEN environment variables be global?

