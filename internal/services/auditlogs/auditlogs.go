package auditlogs

import (
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type AuditLogsService struct {
	Storage *storage.Storage
}

func NewAuditLogsService(storage *storage.Storage) *AuditLogsService {
	return &AuditLogsService{
		Storage: storage,
	}
}

func (s *AuditLogsService) GetById(id uuid.UUID) (*models.AuditLog, error) {
	var auditlog models.AuditLog

	if err := s.Storage.Postgres.
		Where(&models.AuditLog{
			Id: id,
		}).
		Preload("ModifiedByUser").
		Find(&auditlog).Error; err != nil {
		return nil, err
	}

	return &auditlog, nil
}

func (s *AuditLogsService) GetAll(organizationId uuid.UUID, clauses ...clause.Expression) ([]models.AuditLog, error) {
	var auditlogs []models.AuditLog

	if err := s.Storage.Postgres.Clauses(clauses...).Find(&auditlogs).Error; err != nil {
		return nil, err
	}

	return auditlogs, nil
}

func (s *AuditLogsService) GetTotal(organizationId uuid.UUID, clauses ...clause.Expression) (int64, error) {
	var total int64

	if err := s.Storage.Postgres.Clauses(clauses...).Model(&models.AuditLog{}).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}
