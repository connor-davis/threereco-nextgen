package materials

import (
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type MaterialsService struct {
	Storage *storage.Storage
}

func NewMaterialsService(storage *storage.Storage) *MaterialsService {
	return &MaterialsService{
		Storage: storage,
	}
}

func (s *MaterialsService) Create(auditId uuid.UUID, organizationId uuid.UUID, material models.CreateMaterialPayload) error {
	var newMaterial models.Material

	newMaterial.Name = material.Name
	newMaterial.GwCode = material.GwCode
	newMaterial.CarbonFactor = material.CarbonFactor

	newMaterial.ModifiedByUserId = auditId

	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).
		Create(&newMaterial).Error; err != nil {
		return err
	}

	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).
		Model(&models.Organization{
			Id:               organizationId,
			ModifiedByUserId: auditId,
		}).
		Association("Materials").
		Append(&newMaterial); err != nil {
		return err
	}

	return nil
}

func (s *MaterialsService) Update(auditId uuid.UUID, id uuid.UUID, material models.UpdateMaterialPayload) error {
	var existingMaterial models.Material

	if err := s.Storage.Postgres.Where("id = $1", id).Find(&existingMaterial).Error; err != nil {
		return err
	}

	if material.Name != nil {
		existingMaterial.Name = *material.Name
	}

	if material.GwCode != nil {
		existingMaterial.GwCode = *material.GwCode
	}

	if material.CarbonFactor != nil {
		existingMaterial.CarbonFactor = *material.CarbonFactor
	}

	existingMaterial.ModifiedByUserId = auditId

	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).
		Where(&models.Material{
			Id: id,
		}).
		Updates(&map[string]any{
			"name":                existingMaterial.Name,
			"gw_code":             existingMaterial.GwCode,
			"carbon_factor":       existingMaterial.CarbonFactor,
			"modified_by_user_id": existingMaterial.ModifiedByUserId,
		}).Error; err != nil {
		return err
	}

	return nil
}

func (s *MaterialsService) Delete(auditId uuid.UUID, id uuid.UUID) error {
	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).
		Where(&models.Material{
			Id: id,
		}).
		Delete(&models.Material{}).Error; err != nil {
		return err
	}

	return nil
}

func (s *MaterialsService) GetById(id uuid.UUID) (*models.Material, error) {
	var material models.Material

	if err := s.Storage.Postgres.
		Where(&models.Material{
			Id: id,
		}).
		Preload("ModifiedByUser").
		Find(&material).Error; err != nil {
		return nil, err
	}

	return &material, nil
}

func (s *MaterialsService) GetAll(organizationId uuid.UUID, clauses ...clause.Expression) ([]models.Material, error) {
	var materials []models.Material

	if err := s.Storage.Postgres.
		Model(&models.Material{}).
		Joins("JOIN organizations_materials om ON om.material_id = materials.id").
		Where("om.organization_id = ?", organizationId).
		Clauses(clauses...).
		Find(&materials).Error; err != nil {
		return nil, err
	}

	return materials, nil
}

func (s *MaterialsService) GetTotal(organizationId uuid.UUID, clauses ...clause.Expression) (int64, error) {
	var total int64

	if err := s.Storage.Postgres.
		Model(&models.Material{}).
		Joins("JOIN organizations_materials om ON om.material_id = materials.id").
		Where("om.organization_id = ?", organizationId).
		Clauses(clauses...).
		Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}
