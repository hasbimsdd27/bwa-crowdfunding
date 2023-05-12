package helper

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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

func S3ImageUploader(file *multipart.FileHeader) (string, error) {

	bucketName := os.Getenv("S3_BUCKET")
	endpoint := os.Getenv("S3_ENDPOINT")

	creds := credentials.NewStaticCredentials(
		os.Getenv("S3_ID"), os.Getenv("S3_SECRET"), "")

	cfg := aws.NewConfig().WithRegion("SouthJkt-a").WithEndpoint(endpoint).WithCredentials(creds)

	sess := session.Must(session.NewSession(cfg))

	svc := s3.New(sess)

	f, err := file.Open()

	if err != nil {
		return "", err
	}

	defer f.Close()

	size := file.Size
	buffer := make([]byte, size)

	f.Read(buffer)

	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	extensions := strings.Split(file.Filename, ".")
	path := "/images/" + fmt.Sprint(time.Now().UnixMilli()) + "." + extensions[len(extensions)-1]

	params := &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(path),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
		ACL:           aws.String("public-read"),
	}

	_, err = svc.PutObject(params)

	if err != nil {
		return "", err
	}

	return (fmt.Sprintf("https://%s/%s%s", endpoint, bucketName, path)), nil

}
