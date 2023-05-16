package helper

import (
	"bytes"
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
)

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
