package serviceRepo

import "db-forum/internal/app/models"

type ServiceRepository interface {
	GetSystemStatus() models.Service
	ClearSystem() error
}
