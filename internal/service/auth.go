package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/internal/repository"
	"github.com/azmiagr/garudahacks-hackathon/model"
	appbcrypt "github.com/azmiagr/garudahacks-hackathon/pkg/bcrypt"
	"github.com/azmiagr/garudahacks-hackathon/pkg/database/mariadb"
	apperrors "github.com/azmiagr/garudahacks-hackathon/pkg/errors"
	"github.com/azmiagr/garudahacks-hackathon/pkg/hash"
	"github.com/azmiagr/garudahacks-hackathon/pkg/jwt"
	"github.com/azmiagr/garudahacks-hackathon/pkg/mail"
	"github.com/azmiagr/garudahacks-hackathon/pkg/supabase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	adminRoleName              = "admin"
	adminPoskoDisplayRoleName  = "admin_posko"
	donorRoleName              = "donor"
	registerOtpExpiryDuration  = 10 * time.Minute
	registerSessionDuration    = 30 * time.Minute
	registerOtpExpiryInSeconds = 600
	maxStoreKTPImageSize       = 5 * 1024 * 1024
)

type IAuthService interface {
	Login(req model.LoginRequest) (*model.LoginResponse, error)
	RequestRegisterOtp(req model.RequestAdminRegisterOtpRequest) (*model.RequestAdminRegisterOtpResponse, error)
	RequestAdminRegisterOtp(req model.RequestAdminRegisterOtpRequest) (*model.RequestAdminRegisterOtpResponse, error)
	VerifyAdminRegisterOtp(req model.VerifyAdminRegisterOtpRequest) (*model.VerifyAdminRegisterOtpResponse, error)
	SetAdminRegisterPassword(req model.SetAdminRegisterPasswordRequest) (*model.SetAdminRegisterPasswordResponse, error)
	CompleteAdminRegister(req model.CompleteAdminRegisterRequest) (*model.CompleteAdminRegisterResponse, error)
	CompleteDonorRegister(req model.CompleteDonorRegisterRequest) (*model.CompleteDonorRegisterResponse, error)
	CompleteCourierRegister(req model.CompleteCourierRegisterRequest) (*model.CompleteCourierRegisterResponse, error)
	CompleteStoreRegister(req model.CompleteStoreRegisterRequest) (*model.CompleteStoreRegisterResponse, error)
	Logout(token string) (*model.LogoutResponse, error)
	IsTokenRevoked(token string) (bool, error)
}

type AuthService struct {
	db                          *gorm.DB
	userRepository              repository.IUserRepository
	roleRepository              repository.IRoleRepository
	registrationRepository      repository.IRegistrationRepository
	adminPoskoProfileRepository repository.IAdminPoskoProfileRepository
	donorProfileRepository      repository.IDonorProfileRepository
	courierProfileRepository    repository.ICourierProfileRepository
	revokedTokenRepository      repository.IRevokedTokenRepository
	bcrypt                      appbcrypt.Interface
	jwtAuth                     jwt.Interface
	hasher                      hash.Interface
	storage                     supabase.Interface
	storeRepository             repository.IStoreRepository
}

func NewAuthService(
	userRepository repository.IUserRepository,
	roleRepository repository.IRoleRepository,
	registrationRepository repository.IRegistrationRepository,
	adminPoskoProfileRepository repository.IAdminPoskoProfileRepository,
	donorProfileRepository repository.IDonorProfileRepository,
	courierProfileRepository repository.ICourierProfileRepository,
	revokedTokenRepository repository.IRevokedTokenRepository,
	bcrypt appbcrypt.Interface,
	jwtAuth jwt.Interface,
	hasher hash.Interface,
	storage supabase.Interface,
	storeRepository repository.IStoreRepository,
) IAuthService {
	return &AuthService{
		db:                          mariadb.Connection,
		userRepository:              userRepository,
		roleRepository:              roleRepository,
		registrationRepository:      registrationRepository,
		adminPoskoProfileRepository: adminPoskoProfileRepository,
		donorProfileRepository:      donorProfileRepository,
		courierProfileRepository:    courierProfileRepository,
		revokedTokenRepository:      revokedTokenRepository,
		bcrypt:                      bcrypt,
		jwtAuth:                     jwtAuth,
		hasher:                      hasher,
		storage:                     storage,
		storeRepository:             storeRepository,
	}
}

