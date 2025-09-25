package models

import (
	"crypto/rand"
	"math/big"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserType string

const (
	SystemUser    UserType = "system"
	CollectorUser UserType = "collector"
	BusinessUser  UserType = "business"
)

type User struct {
	Base
	Name          string         `json:"name" gorm:"not null"`
	Username      string         `json:"username" gorm:"uniqueIndex;not null"`
	Password      []byte         `json:"-" gorm:"type:bytea"`
	PasswordReset bool           `json:"passwordReset" gorm:"default:false"`
	MfaSecret     []byte         `json:"-" gorm:"type:bytea"`
	MfaEnabled    bool           `json:"mfaEnabled" gorm:"default:false"`
	MfaVerified   bool           `json:"mfaVerified" gorm:"default:false"`
	Permissions   pq.StringArray `json:"permissions" gorm:"type:text[];default:'{}'"`
	Roles         []Role         `json:"roles" gorm:"many2many:users_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Type          UserType       `json:"type" gorm:"type:text;not null;default:'system'"`
	Address       *Address       `json:"address" gorm:"type:jsonb;"`
	BankDetails   *BankDetails   `json:"bankDetails" gorm:"type:jsonb;"`
	IdNumber      *string        `json:"idNumber" gorm:"type:text;"`
	BusinessId    *uuid.UUID     `json:"businessId" gorm:"type:uuid;"`
	Businesses    []Business     `json:"businesses" gorm:"many2many:businesses_users;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-="

	if len(u.Password) == 0 {
		const RandomPasswordLength = 32

		password := make([]byte, RandomPasswordLength)

		for i := range password {
			num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))

			if err != nil {
				return err
			}

			password[i] = charset[num.Int64()]
		}

		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

		if err != nil {
			return err
		}

		u.Password = hashedPassword
		u.PasswordReset = true
	}

	return nil
}
