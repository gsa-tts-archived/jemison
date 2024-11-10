// kv provides an interface to key/value work in S3
// It is specialized to the `jemison` architecture.
package kv

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	minio "github.com/minio/minio-go/v7"
	minio_credentials "github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"go.uber.org/zap"

	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/util"
)

// NewFromBytes(bucket_name string, host string, path string, m []byte) *S3JSON
// NewEmptyS3JSON(bucket_name string, host string, path string) *S3JSON
// (s3json *S3JSON) IsEmpty() bool
// (s3json *S3JSON) Save() error
// (s3json *S3JSON) Load() error

// Only open any given bucket once.
// FIXME: get rid of these long-lived globals.
// Open and close things, until it becomes a performance concern.
// It would be safer if...
// Load() does an open and a close
// Save() does an open and a close
// Then, every object is self-contained. Slower, but self-contained.
// The sync... is hell waiting to happen in terms of debugging.
//var buckets sync.Map

// S3 holds a bucket structure (containing VCAP_SERVICES information)
// and an S3 client connection from the min.io libraries.
type S3 struct {
	Bucket      env.Bucket
	MinioClient *minio.Client
}

// S3JSON structs are JSON documents stored in S3.
// This is because `jemison` shuttles JSON documents in-and-out of S3, and
// we want to be able to find a document representing a host/path in
// multiple, different buckets.
type S3JSON struct {
	Key   util.Key
	raw   []byte
	S3    S3
	empty bool
	open  bool
}

func NewS3JSON(bucket_name string) *S3JSON {
	s3 := newS3FromBucketName(bucket_name)
	return &S3JSON{
		Key:   util.Key{},
		raw:   nil,
		S3:    s3,
		empty: true,
		open:  false,
	}
}

// NewFromBytes takes a []byte representation of a JSON document and constructs
// a S3JSON document from it.
// Inserts _key
func NewFromBytes(bucket_name string, host string, path string, m []byte) *S3JSON {
	s3 := newS3FromBucketName(bucket_name)
	key := util.CreateS3Key(host, path, "json")
	w_key, _ := sjson.SetBytes(m, "_key", key.Render())
	return &S3JSON{
		Key:   key,
		raw:   w_key,
		S3:    s3,
		empty: false,
		open:  true,
	}
}

// Inserts _key
func NewFromMap(bucket_name string, host string, path string, m map[string]string) *S3JSON {
	s3 := newS3FromBucketName(bucket_name)
	key := util.CreateS3Key(host, path, "json")
	m["_key"] = key.Render()
	b, _ := json.Marshal(m)
	return &S3JSON{
		Key:   key,
		raw:   b,
		S3:    s3,
		empty: false,
		open:  true,
	}
}

// Creates a new, empty S3JSON struct, setting it as `empty`.
// `Load()` must be called on it before we can use it.
func NewEmptyS3JSON(bucket_name string, host string, path string) *S3JSON {
	s3 := newS3FromBucketName(bucket_name)
	key := util.CreateS3Key(host, path, "json")
	return &S3JSON{
		Key:   key,
		raw:   nil,
		S3:    s3,
		empty: true,
		open:  true,
	}
}

// IsEmpty() Checks if the S3JSON struct is empty.
// Should be `true` before a call to `Load()`, `false` after.
func (s3json *S3JSON) IsEmpty() bool {
	return s3json.empty
}

func (s3json *S3JSON) IsOpen() bool {
	return s3json.open
}

// Save() will do a `Put` of the JSON to S3.
// BUG(jadudm): handle errors in store gracefully
func (s3json *S3JSON) Save() error {
	if s3json.IsEmpty() {
		return fmt.Errorf("cannot save invalid S3JSON object bucket[%s] host[%s] path[%s]",
			s3json.S3.Bucket.Name, s3json.Key.Host, s3json.Key.Path)
	}

	r := bytes.NewReader(s3json.raw)
	size := int64(len(s3json.raw))
	err := store(&s3json.S3, s3json.Key.Render(), size, r, "application/json")
	if err != nil {
		zap.L().Panic("could not store S3JSON",
			zap.String("bucket_name", s3json.S3.Bucket.Name),
			zap.String("key", s3json.Key.Render()))
	}
	return nil
}

