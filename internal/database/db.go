package database

import (
	"github.com/rs/zerolog/log"
	"github.com/tmc/langchaingo/llms"
)

func (s *Service) GetHistory(
    issueId, repoId int64, model string,
) ([]llms.MessageContent, error) {

    messages := []llms.MessageContent{}
    rows, err := s.DBHandler.DB.Query(`
        SELECT m.message_data
        FROM gh_chat_messages m
        JOIN gh_sessions cs ON m.chat_session = cs.issue_id
        WHERE issue_id = $1 AND
            repo_id = $2 AND
            model = $3
        ORDER BY m.created_at
        `, issueId, repoId, model)
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

func (s *Service) DropHistory(issueId, repoId int64, model string) error {
    res, err := s.DBHandler.DB.Exec(`
        DELETE FROM gh_chat_messages
        WHERE chat_session IN (
            SELECT issue_id
            FROM gh_sessions
            WHERE issue_id = $1 AND
                repo_id = $2 AND
                model = $3
        )
        `, issueId, repoId, model)
    if err != nil {
        log.Warn().Err(err).Msg("failed to execute delete query")
        return err
    }
    rowsAffected, _ := res.RowsAffected()
    log.Info().Int64("rowsAffected", rowsAffected).Msg("Query executed successfully")
    return nil
}

func (s *Service) UpdateHistory(
    issueId, repoId int64, model string, content llms.MessageContent,
) error {
    contentBytes, err := content.MarshalJSON()
    if err != nil {
        return err
    }
    _, err = s.DBHandler.DB.Exec(`
        INSERT INTO gh_chat_messages (chat_session, message_data)
        SELECT $1, $2
        WHERE NOT EXISTS (
            SELECT 1
            FROM gh_chat_messages
            WHERE chat_session = $1 AND message_data = $2
        )
        `, issueId, contentBytes)
    return err
}