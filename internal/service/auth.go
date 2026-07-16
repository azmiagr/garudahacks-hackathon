package service

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/internal/repository"
	"github.com/azmiagr/garudahacks-hackathon/model"
	appbcrypt "github.com/azmiagr/garudahacks-hackathon/pkg/bcrypt"
	"github.com/azmiagr/garudahacks-hackathon/pkg/database/mariadb"
	apperrors "github.com/azmiagr/garudahacks-hackathon/pkg/errors"
	"github.com/azmiagr/garudahacks-hackathon/pkg/jwt"
	"github.com/azmiagr/garudahacks-hackathon/pkg/mail"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	adminRoleName              = "admin"
	adminPoskoDisplayRoleName  = "admin_posko"
	registerOtpExpiryDuration  = 10 * time.Minute
	registerSessionDuration    = 30 * time.Minute
	registerOtpExpiryInSeconds = 600
)

type IAuthService interface {
	RequestAdminRegisterOtp(req model.RequestAdminRegisterOtpRequest) (*model.RequestAdminRegisterOtpResponse, error)
	VerifyAdminRegisterOtp(req model.VerifyAdminRegisterOtpRequest) (*model.VerifyAdminRegisterOtpResponse, error)
	SetAdminRegisterPassword(req model.SetAdminRegisterPasswordRequest) (*model.SetAdminRegisterPasswordResponse, error)
	CompleteAdminRegister(req model.CompleteAdminRegisterRequest) (*model.CompleteAdminRegisterResponse, error)
}

type AuthService struct {
	db                          *gorm.DB
	userRepository              repository.IUserRepository
	roleRepository              repository.IRoleRepository
	registrationRepository      repository.IRegistrationRepository
	adminPoskoProfileRepository repository.IAdminPoskoProfileRepository
	bcrypt                      appbcrypt.Interface
	jwtAuth                     jwt.Interface
}

func NewAuthService(
	userRepository repository.IUserRepository,
	roleRepository repository.IRoleRepository,
	registrationRepository repository.IRegistrationRepository,
	adminPoskoProfileRepository repository.IAdminPoskoProfileRepository,
	bcrypt appbcrypt.Interface,
	jwtAuth jwt.Interface,
) IAuthService {
	return &AuthService{
		db:                          mariadb.Connection,
		userRepository:              userRepository,
		roleRepository:              roleRepository,
		registrationRepository:      registrationRepository,
		adminPoskoProfileRepository: adminPoskoProfileRepository,
		bcrypt:                      bcrypt,
		jwtAuth:                     jwtAuth,
	}
}

