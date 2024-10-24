package reflexia_integration

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	store "github.com/JackBekket/reflexia/pkg"
	runner "github.com/JackBekket/reflexia/pkg/package_runner"
	"github.com/JackBekket/reflexia/pkg/project"
	"github.com/JackBekket/reflexia/pkg/summarize"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms"
)

type Config struct {
	GithubLink                    *string
	GithubUsername                *string
	GithubToken                   *string
	WithConfigFile                *string
	ExactPackages                 *string
	LightCheck                    bool
	WithFileSummary               bool
	UseEmbeddings                 bool
	OverwriteReadme               bool
	OverwriteCache                bool
	EmbeddingsAIURL               *string
	EmbeddingsAIAPIKey            *string
	EmbeddingsDBURL               *string
	EmbeddingsSimSearchTestPrompt *string
	CachePath                     *string
}

func InitPackageRunner(ghLink, ghUsername string) runner.PackageRunnerService {
	initialConfig, err := initConfig(ghLink, ghUsername)
	if err != nil {
		log.Fatalf("initConfig(...) error: %v", err)
	}

	workdir, err := processWorkingDirectory(
		*initialConfig.GithubLink, *initialConfig.GithubUsername, *initialConfig.GithubToken)
	if err != nil {
		log.Fatalf("processWorkingDirectory(...) error: %v", err)
	}

	projectConfigVariants, err := project.GetProjectConfig(workdir, *initialConfig.WithConfigFile, initialConfig.LightCheck)
	if err != nil {
		log.Fatal(err)
	}
	projectConfig, err := chooseProjectConfig(projectConfigVariants)
	if err != nil {
		log.Fatalf("chooseProjectConfig(...) error: %v", err)
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
		OverwriteCache: initialConfig.OverwriteCache,
		CachePath:      *initialConfig.CachePath,
	}

	var embeddingsService *store.EmbeddingsService
	if initialConfig.UseEmbeddings {
		projectName := filepath.Base(projectConfig.RootPath)
		vectorStore, err := store.NewVectorStoreWithPreDelete(
			*initialConfig.EmbeddingsAIURL,
			*initialConfig.EmbeddingsAIAPIKey,
			*initialConfig.EmbeddingsDBURL,
			projectName,
		)
		if err != nil {
			log.Fatalf("store.NewVectorStoreWithPreDelete(...) error: %v", err)
		}

		embeddingsService = &store.EmbeddingsService{
			Store: vectorStore,
		}
		fmt.Printf("Initialized vector store with %s as project name\n", projectName)
	}

	pkgRunner := runner.PackageRunnerService{
		PkgFiles:          pkgFiles,
		ProjectConfig:     projectConfig,
		SummarizeService:  summarizeService,
		EmbeddingsService: embeddingsService,
		ExactPackages:     initialConfig.ExactPackages,
		OverwriteReadme:   initialConfig.OverwriteReadme,
		WithFileSummary:   initialConfig.WithFileSummary,
	}
	return pkgRunner
}

// TODO: I fear no man, but this thing scares me...
func chooseProjectConfig(projectConfigVariants map[string]*project.ProjectConfig) (*project.ProjectConfig, error) {
	switch len(projectConfigVariants) {
	case 0:
		return nil, errors.New(
			"failed to detect project language, available languages: go, python, typescript",
		)
	case 1:
		for _, pc := range projectConfigVariants {
			return pc, nil
		}
	default:
		var filenames []string
		for filename := range projectConfigVariants {
			filenames = append(filenames, filename)
		}
		fmt.Println("Multiple project config matches found!")
		for i, filename := range filenames {
			fmt.Printf("%d. %v\n", i+1, filename)
		}
		fmt.Print("Enter the number or filename: ")
		for {
			var input string
			if _, err := fmt.Scanln(&input); err != nil {
				log.Fatalf("fmt.Scanln(...) error: %v", err)
			}
			if index, err := strconv.Atoi(input); err == nil && index > 0 && index <= len(filenames) {
				return projectConfigVariants[filenames[index-1]], nil
			} else {
				for filename, config := range projectConfigVariants {
					if filename == input || strings.TrimSuffix(filename, ".toml") == input {
						return config, nil
					}
				}
			}
		}
	}
	panic("unreachable")
}

func loadEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("empty environment key %s", key)
	}
	return value
}

func initConfig(ghLink, ghUsername string) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println(err)
	}

	config := Config{}

	config.GithubLink = &ghLink
	config.GithubUsername = &ghUsername

	githubToken := os.Getenv("GH_TOKEN")
	config.GithubToken = &githubToken
	flag.StringVar(config.GithubToken, "t", *config.GithubToken, "github token for ssh auth")
	configFile := "go.toml"
	config.WithConfigFile = &configFile

	cachePath := ".reflexia_cache"
	config.CachePath = &cachePath

	embAIURL := os.Getenv("EMBEDDINGS_AI_URL")
	config.EmbeddingsAIURL = &embAIURL
	flag.StringVar(config.EmbeddingsAIURL, "eu", *config.EmbeddingsAIURL, "Embeddings AI URL")

	embAIAPIKey := os.Getenv("EMBEDDINGS_AI_KEY")
	config.EmbeddingsAIAPIKey = &embAIAPIKey
	flag.StringVar(config.EmbeddingsAIAPIKey, "ea", *config.EmbeddingsAIAPIKey, "Embeddings AI API Key")

	embDBURL := os.Getenv("EMBEDDINGS_DB_URL")
	config.EmbeddingsDBURL = &embDBURL
	flag.StringVar(config.EmbeddingsDBURL, "ed", *config.EmbeddingsDBURL, "Embeddings pgxpool DB connect URL")

	embSimSearchTestPrompt := os.Getenv("EMBEDDINGS_SIM_SEARCH_TEST_PROMPT")
	config.EmbeddingsSimSearchTestPrompt = &embSimSearchTestPrompt
	flag.StringVar(config.EmbeddingsSimSearchTestPrompt, "et", *config.EmbeddingsSimSearchTestPrompt, "Embeddings similarity search validation test prompt")

	config.ExactPackages = flag.String("p", "", "exact package names, ',' delimited")
	config.LightCheck = false
	config.WithFileSummary = false
	config.OverwriteReadme = false
	config.OverwriteCache = false
	config.UseEmbeddings = true
	//TODO: what to set here?
	flag.BoolFunc("c",
		"do not check project root folder specific files such as go.mod or package.json",
		func(_ string) error {
			config.LightCheck = true
			return nil
		})
	flag.BoolFunc("f",
		"Save individual file summary intermediate result to the FILES.md",
		func(_ string) error {
			config.WithFileSummary = true
			return nil
		})
	flag.BoolFunc("r",
		"overwrite README.md instead of README_GENERATED.md creation/overwrite",
		func(_ string) error {
			config.OverwriteReadme = true
			return nil
		})

	flag.Parse()

	return &config, nil
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
