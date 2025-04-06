# pkg/rag/rag_test.go  
## rag_test  
  
### Imports  
  
```  
fmt  
log  
os  
testing  
RAG "github.com/JackBekket/GitHelper/pkg/rag"  
embeddings "github.com/JackBekket/hellper/lib/embeddings"  
godotenv "github.com/joho/godotenv"  
vectorstores "github.com/tmc/langchaingo/vectorstores"  
```  
  
### External Data, Input Sources  
  
The code uses environment variables to get the following information:  
  
- AI_ENDPOINT: URL of the AI endpoint  
- API_TOKEN: API token for the AI endpoint  
- DB_URL: URL of the database  
  
The code also uses the following input sources:  
  
- test_prompts: An array of strings containing the prompts to be used in the tests.  
- repo_names: An array of strings containing the names of the repositories to be used in the tests.  
  
### Major Code Parts  
  
#### Test_RagWithFilteres  
  
This test function calls RAG with optional filters. It first loads the environment variables and then initializes the AI, API token, and database URL. It then defines the test prompts and repository names. The function then calls the getCollection function to get the vector store for the specified repository. Finally, it calls the RAG.RagWithOptions function with the test prompt, AI, API token, and the vector store.  
  
#### Test_RagReflexia  
  
This test function calls RAG with two types of documents (docs and code) and calls the 'stuffed' method of RAG. It follows a similar process as Test_RagWithFilteres, but it calls the RAG.RagReflexia function instead of RAG.RagWithOptions.  
  
#### Test_StuffRag  
  
This test function tests the 'stuffed' method of RAG. It loads the environment variables, initializes the AI, API token, and database URL, and then defines the test prompts and repository names. It calls the getCollection function to get the vector store for the specified repository and then calls the RAG.StuffedQA_Rag function with the test prompt, AI, API token, and the vector store.  
  
#### Test_RefinedQA_RAG  
  
This test function calls the Refined QA method of RAG. It follows a similar process as Test_StuffRag, but it calls the RAG.StuffedQA_Rag function instead of RAG.RagWithOptions.  
  
#### getCollection  
  
This function takes the AI endpoint URL, API token, database URL, and namespace as input and returns a vector store. It uses the embeddings library to get the vector store with the specified options.  
  
  
  
