package s3dir

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type mockTest404 struct {
	s3iface.S3API
	ListObjectsResp s3.ListObjectsOutput
	GetObjectResp   s3.GetObjectOutput
}

func (m mockTest404) ListObjects(in *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	// Only need to return mocked response output
	return &m.ListObjectsResp, nil
}

func (m mockTest404) GetObject(in *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	// Only need to return error
	return nil, awserr.New(s3.ErrCodeNoSuchKey, "something", nil)
}

func Test404(t *testing.T) {

	b := &Bucket{
		s3: mockTest404{
			ListObjectsResp: s3.ListObjectsOutput{},
			GetObjectResp:   s3.GetObjectOutput{},
		},
		bucket: aws.String("bucket"),
	}

	_, err := b.Open("/something")
	if err == nil {
		t.Error("Expected error, but was `nil`")
	} else if !os.IsNotExist(err) {
		t.Error("expected NotFound error")
	}

}
