package requests

type CreateWalletRequest struct {
	Name   string  `json:"name" validate:"required,max=100"`
	Amount float64 `json:"amount" validate:"number,required,min=0"`
}
