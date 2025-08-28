package users

import (
	"errors"

	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

// UsersService provides methods for managing user-related operations,
// utilizing the underlying Storage for persistence and retrieval.
type UsersService struct {
	Storage *storage.Storage
}

// NewUsersService creates and returns a new instance of UsersService using the provided storage.
// It initializes the UsersService with the given storage backend.
//
// Parameters:
//
//	storage - a pointer to a Storage instance used for user data persistence.
//
// Returns:
//
//	A pointer to a newly created UsersService.
func NewUsersService(storage *storage.Storage) *UsersService {
	return &UsersService{
		Storage: storage,
	}
}

// Create adds a new user to the database with the provided audit ID.
// It sets the audit ID in the database context for tracking purposes.
// Returns an error if the user creation fails.
func (s *UsersService) Create(auditId uuid.UUID, organizationId uuid.UUID, user models.CreateUserPayload) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	var newUser models.User

	newUser.Email = user.Email
	newUser.Password = hashedPassword
	newUser.Name = user.Name
	newUser.Phone = user.Phone
	newUser.JobTitle = user.JobTitle
	newUser.Tags = user.Tags

	newUser.ModifiedByUserId = auditId

	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).Create(&newUser).Error; err != nil {
		return err
	}

	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).Model(&models.Organization{Id: organizationId}).Association("Users").Append(&newUser); err != nil {
		return err
	}

	if len(user.Roles) > 0 {
		roles := []models.Role{}

		for _, roleId := range user.Roles {
			roles = append(roles, models.Role{
				Id: roleId,
			})
		}

		if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).Model(&newUser).Association("Roles").Append(roles); err != nil {
			return err
		}
	}

	if user.Address != nil {
		user.Address.UserId = newUser.Id

		if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).Create(user.Address).Error; err != nil {
			return err
		}
	}

	if user.BankDetails != nil {
		user.BankDetails.UserId = newUser.Id

		if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).Create(user.BankDetails).Error; err != nil {
			return err
		}
	}

	return nil
}

// Update updates the user record identified by the given id with the provided user data.
// It also sets the audit user ID for tracking who performed the update.
// Returns an error if the update operation fails.
func (s *UsersService) Update(auditId uuid.UUID, id uuid.UUID, user models.UpdateUserPayload) error {
	var existingUser models.User

	if err := s.Storage.Postgres.Where("id = ?", id).Find(&existingUser).Error; err != nil {
		return err
	}

	if existingUser.Id == uuid.Nil {
		return errors.New("user not found")
	}

	if user.Email != nil {
		existingUser.Email = *user.Email
	}

	if user.Name != nil {
		existingUser.Name = user.Name
	}

	if user.Phone != nil {
		existingUser.Phone = user.Phone
	}

	if user.JobTitle != nil {
		existingUser.JobTitle = user.JobTitle
	}

	if user.Tags != nil {
		existingUser.Tags = user.Tags
	}

	existingUser.ModifiedByUserId = auditId

	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).
		Model(&models.User{}).
		Where("id = ?", id).
		Updates(&map[string]any{
			"email":               existingUser.Email,
			"name":                existingUser.Name,
			"phone":               existingUser.Phone,
			"job_title":           existingUser.JobTitle,
			"tags":                existingUser.Tags,
			"modified_by_user_id": auditId,
		}).Error; err != nil {
		return err
	}

	if len(user.Roles) > 0 {
		roles := []models.Role{}

		for _, roleId := range user.Roles {
			roles = append(roles, models.Role{
				Id: roleId,
			})
		}

		if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).Model(&existingUser).Association("Roles").Replace(roles); err != nil {
			return err
		}
	}

	if user.Address != nil {
		if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).
			Model(&models.Address{}).
			Where("user_id = ?", existingUser.Id).
			FirstOrCreate(map[string]any{
				"line_one":    user.Address.LineOne,
				"line_two":    user.Address.LineTwo,
				"city":        user.Address.City,
				"postal_code": user.Address.PostalCode,
				"state":       user.Address.State,
				"country":     user.Address.Country,
				"user_id":     existingUser.Id,
			}).Error; err != nil {
			return err
		}
	}

	if user.BankDetails != nil {
		if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).
			Model(&models.BankDetails{}).
			Where("user_id = ?", existingUser.Id).
			FirstOrCreate(map[string]any{
				"account_holder": user.BankDetails.AccountHolder,
				"account_number": user.BankDetails.AccountNumber,
				"bank_name":      user.BankDetails.BankName,
				"branch_code":    user.BankDetails.BranchCode,
				"user_id":        existingUser.Id,
			}).Error; err != nil {
			return err
		}
	}

	return nil
}

// Delete removes a user record from the database based on the provided user ID.
// It also sets the audit user ID for tracking purposes.
// Parameters:
//   - auditId: UUID of the user performing the deletion, used for auditing.
//   - id: String identifier of the user to be deleted.
//
// Returns:
//   - error: An error if the deletion fails, otherwise nil.
func (s *UsersService) Delete(auditId uuid.UUID, id uuid.UUID) error {
	var existingUser models.User

	if err := s.Storage.Postgres.
		Where("id = ?", id).
		Find(&existingUser).Error; err != nil {
		return err
	}

	if existingUser.Id == uuid.Nil {
		return nil
	}

	if err := s.Storage.Postgres.
		Set("one:audit_user_id", auditId).
		Delete(&existingUser).Error; err != nil {
		return err
	}

	return nil
}

