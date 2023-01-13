package client

import (
	"context"
	"ecomdream/src/contracts"
	"fmt"
	"google.golang.org/grpc"
	"os"
)

func ProcessImages(images []*contracts.Image) ([]*contracts.Image, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", os.Getenv("IMAGER_HOST"), os.Getenv("IMAGER_PORT")), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()

	c := contracts.NewImageServiceClient(conn)

	res, err := c.ValidateAndResizeImages(context.Background(), &contracts.ValidateAndResizeImagesRequest{Images: images})
	if err != nil {
		return nil, fmt.Errorf("error while calling SendImages: %v", err)
	}

	return res.Images, nil
}
