package reflexia_integration

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	store "github.com/JackBekket/reflexia/pkg"
	runner "github.com/JackBekket/reflexia/pkg/package_runner"
	"github.com/JackBekket/reflexia/pkg/project"
	"github.com/JackBekket/reflexia/pkg/summarize"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/tmc/langchaingo/llms"
)

func InitPackageRunner(ghLink string) (runner.PackageRunnerService, error) {
	ghUsername := loadEnv("GH_USERNAME")
	ghToken := loadEnv("GH_TOKEN")

	workdir, err := processWorkingDirectory(
		ghLink, ghUsername, ghToken)
	if err != nil {
		log.Fatalf("processWorkingDirectory(...) error: %v", err)
	}

	projectConfigVariants, err := project.GetProjectConfig(workdir, "", false)
	if err != nil {
		log.Fatal(err)
	}

	var projectConfig *project.ProjectConfig
	switch len(projectConfigVariants) {
	case 0:
		log.Fatal("no languages detected")
	case 1:
		for _, pc := range projectConfigVariants {
			projectConfig = pc
		}
	default:
		projectConfig = projectConfigVariants["go.toml"]
	}

	pkgFiles, err := projectConfig.BuildPackageFiles()
	if err != nil {
		log.Fatalf("projectConfig.BuildPackageFiles() error: %v", err)
	}
	summarizeService := &summarize.SummarizeService{
		HelperURL: loadEnv("HELPER_URL"),
		Model:     loadEnv("MODEL"),
		ApiToken:  loadEnv("API_TOKEN"),
		Network:   "local",
		LlmOptions: []llms.CallOption{
			llms.WithStopWords(
				projectConfig.StopWords,
			),
			llms.WithRepetitionPenalty(0.7),
		},
		// Tests only
		IgnoreCache:    false,
		OverwriteCache: false,
		CachePath:      ".githelper_cache",
	}

	var embeddingsService *store.EmbeddingsService
	projectName := filepath.Base(projectConfig.RootPath)
	vectorStore, err := store.NewVectorStoreWithPreDelete(
		os.Getenv("EMBEDDINGS_AI_URL"),
		os.Getenv("EMBEDDINGS_AI_KEY"),
		os.Getenv("EMBEDDINGS_DB_URL"),
		projectName,
	)
	if err != nil {
		log.Fatalf("store.NewVectorStoreWithPreDelete(...) error: %v", err)
	}

	embeddingsService = &store.EmbeddingsService{
		Store: vectorStore,
	}
	fmt.Printf("Initialized vector store with %s as project name\n", projectName)

	pkgRunner := runner.PackageRunnerService{
		PkgFiles:          pkgFiles,
		ProjectConfig:     projectConfig,
		SummarizeService:  summarizeService,
		EmbeddingsService: embeddingsService,
		ExactPackages:     nil,
		OverwriteReadme:   false,
		WithFileSummary:   false,
	}
	return pkgRunner, nil
}

func loadEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("empty environment key %s", key)
	}
	return value
}

func processWorkingDirectory(githubLink, githubUsername, githubToken string) (string, error) {
	workdir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if githubLink != "" {
		u, err := url.ParseRequestURI(githubLink)
		if err != nil {
			return "", err
		}

		sPath := strings.Split(strings.TrimPrefix(u.Path, "/"), "/")
		if len(sPath) != 2 {
			return "", errors.New("github repository url does not have two path elements")
		}

		tempDirEl := []string{workdir, "temp"}
		tempDirEl = append(tempDirEl, sPath...)
		tempDir := filepath.Join(tempDirEl...)

		workdir = tempDir

		if _, err := os.Stat(workdir); err != nil {
			if os.IsNotExist(err) {
				if err := os.MkdirAll(workdir, os.FileMode(0755)); err != nil {
					return "", err
				}

				cloneOptions := git.CloneOptions{
					URL:               githubLink,
					RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
					Depth:             1,
				}
				if githubUsername != "" && githubToken != "" {
					cloneOptions.Auth = &http.BasicAuth{
						Username: githubUsername,
						Password: githubToken,
					}
				}

				if _, err := git.PlainClone(workdir, false, &cloneOptions); err != nil {
					if err := os.RemoveAll(workdir); err != nil {
						return "", err
					}
					return "", err
				}

			} else {
				return "", err
			}
		}
	} else if len(flag.Args()) > 0 {
		workdir = flag.Arg(0)
		if _, err := os.Stat(workdir); err != nil {
			return "", err
		}
	}

	return workdir, nil
}
