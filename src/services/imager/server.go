package main

import (
	"context"
	"ecomdream/src/contracts"
)

type imageServiceServer struct{}

func (s *imageServiceServer) ValidateAndResizeImages(ctx context.Context, req *contracts.ValidateAndResizeImagesRequest) (*contracts.ValidateAndResizeImagesResponse, error) {
	result := &contracts.ValidateAndResizeImagesResponse{}
	resultCh := make(chan *contracts.Image)

	for _, image := range req.Images {
		go func(image *contracts.Image) {
			result := CheckConvertAndResizeImage(image)
			resultCh <- result
		}(image)
	}

	for i := 0; i < len(req.Images); i++ {
		result.Images = append(result.Images, <-resultCh)
	}

	return result, nil
}
