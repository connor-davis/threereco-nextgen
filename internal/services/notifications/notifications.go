package notifications

import (
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type NotificationsService struct {
	Storage *storage.Storage
}

func NewNotificationsService(storage *storage.Storage) *NotificationsService {
	return &NotificationsService{
		Storage: storage,
	}
}

func (s *NotificationsService) Create(auditId uuid.UUID, notification models.CreateNotificationPayload) error {
	var newNotification models.Notification

	newNotification.Title = notification.Title
	newNotification.Message = notification.Message

	if notification.Action != nil {
		newNotification.Action = notification.Action
	}

	newNotification.ModifiedByUserId = auditId

	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).Create(&newNotification).Error; err != nil {
		return err
	}

	return nil
}

func (s *NotificationsService) Update(auditId uuid.UUID, id uuid.UUID, notification models.UpdateNotificationPayload) error {
	var existingNotification models.Notification

	if err := s.Storage.Postgres.Where("id = $1", id).Find(&existingNotification).Error; err != nil {
		return err
	}

	if notification.Title != nil {
		existingNotification.Title = *notification.Title
	}

	if notification.Message != nil {
		existingNotification.Message = *notification.Message
	}

	if notification.Action != nil {
		existingNotification.Action = notification.Action
	}

	existingNotification.ModifiedByUserId = auditId

	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).
		Where(&models.Notification{
			Id: id,
		}).
		Updates(&map[string]any{
			"title":               existingNotification.Title,
			"message":             existingNotification.Message,
			"action":              existingNotification.Action,
			"modified_by_user_id": existingNotification.ModifiedByUserId,
		}).Error; err != nil {
		return err
	}

	return nil
}

func (s *NotificationsService) Delete(auditId uuid.UUID, id uuid.UUID) error {
	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).
		Where(&models.Notification{
			Id: id,
		}).
		Delete(&models.Notification{}).Error; err != nil {
		return err
	}

	return nil
}

func (s *NotificationsService) GetById(id uuid.UUID) (*models.Notification, error) {
	var notification models.Notification

	if err := s.Storage.Postgres.
		Where(&models.Notification{
			Id: id,
		}).
		Find(&notification).Error; err != nil {
		return nil, err
	}

	return &notification, nil
}

func (s *NotificationsService) GetAll(clauses ...clause.Expression) ([]models.Notification, error) {
	var notifications []models.Notification

	if err := s.Storage.Postgres.Clauses(clauses...).Find(&notifications).Error; err != nil {
		return nil, err
	}

	return notifications, nil
}

func (s *NotificationsService) GetTotal(clauses ...clause.Expression) (int64, error) {
	var total int64

	if err := s.Storage.Postgres.Clauses(clauses...).Model(&models.Notification{}).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}
