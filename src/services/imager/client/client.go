package client

import (
	"context"
	"ecomdream/src/contracts"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"os"
)

func ProcessAndUploadImages(request *contracts.ValidateAndResizeImagesRequest) (*contracts.ValidateAndResizeImagesResponse, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", os.Getenv("IMAGER_HOST"), os.Getenv("IMAGER_PORT")), grpc.WithInsecure())
	if err != nil {
		logrus.Error(err)
		return nil, errors.New("Please try again later")
	}

	defer conn.Close()
	c := contracts.NewImageServiceClient(conn)

	return c.ValidateAndResizeImages(context.Background(), request)
}
