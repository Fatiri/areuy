package storage

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// Storage Storage service interface.
type Storage interface {
	Put(filename string, file io.Reader) error
	Delete(filename string) error
	URL(filename string) string
	READ(filename string) error
}

type gcsCtx struct {
	gcsClient  *storage.Client
	bucketName string
	configPath string
}

const (
	gcsBaseDomain = "storage.googleapis.com"
)

var (
	googleCloudStorageClientOnce sync.Once
	googleCloudStorageClient     Storage
)

func GoogleCloudStorage(gcs gcsCtx) Storage {
	// GoogleCloudStorage Singletion to get Google cloud storage.
	googleCloudStorageClientOnce.Do(func() {
		googleCloudStorageClient = NewGCS(gcs)
	})

	return googleCloudStorageClient
}

// GCS Construct new GCS implementation object.
func NewGCS(gcs gcsCtx) Storage {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(gcs.configPath))
	if err != nil {
		log.Fatal(err)
	}

	return &gcsCtx{
		gcsClient:  client,
		bucketName: gcs.bucketName,
		configPath: gcs.configPath,
	}
}

// Put Send upload request to google cloud storage.
func (c gcsCtx) Put(filename string, file io.Reader) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	w := c.gcsClient.Bucket(c.bucketName).Object(filename).NewWriter(ctx)
	w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	reader := bufio.NewReader(file)
	if _, err := reader.WriteTo(w); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// Put Send upload request to google cloud storage.
func (c gcsCtx) READ(filename string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	rc, err := c.gcsClient.Bucket("algo-dev").Object("shipment/banner/SDRHP/1585805344").NewReader(ctx)
	if err != nil {
		return err
	}
	defer rc.Close()
	return nil
}

// Delete Delete request to google cloud storage.
func (c *gcsCtx) Delete(filename string) error {
	ctx := context.Background()
	o := c.gcsClient.Bucket(c.bucketName).Object(filename)
	if err := o.Delete(ctx); err != nil {
		return err
	}

	return nil
}

// URL Generate HTTPs url to fetch google clould file
func (c *gcsCtx) URL(filename string) string {
	return fmt.Sprintf("https://%s/%s/%s", gcsBaseDomain, c.bucketName, filename)
}