func (s *AuthService) RequestAdminRegisterOtp(req model.RequestAdminRegisterOtpRequest) (*model.RequestAdminRegisterOtpResponse, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	email := strings.ToLower(strings.TrimSpace(req.Email))

	existingUser, err := s.userRepository.GetUser(tx, model.GetUserParam{
		Email: email,
	})
	if err == nil && existingUser != nil {
		return nil, apperrors.Conflict("email already registered")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	now := time.Now().UTC()
	sessionID := uuid.New()

	existingSession, err := s.registrationRepository.GetRegistrationSession(tx, model.GetRegistrationSessionParam{
		Email: email,
	})
	if err == nil && existingSession != nil {
		sessionID = existingSession.RegistrationID
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	otpCode := mail.GenerateCode()
	session := &entity.RegistrationSession{
		RegistrationID: sessionID,
		Email:          email,
		RoleName:       adminRoleName,
		OtpCode:        otpCode,
		OtpExpiresAt:   now.Add(registerOtpExpiryDuration),
		OtpVerifiedAt:  nil,
		PasswordHash:   "",
		ExpiresAt:      now.Add(registerSessionDuration),
	}

	err = s.registrationRepository.UpsertRegistrationSession(tx, session)
	if err != nil {
		return nil, err
	}

	err = mail.SendVerificationEmail(email, email, otpCode, "")
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &model.RequestAdminRegisterOtpResponse{
		RegistrationID:     session.RegistrationID,
		Email:              session.Email,
		OtpExpiresInSecond: registerOtpExpiryInSeconds,
	}, nil
}

func (s *AuthService) VerifyAdminRegisterOtp(req model.VerifyAdminRegisterOtpRequest) (*model.VerifyAdminRegisterOtpResponse, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	session, err := s.registrationRepository.GetRegistrationSession(tx, model.GetRegistrationSessionParam{
		RegistrationID: req.RegistrationID,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("registration session not found")
		}
		return nil, err
	}

	now := time.Now().UTC()
	if session.ExpiresAt.Before(now) {
		return nil, apperrors.BadRequest("registration session expired")
	}

	if session.OtpExpiresAt.Before(now) {
		return nil, apperrors.BadRequest("otp expired")
	}

	if session.OtpCode != req.OtpCode {
		return nil, apperrors.BadRequest("invalid otp code")
	}

	session.OtpVerifiedAt = &now
	err = s.registrationRepository.UpdateRegistrationSession(tx, session)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &model.VerifyAdminRegisterOtpResponse{
		RegistrationID: session.RegistrationID,
		OtpVerified:    true,
	}, nil
}

func (s *AuthService) SetAdminRegisterPassword(req model.SetAdminRegisterPasswordRequest) (*model.SetAdminRegisterPasswordResponse, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	session, err := s.registrationRepository.GetRegistrationSession(tx, model.GetRegistrationSessionParam{
		RegistrationID: req.RegistrationID,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("registration session not found")
		}
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now().UTC()) {
		return nil, apperrors.BadRequest("registration session expired")
	}

	if session.OtpVerifiedAt == nil {
		return nil, apperrors.BadRequest("otp must be verified before creating password")
	}

	if req.Password != req.ConfirmPassword {
		return nil, apperrors.BadRequest("password confirmation does not match")
	}

	err = validatePassword(req.Password)
	if err != nil {
		return nil, err
	}

	passwordHash, err := s.bcrypt.GenerateFromPassword(req.Password)
	if err != nil {
		return nil, err
	}

	session.PasswordHash = passwordHash
	err = s.registrationRepository.UpdateRegistrationSession(tx, session)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &model.SetAdminRegisterPasswordResponse{
		RegistrationID:  session.RegistrationID,
		PasswordCreated: true,
	}, nil
}

func (s *AuthService) CompleteAdminRegister(req model.CompleteAdminRegisterRequest) (*model.CompleteAdminRegisterResponse, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	session, err := s.registrationRepository.GetRegistrationSession(tx, model.GetRegistrationSessionParam{
		RegistrationID: req.RegistrationID,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("registration session not found")
		}
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now().UTC()) {
		return nil, apperrors.BadRequest("registration session expired")
	}

	if session.OtpVerifiedAt == nil {
		return nil, apperrors.BadRequest("otp must be verified before completing registration")
	}

	if session.PasswordHash == "" {
		return nil, apperrors.BadRequest("password must be created before completing registration")
	}

	err = validateNIK(req.NIK)
	if err != nil {
		return nil, err
	}

	existingUser, err := s.userRepository.GetUser(tx, model.GetUserParam{
		Email: session.Email,
	})
	if err == nil && existingUser != nil {
		return nil, apperrors.Conflict("email already registered")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	existingProfile, err := s.adminPoskoProfileRepository.GetAdminPoskoProfile(tx, model.GetAdminPoskoProfileParam{
		NIK: req.NIK,
	})
	if err == nil && existingProfile != nil {
		return nil, apperrors.Conflict("nik already registered")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	role, err := s.roleRepository.GetRole(tx, model.GetRoleParam{
		RoleName: adminRoleName,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.Wrap(err, http.StatusInternalServerError, "admin role is not seeded")
		}
		return nil, err
	}

	user := &entity.User{
		UserID:    uuid.New(),
		RoleID:    role.RoleID,
		Name:      strings.TrimSpace(req.FullName),
		Email:     session.Email,
		Password:  session.PasswordHash,
		Status:    "active",
		KYCStatus: "pending",
	}

	err = s.userRepository.CreateUser(tx, user)
	if err != nil {
		return nil, err
	}

	profile := &entity.AdminProfile{
		ProfileID:   uuid.New(),
		UserID:      user.UserID,
		NIK:         req.NIK,
		Affiliation: strings.TrimSpace(req.Affiliation),
	}

	err = s.adminPoskoProfileRepository.CreateAdminPoskoProfile(tx, profile)
	if err != nil {
		return nil, err
	}

	err = s.registrationRepository.DeleteRegistrationSession(tx, session)
	if err != nil {
		return nil, err
	}

	token, err := s.jwtAuth.CreateJWTToken(user.UserID, role.RoleName)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &model.CompleteAdminRegisterResponse{
		Token: token,
		User: model.RegisterUserResponse{
			UserID:      user.UserID,
			Role:        role.RoleName,
			DisplayRole: adminPoskoDisplayRoleName,
			Name:        user.Name,
			Email:       user.Email,
			Status:      user.Status,
			KYCStatus:   user.KYCStatus,
		},
	}, nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return apperrors.BadRequest("password must be at least 8 characters")
	}

	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)

	if !hasUppercase || !hasLowercase || !hasNumber || !hasSpecial {
		return apperrors.BadRequest("password must contain uppercase, lowercase, number, and special character")
	}

	return nil
}

func validateNIK(nik string) error {
	if len(nik) != 16 {
		return apperrors.BadRequest("nik must be 16 digits")
	}

	if !regexp.MustCompile(`^[0-9]+$`).MatchString(nik) {
		return apperrors.BadRequest("nik must contain only digits")
	}

	return nil
}
