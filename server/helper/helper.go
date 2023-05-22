package helper

import (
	"crypto/sha512"
	"errors"
	"fmt"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

type GenerateTransactionKeyInput struct {
	Code   string
	UserID int
	Amount int
}

func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	jsonResponse := Response{
		Meta: meta, Data: data,
	}

	return jsonResponse
}

func FormatValidationError(err error) []string {
	var errors []string
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}

	return errors
}

func ImageFileValidator(file *multipart.FileHeader) error {
	f, err := file.Open()

	if err != nil {
		return err
	}

	defer f.Close()

	size := file.Size
	buffer := make([]byte, size)

	f.Read(buffer)

	fileType := http.DetectContentType(buffer)

	maxSize := int64(1024000)

	if size > maxSize {
		return errors.New("filesize too large")
	}

	if !strings.Contains(fileType, "image") {
		return errors.New("unsupported filetype")
	}

	return nil
}

func RandStr(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GenerateTransactionKey(input GenerateTransactionKeyInput) string {

	combinedData := fmt.Sprintf("%s%d%d%s", input.Code, input.Amount, input.UserID, os.Getenv("MIDTRANS_ADAPTER_SERVER_KEY"))

	transactionKey := sha512.New()
	transactionKey.Write([]byte(combinedData))

	return fmt.Sprintf("%x", transactionKey.Sum(nil))
}
