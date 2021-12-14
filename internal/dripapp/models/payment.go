package models

const DateLayout = "2006-01-02T15:04:05.000Z"

const (
	PaymentStatusCanceled       = false
	PaymentStatusCanceledString = "canceled"
	PaymentStatusSuccess        = true
	PaymentStatusSuccessString  = "succeeded"
)

//easyjson:json
type Payment struct {
	Period uint64 `json:"period"`
	Amount string `json:"amount"`
}

type RedirectUrl struct {
	URL string `json:"redirectUrl"`
}

type PaymentInfo struct {
	Amount       map[string]string `json:"amount"`
	Capture      bool              `json:"capture"`
	Confirmation map[string]string `json:"confirmation"`
}

type YooKassaResponse struct {
	Id           string           `json:"id"`
	Status       string           `json:"status"`
	Amount       AmountType       `json:"amount"`
	Recipient    RecipientType    `json:"recipient"`
	CreatedAt    string           `json:"created_at"`
	Confirmation ConfirmationType `json:"confirmation"`
	Test         bool             `json:"test"`
	Paid         bool             `json:"paid"`
	Refundable   bool             `json:"refundable"`
}

type AmountType struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type RecipientType struct {
	AccountId string `json:"account_id"`
	GatewayId string `json:"gateway_id"`
}

type ConfirmationType struct {
	Type            string `json:"type"`
	ConfirmationUrl string `json:"confirmation_url"`
}

type PaymentNotification struct {
	Object ObjectType `json:"object"`
}

type ObjectType struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}
