# pkg/rag/agent/semantic_search_test.go  
package: agent_test  
imports: flag, fmt, log, os, testing, github.com/JackBekket/GitHelper/pkg/rag/agent, github.com/joho/godotenv, github.com/tmc/langchaingo/llms, github.com/tmc/langchaingo/llms/openai  
  
func Test_Search(t *testing.T):  
	- prints "testing one-shot search"  
	- creates a generic LLM model  
	- runs the OneShotRun function with the model and a query  
	- prints the result  
func TestMemory(t *testing.T):  
	- prints "testing with memory"  
	- sets up an initial state for the conversation  
	- creates a generic LLM model  
	- runs the OneShotRun function with the model, initial state, and a query  
	- prints the result  
func TestLongConversation(t *testing.T):  
	- prints "testing with long conversation"  
	- sets the test timeout to 3 minutes  
	- sets up an initial state for the conversation  
	- creates a generic LLM model  
	- runs the OneShotRun function with the model, initial state, and a query  
	- prints the result  
	- appends the result to the initial state  
	- runs the OneShotRun function again with the updated initial state and a new query  
	- prints the result  
func Test5Conversation(t *testing.T):  
	- prints "testing with 5 turns conversation"  
	- sets the test timeout to 6 minutes  
	- sets up an initial state for the conversation  
	- creates a generic LLM model  
	- runs the OneShotRun function with the model, initial state, and a query  
	- prints the result  
	- appends the result to the initial state  
	- runs the OneShotRun function again with the updated initial state and a new query  
	- prints the result  
	- runs the OneShotRun function two more times with updated initial states and new queries  
	- prints the results  
func createGenericLLM():  
	- sets the model name to "big-tiger-gemma-27b-v1"  
	- loads environment variables  
	- sets the base URL and API token  
	- creates a new openai LLM model with the specified parameters  
	- returns the model  
  
  
