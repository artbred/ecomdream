package bucket

import (
	"context"
	"ecomdream/src/pkg/configs"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var (
	client      *minio.Client
	bucketName  string
	cdnEndpoint string
	basePath    = "files"
)

func Upload(fileName string, file io.Reader, size int64, contentType string) (fileUrl string, err error) {
	_, err = client.PutObject(context.Background(), bucketName, fmt.Sprintf("%s/%s", basePath, fileName), file, size,
		minio.PutObjectOptions{
			UserMetadata: map[string]string{"x-amz-acl": "public-read"},
			ContentType:  contentType,
		},
	)

	fileUrl = fmt.Sprintf("%s/%s/%s", cdnEndpoint, basePath, fileName)
	return
}

func DeleteFile(fileName string) (err error) {
	err = client.RemoveObject(context.Background(), bucketName,
		fmt.Sprintf("%s/%s", basePath, fileName),
		minio.RemoveObjectOptions{
			ForceDelete: true,
		},
	)

	if err != nil {
		logrus.Error(err)
	}

	return
}

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Error(err)
	}

	secretKey := os.Getenv("MINIO_SECRET_KEY")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	endpoint := os.Getenv("MINIO_BUCKET_ENDPOINT")

	Client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})

	if err != nil {
		logrus.Error(err)
	}

	cdnEndpoint = os.Getenv("MINIO_CDN_ENDPOINT")
	bucketName = os.Getenv("MINIO_BUCKET_NAME")
	client = Client

	if configs.Debug {
		basePath = "dev"
	}
}
