package service

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/internal/repository"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/azmiagr/garudahacks-hackathon/pkg/database/mariadb"
	apperrors "github.com/azmiagr/garudahacks-hackathon/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	custodyTokenTTLSeconds         = 30
	custodyTokenCacheWindowSeconds = 90
	custodyDevelopmentSecret       = "development-custody-handshake-secret"
)

type IStoreCustodyService interface {
	GetStoreOrders(user *entity.User, param model.StoreOrderListParam) (*model.StoreOrderListResponse, error)
	GetStoreOrderDetail(user *entity.User, orderID uuid.UUID) (*model.StoreOrderDetailResponse, error)
	AcceptOrder(user *entity.User, req model.StoreOrderActionRequest) (*model.StoreOrderActionResponse, error)
	MarkOrderReady(user *entity.User, req model.StoreOrderActionRequest) (*model.StoreHandoffTokenResponse, error)
	GenerateStoreHandoffToken(user *entity.User, req model.StoreOrderActionRequest) (*model.StoreHandoffTokenResponse, error)
	SubmitHandoff(user *entity.User, req model.SubmitCustodyHandshakeRequest) (*model.CustodyHandshakeResponse, error)
}

type StoreCustodyService struct {
	db                       *gorm.DB
	orderRepository          repository.IOrderRepository
	storeRepository          repository.IStoreRepository
	custodyLogRepository     repository.ICustodyLogRepository
	handshakeTokenRepository repository.ICustodyHandshakeTokenRepository
}

func NewStoreCustodyService(
	orderRepository repository.IOrderRepository,
	storeRepository repository.IStoreRepository,
	custodyLogRepository repository.ICustodyLogRepository,
	handshakeTokenRepository repository.ICustodyHandshakeTokenRepository,
) IStoreCustodyService {
	return &StoreCustodyService{
		db:                       mariadb.Connection,
		orderRepository:          orderRepository,
		storeRepository:          storeRepository,
		custodyLogRepository:     custodyLogRepository,
		handshakeTokenRepository: handshakeTokenRepository,
	}
}

