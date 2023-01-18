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

func ProcessImages(images []*contracts.Image) ([]*contracts.Image, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", os.Getenv("IMAGER_HOST"), os.Getenv("IMAGER_PORT")), grpc.WithInsecure())
	if err != nil {
		logrus.Error(err)
		return nil, errors.New("Please try again later")
	}
	defer conn.Close()

	c := contracts.NewImageServiceClient(conn)

	res, err := c.ValidateAndResizeImages(context.Background(), &contracts.ValidateAndResizeImagesRequest{Images: images})
	if err != nil {
		return nil, fmt.Errorf("error while calling SendImages: %v", err)
	}

	return res.Images, nil
}
