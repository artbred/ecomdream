package versions

import (
	"archive/zip"
	"bytes"
	"ecomdream/src/contracts"
	"ecomdream/src/services/imager/client"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
)

func processImagesToZip(form *multipart.Form) (zipOut *bytes.Buffer, err error) {
	var inputImages []*contracts.Image

	for _, fileHeaders := range form.File {
		for _, fileHeader := range fileHeaders {
			imgFile, err := fileHeader.Open(); if err != nil {
				return nil, errors.New(fmt.Sprintf("Image %s is corrupted", fileHeader.Filename))
			}

			defer imgFile.Close()

			imgBytes, err := io.ReadAll(imgFile); if err != nil {
				return nil, errors.New(fmt.Sprintf("Image %s is corrupted", fileHeader.Filename))
			}

			inputImages = append(inputImages, &contracts.Image{
				Id: fileHeader.Filename,
				Data: imgBytes,
				ContentType: fileHeader.Header.Get("Content-Type"),
			})
		}
	}

	if len(inputImages) < 5 {
		return nil, errors.New("You must provide at least 5 images to successfully train AI")
	}

	if len(inputImages) > 30 {
		return nil, errors.New("You can only upload 30 images")
	}

	outImages, err := client.ProcessImages(inputImages); if err != nil {
		return nil, err
	}

	zipOut = new(bytes.Buffer)

	zipWriter := zip.NewWriter(zipOut)
	defer zipWriter.Close()

	var errText string

	for i, img := range outImages {
		if len(img.Error) > 0 {
			errText += img.Error + "\n"
			continue
		}

		f, err := zipWriter.Create(fmt.Sprintf("%d.jpeg", i))
		if err != nil {
			logrus.Error(err)
			return nil, errors.New("Please try again later")
		}

		_, err = f.Write(img.Data); if err != nil {
			logrus.Error(err)
			return nil, errors.New("Please try again later")
		}
	}

	if len(errText) > 0 {
		return nil, errors.New(errText)
	}

	return
}
