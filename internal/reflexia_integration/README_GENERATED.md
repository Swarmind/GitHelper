# reflexia_integration

This package is designed to integrate with Reflexia, a tool for summarizing and analyzing code. It allows users to run and summarize packages in a project, either from a local directory or a GitHub repository.

The package relies on several external data sources, including environment variables, command line arguments, and a GitHub repository. Environment variables such as GH_USERNAME, GH_TOKEN, HELPER_URL, MODEL, API_TOKEN, EMBEDDINGS_AI_URL, EMBEDDINGS_AI_KEY, and EMBEDDINGS_DB_URL are used for configuration. Command line arguments can be a GitHub repository URL or a local project path.

The package's main function, `InitPackageRunner`, initializes a `PackageRunnerService` object. This object is responsible for running and summarizing packages in a project. The function first loads environment variables and command line arguments to determine the working directory and project configuration. Then, it creates a `SummarizeService` object to handle summarization tasks and an `EmbeddingsService` object to store and retrieve embeddings for code snippets. Finally, it returns a `PackageRunnerService` object with the necessary components to run and summarize packages.

The `loadEnv` function loads environment variables and checks if they are empty. The `processWorkingDirectory` function determines the working directory based on command line arguments or a GitHub repository URL. If necessary, it clones the repository and returns the working directory path.

## Project package structure:
- reflexia_integration.go
- internal/reflexia_integration/reflexia_integration.go

## Relations between code entities:
The `InitPackageRunner` function relies on the `loadEnv` and `processWorkingDirectory` functions to set up the necessary environment and working directory. The `PackageRunnerService` object then uses the `SummarizeService` and `EmbeddingsService` objects to perform summarization tasks.

## Unclear places:
It's unclear how the `EmbeddingsService` object is used to store and retrieve embeddings for code snippets.

## Dead code:
None found.

