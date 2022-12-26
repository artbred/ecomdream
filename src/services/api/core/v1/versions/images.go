package versions

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"github.com/h2non/bimg"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"strings"
)

var supportedImages = []string{
	"image/jpeg",
	"image/png",
	"image/jpg",
	"image/heic",
	"image/webp",
}

func arrayContains(arr []string, target string) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}

	return false
}

func checkAndConvertImage(fileHeader *multipart.FileHeader) (img []byte, err error) {
	if !arrayContains(supportedImages, fileHeader.Header.Get("Content-Type")) {
		return nil, errors.New(fmt.Sprintf("File %s has unsupported format", fileHeader.Filename))
	}

	imgFile, err := fileHeader.Open(); if err != nil {
		return nil, errors.New(fmt.Sprintf("Image %s is corrupted", fileHeader.Filename))
	}

	defer imgFile.Close()

	imgBytes, err := io.ReadAll(imgFile); if err != nil {
		return nil, errors.New(fmt.Sprintf("Image %s is corrupted", fileHeader.Filename))
	}

	imgBmg := bimg.NewImage(imgBytes); if imgBmg.Image() == nil {
		return nil, errors.New(fmt.Sprintf("Image %s is corrupted", fileHeader.Filename))
	}

	img, err = imgBmg.Convert(bimg.JPEG); if err != nil {
		return nil, errors.New(fmt.Sprintf("We can't convert image %s to jpeg, please convert it on your side\n", fileHeader.Filename))
	}

	return
}

func convertImagesToZip(form *multipart.Form) (zipOut *bytes.Buffer, err error) {
	var errorText string
	var images [][]byte

	imgCount := 0

	for _, fileHeaders := range form.File {
		for _, fileHeader := range fileHeaders {
			imgCount++

			img, err := checkAndConvertImage(fileHeader); if err != nil {
				errorText += err.Error() + ", "
				continue
			}

			images = append(images, img)
		}
	}

	if imgCount < 4 {
		return nil, errors.New("You must provide at least 4 images to successfully train AI")
	}

	if imgCount > 30 {
		return nil, errors.New("You can only upload 30 images")
	}

	if len(errorText) > 0 {
		i := strings.LastIndex(errorText, ", ")
		return nil, errors.New(errorText[:i] + strings.Replace(errorText[i:], ", ", "", 1))
	}

	zipOut = new(bytes.Buffer)

	zipWriter := zip.NewWriter(zipOut)
	defer zipWriter.Close()

	for i, img := range images {
		f, err := zipWriter.Create(fmt.Sprintf("%d.jpeg", i))
		if err != nil {
			logrus.Error(err)
			return nil, errors.New("Please try again later")
		}

		_, err = f.Write(img); if err != nil {
			logrus.Error(err)
			return nil, errors.New("Please try again later")
		}
	}

	return
}
