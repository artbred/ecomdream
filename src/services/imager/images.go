package main

import (
	"ecomdream/src/contracts"
	"fmt"
	"github.com/h2non/bimg"
	"github.com/sirupsen/logrus"
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

func CheckConvertAndResizeImage(inputImage *contracts.Image) *contracts.Image {
	if !arrayContains(supportedImages, inputImage.ContentType) {
		inputImage.Error = fmt.Sprintf("Image %s has unsupported format\n", inputImage.Id)
		return inputImage
	}

	imgBmg := bimg.NewImage(inputImage.Data); if imgBmg.Image() == nil {
		inputImage.Error = fmt.Sprintf("Image %s is corrupted\n", inputImage.Id)
		return inputImage
	}

	size, err := imgBmg.Size(); if err != nil {
		inputImage.Error = fmt.Sprintf("Image %s is corrupted\n", inputImage.Id)
		return inputImage
	}

	if size.Height < 512 || size.Width < 512 {
		inputImage.Error = fmt.Sprintf("Image %s must be 512x512 or more\n", inputImage.Id)
		return inputImage
	}

	imgJpeg, err := imgBmg.Convert(bimg.JPEG); if err != nil {
		logrus.Error(err)
		inputImage.Error = fmt.Sprintf("We can't convert image %s to jpeg, please convert it on your side\n", inputImage.Id)
		return inputImage
	}

	inputImage.Data, err = bimg.Resize(imgJpeg, bimg.Options{
		Height: 512,
		Width: 512,
	})

	if err != nil {
		logrus.Error(err)
		inputImage.Error = fmt.Sprintf("We can't resize image %s to 512x512, please resize it on your side\n")
		return inputImage
	}

	return inputImage
}
