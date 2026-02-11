package requests

type TransferRequest struct {
	SourceWalletId      uint    `json:"source_wallet_id" validate:"required,number"`
	DestinationWalletId uint    `json:"destination_wallet_id" validate:"required,number"`
	Amount              float64 `json:"amount" validate:"required,numeric"`
}
