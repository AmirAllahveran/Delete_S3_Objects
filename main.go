package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"

	"fmt"
	"os"
)

//    go run main.go BUCKET

func main() {
	if len(os.Args) != 2 {
		exitErrorf("Bucket name required\nUsage: %s BUCKET", os.Args[0])
	}

	bucket := os.Args[1]

	// load .env file
	err := godotenv.Load(".env")
	exitErrorf("Error loading .env file", err)

	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), ""),
	})

	// Create S3 service client
	svc := s3.New(sess, &aws.Config{Endpoint: aws.String(os.Getenv("ENDPOINT_URL"))})

	// Setup BatchDeleteIterator to iterate through a list of objects.
	iter := s3manager.NewDeleteListIterator(svc, &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	})

	// Traverse iterator deleting each object
	if err := s3manager.NewBatchDeleteWithClient(svc).Delete(aws.BackgroundContext(), iter); err != nil {
		exitErrorf("Unable to delete objects from bucket %q, %v", bucket, err)
	}

	fmt.Printf("Deleted object(s) from bucket: %s", bucket)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
