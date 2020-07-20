package s3

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	as "github.com/clok/awssession"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

func Get(src string, dest string) error {
	var dstFile *os.File
	var err error
	dstFile, err = os.Create(dest)
	if err != nil {
		return cli.NewExitError(err, 2)
	}
	defer dstFile.Close()

	sess, _ := as.New()
	downloader := s3manager.NewDownloader(sess)

	bucket, object, err := getBucketAndObject(src)
	if err != nil {
		return cli.NewExitError(err, 2)
	}

	_, err = downloader.Download(dstFile, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(object),
	})

	if err != nil {
		return err
	}

	return nil
}

func getBucketAndObject(path string) (bucket string, object string, err error) {
	// error if doesn't start with s3://
	parts := strings.SplitN(path, "s3://", 2)
	if len(parts) > 1 {
		b := strings.SplitN(parts[1], "/", 2)
		if len(b) < 2 {
			return "", "", errors.New("ERROR unable to parse bucket path")
		}
		bucket, object = b[0], b[1]
	} else {
		return "", "", errors.New("ERROR unable to parse bucket path")
	}

	return
}
