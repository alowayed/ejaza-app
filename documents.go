package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var sess = session.Must(session.NewSession())
var uploader = s3manager.NewUploader(sess)

func uploadDocument() {
	// f, err := os.Open(filename)
	// if err != nil {
	// 	return fmt.Errorf("failed to open file %q, %v", filename, err)
	// }

	// // Upload the file to S3.
	// result, err := uploader.Upload(&s3manager.UploadInput{
	// 	Bucket: aws.String(myBucket),
	// 	Key:    aws.String(myString),
	// 	Body:   f,
	// })
	// if err != nil {
	// 	return fmt.Errorf("failed to upload file, %v", err)
	// }
	// fmt.Printf("file uploaded to, %s\n", aws.StringValue(result.Location))
}
