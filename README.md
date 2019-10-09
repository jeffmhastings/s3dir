# s3dir

Implements [http.Dir](https://golang.org/pkg/net/http/#Dir).

This can be used to expose a s3 bucket (ex: for static assets).

## Usage

```go
package main
import (
  "net/http"

  "github.com/jeffmhastings/s3dir"
)

func main() {
	cfg := s3dir.BucketConfig{
		BucketName: "mah-bucket",
		Region:     "us-east-1",
	}

	bucket, err := s3dir.NewBucket(cfg)
	if err != nil {
		panic(err)
	}

  // note the trailing '/' when combining with http.StripPrefix
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(bucket)))
	http.ListenAndServe(":8080", nil)
}
```

## License

MIT
