package users

import (
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
	"github.com/google/uuid"
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
func (s *UsersService) Create(auditId uuid.UUID, user *models.User) error {
	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).Create(&user).Error; err != nil {
		return err
	}

	return nil
}

// Update updates the user record identified by the given id with the provided user data.
// It also sets the audit user ID for tracking who performed the update.
// Returns an error if the update operation fails.
func (s *UsersService) Update(auditId uuid.UUID, id uuid.UUID, user *models.User) error {
	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).
		Where(&models.User{
			Id: id,
		}).
		Updates(&user).Error; err != nil {
		return err
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
	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).
		Where(&models.User{
			Id: id,
		}).
		Delete(&models.User{}).Error; err != nil {
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
		Find(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
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
		Where(&models.User{
			Email: email,
		}).
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
func (s *UsersService) GetAll(clauses ...clause.Expression) ([]models.User, error) {
	var users []models.User

	if err := s.Storage.Postgres.Clauses(clauses...).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
