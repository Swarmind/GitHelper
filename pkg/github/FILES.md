# pkg/github/github.go  
package: github  
imports: context, fmt, net/http, slices, time, github.com/bradleyfalzon/ghinstallation/v2, github.com/google/go-github/v65/github, github.com/rs/zerolog/log  
  
var: IssueStateClosed, IssueClosedReasonCompleted, IssueClosedReasonNotPlanned  
  
type: Service  
func: NewGHService  
func: CloseIssue  
func: CreateIssue  
func: CommentIssue  
func: GetClientByRepoOwner  
  
  
