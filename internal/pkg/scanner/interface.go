package scanner

import "main/internal/models"

type Repository interface {
	GetRequest(id int) (models.Request, error)
	GetResponse(id int) (models.Response, error)
	GetAllRequests() ([]models.Request, error)
}
