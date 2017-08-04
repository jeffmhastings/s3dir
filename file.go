package s3dir

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
)

type File struct {
	s3Object *s3.GetObjectOutput
	bucket   *Bucket
	path     string
}

func NewFile(bucket *Bucket, path string, object *s3.GetObjectOutput) *File {
	return &File{
		bucket:   bucket,
		path:     path,
		s3Object: object,
	}
}

func (f *File) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	// since s3 doesn't support seeking, we can only seek from the start of the file
	if whence != io.SeekStart {
		return 0, ErrNotImplemented
	}

	// reopen the file
	nf, err := f.bucket.Open(f.path)
	if err != nil {
		return 0, err
	}

	f.s3Object = nf.(*File).s3Object

	if _, err := f.Read(make([]byte, offset)); err != nil {
		return 0, err
	}

	return 0, nil
}

func (f *File) Read(p []byte) (int, error) {
	i, err := f.s3Object.Body.Read(p)
	return i, err
}

func (f *File) Close() error {
	return f.s3Object.Body.Close()
}

func (f *File) Name() string {
	return strings.TrimRight(f.path, "/")
}

func (f *File) Size() int64 {
	return *f.s3Object.ContentLength
}

func (f *File) Mode() os.FileMode {
	return 0444
}

func (f *File) ModTime() time.Time {
	if f.s3Object.LastModified == nil {
		return time.Time{}
	}
	return *f.s3Object.LastModified
}

func (f *File) IsDir() bool {
	return false
}

func (f *File) Sys() interface{} {
	return f.s3Object
}

func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	// Not a directory
	return nil, nil
}
