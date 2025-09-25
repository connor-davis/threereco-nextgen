package storage

import (
	"errors"

	"github.com/connor-davis/threereco-nextgen/common"
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage interface {
	Database() *gorm.DB
	Migrate() error
	SeedAdmin() error
	SeedDefaultBusiness() error
}

type storage struct {
	db *gorm.DB
}

func New() Storage {
	dsn := common.EnvString("APP_DSN", "host=localhost user=postgres password=postgres dbname=kalimbu port=5432 sslmode=disable TimeZone=Africa/Johannesburg")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Errorf("failed to connect database: %s", err.Error())
	}

	return &storage{
		db: db,
	}
}

func (s *storage) Database() *gorm.DB {
	return s.db
}

func (s *storage) Migrate() error {
	if err := s.db.AutoMigrate(
		&models.Business{},
		&models.User{},
		&models.Role{},
		&models.Material{},
		&models.Collection{},
		&models.CollectionMaterial{},
		&models.Transaction{},
		&models.TransactionMaterial{},
	); err != nil {
		log.Errorf("failed to migrate database: %s", err.Error())

		return err
	}

	return nil
}

func (s *storage) SeedAdmin() error {
	adminUserId := uuid.New()
	adminRoleId := uuid.New()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(common.EnvString("APP_ADMIN_PASSWORD", "password")), bcrypt.DefaultCost)

	if err != nil {
		log.Errorf("failed to hash admin password: %s", err.Error())
	}

	adminRoleName := "Administrator"
	adminRoleDescription := "Full access to all system resources."

	adminRole := models.Role{
		Base: models.Base{
			Id: adminRoleId,
		},
		Name:        adminRoleName,
		Description: &adminRoleDescription,
		Permissions: []string{
			"*",
		},
		Default: true,
	}

	adminUser := models.User{
		Base: models.Base{
			Id: adminUserId,
		},
		Name:     common.EnvString("APP_ADMIN_NAME", "Admin User"),
		Username: common.EnvString("APP_ADMIN_EMAIL", "admin@example.com"),
		Password: hashedPassword,
		Roles: []models.Role{
			adminRole,
		},
	}

	var existingAdminRole models.Role
	var existingAdminUser models.User

	if err := s.db.Where("name = ?", adminRoleName).First(&existingAdminRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := s.db.Create(&adminRole).Error; err != nil {
				log.Errorf("failed to create admin role: %s", err.Error())
			}
		} else {
			log.Errorf("failed to query admin role: %s", err.Error())
		}
	}

	if err := s.db.Where("username = ?", adminUser.Username).First(&existingAdminUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := s.db.Create(&adminUser).Error; err != nil {
				log.Errorf("failed to create admin user: %s", err.Error())
			}
		} else {
			log.Errorf("failed to query admin user: %s", err.Error())
		}
	}

	return nil
}

