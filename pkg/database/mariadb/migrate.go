package mariadb

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if db.Migrator().HasTable(&entity.RegistrationSession{}) &&
		db.Migrator().HasIndex(&entity.RegistrationSession{}, "idx_registration_sessions_email") {
		if err := db.Migrator().DropIndex(&entity.RegistrationSession{}, "idx_registration_sessions_email"); err != nil {
			return err
		}
	}

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
		&entity.CustodyHandshakeToken{},
		&entity.OrderItems{},
		&entity.RegistrationSession{},
		&entity.RevokedToken{},
		&entity.AdminProfile{},
		&entity.DonorProfile{},
		&entity.PaymentTransactions{},
		&entity.PointAccount{},
		&entity.PointTransaction{},
		&entity.Reward{},
		&entity.RewardClaim{},
	)

	if err != nil {
		return err
	}

	return nil
}
