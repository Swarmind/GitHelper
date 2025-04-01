package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tmc/langchaingo/llms"
)

// memory
type AiSession struct {
	AIToken      string
	Model     string
	DialogThread ChatSessionGraph
}



// langgraph doesn't work with same types as langchain, so we have to improvise here.
type ChatSessionGraph struct {
	ConversationBuffer []llms.MessageContent
	//DialogThread string

}


// check if collection is exist
func (s *Service) CheckCollection(repo_name string) (bool) {
    //var repoNameFromCollection sql.NullString
    err := s.DBHandler.DB.QueryRow(`
        SELECT name
        FROM langchain_pg_collection
        WHERE name = $1`, repo_name)
    if err != nil {
        return false
    } else {
		return true
	}
}



// check if session exists
func (s *Service) CheckSession(userID int64) bool {
	_, err := s.GetAISession(userID)
	if err != nil {
		log.Warn().Err(err).Msg("error encountered")
		return false
	} else {
		return true
	}
}

func (s *Service) CreateGHSession(issueID int8, repo_name string) error {
	check := s.CheckCollection(repo_name)
	if check == true{
    res, err := s.DBHandler.DB.Exec("INSERT INTO gh_session (issue_id, repo_name) VALUES ($1, $2)", issueID, repo_name)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	log.Printf("Query executed successfully. Rows affected: %d", rowsAffected)
	return nil
	} else {
    e:= fmt.Sprintf(" collection does not exist")
	return fmt.Errorf(e)
	}
}



func (s *Service) CreateGHChatMessage(id int8, messageData string, createdAt time.Time) error {
    _, err := s.DBHandler.DB.Exec("INSERT INTO gh_chat_messages (id, message_data, created_at) VALUES (?, ?, ?, ?)", id, messageData, createdAt)
    if err != nil {
        log.Printf("Error executing query: %v", err)
        return err
    }
    return nil
}





func (s *Service) CreateAISession(userID int64, model string, providerID int64) error {
	log.Info().Int64("userID", userID).Str("model", model).Int64("providerID", providerID).Msg("CreateLSession called")

	res, err := s.DBHandler.DB.Exec(`
		INSERT INTO gh_sessions (tg_user_id, model, endpoint)
		VALUES ($1, $2, $3)
		ON CONFLICT(tg_user_id) DO UPDATE SET
		model = $2
		RETURNING tg_user_id
	`, userID, model, providerID)

	if err != nil {
		log.Printf("Error executing query: %v", err)
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	log.Printf("Query executed successfully. Rows affected: %d", rowsAffected)

	return nil
}

func (s *Service) DeleteAISession(userID int64) error {
	_, err := s.DBHandler.DB.Exec(`
    DELETE FROM ai_sessions
    WHERE tg_user_id = $1
    `, userID)
	return err
}

func (s *Service) GetAISession(userID int64) (AiSession, error) {
	var (
		model, endpointName, endpointURL sql.NullString
		endpointId, endpointAuthMethod   sql.NullInt64
		session                          AiSession
	)

	err := s.DBHandler.DB.QueryRow(`
		SELECT lt.model,
		       rt.id,
		       rt.name,
		       rt.url,
		       rt.auth_method
		FROM ai_sessions lt
		LEFT JOIN endpoints rt ON lt.endpoint = rt.id
		WHERE tg_user_id = $1`, userID).Scan(
		&model, &endpointId, &endpointName,
		&endpointURL, &endpointAuthMethod,
	)
	if err != nil {
		return session, err
	}

	if model.Valid {
		session.Model = model.String
	}

	/*
	if endpointId.Valid && endpointName.Valid && endpointURL.Valid && endpointAuthMethod.Valid {
		session.AIProvider = &Endpoint{
			ID:         endpointId.Int64,
			Name:       endpointName.String,
			BaseURL:    endpointURL.String,
			AuthMethod: endpointAuthMethod.Int64,
		}
	}
	*/

	return session, nil
}




func (s *Service) UpdateModelInAISession(userID int64, model *string) error {
	var modelValue interface{}
	if model != nil {
		modelValue = *model
	} else {
		modelValue = nil
	}
	_, err := s.DBHandler.DB.Exec(`INSERT INTO ai_sessions
			(tg_user_id, model)
		VALUES
			($1, $2)
		ON CONFLICT(tg_user_id) DO UPDATE SET
			model = $2
		`, userID, modelValue)
	return err
}

// also create new session?
func (s *Service) UpdateEndpointInAISession(userID int64, endpointId *int64) error {
	var endpointIdValue interface{}
	if endpointId != nil {
		endpointIdValue = *endpointId
	} else {
		endpointIdValue = nil
	}
	_, err := s.DBHandler.DB.Exec(`INSERT INTO ai_sessions
			(tg_user_id, endpoint)
		VALUES
			($1, $2)
		ON CONFLICT(tg_user_id) DO UPDATE SET
			endpoint = $2
		`, userID, endpointIdValue)
	return err
}
