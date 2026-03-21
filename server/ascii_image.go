package main

import (
	"bytes"
	"fmt"
	"image/jpeg"

	"github.com/qeesung/image2ascii/convert"
)

//convert an image to ascii characters
func imageToASCII(inputImage []byte) ([]byte, error) {
    // 1. Preprocess with gocv: enhance contrast for better ASCII output
    // preprocessed, err := preprocessForASCII(inputImage)
    // if err != nil {
    //     // Fall back to original image if preprocessing fails
    //     log.Printf("preprocessing failed, using original: %v", err)
    //     preprocessed = inputImage
    // }

    readerPtr := bytes.NewReader(inputImage)
    //decode image bytes to image.Image using image/jpeg
    jpegImage, err := jpeg.Decode(readerPtr)
    if err != nil {
        return nil, fmt.Errorf("imageToASCII %w", err)
    }
    fmt.Printf("width: %d, height: %d\n", jpegImage.Bounds().Dx(), jpegImage.Bounds().Dy())
    //convert image.Image to ASCII using image2ascii library
    image2ASCIIConverter := convert.NewImageConverter()
    convertOptions := convert.DefaultOptions
    convertOptions.Colored = false
    //width/height = 2.7
    convertOptions.FixedWidth = 40
    convertOptions.FixedHeight = 15
    asciiString := image2ASCIIConverter.Image2ASCIIString(jpegImage, &convertOptions)
    //send ASCII string as byte slice
    return []byte(asciiString), nil
}