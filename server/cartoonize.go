package main

import (
	"fmt"
	"image"

	"gocv.io/x/gocv"
)

// cartoonize image function
func cartoonizeImage(inputPath, outputPath string) bool {
    // Read the image
    img := gocv.IMRead(inputPath, gocv.IMReadColor)
    if img.Empty() {
        fmt.Println("Error: Could not read image")
        return false
    }
    defer img.Close()

    // Step 1: Downscale for faster processing and smoother result
    // then upscale back - this helps reduce noise
    smallImg := gocv.NewMat()
    defer smallImg.Close()
    
    // Step 2: Apply bilateral filter multiple times for strong color smoothing
    // while preserving edges - this creates the "painted" look
    colorReduced := gocv.NewMat()
    defer colorReduced.Close()
    
    temp1 := gocv.NewMat()
    defer temp1.Close()
    temp2 := gocv.NewMat()
    defer temp2.Close()
    
    // Apply bilateral filter 4 times for stronger posterization effect
    // diameter=9, sigmaColor=150, sigmaSpace=150 for aggressive smoothing
    gocv.BilateralFilter(img, &temp1, 9, 150, 150)
    gocv.BilateralFilter(temp1, &temp2, 9, 150, 150)
    gocv.BilateralFilter(temp2, &temp1, 9, 150, 150)
    gocv.BilateralFilter(temp1, &colorReduced, 9, 150, 150)

    // Step 3: Convert to grayscale for edge detection
    gray := gocv.NewMat()
    defer gray.Close()
    gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

    // Step 4: Apply stronger median blur to reduce noise before edge detection
    // Larger kernel size removes more noise/spots
    gocv.MedianBlur(gray, &gray, 11)

    // Step 5: Detect edges using adaptive threshold with larger block size
    // Larger block size (15 instead of 9) creates thicker, cleaner edges
    // Higher constant (5 instead of 2) removes more weak edges
    edges := gocv.NewMat()
    defer edges.Close()
    gocv.AdaptiveThreshold(gray, &edges, 255, gocv.AdaptiveThresholdMean, gocv.ThresholdBinary, 15, 5)

    // Step 6: Apply morphological operations to clean up edges
    // Dilate to thicken edges, then erode to remove small noise
    kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Point{X: 2, Y: 2})
    defer kernel.Close()
    
    cleanedEdges := gocv.NewMat()
    defer cleanedEdges.Close()
    gocv.Dilate(edges, &cleanedEdges, kernel)
    gocv.Erode(cleanedEdges, &cleanedEdges, kernel)

    // Step 7: Convert edges to 3-channel for combining with color image
    edgesColor := gocv.NewMat()
    defer edgesColor.Close()
    gocv.CvtColor(cleanedEdges, &edgesColor, gocv.ColorGrayToBGR)

    // Step 8: Combine edges with color-reduced image
    cartoon := gocv.NewMat()
    defer cartoon.Close()
    gocv.BitwiseAnd(colorReduced, edgesColor, &cartoon)

    // Save the cartoonized image
    return gocv.IMWrite(outputPath, cartoon)
}