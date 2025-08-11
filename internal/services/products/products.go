package products

import (
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type ProductsService struct {
	Storage *storage.Storage
}

func NewProductsService(storage *storage.Storage) *ProductsService {
	return &ProductsService{
		Storage: storage,
	}
}

func (s *ProductsService) Create(auditId uuid.UUID, product models.CreateProductPayload) error {
	var newProduct models.Product

	newProduct.Name = product.Name
	newProduct.Value = product.Value

	newProduct.ModifiedByUserId = auditId

	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).Create(&newProduct).Error; err != nil {
		return err
	}

	if product.Materials != nil {
		materials := []models.Material{}

		for _, materialId := range product.Materials {
			materials = append(materials, models.Material{
				Id: materialId,
			})
		}

		if err := s.Storage.Postgres.Model(&newProduct).Association("Materials").Append(materials); err != nil {
			return err
		}
	}

	return nil
}

func (s *ProductsService) Update(auditId uuid.UUID, id uuid.UUID, product models.UpdateProductPayload) error {
	var existingProduct models.Product

	if err := s.Storage.Postgres.Where("id = $1", id).Find(&existingProduct).Error; err != nil {
		return err
	}

	if product.Name != nil {
		existingProduct.Name = *product.Name
	}

	if product.Value != nil {
		existingProduct.Value = *product.Value
	}

	existingProduct.ModifiedByUserId = auditId

	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).
		Where(&models.Product{
			Id: id,
		}).
		Updates(&map[string]any{
			"name":                existingProduct.Name,
			"value":               existingProduct.Value,
			"modified_by_user_id": existingProduct.ModifiedByUserId,
		}).Error; err != nil {
		return err
	}

	if product.Materials != nil {
		materials := []models.Material{}

		for _, materialId := range product.Materials {
			materials = append(materials, models.Material{
				Id: materialId,
			})
		}

		if err := s.Storage.Postgres.Model(&existingProduct).Association("Materials").Replace(materials); err != nil {
			return err
		}
	}

	return nil
}

func (s *ProductsService) Delete(auditId uuid.UUID, id uuid.UUID) error {
	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).
		Where(&models.Product{
			Id: id,
		}).
		Delete(&models.Product{}).Error; err != nil {
		return err
	}

	return nil
}

func (s *ProductsService) GetById(id uuid.UUID) (*models.Product, error) {
	var product models.Product

	if err := s.Storage.Postgres.
		Where(&models.Product{
			Id: id,
		}).
		Find(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (s *ProductsService) GetAll(clauses ...clause.Expression) ([]models.Product, error) {
	var products []models.Product

	if err := s.Storage.Postgres.Clauses(clauses...).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (s *ProductsService) GetTotal(clauses ...clause.Expression) (int64, error) {
	var total int64

	if err := s.Storage.Postgres.Clauses(clauses...).Model(&models.Product{}).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}
