package s3dir

import (
	"errors"
	"log"
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
	log.Printf("S3 Open: %s", path)
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
	log.Printf("S3 Response: %v", resp)

	s3File.s3Object = resp

	return s3File, nil
}

// http.File
func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	log.Printf("Readdir: %d. %s", count, f.path)
	return nil, ErrNotImplemented
}

// http.File
func (f *File) Stat() (os.FileInfo, error) {
	log.Printf("Stat: %s", f.path)
	return f, nil
}

// io.Seeker
func (f *File) Seek(offset int64, whence int) (int64, error) {
	log.Printf("seek, %v %v ", offset, whence)
	return 0, ErrNotImplemented
}

// io.Reader
func (f *File) Read(p []byte) (int, error) {

	log.Printf("read")
	return f.s3Object.Body.Read(p)
}

// io.Closer
func (f *File) Close() error {
	log.Printf("close")
	return f.s3Object.Body.Close()
}

// base name of the file
func (f *File) Name() string {
	log.Printf("name")
	name := strings.TrimLeft(f.path, "/")
	return name
}

// length in bytes for regular files; system-dependent for others
func (f *File) Size() int64 {
	log.Printf("size")
	return *f.s3Object.ContentLength
}

// file mode bits
func (f *File) Mode() os.FileMode {
	log.Printf("mode")
	return 0444
}

// modification time
func (f *File) ModTime() time.Time {
	log.Printf("modtime")
	return *f.s3Object.LastModified
}

// abbreviation for Mode().IsDir()
func (f *File) IsDir() bool {
	log.Printf("isdir")
	return strings.HasSuffix(f.path, "/")
}

func (f *File) Sys() interface{} {
	log.Printf("Sys")
	return f.s3Object
}
