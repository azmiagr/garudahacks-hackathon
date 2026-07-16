package mariadb

import (
	"errors"

	"github.com/azmiagr/garudahacks-hackathon/entity"
	appbcrypt "github.com/azmiagr/garudahacks-hackathon/pkg/bcrypt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	adminRoleName   = "admin"
	donorRoleName   = "donor"
	storeRoleName   = "store"
	courierRoleName = "relawan"
	adminEmail      = "admin@example.com"
	adminPassword   = "admin123"
)

func Seed(db *gorm.DB) error {
	if err := seedRoles(db); err != nil {
		return err
	}

	return seedAdmin(db)
}

func seedRoles(db *gorm.DB) error {
	roleNames := []string{
		adminRoleName,
		donorRoleName,
		storeRoleName,
		courierRoleName,
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, roleName := range roleNames {
			if err := seedRole(tx, roleName); err != nil {
				return err
			}
		}

		return nil
	})
}

func seedRole(tx *gorm.DB, roleName string) error {
	var role entity.Role
	err := tx.Where("role_name = ?", roleName).First(&role).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	role = entity.Role{
		RoleID:   uuid.New(),
		RoleName: roleName,
	}

	return tx.Create(&role).Error
}

func seedAdmin(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var role entity.Role
		err := tx.Where("role_name = ?", adminRoleName).First(&role).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			role = entity.Role{
				RoleID:   uuid.New(),
				RoleName: adminRoleName,
			}
			if err := tx.Create(&role).Error; err != nil {
				return err
			}
		}

		var user entity.User
		err = tx.Where("email = ?", adminEmail).First(&user).Error
		if err == nil {
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		hashedPassword, err := appbcrypt.Init().GenerateFromPassword(adminPassword)
		if err != nil {
			return err
		}

		user = entity.User{
			UserID:    uuid.New(),
			RoleID:    role.RoleID,
			Name:      "Admin",
			Email:     adminEmail,
			Password:  hashedPassword,
			Status:    "active",
			KYCStatus: "approved",
		}

		return tx.Create(&user).Error
	})
}
