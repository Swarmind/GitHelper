package database_test

import (
	"database/sql"
	"testing"

	"github.com/JackBekket/GitHelper/internal/database"
	_ "github.com/lib/pq"
	"github.com/rubenv/pgtest"
	"github.com/tmc/langchaingo/llms"
)

func expectMessage(aiService *database.Service, t *testing.T, issueId int64, repoName, model, message string) {
	content, err := aiService.GetHistory(issueId, repoName, model)
	if err != nil {
		t.Fatal(err)
	}
	if len(content) > 0 {
		msg := content[0]
		if msg.Role != llms.ChatMessageTypeAI {
			t.Fatal("unexpected message role")
		}
		if len(msg.Parts) > 0 {
			part, ok := msg.Parts[0].(llms.TextContent)
			if !ok {
				t.Fatal("failed to cast message part to the TextContent part type")
			}
			if part.Text != message {
				t.Fatal("unexpected message!")
			}
		} else {
			t.Fatal("expected non-empty message parts")
		}
	} else {
		t.Fatal("expected non-empty content")
	}
}

func Test_DB(t *testing.T) {
	pg, err := pgtest.Start()
	defer pg.Stop()

	aiService, err := database.NewAIService(
		&database.Handler{
			DB: pg.DB,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := aiService.UpdateHistory(1, "repo1", "model_test1", llms.MessageContent{
		Role: llms.ChatMessageTypeAI,
		Parts: []llms.ContentPart{
			llms.TextPart("test1"),
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := aiService.UpdateHistory(1, "repo2", "model_test2", llms.MessageContent{
		Role: llms.ChatMessageTypeAI,
		Parts: []llms.ContentPart{
			llms.TextPart("test2"),
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := aiService.UpdateHistory(2, "repo1", "model_test3", llms.MessageContent{
		Role: llms.ChatMessageTypeAI,
		Parts: []llms.ContentPart{
			llms.TextPart("test3"),
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := aiService.UpdateHistory(1, "repo1", "model_test4", llms.MessageContent{
		Role: llms.ChatMessageTypeAI,
		Parts: []llms.ContentPart{
			llms.TextPart("test4"),
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := aiService.DropHistory(1, "repo2", "model_test2"); err != nil {
		t.Fatal(err)
	}
	if err := aiService.DropHistory(1, "repo1", "model_test999"); err != nil {
		t.Fatal(err)
	}

	expectMessage(aiService, t, 1, "repo1", "model_test1", "test1")
	expectMessage(aiService, t, 2, "repo1", "model_test3", "test3")
	expectMessage(aiService, t, 1, "repo1", "model_test4", "test4")

	content, err := aiService.GetHistory(1, "repo2", "model_test2")
	if err != nil && err != sql.ErrNoRows {
		t.Fatal(err)
	}
	if len(content) > 0 {
		t.Fatal("unexpected content for that case")
	}
}
