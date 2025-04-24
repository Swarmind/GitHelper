# tools

This package provides a set of tools for semantic search and web search.

## Project package structure:
- pkg/agent/rag/tools/tools.go

## Code summary:
The `tools` package provides a set of tools for semantic search and web search. The `GetTools` function returns a list of available tools, including a semantic search tool and a web search tool. The `SemanticSearchTool` struct implements the `Tool` interface and is responsible for executing the semantic search. The `Execute` method of the `SemanticSearchTool` struct takes a context and a list of messages as input and returns a list of messages and an error. The method extracts the search query and collection name from the input message, retrieves the vector store, performs the semantic search, and returns the search results.

## Environment variables:
- AI_URL
- API_TOKEN
- DB_LINK

## Edge cases:
- The application can be launched by running the `main` function in the `tools.go` file.

## Unclear places:
- It is unclear how the `SemanticSearchTool` struct retrieves the vector store.
- It is unclear how the `SemanticSearchTool` struct performs the semantic search.

