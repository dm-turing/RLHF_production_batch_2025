package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3BucketPool provides a synchronized way to manage S3 buckets and file handles within them.
type S3BucketPool struct {
	mu         sync.Mutex
	release    chan *BucketHandle
	buckets    map[string]*BucketHandle
	newBucket  func() (*BucketHandle, error)
	s3         *s3.S3
	uploadBucket string
}

type BucketHandle struct {
	bucket  *s3.Bucket
	file    *os.File
	wc      io.WriteCloser
}

func (p *S3BucketPool) NewBucketHandle() (*BucketHandle, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check if there are existing buckets
	for _, handle := range p.buckets {
		if handle.wc != nil {
			return handle, nil
		}
	}

	// If no available bucket, create a new one
	bucketHandle, err := p.newBucket()
	if err != nil {
		return nil, err
	}

	p.buckets[bucketHandle.bucket.Name] = bucketHandle
	return bucketHandle, nil
}

func (p *S3BucketPool) ReleaseBucketHandle(bucketHandle *BucketHandle) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.buckets[bucketHandle.bucket.Name]; ok {
		bucketHandle.wc.Close()
		bucketHandle.file.Close()
		p.buckets[bucketHandle.bucket.Name] = nil
	}

	p.release <- bucketHandle
}

func (p *S3BucketPool) Close() {
	close(p.release)
	for _, bucketHandle := range p.buckets {
		bucketHandle.wc.Close()
		bucketHandle.file.Close()
	}
}

// NewS3BucketPool initializes the S3BucketPool with the given configuration.
func NewS3BucketPool(s3 *s3.S3, capacity int, uploadBucket string) *S3BucketPool {
	return &S3BucketPool{
		mu:      sync.Mutex{},
		release: make(chan *BucketHandle, capacity),
		buckets:  make(map[string]*BucketHandle, capacity),
		newBucket: func() (*BucketHandle, error) {
			// Generate a new bucket name or use an existing one with proper naming convention
			bucketName := "your-bucket-prefix-" + aws.StringValue(s3.Client.Config.Region)

			_, err := s3.CreateBucket(&s3.CreateBucketInput{
				Bucket: aws.String(bucketName),
			})
			if err != nil {
				return nil, err
			}

			f, err := os.CreateTemp("", "file-chunk-")
			if err != nil {
				return nil, err
			}

			wc, err := s3.NewUploader(s3.Client.Config).Upload(&s3manager.UploadInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String("file.dat"),
				Body:   f,
			})
			if err != nil {
				return nil, err
			}

			return &BucketHandle{
				bucket: &s3.Bucket{
					Name: aws.String(bucketName),
				},