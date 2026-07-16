package model

import "github.com/google/uuid"

type GetPaymentTransactionParam struct {
	PaymentTransactionID uuid.UUID
	OrderID              string
	DonationID           uuid.UUID
	WalletTransactionID  uuid.UUID
	UserID               uuid.UUID
	RequestID            uuid.UUID
}

type GetWalletParam struct {
	WalletID uuid.UUID
	UserID   uuid.UUID
}

type GetWalletTransactionParam struct {
	WalletTransactionID uuid.UUID
	WalletID            uuid.UUID
}

type GetDonationParam struct {
	DonationID          uuid.UUID
	RequestID           uuid.UUID
	DonatedBy           uuid.UUID
	WalletTransactionID uuid.UUID
}
