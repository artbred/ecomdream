package versions

import (
	"ecomdream/src/contracts"
	"ecomdream/src/services/imager/client"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
)

func validateAndUploadImages(paymentID string, form *multipart.Form) (imagerResponse *contracts.ValidateAndResizeImagesResponse, err error) {
	var inputImages []*contracts.Image

	for _, fileHeaders := range form.File {
		for _, fileHeader := range fileHeaders {
			imgFile, err := fileHeader.Open(); if err != nil {
				return nil, fmt.Errorf("Image %s is corrupted", fileHeader.Filename)
			}

			defer imgFile.Close()

			imgBytes, err := io.ReadAll(imgFile); if err != nil {
				return nil, fmt.Errorf("Image %s is corrupted", fileHeader.Filename)
			}

			inputImages = append(inputImages, &contracts.Image{
				Id: strings.ToValidUTF8(fileHeader.Filename, " "),
				Data: imgBytes,
			})
		}
	}

	if len(inputImages) < 5 {
		return nil, errors.New("You must provide at least 5 images to successfully train AI")
	}

	if len(inputImages) > 30 {
		return nil, errors.New("You can only upload 30 images")
	}

	req := &contracts.ValidateAndResizeImagesRequest{
		Images: inputImages,
		PaymentID: paymentID,
	}

	imagerResponse, err = client.ProcessAndUploadImages(req)
	return
}
