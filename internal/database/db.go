package database

import (
	"github.com/tmc/langchaingo/llms"
)

func (s *Service) CreateTables() error {
	_, err := s.DBHandler.DB.Exec(`
		CREATE TABLE IF NOT EXISTS gh_sessions (
			id SERIAL PRIMARY KEY,
			repo_name TEXT NOT NULL,
			issue_id BIGINT NOT NULL,
			model TEXT NOT NULL
		)`)
	if err != nil {
		return err
	}

	_, err = s.DBHandler.DB.Exec(`
		CREATE TABLE IF NOT EXISTS gh_chat_messages (
			id BIGSERIAL PRIMARY KEY,
			gh_session INT NOT NULL REFERENCES gh_sessions(id),
			message_data BYTEA NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_gh_chat_session_timestamp ON gh_chat_messages (gh_session, created_at);
		`)
	return err
}

func (s *Service) GetHistory(
	issueId int64, repoName, model string,
) ([]llms.MessageContent, error) {

	messages := []llms.MessageContent{}
	rows, err := s.DBHandler.DB.Query(`
        SELECT m.message_data
        FROM gh_chat_messages m
        JOIN gh_sessions gs ON m.gh_session = gs.id
        WHERE issue_id = $1 AND
            repo_name = $2 AND
            model = $3
        ORDER BY m.created_at
        `, issueId, repoName, model)
	if err != nil {
		return messages, err
	}
	defer rows.Close()

	for rows.Next() {
		contentBytes := []byte{}
		err := rows.Scan(&contentBytes)
		if err != nil {
			return messages, err
		}
		content := llms.MessageContent{}
		err = content.UnmarshalJSON(contentBytes)
		if err != nil {
			return messages, err
		}
		messages = append(messages, content)
	}

	return messages, err
}

func (s *Service) DropHistory(issueId int64, repoName, model string) error {
	_, err := s.DBHandler.DB.Exec(`
        DELETE FROM gh_chat_messages
        WHERE gh_session IN (
            SELECT id
            FROM gh_sessions
            WHERE issue_id = $1 AND
                repo_name = $2 AND
                model = $3
        )
        `, issueId, repoName, model)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateHistory(
	issueId int64, repoName, model string, content llms.MessageContent,
) error {
	contentBytes, err := content.MarshalJSON()
	if err != nil {
		return err
	}
	_, err = s.DBHandler.DB.Exec(`
        WITH
		GHSessionCheck AS (
			SELECT id
			FROM gh_sessions
			WHERE issue_id = $1 AND
                repo_name = $2 AND
                model = $3
        ),
		InsertGHSession AS (
			INSERT INTO gh_sessions
			(issue_id, repo_name, model)
			SELECT $1, $2, $3
			WHERE NOT EXISTS (SELECT 1 FROM GHSessionCheck)
			RETURNING id
		),
		GHSessionID AS (
			SELECT id FROM GHSessionCheck
			UNION ALL
			SELECT id FROM InsertGHSession
		)
		INSERT INTO gh_chat_messages (gh_session, message_data)
		SELECT (SELECT id FROM GHSessionID), $4
        `, issueId, repoName, model, contentBytes)
	return err
}
