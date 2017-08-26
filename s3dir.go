package s3dir

import (
	"errors"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	ErrNotImplemented = errors.New("Not implemented")
)

type Bucket struct {
	s3     *s3.S3
	bucket *string
	config BucketConfig
}

type BucketConfig struct {
	Region     string
	BucketName string
}

func NewBucket(cfg BucketConfig) (*Bucket, error) {
	awsCfg := &aws.Config{
		Region:   aws.String(cfg.Region),
		LogLevel: aws.LogLevel(aws.LogDebug),
	}

	b := &Bucket{
		s3:     s3.New(session.New(awsCfg)),
		bucket: aws.String(cfg.BucketName),
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
		return nil, err
	}

	return NewFile(b, path, resp), nil
}
