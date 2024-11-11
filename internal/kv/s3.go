package kv

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/util"
	minio "github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

// S3 holds a bucket structure (containing VCAP_SERVICES information)
// and an S3 client connection from the min.io libraries.
type S3 struct {
	Bucket      env.Bucket
	MinioClient *minio.Client
}

// /////////////////////////////
// FILE INTERACTIONS

// NewS3 creates a new S3 object for the bucket given.
// Lets us copy files to/from the bucket.
func NewS3(bucket_name string) *S3 {
	s3 := newS3FromBucketName(bucket_name)
	return &s3
}

func (s3 *S3) FileToS3(key *util.Key, local_filename string, mime_type string) error {
	reader, err := os.Open(local_filename)
	if err != nil {
		log.Fatal("KV cannot open file", local_filename)
	}
	fi, err := reader.Stat()
	if err != nil {
		log.Println("KV could not stat file")
		log.Fatal(err)
	}

	return store(s3, key.Render(), fi.Size(), reader, mime_type)
}

func (s3 *S3) FileToS3Path(key string, local_filename string, mime_type string) error {
	reader, err := os.Open(local_filename)
	if err != nil {
		log.Fatal("KV cannot open file", local_filename)
	}
	fi, err := reader.Stat()
	if err != nil {
		log.Println("KV could not stat file")
		log.Fatal(err)
	}

	return store(s3, key, fi.Size(), reader, mime_type)
}

func (s3 *S3) S3ToFile(key *util.Key, local_filename string) error {
	ctx := context.Background()
	err := s3.MinioClient.FGetObject(
		ctx,
		s3.Bucket.CredentialString("bucket"),
		key.Render(),
		local_filename,
		minio.GetObjectOptions{})

	if err != nil {
		zap.L().Error("could not FGetObject",
			zap.String("bucket", s3.Bucket.Name),
			zap.String("key", key.Render()),
			zap.String("local_filename", local_filename),
		)
		return err
	}
	return nil
}

func (s3 *S3) S3PathToFile(path string, local_filename string) error {
	ctx := context.Background()
	err := s3.MinioClient.FGetObject(
		ctx,
		s3.Bucket.CredentialString("bucket"),
		path,
		local_filename,
		minio.GetObjectOptions{})

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// Lists objects in the bucket, returning keys and sizes.
func (s3 *S3) List(prefix string) ([]*ObjInfo, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	objectCh := s3.MinioClient.ListObjects(
		ctx,
		s3.Bucket.CredentialString("bucket"),
		minio.ListObjectsOptions{
			Prefix:    prefix,
			Recursive: false,
		})

	objects := make([]*ObjInfo, 0)
	for object := range objectCh {
		if object.Err != nil {
			fmt.Println(object.Err)
			return nil, object.Err
		}
		objects = append(objects, NewObjInfo(object.Key, object.Size))
	}
	return objects, nil
}
