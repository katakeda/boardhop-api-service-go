package utils

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"cloud.google.com/go/storage"

	firebase "firebase.google.com/go"
)

func GetDefaultBucket(ctx context.Context) (bucket *storage.BucketHandle, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to get default bucket | %w", err)
		}
	}()

	app, err := firebase.NewApp(ctx, &firebase.Config{
		StorageBucket: os.Getenv("FIREBASE_DEFAULT_BUCKET_NAME"),
	})

	client, err := app.Storage(ctx)

	bucket, err = client.DefaultBucket()

	return
}

func UploadFile(ctx context.Context, bucket *storage.BucketHandle, object string, file *multipart.FileHeader) error {
	f, err := file.Open()
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	objHandle := bucket.Object(object)

	wc := objHandle.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}
