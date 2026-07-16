package service

import (
	"errors"
	"time"

	"github.com/azmiagr/garudahacks-hackathon/internal/repository"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/azmiagr/garudahacks-hackathon/pkg/database/mariadb"
	"github.com/azmiagr/garudahacks-hackathon/pkg/mail"
	"gorm.io/gorm"
)

type IOtpService interface {
	ResendOtp(param model.GetOtp) error
	ResendOtpChangePassword(param model.GetOtp) error
}

type OtpService struct {
	db             *gorm.DB
	OtpRepository  repository.IOtpRepository
	UserRepository repository.IUserRepository
}

func NewOtpService(OtpRepository repository.IOtpRepository, UserRepository repository.IUserRepository) IOtpService {
	return &OtpService{
		db:             mariadb.Connection,
		OtpRepository:  OtpRepository,
		UserRepository: UserRepository,
	}
}

func (o *OtpService) ResendOtp(param model.GetOtp) error {
	tx := o.db.Begin()
	defer tx.Rollback()

	user, err := o.UserRepository.GetUser(tx, model.GetUserParam{
		UserID: param.UserID,
	})
	if err != nil {
		return err
	}

	if user.Status == "active" {
		return errors.New("your account is already active")
	}

	otp, err := o.OtpRepository.GetOtp(tx, model.GetOtp{
		UserID: user.UserID,
	})
	if err != nil {
		return err
	}

	if otp.UpdatedAt.After(time.Now().UTC().Add(-5 * time.Minute)) {
		return errors.New("you can only resend otp every 5 minutes")
	}

	otp.Code = mail.GenerateCode()

	err = mail.SendEmail(user.Email, "OTP Verification", otp.Code)
	if err != nil {
		return err
	}

	err = o.OtpRepository.UpdateOtp(tx, otp)
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (o *OtpService) ResendOtpChangePassword(param model.GetOtp) error {
	tx := o.db.Begin()
	defer tx.Rollback()

	user, err := o.UserRepository.GetUser(tx, model.GetUserParam{
		UserID: param.UserID,
	})
	if err != nil {
		return err
	}

	otp, err := o.OtpRepository.GetOtp(tx, model.GetOtp{
		UserID: user.UserID,
	})
	if err != nil {
		return err
	}

	if otp.UpdatedAt.After(time.Now().UTC().Add(-5 * time.Minute)) {
		return errors.New("you can only resend otp every 5 minutes")
	}

	otp.Code = mail.GenerateCode()

	err = mail.SendEmail(user.Email, "Reset Password Token", "Your Reset Password Code is "+otp.Code+".")
	if err != nil {
		return err
	}

	err = o.OtpRepository.UpdateOtp(tx, otp)
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil

}