func (s *StoreCustodyService) GetStoreOrders(user *entity.User, param model.StoreOrderListParam) (*model.StoreOrderListResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	store, err := s.getOwnedStore(s.db, user.UserID)
	if err != nil {
		return nil, err
	}

	status := normalizeStoreOrderStatus(param.Status)
	limit := normalizeStoreOrderServiceLimit(param.Limit)
	offset := normalizeStoreOrderServiceOffset(param.Offset)
	rows, err := s.orderRepository.GetStoreOrders(s.db, model.StoreOrderListRepositoryParam{
		StoreID: store.StoreID,
		Status:  status,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, err
	}

	return &model.StoreOrderListResponse{
		Items:  buildStoreOrderListItems(rows),
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *StoreCustodyService) GetStoreOrderDetail(user *entity.User, orderID uuid.UUID) (*model.StoreOrderDetailResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	store, err := s.getOwnedStore(s.db, user.UserID)
	if err != nil {
		return nil, err
	}

	row, err := s.orderRepository.GetStoreOrderDetail(s.db, model.StoreOrderDetailRepositoryParam{
		OrderID: orderID,
		StoreID: store.StoreID,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("order not found")
		}
		return nil, err
	}

	items, err := s.orderRepository.GetStoreOrderItems(s.db, orderID)
	if err != nil {
		return nil, err
	}

	return buildStoreOrderDetailResponse(*row, items), nil
}

func (s *StoreCustodyService) AcceptOrder(user *entity.User, req model.StoreOrderActionRequest) (*model.StoreOrderActionResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	tx := s.db.Begin()
	defer tx.Rollback()

	store, err := s.getOwnedStore(tx, user.UserID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	err = s.orderRepository.AcceptOrderForStore(tx, req.OrderID, store.StoreID, now)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.Conflict("order is not available for this store")
		}
		return nil, err
	}

	order, err := s.orderRepository.GetOrder(tx, req.OrderID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return buildStoreOrderActionResponse(order), nil
}

func (s *StoreCustodyService) MarkOrderReady(user *entity.User, req model.StoreOrderActionRequest) (*model.StoreHandoffTokenResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	tx := s.db.Begin()
	defer tx.Rollback()

	store, err := s.getOwnedStore(tx, user.UserID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	err = s.orderRepository.MarkReadyForPickup(tx, req.OrderID, store.StoreID, now)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.BadRequest("order cannot be marked ready by this store")
		}
		return nil, err
	}

	resp, err := s.createStoreHandoffToken(tx, req.OrderID, store.OwnerID, now)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *StoreCustodyService) GenerateStoreHandoffToken(user *entity.User, req model.StoreOrderActionRequest) (*model.StoreHandoffTokenResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	tx := s.db.Begin()
	defer tx.Rollback()

	store, err := s.getOwnedStore(tx, user.UserID)
	if err != nil {
		return nil, err
	}

	order, err := s.orderRepository.GetOrderForUpdate(tx, req.OrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("order not found")
		}
		return nil, err
	}
	if order.StoreID != store.StoreID {
		return nil, apperrors.Forbidden("order does not belong to this store")
	}
	if order.OrderStatus != entity.OrderStatusReadyForPickup {
		return nil, apperrors.BadRequest("order is not ready for pickup")
	}

	resp, err := s.createStoreHandoffToken(tx, order.OrderID, store.OwnerID, time.Now().UTC())
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *StoreCustodyService) SubmitHandoff(user *entity.User, req model.SubmitCustodyHandshakeRequest) (*model.CustodyHandshakeResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	method := strings.ToLower(strings.TrimSpace(req.Method))
	if method != entity.HandshakeMethodQR && method != entity.HandshakeMethodPIN {
		return nil, apperrors.BadRequest("handshake method must be qr or pin")
	}

	tx := s.db.Begin()
	defer tx.Rollback()

	if strings.TrimSpace(req.IdempotencyKey) != "" {
		exists, err := s.custodyLogRepository.ExistsIdempotencyKey(tx, req.IdempotencyKey)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, apperrors.Conflict("handshake idempotency key has already been used")
		}
	}

	token, err := s.getHandshakeTokenForUpdate(tx, req, method)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.BadRequest("handshake token is invalid or expired")
		}
		return nil, err
	}
	if token.HandoffStage != entity.CustodyStageStoreToCourier {
		return nil, apperrors.BadRequest("handshake token stage is not supported")
	}

	order, err := s.orderRepository.GetOrderForUpdate(tx, token.OrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("order not found")
		}
		return nil, err
	}
	if order.OrderStatus != entity.OrderStatusReadyForPickup {
		return nil, apperrors.BadRequest("order is not ready for pickup")
	}

	store, err := s.storeRepository.GetStore(tx, model.GetStoreParam{StoreID: order.StoreID})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("store not found")
		}
		return nil, err
	}
	if store.OwnerID != token.PresentedBy {
		return nil, apperrors.BadRequest("handshake token does not match order store")
	}

	if order.CourierID != uuid.Nil && order.CourierID != user.UserID {
		return nil, apperrors.Conflict("order is already assigned to another courier")
	}

	capturedAt := req.CapturedAt.UTC()
	if capturedAt.IsZero() {
		capturedAt = time.Now().UTC()
	}

	log, err := s.appendCustodyLog(tx, appendCustodyLogParam{
		OrderID:           order.OrderID,
		HandoffStage:      entity.CustodyStageStoreToCourier,
		HandshakeMethod:   method,
		FromActorID:       store.OwnerID,
		ToActorID:         user.UserID,
		ScannedBy:         user.UserID,
		Latitude:          req.Latitude,
		Longitude:         req.Longitude,
		IdempotencyKey:    strings.TrimSpace(req.IdempotencyKey),
		CapturedAt:        capturedAt,
		GPSDistanceMeters: calculateDistanceMeters(req.Latitude, req.Longitude, store.Latitude, store.Longitude),
	})
	if err != nil {
		return nil, err
	}

	err = s.handshakeTokenRepository.MarkTokenUsed(tx, token.TokenID, user.UserID, capturedAt)
	if err != nil {
		return nil, err
	}

	order.CourierID = user.UserID
	order.OrderStatus = entity.OrderStatusInTransit
	order.PickedUpAt = &capturedAt
	order.UpdatedAt = time.Now().UTC()
	if err := s.orderRepository.UpdateOrder(tx, order); err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &model.CustodyHandshakeResponse{
		OrderID:          order.OrderID,
		LogID:            log.LogsID,
		OrderStatus:      order.OrderStatus,
		HandoffStage:     log.HandoffStage,
		HandshakeMethod:  log.HandshakeMethod,
		Sequence:         log.Sequence,
		CurrentHash:      log.CurrentHash,
		ShortCurrentHash: shortenLedgerHash(log.CurrentHash),
		CapturedAt:       capturedAt,
	}, nil
}

