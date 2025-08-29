package storage

import (
	"github.com/connor-davis/threereco-nextgen/env"
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Storage encapsulates the application's database connection.
// It holds a reference to a GORM-managed PostgreSQL database instance.
type Storage struct {
	Postgres *gorm.DB
}

// New creates and returns a new instance of Storage.
func New() *Storage {
	return &Storage{}
}

// ConnectPostgres establishes a connection to a PostgreSQL database using the DSN provided in the environment.
// If the connection is successful, it assigns the database instance to the Storage struct's Postgres field.
// In case of failure, it logs an error message and returns without setting the Postgres field.
func (s *Storage) ConnectPostgres() {
	database, err := gorm.Open(postgres.Open(string(env.POSTGRES_DSN)), &gorm.Config{})

	if err != nil {
		log.Infof("🔥 Failed to connect to Postgres: %v", err)
		return
	}

	log.Info("✅ Successfully connected to Postgres")

	s.Postgres = database
}

// MigratePostgres performs database migrations for the Postgres connection.
// It ensures the required schema and extensions are created, then runs GORM
// auto-migrations for core models. If the Postgres connection is not established,
// or if any migration step fails, appropriate error messages are logged.
func (s *Storage) MigratePostgres() {
	if s.Postgres == nil {
		log.Error("❌ Postgres connection is not established, cannot migrate")

		return
	}

	log.Info("🔄 Running Postgres migrations...")

	if err := s.Postgres.Exec(`
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	`).Error; err != nil {
		log.Errorf("❌ Failed to create schema or extensions: %v", err)

		return
	} else {
		log.Info("✅ Extensions created successfully")
	}

	log.Info("🔃 Running GORM migrations...")

	if err := s.Postgres.AutoMigrate(
		&models.AuditLog{},
		&models.User{},
		&models.Organization{},
		&models.Role{},
		&models.Product{},
		&models.Material{},
		&models.Transaction{},
		&models.Notification{},
		&models.Address{},
		&models.BankDetails{},
	); err != nil {
		log.Errorf("❌ AutoMigrate failed: %v", err)

		return
	}

	log.Info("✅ Postgres migrations completed successfully")
}

// SeedPostgres seeds the PostgreSQL database with default organization, user, and role data.
// It checks if the Postgres connection is established, then attempts to find or create:
//   - An "Organization Admin" role with full permissions.
//   - An organization admin user with default credentials.
//   - A default organization associated with the admin user and role.
//
// The function ensures that these entities exist, creating them if necessary, and logs progress and errors.
// This is typically used for initial setup or testing environments.
func (s *Storage) SeedPostgres() {
	if s.Postgres == nil {
		log.Error("❌ Postgres connection is not established, cannot seed")

		return
	}

	log.Info("🔄 Seeding Postgres...")

	organizationId := uuid.New()
	organizationAdminUserId := uuid.New()
	organizationAdminRoleId := uuid.New()
	organizationUserRoleId := uuid.New()

	organizationAdminRoleDescription := "Has full access to the organization, including user management, settings, and data."
	organizationUserRoleDescription := "Has limited access to the organization, primarily for their own profile information or data related to them that has been created by other users with their information."

	organizationAdminRole := &models.Role{
		Name:             "Administrator",
		Description:      &organizationAdminRoleDescription,
		Permissions:      []string{"*"},
		IsDefault:        true,
		ModifiedByUserId: organizationAdminUserId,
	}

	if err := s.Postgres.
		Where("name = ?", organizationAdminRole.Name).
		Find(&organizationAdminRole).Error; err != nil {
		log.Errorf("❌ Failed to find organization admin role: %v", err)

		return
	}

	if organizationAdminRole.Id == uuid.Nil {
		log.Info("🔄 Creating organization admin role...")

		organizationAdminRole.Id = organizationAdminRoleId
	}

	organizationUserRole := &models.Role{
		Name:             "User",
		Description:      &organizationUserRoleDescription,
		Permissions:      []string{"users.view.self", "users.update.self", "users.delete.self"},
		IsDefault:        true,
		ModifiedByUserId: organizationAdminUserId,
	}

	if err := s.Postgres.
		Where("name = ?", organizationUserRole.Name).
		Find(&organizationUserRole).Error; err != nil {
		log.Errorf("❌ Failed to find organization user role: %v", err)

		return
	}

	if organizationUserRole.Id == uuid.Nil {
		log.Info("🔄 Creating organization user role...")

		organizationUserRole.Id = organizationUserRoleId
	}

	organizationAdminUserName := string(env.DEFAULT_ORGANIZATION_ADMIN_NAME)
	organizationAdminUserPhone := string(env.DEFAULT_ORGANIZATION_ADMIN_PHONE)

	organizationAdminUserHashedPassword, err := bcrypt.GenerateFromPassword([]byte(string(env.DEFAULT_ORGANIZATION_ADMIN_PASSWORD)), bcrypt.DefaultCost)

	if err != nil {
		log.Errorf("❌ Failed to hash organization admin user password: %v", err)

		return
	}

	organizationAdminUser := &models.User{
		Email:                 string(env.DEFAULT_ORGANIZATION_ADMIN_EMAIL),
		Password:              organizationAdminUserHashedPassword,
		Name:                  &organizationAdminUserName,
		Phone:                 &organizationAdminUserPhone,
		Image:                 nil,
		ModifiedByUserId:      organizationAdminUserId,
		PrimaryOrganizationId: &organizationId,
		Roles:                 []models.Role{*organizationAdminRole},
	}

	if err := s.Postgres.
		Where("email = ?", string(env.DEFAULT_ORGANIZATION_ADMIN_EMAIL)).
		Find(&organizationAdminUser).Error; err != nil {
		log.Errorf("❌ Failed to find organization admin user: %v", err)

		return
	}

	if organizationAdminUser.Id == uuid.Nil {
		log.Info("🔄 Creating organization admin user...")

		organizationAdminUser.Id = organizationAdminUserId
	}

	organization := &models.Organization{
		Name:             string(env.DEFAULT_ORGANIZATION_NAME),
		Domain:           string(env.DEFAULT_ORGANIZATION_DOMAIN),
		OwnerId:          organizationAdminUser.Id,
		ModifiedByUserId: organizationAdminUser.Id,
		ModifiedByUser:   organizationAdminUser,
		Users:            []models.User{*organizationAdminUser},
		Roles:            []models.Role{*organizationAdminRole, *organizationUserRole},
	}

	if err := s.Postgres.
		Where("name = ?", string(env.DEFAULT_ORGANIZATION_NAME)).
		Find(&organization).Error; err != nil {
		log.Errorf("❌ Failed to find organization: %v", err)

		return
	}

	if organization.Id == uuid.Nil {
		log.Info("🔄 Creating organization...")

		organization.Id = organizationId
	}

	if err := s.Postgres.
		Set("one:ignore_audit_log", true).
		Assign(&organization).
		FirstOrCreate(&organization).Error; err != nil {
		log.Errorf("❌ Failed to create organization admin user: %v", err)

		return
	}

	log.Info("✅ Organization, users and roles seeded successfully")
}
