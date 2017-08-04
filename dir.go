package s3dir

import (
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
)

type Dir struct {
	bucket *Bucket
	data   *s3.ListObjectsOutput
	path   string
}

func NewDir(bucket *Bucket, path string, listObjects *s3.ListObjectsOutput) *Dir {
	return &Dir{
		bucket: bucket,
		data:   listObjects,
		path:   path,
	}
}

func (d *Dir) Name() string {
	return strings.TrimSuffix(d.path, "/")
}

func (d *Dir) Stat() (os.FileInfo, error) {
	return d, nil
}

func (*Dir) Size() int64 {
	return 0
}

func (*Dir) IsDir() bool {
	return true
}

func (*Dir) ModTime() time.Time {
	return time.Time{}
}

func (*Dir) Mode() os.FileMode {
	return 0755
}

func (d *Dir) Sys() interface{} {
	return d.data
}

func (*Dir) Close() error {
	// nothing to do
	return nil
}

func (*Dir) Read(p []byte) (int, error) {
	return 0, nil
}

func (*Dir) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (d *Dir) Readdir(count int) ([]os.FileInfo, error) {
	// if there are more than one common prefix, treat as directory
	info := make([]os.FileInfo, len(d.data.CommonPrefixes)+len(d.data.Contents))
	for i := 0; i < len(d.data.CommonPrefixes); i++ {
		info[i] = &Dir{
			bucket: d.bucket,
			path:   *d.data.CommonPrefixes[i].Prefix,
		}

	}

	for i := 0; i < len(d.data.Contents); i++ {
		info[i] = &File{
			bucket: d.bucket,
			path:   strings.TrimPrefix(*d.data.Contents[i].Key, d.path),
		}
	}
	return info, nil
}