func (s *StoreCustodyService) getOwnedStore(tx *gorm.DB, ownerID uuid.UUID) (*entity.Stores, error) {
	store, err := s.storeRepository.GetStore(tx, model.GetStoreParam{OwnerID: ownerID})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.Forbidden("store profile is required")
		}
		return nil, err
	}

	return store, nil
}

func (s *StoreCustodyService) createStoreHandoffToken(tx *gorm.DB, orderID uuid.UUID, presentedBy uuid.UUID, now time.Time) (*model.StoreHandoffTokenResponse, error) {
	tokenID := uuid.New()
	nonce := strings.ReplaceAll(uuid.NewString(), "-", "")
	expiresAt := now.Add(custodyTokenTTLSeconds * time.Second)
	cacheValidUntil := now.Add(custodyTokenCacheWindowSeconds * time.Second)
	pin, err := generateNumericPIN(6)
	if err != nil {
		return nil, err
	}

	payload := buildStoreHandoffQRPayload(tokenID, orderID, presentedBy, nonce, expiresAt)
	token := &entity.CustodyHandshakeToken{
		TokenID:         tokenID,
		OrderID:         orderID,
		HandoffStage:    entity.CustodyStageStoreToCourier,
		PresentedBy:     presentedBy,
		QRPayloadHash:   hashTokenValue(payload),
		PINHash:         hashTokenValue(pin),
		Nonce:           nonce,
		Status:          entity.CustodyTokenStatusActive,
		ExpiresAt:       expiresAt,
		CacheValidUntil: cacheValidUntil,
		CreatedAt:       now,
	}
	if err := s.handshakeTokenRepository.CreateToken(tx, token); err != nil {
		return nil, err
	}

	return &model.StoreHandoffTokenResponse{
		OrderID:              orderID,
		TokenID:              tokenID,
		HandoffStage:         entity.CustodyStageStoreToCourier,
		QRPayload:            payload,
		FallbackPIN:          pin,
		ExpiresAt:            expiresAt,
		CacheValidUntil:      cacheValidUntil,
		RefreshInSeconds:     custodyTokenTTLSeconds,
		CacheWindowInSeconds: custodyTokenCacheWindowSeconds,
	}, nil
}

func (s *StoreCustodyService) getHandshakeTokenForUpdate(tx *gorm.DB, req model.SubmitCustodyHandshakeRequest, method string) (*entity.CustodyHandshakeToken, error) {
	now := time.Now().UTC()
	if method == entity.HandshakeMethodQR {
		payload := strings.TrimSpace(req.QRPayload)
		if payload == "" {
			return nil, apperrors.BadRequest("qr_payload is required")
		}
		return s.handshakeTokenRepository.GetActiveTokenByQRHashForUpdate(tx, hashTokenValue(payload), now)
	}

	pin := strings.TrimSpace(req.FallbackPIN)
	if pin == "" || req.OrderID == uuid.Nil {
		return nil, apperrors.BadRequest("order_id and fallback_pin are required")
	}
	return s.handshakeTokenRepository.GetActiveTokenByPINHashForUpdate(tx, req.OrderID, entity.CustodyStageStoreToCourier, hashTokenValue(pin), now)
}

