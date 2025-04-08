# legacy/rag.go  
package: rag  
imports: context, fmt, github.com/JackBekket/hellper/lib/embeddings, github.com/tmc/langchaingo/chains, github.com/tmc/langchaingo/llms/openai, github.com/tmc/langchaingo/vectorstores  
  
func RagReflexia(question string, ai_url string, api_token string, numOfResults int, store vectorstores.VectorStore) (result string, err error):  
	- creates an embeddings client using the specified API and embedding model  
	- retrieves documentation and code metadata based on the question  
	- combines the documentation and code metadata into a single input for the LLM  
	- uses the stuffed QA chain to generate an answer to the question  
	- returns the answer and any errors encountered  
func RagWithOptions(question string, ai_url string, api_token string, numOfResults int, store vectorstores.VectorStore, option ...vectorstores.Option) (result string, err error):  
	- creates an embeddings client using the specified API and embedding model  
	- retrieves relevant documents from the vector store  
	- uses the retrieval QA chain to generate an answer to the question  
	- returns the answer and any errors encountered  
func StuffedQA_Rag(question string, ai_url string, api_token string, numOfResults int, store vectorstores.VectorStore) (result string, err error):  
	- creates an embeddings client using the specified API and embedding model  
	- retrieves relevant documents from the vector store  
	- uses the stuffed QA chain to generate an answer to the question  
	- returns the answer and any errors encountered  
func RefinedQA_RAG(question string, ai_url string, api_token string, numOfResults int, store vectorstores.VectorStore) (result string, err error):  
	- creates an embeddings client using the specified API and embedding model  
	- retrieves relevant documents from the vector store  
	- uses the refine documents chain for question answering to generate an answer to the question  
	- returns the answer and any errors encountered  
  
  
