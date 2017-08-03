package s3dir

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

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

type File struct {
	s3Object *s3.GetObjectOutput
	bucket   *Bucket
	path     string
}

func NewBucket(cfg BucketConfig) (*Bucket, error) {
	awsCfg := &aws.Config{
		Region:   aws.String(cfg.Region),
		LogLevel: 5,
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
	s3File := &File{
		bucket: b,
		path:   path,
	}

	params := &s3.GetObjectInput{
		Bucket: b.bucket,
		Key:    aws.String(path),
	}

	resp, err := b.s3.GetObject(params)
	if err != nil {
		return s3File, err
	}

	s3File.s3Object = resp

	return s3File, nil
}

// http.File
func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	return nil, ErrNotImplemented
}

// http.File
func (f *File) Stat() (os.FileInfo, error) {
	return f, nil
}

// io.Seeker
func (f *File) Seek(offset int64, whence int) (int64, error) {
	return 0, ErrNotImplemented
}

// io.Reader
func (f *File) Read(p []byte) (int, error) {
	return f.s3Object.Body.Read(p)
}

// io.Closer
func (f *File) Close() error {
	return f.s3Object.Body.Close()
}

// base name of the file
func (f *File) Name() string {
	name := strings.TrimLeft(f.path, "/")
	return name
}

// length in bytes for regular files; system-dependent for others
func (f *File) Size() int64 {
	return *f.s3Object.ContentLength
}

// file mode bits
func (f *File) Mode() os.FileMode {
	return 000
}

// modification time
func (f *File) ModTime() time.Time {
	return *f.s3Object.LastModified
}

// abbreviation for Mode().IsDir()
func (f *File) IsDir() bool {
	return false
}

func (f *File) Sys() interface{} {
	return f.s3Object
}
