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

	if err := seedAdmin(db); err != nil {
		return err
	}

	return seedRewards(db)
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

func seedRewards(db *gorm.DB) error {
	rewards := []entity.Reward{
		{
			RewardID:     uuid.MustParse("aaaaaaaa-1111-4111-8111-aaaaaaaaaaaa"),
			Name:         "Pulsa Rp25.000",
			Description:  "Semua operator",
			RewardType:   "pulsa",
			PointsCost:   1000,
			Stock:        100,
			IsActive:     true,
			ValidityDays: 0,
		},
		{
			RewardID:     uuid.MustParse("bbbbbbbb-2222-4222-8222-bbbbbbbbbbbb"),
			Name:         "Voucher e-commerce Rp50.000",
			Description:  "Berlaku 90 hari",
			RewardType:   "voucher",
			PointsCost:   1800,
			Stock:        50,
			IsActive:     true,
			ValidityDays: 90,
		},
		{
			RewardID:     uuid.MustParse("cccccccc-3333-4333-8333-cccccccccccc"),
			Name:         "Donasikan kembali poin",
			Description:  "1 poin = Rp10 ke posko pilihan",
			RewardType:   "donation",
			PointsCost:   1,
			Stock:        999999,
			IsActive:     true,
			ValidityDays: 0,
		},
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, reward := range rewards {
			if err := seedReward(tx, reward); err != nil {
				return err
			}
		}

		return nil
	})
}

func seedReward(tx *gorm.DB, reward entity.Reward) error {
	var existingReward entity.Reward
	err := tx.Where("reward_id = ?", reward.RewardID).First(&existingReward).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return tx.Create(&reward).Error
}
