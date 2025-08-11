package transactions

import (
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type TransactionsService struct {
	Storage *storage.Storage
}

func NewTransactionsService(storage *storage.Storage) *TransactionsService {
	return &TransactionsService{
		Storage: storage,
	}
}

func (s *TransactionsService) Create(auditId uuid.UUID, organizationId uuid.UUID, transaction models.CreateTransactionPayload) error {
	var newTransaction models.Transaction

	newTransaction.Type = transaction.Type
	newTransaction.Weight = transaction.Weight
	newTransaction.Amount = transaction.Amount
	newTransaction.SellerID = transaction.SellerID
	newTransaction.BuyerID = transaction.BuyerID

	newTransaction.ModifiedByUserId = auditId

	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).Create(&newTransaction).Error; err != nil {
		return err
	}

	if transaction.Products != nil {
		products := []models.Product{}

		for _, productId := range transaction.Products {
			products = append(products, models.Product{
				Id: productId,
			})
		}

		if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).Model(&newTransaction).Association("Products").Append(products); err != nil {
			return err
		}
	}

	if newTransaction.SellerID == organizationId {
		if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).
			Model(&models.Organization{
				Id:               organizationId,
				ModifiedByUserId: auditId,
			}).
			Association("Sales").
			Append(&newTransaction); err != nil {
			return err
		}
	}

	if newTransaction.BuyerID == organizationId {
		if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).
			Model(&models.Organization{
				Id:               organizationId,
				ModifiedByUserId: auditId,
			}).
			Association("Purchases").
			Append(&newTransaction); err != nil {
			return err
		}
	}

	return nil
}

func (s *TransactionsService) Update(auditId uuid.UUID, id uuid.UUID, transaction models.UpdateTransactionPayload) error {
	var existingTransaction models.Transaction

	if err := s.Storage.Postgres.Where("id = $1", id).Find(&existingTransaction).Error; err != nil {
		return err
	}

	if transaction.Type != nil {
		existingTransaction.Type = *transaction.Type
	}

	if transaction.Weight != nil {
		existingTransaction.Weight = *transaction.Weight
	}

	if transaction.Amount != nil {
		existingTransaction.Amount = *transaction.Amount
	}

	if transaction.SellerID != uuid.Nil {
		existingTransaction.SellerID = transaction.SellerID
	}

	if transaction.BuyerID != uuid.Nil {
		existingTransaction.BuyerID = transaction.BuyerID
	}

	if transaction.SellerAccepted != nil {
		existingTransaction.SellerAccepted = *transaction.SellerAccepted
	}

	if transaction.SellerDeclined != nil {
		existingTransaction.SellerDeclined = *transaction.SellerDeclined
	}

	existingTransaction.ModifiedByUserId = auditId

	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).
		Where(&models.Transaction{
			Id: id,
		}).
		Updates(&map[string]any{
			"type":                existingTransaction.Type,
			"weight":              existingTransaction.Weight,
			"amount":              existingTransaction.Amount,
			"seller_id":           existingTransaction.SellerID,
			"buyer_id":            existingTransaction.BuyerID,
			"seller_accepted":     existingTransaction.SellerAccepted,
			"seller_declined":     existingTransaction.SellerDeclined,
			"modified_by_user_id": existingTransaction.ModifiedByUserId,
		}).Error; err != nil {
		return err
	}

	if transaction.Products != nil {
		products := []models.Product{}

		for _, productId := range transaction.Products {
			products = append(products, models.Product{
				Id: productId,
			})
		}

		if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).Model(&existingTransaction).Association("Products").Replace(products); err != nil {
			return err
		}
	}

	return nil
}

func (s *TransactionsService) Delete(auditId uuid.UUID, id uuid.UUID) error {
	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).
		Where(&models.Transaction{
			Id: id,
		}).
		Delete(&models.Transaction{}).Error; err != nil {
		return err
	}

	return nil
}

func (s *TransactionsService) GetById(id uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction

	if err := s.Storage.Postgres.
		Where(&models.Transaction{
			Id: id,
		}).
		Preload("Products.Materials").
		Preload("ModifiedBy").
		Find(&transaction).Error; err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (s *TransactionsService) GetAll(clauses ...clause.Expression) ([]models.Transaction, error) {
	var transactions []models.Transaction

	if err := s.Storage.Postgres.Clauses(clauses...).Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *TransactionsService) GetTotal(clauses ...clause.Expression) (int64, error) {
	var total int64

	if err := s.Storage.Postgres.Clauses(clauses...).Model(&models.Transaction{}).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}
