# internal/reflexia_integration/reflexia_integration.go  
package: reflexia_integration  
imports: errors, flag, fmt, log, net/url, os, path/filepath, strings, store, runner, project, summarize, git, github.com/go-git/go-git/v5/plumbing/transport/http, github.com/tmc/langchaingo/llms  
  
func InitPackageRunner(ghLink string):  
	- loads environment variables GH_USERNAME and GH_TOKEN  
	- calls processWorkingDirectory(...) to get the working directory  
	- calls GetProjectConfig(...) to get the project configuration  
	- checks the number of project configuration variants and sets the projectConfig accordingly  
	- calls BuildPackageFiles() to get the package files  
	- creates a SummarizeService object with the loaded environment variables and LlmOptions  
	- creates a vector store with the loaded environment variables  
	- creates an EmbeddingsService object with the vector store  
	- creates a PackageRunnerService object with the pkgFiles, ProjectConfig, SummarizeService, EmbeddingsService, ExactPackages, OverwriteReadme, and WithFileSummary  
	- returns the PackageRunnerService object  
func loadEnv(key string):  
	- loads the environment variable with the given key  
	- if the environment variable is empty, it logs a fatal error  
	- returns the environment variable value  
func processWorkingDirectory(githubLink, githubUsername, githubToken string):  
	- gets the current working directory  
	- if githubLink is not empty, it parses the URL and creates a temporary directory to clone the repository  
	- if flag.Args() is not empty, it sets the working directory to the first argument  
	- returns the working directory  
  
  
