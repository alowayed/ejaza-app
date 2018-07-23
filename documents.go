package main

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	documentBucket = "ejaza-documents"
	awsRegion      = "us-east-2"
	awsKeyID       = "AKIAIBNG5PXNNNMC54LQ"
	awsKeySecret   = "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"
)

// -------------------- Global Variables

var sess = session.Must(session.NewSession(
	&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsKeyID, awsKeySecret, ""),
	},
))
var uploader = s3manager.NewUploader(sess)
var downloader = s3manager.NewDownloader(sess)

/**
 * Upload a document to S3
 */
func upload(file io.Reader, id string) error {
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(documentBucket),
		Key:    aws.String(id),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to, %s\n", result.Location)
	return nil
}

// func download(string id) error {
// 	file, err := os.Create(id)
// 	numBytes, err := downloader.Download(file,
// 		&s3.GetObjectInput{
// 			Bucket: aws.String(bucket),
// 			Key:    aws.String(item),
// 		})
// 	if err != nil {
// 		exitErrorf("Unable to download item %q, %v", item, err)
// 	}

// 	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
// }
