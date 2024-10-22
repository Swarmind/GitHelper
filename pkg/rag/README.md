## Package: rag

### Imports:

- context
- fmt
- github.com/JackBekket/hellper/lib/embeddings
- github.com/tmc/langchaingo/chains
- github.com/tmc/langchaingo/llms/openai
- github.com/tmc/langchaingo/vectorstores

### External Data, Input Sources:

- vectorstores.VectorStore: This is an interface for interacting with a vector store, which is used to store and retrieve embeddings of documents.

### Code Summary:

#### RagReflexia:

This function takes a question, API URL, API token, number of results, and a vector store as input. It first creates an embeddings client using the specified API and embedding model. Then, it performs semantic search on the vector store to retrieve relevant documents and code snippets based on the question. The retrieved documents and code snippets are combined into a single prompt for the LLM. Finally, it calls a stuffed QA chain to generate an answer to the question using the combined prompt.

#### RagWithOptions:

This function takes a question, API URL, API token, number of results, a vector store, and optional vectorstore.Option as input. It creates an embeddings client and then calls a retrieval QA chain to generate an answer to the question using the specified vector store and options.

#### StuffedQA_Rag:

This function takes a question, API URL, API token, number of results, and a vector store as input. It creates an embeddings client and performs semantic search on the vector store to retrieve relevant documents. Then, it calls a stuffed QA chain to generate an answer to the question using the retrieved documents.

#### RefinedQA_RAG:

This function takes a question, API URL, API token, number of results, and a vector store as input. It creates an embeddings client and performs semantic search on the vector store to retrieve relevant documents. Then, it calls a refine QA chain to generate an answer to the question using the retrieved documents.



