package domain

import (
	"time"

	"github.com/google/uuid"
)


type TransferListingStatus string

const (
	TransferStatusActive    TransferListingStatus = "active"
	TransferStatusSold      TransferListingStatus = "sold"
	TransferStatusCancelled TransferListingStatus = "cancelled"
)


type TransferListing struct {
	ID          uuid.UUID             `json:"id" db:"id"`
	PlayerID    uuid.UUID             `json:"player_id" db:"player_id"`
	AskingPrice float64               `json:"asking_price" db:"asking_price"`
	Status      TransferListingStatus `json:"status" db:"status"`
	ListedAt    time.Time             `json:"listed_at" db:"listed_at"`
}


type TransferListingWithPlayer struct {
	TransferListing
	Player Player `json:"player"`
}


func NewTransferListing(playerID uuid.UUID, askingPrice float64) *TransferListing {
	return &TransferListing{
		ID:          uuid.New(),
		PlayerID:    playerID,
		AskingPrice: askingPrice,
		Status:      TransferStatusActive,
		ListedAt:    time.Now(),
	}
}


func (tl *TransferListing) MarkAsSold() {
	tl.Status = TransferStatusSold
}


func (tl *TransferListing) Cancel() {
	tl.Status = TransferStatusCancelled
}


func (tl *TransferListing) IsActive() bool {
	return tl.Status == TransferStatusActive
}


type Transfer struct {
	ID            uuid.UUID `json:"id" db:"id"`
	PlayerID      uuid.UUID `json:"player_id" db:"player_id"`
	SellerTeamID  uuid.UUID `json:"seller_team_id" db:"seller_team_id"`
	BuyerTeamID   uuid.UUID `json:"buyer_team_id" db:"buyer_team_id"`
	TransferPrice float64   `json:"transfer_price" db:"transfer_price"`
	TransferredAt time.Time `json:"transferred_at" db:"transferred_at"`
}


func NewTransfer(playerID, sellerTeamID, buyerTeamID uuid.UUID, transferPrice float64) *Transfer {
	return &Transfer{
		ID:            uuid.New(),
		PlayerID:      playerID,
		SellerTeamID:  sellerTeamID,
		BuyerTeamID:   buyerTeamID,
		TransferPrice: transferPrice,
		TransferredAt: time.Now(),
	}
}
