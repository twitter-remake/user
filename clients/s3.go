package clients

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3 struct {
	sess     *session.Session
	client   *s3.S3
	uploader *s3manager.Uploader
	bucket   string
}

func NewS3(sess *session.Session, bucket string) *S3 {
	return &S3{
		sess:     sess,
		bucket:   bucket,
		client:   s3.New(sess),
		uploader: s3manager.NewUploader(sess),
	}
}

func (s *S3) UploadFile(ctx context.Context, key string, file io.Reader) (*s3manager.UploadOutput, error) {
	input := &s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   file,
	}

	return s.uploader.UploadWithContext(ctx, input)
}

func (s *S3) FindFile(ctx context.Context, key string) (*s3.GetObjectOutput, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	return s.client.GetObject(input)
}

func (s *S3) DeleteFile(ctx context.Context, key string) (*s3.DeleteObjectOutput, error) {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	return s.client.DeleteObject(input)
}