func (s *storage) SeedDefaultBusiness() error {
	businessOwnerId := uuid.New()
	businessId := uuid.New()
	businessOwnerRoleId := uuid.New()
	businessStaffRoleId := uuid.New()
	businessUserRoleId := uuid.New()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(common.EnvString("APP_DEFAULT_BUSINESS_PASSWORD", "password")), bcrypt.DefaultCost)

	if err != nil {
		log.Errorf("failed to hash default business owner password: %s", err.Error())
	}

	businessOwnerRoleName := "Business Owner"
	businessOwnerRoleDescription := "Owner of the business with full access to business resources."

	businessStaffRoleName := "Business Staff"
	businessStaffRoleDescription := "Staff member of the business with limited access to business resources."

	businessUserRoleName := "Business User"
	businessUserRoleDescription := "User of the business with minimal access to business resources."

	businessOwnerRole := models.Role{
		Base: models.Base{
			Id: businessOwnerRoleId,
		},
		Name:        businessOwnerRoleName,
		Description: &businessOwnerRoleDescription,
		Permissions: []string{
			"materials.view",
			"collections.*",
			"transactions.*",
			"users.view.self",
			"users.update.self",
			"users.delete.self",
			"businesses.view",
			"businesses.update.self",
			"businesses.delete.self",
			"businesses.roles.assign",
			"businesses.roles.unassign",
			"businesses.roles.view",
			"businesses.users.assign",
			"businesses.users.unassign",
			"businesses.users.view",
		},
		Default: false,
	}

	businessStaffRole := models.Role{
		Base: models.Base{
			Id: businessStaffRoleId,
		},
		Name:        businessStaffRoleName,
		Description: &businessStaffRoleDescription,
		Permissions: []string{
			"materials.view",
			"collections.view",
			"collections.create",
			"collections.update",
			"transactions.view",
			"transactions.create",
			"transactions.update",
			"users.view.self",
			"users.update.self",
			"users.delete.self",
			"businesses.view",
			"businesses.users.view",
		},
		Default: false,
	}

	businessUserRole := models.Role{
		Base: models.Base{
			Id: businessUserRoleId,
		},
		Name:        businessUserRoleName,
		Description: &businessUserRoleDescription,
		Permissions: []string{
			"materials.view",
			"collections.view",
			"transactions.view",
			"users.view.self",
			"users.update.self",
			"users.delete.self",
			"businesses.view",
		},
		Default: false,
	}

	businessOwner := models.User{
		Base: models.Base{
			Id: businessOwnerId,
		},
		Name:     common.EnvString("APP_DEFAULT_BUSINESS_NAME", "Demo Business"),
		Username: common.EnvString("APP_DEFAULT_BUSINESS_EMAIL", "demo@3reco.co.za"),
		Password: hashedPassword,
		Roles: []models.Role{
			businessOwnerRole,
		},
		Type:       models.BusinessUser,
		BusinessId: &businessId,
	}

	business := models.Business{
		Base: models.Base{
			Id: businessId,
		},
		Name:    common.EnvString("APP_DEFAULT_BUSINESS_NAME", "Demo Business"),
		OwnerId: businessOwnerId,
		Users: []models.User{
			businessOwner,
		},
	}

	var existingBusinessOwner models.User
	var existingBusinessOwnerRole models.Role
	var existingBusinessStaffRole models.Role
	var existingBusinessUserRole models.Role
	var existingBusiness models.Business

	if err := s.db.Where("name = ?", businessOwnerRoleName).First(&existingBusinessOwnerRole).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := s.db.Create(&businessOwnerRole).Error; err != nil {
				log.Errorf("failed to create business role: %s", err.Error())
			}
		} else {
			log.Errorf("failed to query business role: %s", err.Error())
		}
	}

	if err := s.db.Where("name = ?", businessStaffRoleName).First(&existingBusinessStaffRole).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := s.db.Create(&businessStaffRole).Error; err != nil {
				log.Errorf("failed to create business staff role: %s", err.Error())
			}
		} else {
			log.Errorf("failed to query business staff role: %s", err.Error())
		}
	}

	if err := s.db.Where("name = ?", businessUserRoleName).First(&existingBusinessUserRole).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := s.db.Create(&businessUserRole).Error; err != nil {
				log.Errorf("failed to create business user role: %s", err.Error())
			}
		} else {
			log.Errorf("failed to query business user role: %s", err.Error())
		}
	}

	if err := s.db.Where("username = ?", businessOwner.Username).First(&existingBusinessOwner).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := s.db.Create(&businessOwner).Error; err != nil {
				log.Errorf("failed to create business owner: %s", err.Error())
			}
		} else {
			log.Errorf("failed to query business owner: %s", err.Error())
		}
	}

	if err := s.db.Where("name = ?", business.Name).First(&existingBusiness).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := s.db.Create(&business).Error; err != nil {
				log.Errorf("failed to create default business: %s", err.Error())
			}
		} else {
			log.Errorf("failed to query default business: %s", err.Error())
		}
	}

	return nil
}