// GetById retrieves a user from the database by their unique ID.
// It returns a pointer to the User model if found, or an error if the operation fails.
//
// Parameters:
//   - id: The unique identifier of the user to retrieve.
//
// Returns:
//   - *models.User: Pointer to the retrieved user model.
//   - error: Error encountered during the database query, or nil if successful.
func (s *UsersService) GetById(id uuid.UUID) (*models.User, error) {
	var user models.User

	if err := s.Storage.Postgres.
		Where(&models.User{
			Id: id,
		}).
		Preload("Roles").
		Preload("Organizations.Owner").
		Preload("ModifiedByUser").
		Preload("Address").
		Preload("BankDetails").
		Find(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// GetSales returns all transactions where the provided userId matches the SellerID.
// Optional GORM clause expressions may be supplied to customize the query (e.g., OrderBy, Limit, Locking).
// It returns the matching transactions on success, or a non-nil error if the query fails.
func (s *UsersService) GetSales(userId uuid.UUID, clauses ...clause.Expression) ([]models.Transaction, error) {
	var sales []models.Transaction

	if err := s.Storage.Postgres.
		Clauses(clauses...).
		Where(&models.Transaction{
			SellerID: userId,
		}).
		Find(&sales).Error; err != nil {
		return nil, err
	}

	return sales, nil
}

// GetPurchases retrieves all transactions where the given userId matches Transaction.BuyerID.
// Optional GORM clause.Expression values can be provided to refine the query (e.g., ORDER BY,
// LIMIT, or locking clauses). It returns an empty slice and a nil error when no records are found,
// and a non-nil error only if the query execution fails.
func (s *UsersService) GetPurchases(userId uuid.UUID, clauses ...clause.Expression) ([]models.Transaction, error) {
	var purchases []models.Transaction

	if err := s.Storage.Postgres.
		Clauses(clauses...).
		Where(&models.Transaction{
			BuyerID: userId,
		}).
		Find(&purchases).Error; err != nil {
		return nil, err
	}

	return purchases, nil
}

// GetNotifications returns notifications that belong to the provided userId.
// Optional GORM clause expressions can be supplied to customize the query
// (for example: Order("created_at DESC"), Limit(n), Offset(n), Preload("...")).
// It returns the matching notifications, or an error if the query fails.
// If no records match, the returned slice may be empty.
func (s *UsersService) GetNotifications(userId uuid.UUID, clauses ...clause.Expression) ([]models.Notification, error) {
	var notifications []models.Notification

	if err := s.Storage.Postgres.
		Clauses(clauses...).
		Where(&models.Notification{
			UserId: userId,
		}).
		Find(&notifications).Error; err != nil {
		return nil, err
	}

	return notifications, nil
}

// GetByEmail retrieves a user from the database by their email address.
// It returns a pointer to the User model and an error if the operation fails.
// If no user is found with the given email, the returned error will indicate the failure.
//
// Parameters:
//   - email: The email address of the user to retrieve.
//
// Returns:
//   - *models.User: Pointer to the retrieved user, or nil if not found.
//   - error: Error encountered during the database query, or nil if successful.
func (s *UsersService) GetByEmail(email string) (*models.User, error) {
	var user models.User

	if err := s.Storage.Postgres.
		Where("email = ?", email).
		Preload("Roles").
		Preload("Organizations.Owner").
		Preload("ModifiedByUser").
		Find(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// GetAll retrieves all users from the database, applying the provided GORM clause expressions as query modifiers.
// It returns a slice of models.User and an error if the query fails.
//
// Parameters:
//
//	clauses - Optional GORM clause expressions to customize the query (e.g., filtering, ordering).
//
// Returns:
//
//	[]models.User - A slice containing the retrieved user records.
//	error         - An error if the database query fails, otherwise nil.
func (s *UsersService) GetAll(organizationId uuid.UUID, clauses ...clause.Expression) ([]models.User, error) {
	var users []models.User

	usersClauses := []clause.Expression{
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

	usersClauses = append(usersClauses, clauses...)

	if err := s.Storage.Postgres.
		Model(&models.User{}).
		Joins("JOIN organizations_users ou ON ou.user_id = users.id").
		Where("ou.organization_id = ?", organizationId).
		Clauses(usersClauses...).
		Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// GetTotal returns the total number of User records in the database that match the provided GORM clause expressions.
// It accepts a variadic number of clause.Expression arguments to filter the query.
// Returns the count as int64 and an error if the query fails.
func (s *UsersService) GetTotal(organizationId uuid.UUID, clauses ...clause.Expression) (int64, error) {
	var total int64

	if err := s.Storage.Postgres.
		Model(&models.User{}).
		Joins("JOIN organizations_users ou ON ou.user_id = users.id").
		Where("ou.organization_id = ?", organizationId).
		Clauses(clauses...).
		Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}