const storeRoleName = "store"
const courierRoleName = "relawan"

func (s *AuthService) Login(req model.LoginRequest) (*model.LoginResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.userRepository.GetUser(s.db, model.GetUserParam{
		Email: email,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.Unauthorized("invalid email or password")
		}
		return nil, err
	}

	err = s.bcrypt.CompareAndHashPassword(user.Password, req.Password)
	if err != nil {
		return nil, apperrors.Unauthorized("invalid email or password")
	}

	if user.Status != "active" {
		return nil, apperrors.Unauthorized("your account has been deactivated. Please contact administrator")
	}

	role, err := s.roleRepository.GetRole(s.db, model.GetRoleParam{
		RoleID: user.RoleID,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.Wrap(err, http.StatusInternalServerError, "user role is not found")
		}
		return nil, err
	}

	token, err := s.jwtAuth.CreateJWTToken(user.UserID, role.RoleName)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token: token,
	}, nil
}

func (s *AuthService) Logout(token string) (*model.LogoutResponse, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, apperrors.Unauthorized("token is required")
	}

	expiresAt, err := s.jwtAuth.GetTokenExpiresAt(token)
	if err != nil {
		return nil, apperrors.Unauthorized("failed to validate token")
	}

	now := time.Now().UTC()
	if !expiresAt.After(now) {
		return nil, apperrors.Unauthorized("token is expired")
	}

	tx := s.db.Begin()
	defer tx.Rollback()

	err = s.revokedTokenRepository.DeleteExpiredTokens(tx, now)
	if err != nil {
		return nil, err
	}

	err = s.revokedTokenRepository.CreateRevokedToken(tx, &entity.RevokedToken{
		TokenHash: hashToken(token),
		ExpiresAt: expiresAt.UTC(),
	})
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &model.LogoutResponse{LoggedOut: true}, nil
}

func (s *AuthService) IsTokenRevoked(token string) (bool, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return false, nil
	}

	now := time.Now().UTC()
	return s.revokedTokenRepository.ExistsActiveTokenHash(s.db, hashToken(token), now)
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (s *AuthService) RequestAdminRegisterOtp(req model.RequestAdminRegisterOtpRequest) (*model.RequestAdminRegisterOtpResponse, error) {
	if strings.TrimSpace(req.Role) == "" {
		req.Role = adminRoleName
	}

	return s.RequestRegisterOtp(req)
}

func (s *AuthService) RequestRegisterOtp(req model.RequestAdminRegisterOtpRequest) (*model.RequestAdminRegisterOtpResponse, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	email := strings.ToLower(strings.TrimSpace(req.Email))
	roleName, err := normalizeRegisterRole(req.Role)
	if err != nil {
		return nil, err
	}

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
		Email:    email,
		RoleName: roleName,
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
		RoleName:       roleName,
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

	if session.RoleName != adminRoleName {
		return nil, apperrors.BadRequest("registration session is not for admin")
	}

	err = validateNIK(req.NIK)
	if err != nil {
		return nil, err
	}

	hashedNIK := s.hasher.HashNIK(req.NIK)

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
		NIK: hashedNIK,
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
		KYCStatus: "approved",
	}

	err = s.userRepository.CreateUser(tx, user)
	if err != nil {
		return nil, err
	}

	profile := &entity.AdminProfile{
		ProfileID:   uuid.New(),
		UserID:      user.UserID,
		NIK:         hashedNIK,
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
			DisplayRole: resolveDisplayRole(role.RoleName),
			Name:        user.Name,
			Email:       user.Email,
			Status:      user.Status,
			KYCStatus:   user.KYCStatus,
		},
	}, nil
}

