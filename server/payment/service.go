package payment

import (
	"bwastartup/user"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type service struct {
}

type Service interface {
	GetPaymentUrl(transaction Transaction, user user.User) (string, error)
}

func NewService() *service {
	return &service{}
}

func (s *service) GetPaymentUrl(transaction Transaction, user user.User) (string, error) {

	customerDetail := RequestPayloadCustomerDetail{}
	customerDetail.FirstName = user.Name
	customerDetail.Email = user.Email

	transactionDetail := RequestPayloadTransactionDetail{}
	transactionDetail.GrossAmount = transaction.Amount
	transactionDetail.OrderID = transaction.Code

	payloadBody := RequestPayload{}
	payloadBody.CustomerDetail = customerDetail
	payloadBody.TransactionDetail = transactionDetail
	payloadBody.WebhookUrl = "http://localhost:7350/api/v1/webhook/midtrans"

	postBody, _ := json.Marshal(payloadBody)
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(fmt.Sprintf("%s/create", os.Getenv("MIDTRANS_ADAPTER_URL")), "application/json", responseBody)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var respBody ResponseBody
	bodyBytes, err := io.ReadAll(resp.Body)
	fmt.Println(fmt.Println(string(bodyBytes)))
	if err != nil {
		return "", err
	}

	if err = json.Unmarshal(bodyBytes, &respBody); err != nil {
		return "", err
	}
	if resp.StatusCode == http.StatusOK {

		return respBody.Data.PaymentUrl, nil
	} else {
		return "", errors.New(respBody.Message)
	}
}
