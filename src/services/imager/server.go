package main

import (
	"archive/zip"
	"bytes"
	"context"
	"ecomdream/src/contracts"
	"ecomdream/src/pkg/storages/bucket"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
)

type imageServiceServer struct{}

func (s *imageServiceServer) ValidateAndResizeImages(ctx context.Context, req *contracts.ValidateAndResizeImagesRequest) (*contracts.ValidateAndResizeImagesResponse, error) {
	imagesCh := make(chan *contracts.Image)
	defer close(imagesCh)

	var imagesErrorText string

	for _, image := range req.Images {
		go func(image *contracts.Image) {
			result := CheckConvertAndResizeImage(image)
			imagesCh <- result
		}(image)
	}

	zipArchive := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipArchive)
	defer zipWriter.Close()

	for i := 0; i < len(req.Images); i++ {
		image := <-imagesCh; if len(image.Error) > 0 {
			imagesErrorText += image.Error + "\n"
			continue
		}

		f, err := zipWriter.Create(fmt.Sprintf("%d.jpeg", i))
		if err != nil {
			logrus.Error(err)
			return nil, errors.New("Please try again later")
		}

		_, err = f.Write(image.Data); if err != nil {
			logrus.Error(err)
			return nil, errors.New("Please try again later")
		}
	}

	if len(imagesErrorText) > 0 {
		return nil, errors.New(imagesErrorText)
	}

	zipURL, err := bucket.Upload(fmt.Sprintf("%s.zip", req.PaymentID), zipArchive, int64(zipArchive.Len()), "application/zip")
	if err != nil {
		logrus.WithError(err).Errorf("can't upload zip images for payment %s", req.PaymentID)
		return nil, errors.New("Please try again later")
	}

	return &contracts.ValidateAndResizeImagesResponse{ZipData: &contracts.ZipData{Url: zipURL}}, nil
}
