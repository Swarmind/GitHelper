# rag

This package provides a RAG (Retrieval Augmented Generation) system for generating responses to prompts based on a given set of documents. The system uses an AI endpoint, an API token, and a database to store and retrieve documents. The package includes several test functions that demonstrate how to use the RAG system with different options and configurations.

## File Structure

- rag.go
- rag_test.go
- pkg/rag/rag_test.go

## Environment Variables

- AI_ENDPOINT: URL of the AI endpoint
- API_TOKEN: API token for the AI endpoint
- DB_URL: URL of the database

## Input Sources

- test_prompts: An array of strings containing the prompts to be used in the tests.
- repo_names: An array of strings containing the names of the repositories to be used in the tests.

## Major Code Parts

1. Test_RagWithFilteres: This test function calls RAG with optional filters. It loads the environment variables, initializes the AI, API token, and database URL, and then calls the RAG.RagWithOptions function with the test prompt, AI, API token, and the vector store.

2. Test_RagReflexia: This test function calls RAG with two types of documents (docs and code) and calls the 'stuffed' method of RAG. It follows a similar process as Test_RagWithFilteres, but it calls the RAG.RagReflexia function instead of RAG.RagWithOptions.

3. Test_StuffRag: This test function tests the 'stuffed' method of RAG. It loads the environment variables, initializes the AI, API token, and database URL, and then calls the RAG.StuffedQA_Rag function with the test prompt, AI, API token, and the vector store.

4. Test_RefinedQA_RAG: This test function calls the Refined QA method of RAG. It follows a similar process as Test_StuffRag, but it calls the RAG.StuffedQA_Rag function instead of RAG.RagWithOptions.

5. getCollection: This function takes the AI endpoint URL, API token, database URL, and namespace as input and returns a vector store. It uses the embeddings library to get the vector store with the specified options.