type appendCustodyLogParam struct {
	OrderID           uuid.UUID
	HandoffStage      string
	HandshakeMethod   string
	FromActorID       uuid.UUID
	ToActorID         uuid.UUID
	ScannedBy         uuid.UUID
	Latitude          float64
	Longitude         float64
	IdempotencyKey    string
	CapturedAt        time.Time
	GPSDistanceMeters float64
}

func (s *StoreCustodyService) appendCustodyLog(tx *gorm.DB, param appendCustodyLogParam) (*entity.CustodyLogs, error) {
	latestLog, err := s.custodyLogRepository.GetLatestCustodyLogByOrderForUpdate(tx, param.OrderID)

	sequence := 1
	prevHash := "GENESIS"
	if err == nil && latestLog != nil {
		sequence = latestLog.Sequence + 1
		prevHash = latestLog.CurrentHash
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	idempotencyKey := strings.TrimSpace(param.IdempotencyKey)
	var idempotencyKeyPtr *string
	if idempotencyKey != "" {
		idempotencyKeyPtr = &idempotencyKey
	}

	scannedBy := param.ScannedBy
	gpsDistanceMeters := param.GPSDistanceMeters
	log := &entity.CustodyLogs{
		LogsID:            uuid.NewString(),
		OrderID:           param.OrderID.String(),
		HandoffStage:      param.HandoffStage,
		HandshakeMethod:   param.HandshakeMethod,
		FromActorID:       param.FromActorID,
		ToActorID:         param.ToActorID,
		ScannedBy:         &scannedBy,
		Sequence:          sequence,
		Latitude:          param.Latitude,
		Longitude:         param.Longitude,
		GPSDistanceMeters: &gpsDistanceMeters,
		IsGPSAnomaly:      gpsDistanceMeters > 300,
		IdempotencyKey:    idempotencyKeyPtr,
		PrevHash:          prevHash,
		CapturedAt:        &param.CapturedAt,
		CreatedAt:         time.Now().UTC(),
	}
	log.CurrentHash = buildCustodyServiceHash(*log)

	if err := s.custodyLogRepository.CreateCustodyLog(tx, log); err != nil {
		return nil, err
	}

	return log, nil
}

func buildStoreOrderActionResponse(order *entity.Orders) *model.StoreOrderActionResponse {
	return &model.StoreOrderActionResponse{
		OrderID:     order.OrderID,
		StoreID:     order.StoreID,
		OrderStatus: order.OrderStatus,
		UpdatedAt:   order.UpdatedAt,
	}
}

func buildStoreOrderListItems(rows []model.StoreOrderRow) []model.StoreOrderListItem {
	items := make([]model.StoreOrderListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, buildStoreOrderListItem(row))
	}
	return items
}

func buildStoreOrderListItem(row model.StoreOrderRow) model.StoreOrderListItem {
	return model.StoreOrderListItem{
		OrderID:      row.OrderID,
		OrderCode:    row.OrderCode,
		OrderStatus:  row.OrderStatus,
		TotalAmount:  row.TotalAmount,
		RequestTitle: row.RequestTitle,
		PostName:     row.PostName,
		PostAddress:  row.PostAddress,
		StoreName:    row.StoreName,
		CourierName:  row.CourierName,
		UpdatedAt:    row.UpdatedAt,
	}
}

