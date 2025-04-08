# internal/database/db_test.go  
package: database_test  
imports: database/sql, testing, github.com/JackBekket/GitHelper/internal/database, github.com/lib/pq, github.com/rubenv/pgtest, github.com/tmc/langchaingo/llms  
  
func expectMessage(aiService *database.Service, t *testing.T, issueId int64, repoName, model, message string):  
	- retrieves the history for the given issueId, repoName, and model  
	- checks if the message role is AI  
	- checks if the message part is of type TextContent  
	- checks if the message text matches the expected message  
	- if any of the checks fail, it prints a fatal error  
  
func Test_DB(t *testing.T):  
	- starts a PostgreSQL test instance  
	- creates a new AI service instance  
	- updates the history for several issueId, repoName, and model combinations  
	- drops the history for specific issueId, repoName, and model combinations  
	- calls expectMessage to verify the expected messages  
	- retrieves the history for a specific issueId, repoName, and model combination  
	- checks if the content is empty or not  
	- if any of the checks fail, it prints a fatal error  
  
  