func (s *AuthService) CompleteDonorRegister(req model.CompleteDonorRegisterRequest) (*model.CompleteDonorRegisterResponse, error) {
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

	if session.RoleName != donorRoleName {
		return nil, apperrors.BadRequest("registration session is not for donor")
	}

	if strings.TrimSpace(req.FullName) == "" {
		return nil, apperrors.BadRequest("full name is required")
	}

	phoneNumber := normalizePhoneNumber(req.PhoneNumber)
	if phoneNumber == "" {
		return nil, apperrors.BadRequest("phone number is required")
	}

	if !req.ConsentAccepted {
		return nil, apperrors.BadRequest("consent must be accepted")
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

	existingProfile, err := s.donorProfileRepository.GetDonorProfile(tx, model.GetDonorProfileParam{
		PhoneNumber: phoneNumber,
	})
	if err == nil && existingProfile != nil {
		return nil, apperrors.Conflict("phone number already registered")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	role, err := s.roleRepository.GetRole(tx, model.GetRoleParam{
		RoleName: donorRoleName,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.Wrap(err, http.StatusInternalServerError, "donor role is not seeded")
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
		KYCStatus: "approved",
	}

	err = s.userRepository.CreateUser(tx, user)
	if err != nil {
		return nil, err
	}

	preferenceJSON, err := json.Marshal(req.DonationPreferences)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	profile := &entity.DonorProfile{
		ProfileID:         uuid.New(),
		UserID:            user.UserID,
		PhoneNumber:       phoneNumber,
		PreferenceJSON:    string(preferenceJSON),
		ConsentAccepted:   true,
		ConsentAcceptedAt: now,
	}

	err = s.donorProfileRepository.CreateDonorProfile(tx, profile)
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

	return &model.CompleteDonorRegisterResponse{
		Token: token,
		User: model.RegisterUserResponse{
			UserID:      user.UserID,
			Role:        role.RoleName,
			DisplayRole: resolveDisplayRole(role.RoleName),
			Name:        user.Name,
			Email:       user.Email,
			Status:      user.Status,
			KYCStatus:   user.KYCStatus,
		},
	}, nil
}

func (s *AuthService) CompleteCourierRegister(req model.CompleteCourierRegisterRequest) (*model.CompleteCourierRegisterResponse, error) {
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

	if session.RoleName != courierRoleName {
		return nil, apperrors.BadRequest("registration session is not for courier")
	}

	fullName := strings.TrimSpace(req.FullName)
	if fullName == "" {
		return nil, apperrors.BadRequest("full name is required")
	}

	err = validateNIK(req.NIK)
	if err != nil {
		return nil, err
	}

	vehicleType, err := normalizeCourierVehicleType(req.VehicleType)
	if err != nil {
		return nil, err
	}

	vehicleCapacityKG := req.VehicleCapacityKG
	if vehicleCapacityKG < 0 {
		return nil, apperrors.BadRequest("vehicle capacity must not be negative")
	}
	if vehicleCapacityKG == 0 {
		vehicleCapacityKG = defaultCourierVehicleCapacityKG(vehicleType)
	}

	operationalArea := strings.TrimSpace(req.OperationalArea)
	if operationalArea == "" {
		return nil, apperrors.BadRequest("operational area is required")
	}

	if req.OperationRadiusKM <= 0 {
		return nil, apperrors.BadRequest("operation radius must be greater than 0")
	}

	if !req.WaiverAccepted {
		return nil, apperrors.BadRequest("waiver must be accepted")
	}

	hashedNIK := s.hasher.HashNIK(req.NIK)

	existingUser, err := s.userRepository.GetUser(tx, model.GetUserParam{
		Email: session.Email,
	})
	if err == nil && existingUser != nil {
		return nil, apperrors.Conflict("email already registered")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	existingProfile, err := s.courierProfileRepository.GetCourierProfile(tx, model.GetCourierProfileParam{
		NIK: hashedNIK,
	})
	if err == nil && existingProfile != nil {
		return nil, apperrors.Conflict("nik already registered")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	role, err := s.roleRepository.GetRole(tx, model.GetRoleParam{
		RoleName: courierRoleName,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.Wrap(err, http.StatusInternalServerError, "courier role is not seeded")
		}
		return nil, err
	}

	user := &entity.User{
		UserID:    uuid.New(),
		RoleID:    role.RoleID,
		Name:      fullName,
		Email:     session.Email,
		Password:  session.PasswordHash,
		Status:    "active",
		KYCStatus: "approved",
	}

	err = s.userRepository.CreateUser(tx, user)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	profile := &entity.CourierProfile{
		ProfileID:                     uuid.New(),
		UserID:                        user.UserID,
		NIK:                           hashedNIK,
		VehicleType:                   vehicleType,
		VehicleCapacityKG:             vehicleCapacityKG,
		OperationalArea:               operationalArea,
		OperationRadiusKM:             req.OperationRadiusKM,
		WaiverAccepted:                true,
		WaiverAcceptedAt:              now,
		IsAvailable:                   true,
		UrgentTaskNotificationEnabled: true,
	}

	err = s.courierProfileRepository.CreateCourierProfile(tx, profile)
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

	return &model.CompleteCourierRegisterResponse{
		Token: token,
		User: model.RegisterUserResponse{
			UserID:      user.UserID,
			Role:        role.RoleName,
			DisplayRole: resolveDisplayRole(role.RoleName),
			Name:        user.Name,
			Email:       user.Email,
			Status:      user.Status,
			KYCStatus:   user.KYCStatus,
		},
		Courier: model.CourierRegisterResponse{
			ProfileID:         profile.ProfileID,
			UserID:            profile.UserID,
			VehicleType:       profile.VehicleType,
			VehicleCapacityKG: profile.VehicleCapacityKG,
			OperationalArea:   profile.OperationalArea,
			OperationRadiusKM: profile.OperationRadiusKM,
			WaiverAccepted:    profile.WaiverAccepted,
		},
	}, nil
}

func (s *AuthService) CompleteStoreRegister(req model.CompleteStoreRegisterRequest) (*model.CompleteStoreRegisterResponse, error) {
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

	if session.RoleName != storeRoleName {
		return nil, apperrors.BadRequest("registration session is not for store")
	}

	storeName := strings.TrimSpace(req.StoreName)
	ownerName := strings.TrimSpace(req.OwnerName)
	nib := normalizeBusinessNumber(req.NIB)

	if storeName == "" {
		return nil, apperrors.BadRequest("store name is required")
	}
	if ownerName == "" {
		return nil, apperrors.BadRequest("owner name is required")
	}
	if nib == "" {
		return nil, apperrors.BadRequest("nib is required")
	}
	if strings.TrimSpace(req.Address) == "" {
		return nil, apperrors.BadRequest("address is required")
	}
	if req.KTPImage == nil {
		return nil, apperrors.BadRequest("ktp_image is required")
	}
	if req.KTPImage.Size > maxStoreKTPImageSize {
		return nil, apperrors.BadRequest("ktp_image size must not exceed 5MB")
	}

	var ktpImageURL string
	if strings.EqualFold(filepath.Ext(req.KTPImage.Filename), ".pdf") {
		ktpImageURL, err = s.storage.UploadPDF(req.KTPImage)
	} else {
		ktpImageURL, err = s.storage.UploadFile(req.KTPImage)
	}
	if err != nil {
		return nil, err
	}

	committed := false
	defer func() {
		if !committed && ktpImageURL != "" {
			_ = supabase.DeleteFileIfPresent(s.storage, ktpImageURL)
		}
	}()

	existingUser, err := s.userRepository.GetUser(tx, model.GetUserParam{
		Email: session.Email,
	})
	if err == nil && existingUser != nil {
		return nil, apperrors.Conflict("email already registered")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	existingStore, err := s.storeRepository.GetStore(tx, model.GetStoreParam{
		BusinessNumber: nib,
	})
	if err == nil && existingStore != nil {
		return nil, apperrors.Conflict("nib already registered")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	role, err := s.roleRepository.GetRole(tx, model.GetRoleParam{
		RoleName: storeRoleName,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.Wrap(err, http.StatusInternalServerError, "store role is not seeded")
		}
		return nil, err
	}

	user := &entity.User{
		UserID:    uuid.New(),
		RoleID:    role.RoleID,
		Name:      ownerName,
		Email:     session.Email,
		Password:  session.PasswordHash,
		Status:    "active",
		KYCStatus: "pending",
	}

	err = s.userRepository.CreateUser(tx, user)
	if err != nil {
		return nil, err
	}

	categories, err := normalizeStoreCategories(req.Categories, req.CategoriesJSON)
	if err != nil {
		return nil, err
	}

	categoriesJSON, err := json.Marshal(categories)
	if err != nil {
		return nil, err
	}

	store := &entity.Stores{
		StoreID:         uuid.New(),
		OwnerID:         user.UserID,
		Name:            storeName,
		BusinessNumber:  nib,
		NPWP:            strings.TrimSpace(req.NPWP),
		KTPImageURL:     ktpImageURL,
		BankName:        strings.TrimSpace(req.BankName),
		BankAccountNo:   strings.TrimSpace(req.BankAccountNo),
		BankAccountName: strings.TrimSpace(req.BankAccountName),
		CategoriesJSON:  string(categoriesJSON),
		Address:         strings.TrimSpace(req.Address),
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
	}

	err = s.storeRepository.CreateStore(tx, store)
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
	committed = true

	return &model.CompleteStoreRegisterResponse{
		Token: token,
		User: model.RegisterUserResponse{
			UserID:      user.UserID,
			Role:        role.RoleName,
			DisplayRole: resolveDisplayRole(role.RoleName),
			Name:        user.Name,
			Email:       user.Email,
			Status:      user.Status,
			KYCStatus:   user.KYCStatus,
		},
		Store: model.StoreRegisterResponse{
			StoreID:        store.StoreID,
			OwnerID:        store.OwnerID,
			Name:           store.Name,
			BusinessNumber: store.BusinessNumber,
			Address:        store.Address,
			Latitude:       store.Latitude,
			Longitude:      store.Longitude,
		},
	}, nil
}

func normalizeBusinessNumber(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, " ", "")
	value = strings.ReplaceAll(value, "-", "")
	value = strings.ReplaceAll(value, ".", "")
	return value
}

func normalizeStoreCategories(values []string, rawJSON string) ([]string, error) {
	if strings.TrimSpace(rawJSON) != "" {
		var parsed []string
		if err := json.Unmarshal([]byte(rawJSON), &parsed); err != nil {
			return nil, apperrors.BadRequest("categories_json must be a valid JSON array")
		}
		values = append(values, parsed...)
	}

	normalized := make([]string, 0, len(values))
	seen := map[string]struct{}{}
	for _, value := range values {
		for _, part := range strings.Split(value, ",") {
			category := strings.TrimSpace(part)
			if category == "" {
				continue
			}

			key := strings.ToLower(category)
			if _, exists := seen[key]; exists {
				continue
			}

			seen[key] = struct{}{}
			normalized = append(normalized, category)
		}
	}

	return normalized, nil
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

func resolveDisplayRole(roleName string) string {
	if roleName == adminRoleName {
		return adminPoskoDisplayRoleName
	}

	if roleName == storeRoleName {
		return "toko_mitra"
	}

	if roleName == courierRoleName {
		return "relawan_kurir"
	}

	return roleName
}

func normalizeCourierVehicleType(vehicleType string) (string, error) {
	vehicleType = strings.ToLower(strings.TrimSpace(vehicleType))
	vehicleType = strings.ReplaceAll(vehicleType, "_", "-")

	switch vehicleType {
	case "motor", "mobil", "pick-up", "pickup":
		if vehicleType == "pickup" {
			return "pick-up", nil
		}
		return vehicleType, nil
	default:
		return "", apperrors.BadRequest("unsupported vehicle type")
	}
}

func defaultCourierVehicleCapacityKG(vehicleType string) int {
	switch vehicleType {
	case "motor":
		return 50
	case "mobil":
		return 300
	case "pick-up":
		return 1000
	default:
		return 0
	}
}

func normalizeRegisterRole(roleName string) (string, error) {
	roleName = strings.ToLower(strings.TrimSpace(roleName))
	if roleName == "" {
		return adminRoleName, nil
	}

	switch roleName {
	case adminRoleName, adminPoskoDisplayRoleName:
		return adminRoleName, nil
	case donorRoleName, "donatur":
		return donorRoleName, nil
	case "store", "toko_mitra":
		return "store", nil
	case "relawan", "courier", "kurir":
		return "relawan", nil
	default:
		return "", apperrors.BadRequest("unsupported registration role")
	}
}

func normalizePhoneNumber(phoneNumber string) string {
	phoneNumber = strings.TrimSpace(phoneNumber)
	phoneNumber = strings.ReplaceAll(phoneNumber, " ", "")
	phoneNumber = strings.ReplaceAll(phoneNumber, "-", "")

	return phoneNumber
}