func buildStoreOrderDetailResponse(row model.StoreOrderRow, itemRows []model.StoreOrderItemRow) *model.StoreOrderDetailResponse {
	items := make([]model.StoreOrderItemItem, 0, len(itemRows))
	for _, itemRow := range itemRows {
		items = append(items, model.StoreOrderItemItem{
			ItemID:    itemRow.ItemID,
			Name:      itemRow.Name,
			Quantity:  itemRow.Quantity,
			Unit:      itemRow.Unit,
			UnitPrice: itemRow.UnitPrice,
			Subtotal:  itemRow.Subtotal,
		})
	}

	return &model.StoreOrderDetailResponse{
		StoreOrderListItem: buildStoreOrderListItem(row),
		RequestID:          row.RequestID,
		StoreID:            row.StoreID,
		CourierID:          row.CourierID,
		PostLatitude:       row.PostLatitude,
		PostLongitude:      row.PostLongitude,
		AcceptedAt:         row.AcceptedAt,
		ReadyAt:            row.ReadyAt,
		PickedUpAt:         row.PickedUpAt,
		CreatedAt:          row.CreatedAt,
		Items:              items,
	}
}

func normalizeStoreOrderStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "mine", "accepted", "ready", "in_transit":
		return strings.ToLower(strings.TrimSpace(status))
	default:
		return "available"
	}
}

func normalizeStoreOrderServiceLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func normalizeStoreOrderServiceOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}

func buildStoreHandoffQRPayload(tokenID uuid.UUID, orderID uuid.UUID, presentedBy uuid.UUID, nonce string, expiresAt time.Time) string {
	body := map[string]string{
		"token_id":     tokenID.String(),
		"order_id":     orderID.String(),
		"presented_by": presentedBy.String(),
		"stage":        entity.CustodyStageStoreToCourier,
		"nonce":        nonce,
		"expires_at":   expiresAt.UTC().Format(time.RFC3339),
	}
	signature := signCustodyPayload(body)
	body["signature"] = signature

	raw, _ := json.Marshal(body)
	return string(raw)
}

func signCustodyPayload(body map[string]string) string {
	payload := fmt.Sprintf("%s|%s|%s|%s|%s|%s",
		body["token_id"],
		body["order_id"],
		body["presented_by"],
		body["stage"],
		body["nonce"],
		body["expires_at"],
	)
	mac := hmac.New(sha256.New, []byte(custodyHandshakeSecret()))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

func custodyHandshakeSecret() string {
	secret := strings.TrimSpace(os.Getenv("CUSTODY_HANDSHAKE_SECRET"))
	if secret == "" {
		return custodyDevelopmentSecret
	}
	return secret
}

func hashTokenValue(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func generateNumericPIN(length int) (string, error) {
	var builder strings.Builder
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		builder.WriteString(n.String())
	}
	return builder.String(), nil
}

func buildCustodyServiceHash(log entity.CustodyLogs) string {
	capturedAt := ""
	if log.CapturedAt != nil {
		capturedAt = log.CapturedAt.UTC().Format(time.RFC3339Nano)
	}
	idempotencyKey := ""
	if log.IdempotencyKey != nil {
		idempotencyKey = *log.IdempotencyKey
	}
	payload := fmt.Sprintf(
		"%s|%s|%s|%s|%s|%s|%d|%.8f|%.8f|%s|%s|%s|%s",
		log.LogsID,
		log.OrderID,
		log.HandoffStage,
		log.HandshakeMethod,
		log.FromActorID.String(),
		log.ToActorID.String(),
		log.Sequence,
		log.Latitude,
		log.Longitude,
		idempotencyKey,
		log.PrevHash,
		capturedAt,
		log.CreatedAt.UTC().Format(time.RFC3339Nano),
	)

	sum := sha256.Sum256([]byte(payload))
	return "0x" + hex.EncodeToString(sum[:])
}

func calculateDistanceMeters(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	const earthRadiusMeters = 6371000
	degToRad := func(deg float64) float64 {
		return deg * 0.017453292519943295
	}

	lat1Rad := degToRad(lat1)
	lat2Rad := degToRad(lat2)
	deltaLat := degToRad(lat2 - lat1)
	deltaLon := degToRad(lon2 - lon1)

	sinLat := math.Sin(deltaLat / 2)
	sinLon := math.Sin(deltaLon / 2)
	a := sinLat*sinLat + math.Cos(lat1Rad)*math.Cos(lat2Rad)*sinLon*sinLon
	if a > 1 {
		a = 1
	}
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusMeters * c
}