// Load() uses the bucket/path information in the underlying S3 struct
// to do a `Get` against S3 and retrieve the JSON document.
func (s3json *S3JSON) Load() error {
	if s3json.IsEmpty() {
		return fmt.Errorf("will not save empty object bucket[%s] host[%s] path[%s]",
			s3json.S3.Bucket.Name, s3json.Key.Host, s3json.Key.Path)
	}

	// GetObject(ctx context.Context, bucketName, objectName string, opts GetObjectOptions) (*Object, error)
	// func (s3 *S3) Get(key string) (Object, error) {
	ctx := context.Background()
	key := s3json.Key.Render()
	// The object has a channel interface that we have to empty.
	object, err := s3json.S3.MinioClient.GetObject(
		ctx,
		s3json.S3.Bucket.CredentialString("bucket"),
		key,
		minio.GetObjectOptions{})
	// https://rezakhademix.medium.com/defer-functions-in-golang-common-mistakes-and-best-practices-96eacdb551f0
	defer func(obj *minio.Object) {
		err := obj.Close()
		if err != nil {
			zap.L().Error("deferred close on S3 object encountered error",
				zap.String("key", key))
		}
	}(object)

	zap.L().Debug("retrieved S3 object", zap.String("key", key))

	if err != nil {
		zap.L().Error("could not retrieve object",
			zap.String("bucket_name", s3json.S3.Bucket.CredentialString("bucket")),
			zap.String("key", key),
			zap.String("error", err.Error()))
		return err
	}

	raw, err := io.ReadAll(object)
	if err != nil {
		zap.L().Error("could not read object bytes",
			zap.String("bucket_name", s3json.S3.Bucket.CredentialString("bucket")),
			zap.String("key", key),
			zap.String("error", err.Error()))
		return err
	}

	s3json.raw = raw
	current_mime_type := s3json.GetString("content-type")
	sjson.SetBytes(s3json.raw, "content-type", util.CleanMimeType(current_mime_type))
	s3json.empty = false
	return nil
}

func (s3json *S3JSON) GetString(gjson_path string) string {
	r := gjson.GetBytes(s3json.raw, gjson_path)
	return r.String()
}

func (s3json *S3JSON) GetInt64(gjson_path string) int64 {
	r := gjson.GetBytes(s3json.raw, gjson_path)
	return int64(r.Int())
}

func (s3json *S3JSON) GetBool(gjson_path string) bool {
	r := gjson.GetBytes(s3json.raw, gjson_path)
	return r.Bool()
}

// NewS3FromBucketName creates an S3 object containing bucket information
// from VCAP and a minio client ready to talk to the bucket. S3JSON objects
// carry the information so they can load/save.
func newS3FromBucketName(bucket_name string) S3 {
	if !env.IsValidBucketName(bucket_name) {
		log.Fatal("KV INVALID BUCKET NAME ", bucket_name)
	}

	// Check if we already have this in the map, so reconnects don't create
	// new S3 objects/etc.
	// if s3, ok := buckets.Load(bucket_name); ok {
	// 	zap.L().Debug("in the sync map", zap.String("bucket_name", bucket_name))
	// 	return s3.(S3)
	// } else {
	// 	zap.L().Debug("not in the sync map", zap.String("bucket_name", bucket_name))
	// }

	s3 := S3{}

	// Grab a reference to our bucket from the config.
	b, err := env.Env.GetObjectStore(bucket_name)

	if err != nil {
		zap.L().Error("could not get bucket from config", zap.String("bucket_name", bucket_name))
		os.Exit(1)
	}

	zap.L().Debug("got reference to bucket from vcap",
		zap.String("name", b.Name),
		zap.String("bucket", b.CredentialString("bucket")),
		zap.String("region", b.CredentialString("region")))

	s3.Bucket = b

	// Initialize minio client object.
	useSSL := true
	if env.IsContainerEnv() || env.IsLocalTestEnv() {
		// log.Println("ENV disabling SSL in containerized environment")
		useSSL = false
	}

	options := minio.Options{
		Creds: minio_credentials.NewStaticV4(
			b.CredentialString("access_key_id"),
			b.CredentialString("secret_access_key"), ""),
		Secure: useSSL,
	}

	minioClient, err := minio.New(b.CredentialString("endpoint"), &options)
	if err != nil {
		log.Fatalln(err)
	}
	s3.MinioClient = minioClient
	ctx := context.Background()

	found, err := minioClient.BucketExists(ctx, s3.Bucket.CredentialString("bucket"))
	if err != nil {
		//log.Println("KV could not check if bucket exists ", bucket_name)
		//log.Fatal(err)
		zap.L().Fatal("could not check if bucket exists", zap.String("bucket_name", bucket_name))
	}

	if found {
		zap.L().Debug("pre-existing bucket in S3",
			zap.String("bucket_name", bucket_name))
		// Make sure to insert the metadata into the sync.Map
		// when we find a bucket that already exists.
		// buckets.Store(bucket_name, s3)
		return s3
	}

	if env.IsContainerEnv() {
		log.Println("KV creating new bucket ", bucket_name)
		// Try and make the bucket; if we're local, this is necessary.
		ctx := context.Background()
		err = minioClient.MakeBucket(
			ctx,
			s3.Bucket.CredentialString("bucket"),
			minio.MakeBucketOptions{Region: b.CredentialString("region")})

		if err != nil {
			log.Println(err)
			log.Fatal("KV could not create bucket ", bucket_name)
		}
	} // Skip container creation in CF

	// Put a pointer to this object in our syncmap.
	// buckets.Store(bucket_name, &s3)

	// loaded, _ := buckets.Load(bucket_name)
	//zap.L().Info("bucket ready", zap.String("bucket_name", loaded.(*S3).Bucket.Name))

	return s3
}

