package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/jeffmhastings/s3dir"
)

var (
	bucketName = flag.String("bucket", "", "S3 Bucket to expose")
	region     = flag.String("region", "us-east-1", "S3 Region")
	bind       = flag.String("bind", ":8080", "http listener")
)

func main() {
	flag.Parse()

	if *bucketName == "" {
		log.Println("Missing bucket")
		flag.PrintDefaults()
		os.Exit(1)
	}

	cfg := s3dir.BucketConfig{
		BucketName: *bucketName,
		Region:     *region,
	}

	bucket, err := s3dir.NewBucket(cfg)
	if err != nil {
		panic(err)
	}

	dir := http.Dir("./")
	http.Handle("/s3/", http.StripPrefix("/s3/", http.FileServer(bucket)))
	http.Handle("/local/", http.StripPrefix("/local/", http.FileServer(dir)))

	log.Println("Starting HTTP interface on", *bind)
	http.ListenAndServe(*bind, nil)
}
