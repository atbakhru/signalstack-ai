package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Uploader interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

type S3Service struct {
	Client S3Uploader
	Bucket string
}

func (s S3Service) UploadJSON(ctx context.Context, sourceName, externalID string, payload any) (string, error) {
	if s.Client == nil || s.Bucket == "" {
		return "", fmt.Errorf("s3 client or bucket not configured")
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	key := path.Join("raw", sourceName, time.Now().UTC().Format("2006/01/02"), fmt.Sprintf("%s.json", externalID))
	_, err = s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(encoded),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		return "", err
	}

	return key, nil
}
