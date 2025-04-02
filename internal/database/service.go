package database

type Service struct {
	DBHandler *Handler
}

func NewAIService(dbHandler *Handler) (*Service, error) {
	service := Service{
		DBHandler: dbHandler,
	}
	err := service.CreateTables()
	return &service, err
}
