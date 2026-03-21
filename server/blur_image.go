package main

import (
	"fmt"
	"image"

	"gocv.io/x/gocv"
)

// blur image function
func blurImage(inputPath, outputPath string) bool {
    // Read the image
    img := gocv.IMRead(inputPath, gocv.IMReadColor)
    if img.Empty() {
        fmt.Println("Error: Could not read image")
        return false
    }
    defer img.Close()

    // Create a Mat to hold the blurred image
    blurred := gocv.NewMat()
    defer blurred.Close()

    // Apply Gaussian blur
    // Parameters:
    // - img: source image
    // - &blurred: destination image
    // - image.Point{X: 201, Y: 201}: kernel size (must be odd numbers)
    // - 0: sigmaX (0 means calculate from kernel size)
    // - 0: sigmaY (0 means calculate from kernel size)
    gocv.GaussianBlur(img, &blurred, image.Point{X: 201, Y: 201}, 0, 0, gocv.BorderDefault)

    // Save the blurred image
    return gocv.IMWrite(outputPath, blurred)
}