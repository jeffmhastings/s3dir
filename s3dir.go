package s3dir

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

var (
	ErrNotImplemented = errors.New("Not implemented")
)

type Bucket struct {
	s3     s3iface.S3API
	bucket *string
	config BucketConfig
}

type BucketConfig struct {
	Region       string
	BucketName   string
	BucketPrefix string
}

func NewBucket(cfg BucketConfig) (*Bucket, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(cfg.Region),
		LogLevel: aws.LogLevel(aws.LogDebug),
	}))

	b := &Bucket{
		s3:     s3.New(sess),
		bucket: aws.String(cfg.BucketName),
		config: cfg,
	}

	params := &s3.HeadBucketInput{
		Bucket: b.bucket,
	}

	if _, err := b.s3.HeadBucket(params); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *Bucket) Open(path string) (http.File, error) {
	if b.config.BucketPrefix != "" && !strings.HasPrefix(path, b.config.BucketPrefix) {
		log.Printf("[INFO] denying access because %s doesn't begin with %s", path, b.config.BucketPrefix)
		return nil, os.ErrNotExist
	}

	prefix := strings.TrimPrefix(path, "/")
	if len(prefix) > 0 && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	input := &s3.ListObjectsInput{
		Bucket:    b.bucket,
		Prefix:    &prefix,
		Delimiter: aws.String("/"),
	}
	l, err := b.s3.ListObjects(input)
	if err != nil {
		return nil, err
	}

	if len(l.CommonPrefixes) > 0 || len(l.Contents) > 0 {
		return NewDir(b, prefix, l), nil
	}

	// This is an actual object. Get it
	params := &s3.GetObjectInput{
		Bucket: b.bucket,
		Key:    aws.String(path),
	}

	resp, err := b.s3.GetObject(params)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case s3.ErrCodeNoSuchKey:
				return nil, os.ErrNotExist
			}
		}
		return nil, err
	}

	return NewFile(b, path, resp), nil
}
