package main

import (
	"context"
	"ecomdream/src/contracts"
)

type imageServiceServer struct{}

func (s *imageServiceServer) ValidateAndResizeImages(ctx context.Context, req *contracts.ValidateAndResizeImagesRequest) (*contracts.ValidateAndResizeImagesResponse, error) {
	var results []*contracts.Image

	for _, image := range req.Images {
		result := CheckConvertAndResizeImage(image)
		results = append(results, result)
	}

	return &contracts.ValidateAndResizeImagesResponse{Images: results}, nil
}
