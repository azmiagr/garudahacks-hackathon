package mariadb

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&entity.Role{},
		&entity.User{},
		&entity.Stores{},
		&entity.Post{},
		&entity.DisasterEvent{},
		&entity.DisasterReport{},
		&entity.Wallets{},
		&entity.WalletTransactions{},
		&entity.Requests{},
		&entity.Items{},
		&entity.Donations{},
		&entity.Orders{},
		&entity.DeliveryVerification{},
		&entity.Disbursements{},
		&entity.CustodyLogs{},
		&entity.OrderItems{},
		&entity.RegistrationSession{},
		&entity.AdminProfile{},
	)

	if err != nil {
		return err
	}

	return nil
}