// store saves things to S3
func store(s3 *S3, destination_key string, size int64, reader io.Reader, mime_type string) error {
	// mime := "octet/binary"

	// if mime_type != nil {
	// 	mime = "application/json"
	// }

	ctx := context.Background()
	_, err := s3.MinioClient.PutObject(
		ctx,
		s3.Bucket.CredentialString("bucket"),
		destination_key,
		reader,
		size,
		minio.PutObjectOptions{
			ContentType: mime_type,
			// This seems to set the *minimum* partsize for multipart uploads.
			// Which... makes writing JSON objects impossible.
			// PartSize:    5000000
		},
	)
	if err != nil {
		zap.L().Warn("S3JSON could not PUT object",
			zap.String("destination_key", destination_key),
			zap.String("error", err.Error()))
	}
	return err
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

// type Storage interface {
// 	Store(string, JSON) error
// 	List(string) ([]*ObjInfo, error)
// 	Get(string) (Object, error)
// }

// func (s3 *S3) GetFile(key string, dest_filename string) error {
// 	ctx := context.Background()
// 	err := s3.MinioClient.FGetObject(
// 		ctx,
// 		s3.Bucket.CredentialString("bucket"),
// 		key,
// 		dest_filename,
// 		minio.GetObjectOptions{})

// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	return nil
// }

// // //////////////////
// // LIST
// // Lists objects in the bucket, returning keys and sizes.
// func (s3 *S3) List(prefix string) ([]*ObjInfo, error) {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	objectCh := s3.MinioClient.ListObjects(
// 		ctx,
// 		s3.Bucket.CredentialString("bucket"),
// 		minio.ListObjectsOptions{
// 			Prefix:    prefix,
// 			Recursive: false,
// 		})

// 	objects := make([]*ObjInfo, 0)
// 	for object := range objectCh {
// 		if object.Err != nil {
// 			fmt.Println(object.Err)
// 			return nil, object.Err
// 		}
// 		objects = append(objects, NewObjInfo(object.Key, object.Size))
// 	}
// 	return objects, nil
// }

// func (s3 *S3) StoreFile(destination_key string, source_filename string) error {
// 	reader, err := os.Open(source_filename)
// 	if err != nil {
// 		log.Fatal("KV cannot open file", source_filename)
// 	}
// 	fi, err := reader.Stat()
// 	if err != nil {
// 		log.Println("KV could not stat file")
// 		log.Fatal(err)
// 	}

// 	return store(s3, destination_key, fi.Size(), make(JSON, 0), reader)
// }

// ////////////////////////////
// // SUPPORT

// func mapToReader(json_map JSON) (io.Reader, int64) {
// 	b, _ := json.Marshal(json_map)
// 	r := bytes.NewReader(b)
// 	return r, int64(len(b))
// }
