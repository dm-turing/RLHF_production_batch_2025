package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// BucketChunk contains the S3 bucket name and key for each chunk.
type BucketChunk struct {
	BucketName string `json:"bucket_name"`
	Key        string `json:"key"`
}

// BucketPool manages a pool of S3 bucket chunks.
type BucketPool struct {
	mu      sync.Mutex
	release chan *BucketChunk
	chunks  []*BucketChunk
	cap     int
	count   int
}

func (f *BucketPool) NewBucket() (*BucketChunk, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.count < f.cap {
		// Create an S3 bucket and get a chunk key
		// For demonstration, we'll simply generate names
		bucketName := fmt.Sprintf("ephemeral-bucket-%03d", f.count)
		key := fmt.Sprintf("chunk-%03d.txt", f.count)
		chunk := &BucketChunk{BucketName: bucketName, Key: key}
		f.chunks = append(f.chunks, chunk)
		f.count++
		return chunk, nil
	}

	// Wait for a bucket to be returned
	select {
	case chunk := <-f.release:
		return chunk, nil
	default:
		log.Fatalf("bucket pool exhausted")
	}
	return nil, nil
}

func (f *BucketPool) Release(chunk *BucketChunk) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.count > 0 {
		f.chunks[f.count-1].BucketName = chunk.BucketName
		f.chunks[f.count-1].Key = chunk.Key
		f.count--
	}

	select {
	case f.release <- chunk:
	default:
		// If no one is waiting, delete the bucket
		deleteBucket(chunk.BucketName)
		f.release <- chunk
	}
}

func (f *BucketPool) Close() {
	close(f.release)
	for _, chunk := range f.chunks {
		deleteBucket(chunk.BucketName)
	}
}

func deleteBucket(name string) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"), // Change to your region
		Credentials: credentials.NewStaticCredentials("YOUR_ACCESS_KEY_ID", "YOUR_SECRET_ACCESS_KEY", ""),
	}))
	svc := s3.New(sess)
	input := &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	}
	_, err := svc.DeleteBucket(input)
	if err != nil {
		log.Fatalf("Failed to delete bucket %q: %v", name, err)
	}
	log.Printf("Deleted bucket %q", name)
}

func uploadChunkToS3(chunk *BucketChunk, data []byte) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"), // Change to your region
		Credentials: credentials.NewStaticCredentials("YOUR_ACCESS_KEY_ID", "YOUR_SECRET_ACCESS_KEY", ""),
	}))
	svc := s3.New(sess)
	input := &s3.PutObjectInput{
		Bucket: aws.String(chunk.BucketName),
		Key:    aws.String(chunk.Key),
		Body:   aws.NewReadSeekCloser(io.NopCloser(bytes.NewReader(data))),
	}
	_, err := svc.PutObject(input)
	if err != nil {
		log.Fatalf("Failed to upload chunk to S3: %v", err)
	}
	log.Printf("Uploaded chunk to S3: %q", chunk.Key)
}

func downloadChunkFromS3(chunk *BucketChunk) ([]byte, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"), // Change to your region
		Credentials: credentials.NewStaticCredentials("YOUR_ACCESS_KEY_ID", "YOUR_SECRET_ACCESS_KEY", ""),
	}))
	svc := s3.New(sess)
	input := &s3.GetObjectInput{
		Bucket: aws.String(chunk.BucketName),
		Key:    aws.String(chunk.Key),
	}
	resp, err := svc.GetObject(input)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func composeFile(chunks []*BucketChunk) ([]byte, error) {
	var composedData []byte
	for _, chunk := range chunks {
		data, err := downloadChunkFromS3(chunk)
		if err != nil {
			return nil, err
		}
		composedData = append(composedData, data...)
	}
	return composedData, nil
}

func main() {
	pool := BucketPool{
		cap: 100, // Set the maximum number of S3 buckets
		newBucket: func() (*BucketChunk, error) {
			return newBucket()
		},
		release: make(chan *BucketChunk, 1),
	}
	defer pool.Close()

	// Read a large file and split it into chunks
	file, err := os.Open("large_file.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	chunks := make([]*BucketChunk, 0)
	for scanner.Scan() {
		chunk, err := pool.NewBucket()
		if err != nil {
			log.Fatalf("Failed to get new bucket: %v", err)
		}
		data := []byte(scanner.Text() + "\n")
		uploadChunkToS3(chunk, data)
		chunks = append(chunks, chunk)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Compose the larger file from the chunks
	composedData, err := composeFile(chunks)
	if err != nil {
		log.Fatalf("Failed to compose file: %v", err)
	}

	// Upload the composed file to the final S3 bucket
	finalBucketName := "final-bucket"
	finalKey := "composed_file.txt"
	finalSess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"), // Change to your region
		Credentials: credentials.NewStaticCredentials("YOUR_ACCESS_KEY_ID", "YOUR_SECRET_ACCESS_KEY", ""),
	}))
	finalSvc := s3.New(finalSess)
	finalInput := &s3.PutObjectInput{
		Bucket: aws.String(finalBucketName),
		Key:    aws.String(finalKey),
		Body:   aws.NewReadSeekCloser(io.NopCloser(bytes.NewReader(composedData))),
	}
	_, err = finalSvc.PutObject(finalInput)
	if err != nil {
		log.Fatalf("Failed to upload final file to S3: %v", err)
	}
	log.Printf("Uploaded final file to S3: %q", finalKey)

	// Delete ephemeral buckets
	for _, chunk := range chunks {
		pool.Release(chunk)
	}
}
