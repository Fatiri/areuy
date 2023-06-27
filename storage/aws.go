package storage

import (
	"fmt"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AwsStorage interface {
	Upload(bucketName string, files *multipart.FileHeader, key string) (string, error)
	Delete(bucketName string, key string) error
	GetFile(bucketName, key string) (string, error)
}

type awsStorageCtx struct {
	session *session.Session
	err     error
}

func NewAwsStorage(accessID, secretKey, region, endpoint string) AwsStorage {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Endpoint:    aws.String(endpoint),
		Credentials: credentials.NewStaticCredentials(accessID, secretKey, ""),
	})
	return &awsStorageCtx{
		session: sess,
		err:     err,
	}
}

func (a *awsStorageCtx) Upload(bucketName string, files *multipart.FileHeader, key string) (string, error) {
	fileRes, err := files.Open()
	if err != nil {
		return "", err
	}
	uploader := s3manager.NewUploader(a.session)
	_, a.err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		ContentType: aws.String(files.Header.Get("contain-type")),
		Body:        fileRes,
	})
	fmt.Println(a.err)
	if a.err != nil {
		return "", a.err
	}

	return key, nil
}

func (a *awsStorageCtx) Delete(bucketName string, key string) error {
	svc := s3.New(a.session)

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *awsStorageCtx) GetFile(bucketName, key string) (string, error) {
	// Create S3 service client
	svc := s3.New(a.session)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	urlStr, err := req.Presign(20 * time.Minute)
	if err != nil {
		return "", err
	}
	return urlStr, nil
}
