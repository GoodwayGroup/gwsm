package s3

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/clok/kemba"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strings"
)

var (
	k    = kemba.New("s3")
	kget = k.Extend("Get")
	kgbo = k.Extend("getBucketAndObject")
)

func Get(src string, dest string) error {
	var dstFile *os.File
	var err error
	dstFile, err = os.Create(dest)
	if err != nil {
		return cli.Exit(err, 2)
	}
	defer dstFile.Close()
	kget.Printf("writing to path: %s", dest)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)
	downloader := manager.NewDownloader(client)

	bucket, object, err := getBucketAndObject(src)
	if err != nil {
		return cli.Exit(err, 2)
	}

	kget.Printf("downloading - bucket: %s object: %s", bucket, object)
	_, err = downloader.Download(context.TODO(), dstFile, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(object),
	})
	kget.Printf("download complete - bucket: %s object: %s", bucket, object)

	if err != nil {
		return err
	}

	return nil
}

func getBucketAndObject(path string) (bucket string, object string, err error) {
	kgbo.Printf("path: %s", path)
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
	kgbo.Printf("bucket: %s object: %s", bucket, object)

	return
}
