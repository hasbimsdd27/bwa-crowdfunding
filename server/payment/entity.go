package payment

type Transaction struct {
	ID     int
	Amount int
	Code   string
}

type RequestPayloadTransactionDetail struct {
	OrderID     string `json:"order_id" `
	GrossAmount int    `json:"gross_amount" `
}

type RequestPayloadCustomerDetail struct {
	FirstName string `json:"first_name" `
	LastName  string `json:"last_name"`
	Email     string `json:"email" `
	Phone     string `json:"phone"`
}

type RequestPayload struct {
	TransactionDetail RequestPayloadTransactionDetail `json:"transaction_details"`
	CustomerDetail    RequestPayloadCustomerDetail    `json:"customer_detail"`
	WebhookUrl        string                          `json:"webhook_url"`
}

type ResponseBodyData struct {
	PaymentUrl string `json:"payment_url"`
}

type ResponseBody struct {
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Data    ResponseBodyData `json:"data"`
}
