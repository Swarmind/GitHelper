package database

import (
	"errors"
)

type Service struct {
	DBHandler *Handler
}

var ErrHandlerNotFound = errors.New("handler for that user id is not found")
var ErrHandlerCast = errors.New("failed to cast LLM handler")

func NewAIService(dbHandler *Handler) (*Service, error) {
	service := Service{
		DBHandler: dbHandler,
	}
	err := service.CreateTables()
	return &service, err
}
