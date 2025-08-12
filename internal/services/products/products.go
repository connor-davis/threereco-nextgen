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

func (s *ProductsService) Create(auditId uuid.UUID, organizationId uuid.UUID, product models.CreateProductPayload) error {
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

		if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).Model(&newProduct).Association("Materials").Append(materials); err != nil {
			return err
		}
	}

	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).
		Model(&models.Organization{
			Id:               organizationId,
			ModifiedByUserId: auditId,
		}).
		Association("Products").
		Append(&newProduct); err != nil {
		return err
	}

	return nil
}

func (s *ProductsService) Update(auditId uuid.UUID, id uuid.UUID, product models.UpdateProductPayload) error {
	var existingProduct models.Product

	if err := s.Storage.Postgres.Where("id = ?", id).Find(&existingProduct).Error; err != nil {
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
		Model(&models.Product{}).
		Where("id = ?", id).
		Updates(&map[string]any{
			"name":                existingProduct.Name,
			"value":               existingProduct.Value,
			"modified_by_user_id": auditId,
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

		if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).Model(&existingProduct).Association("Materials").Replace(materials); err != nil {
			return err
		}
	}

	return nil
}

func (s *ProductsService) Delete(auditId uuid.UUID, id uuid.UUID) error {
	var existingProduct models.Product

	if err := s.Storage.Postgres.
		Where("id = ?", id).
		Find(&existingProduct).Error; err != nil {
		return err
	}

	if existingProduct.Id == uuid.Nil {
		return nil
	}

	if err := s.Storage.Postgres.
		Set("one:audit_user_id", auditId).
		Delete(&existingProduct).Error; err != nil {
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
		Preload("Materials").
		Preload("ModifiedByUser").
		Find(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (s *ProductsService) GetAll(organizationId uuid.UUID, clauses ...clause.Expression) ([]models.Product, error) {
	var products []models.Product

	productsClauses := []clause.Expression{
		clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "created_at",
					},
					Desc: true,
				},
			},
		},
	}

	productsClauses = append(productsClauses, clauses...)

	if err := s.Storage.Postgres.
		Model(&models.Product{}).
		Joins("JOIN organizations_products op ON op.product_id = products.id").
		Where("op.organization_id = ?", organizationId).
		Clauses(productsClauses...).
		Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (s *ProductsService) GetTotal(organizationId uuid.UUID, clauses ...clause.Expression) (int64, error) {
	var total int64

	if err := s.Storage.Postgres.
		Model(&models.Product{}).
		Joins("JOIN organizations_products op ON op.product_id = products.id").
		Where("op.organization_id = ?", organizationId).
		Clauses(clauses...).
		Model(&models.Product{}).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}
