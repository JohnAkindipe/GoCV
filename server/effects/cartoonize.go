package effects

import (
	"fmt"
	"image"

	"gocv.io/x/gocv"
)

// CartoonizeImage applies a cartoon effect to an in-memory JPEG image.
func CartoonizeImage(inputImage []byte) ([]byte, error) {
	// 1. Decode
	img, err := gocv.IMDecode(inputImage, gocv.IMReadColor)
	if img.Empty() || err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}
	defer img.Close()

	// 2. Apply bilateral filter 4 times for strong color smoothing
	// while preserving edges — creates the "painted" look
	colorReduced := gocv.NewMat()
	defer colorReduced.Close()
	temp1 := gocv.NewMat()
	defer temp1.Close()
	temp2 := gocv.NewMat()
	defer temp2.Close()

	gocv.BilateralFilter(img, &temp1, 9, 150, 150)
	gocv.BilateralFilter(temp1, &temp2, 9, 150, 150)
	gocv.BilateralFilter(temp2, &temp1, 9, 150, 150)
	gocv.BilateralFilter(temp1, &colorReduced, 9, 150, 150)

	// 3. Convert to grayscale for edge detection
	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

	// 4. Median blur to reduce noise before edge detection
	gocv.MedianBlur(gray, &gray, 11)

	// 5. Detect edges via adaptive threshold
	edges := gocv.NewMat()
	defer edges.Close()
	gocv.AdaptiveThreshold(gray, &edges, 255, gocv.AdaptiveThresholdMean, gocv.ThresholdBinary, 15, 5)

	// 6. Morphological cleanup: dilate then erode
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Point{X: 2, Y: 2})
	defer kernel.Close()
	cleanedEdges := gocv.NewMat()
	defer cleanedEdges.Close()
	gocv.Dilate(edges, &cleanedEdges, kernel)
	gocv.Erode(cleanedEdges, &cleanedEdges, kernel)

	// 7. Convert edges to 3-channel for blending
	edgesColor := gocv.NewMat()
	defer edgesColor.Close()
	gocv.CvtColor(cleanedEdges, &edgesColor, gocv.ColorGrayToBGR)

	// 8. Combine edges with color-reduced image
	cartoon := gocv.NewMat()
	defer cartoon.Close()
	gocv.BitwiseAnd(colorReduced, edgesColor, &cartoon)

	// 9. Encode
	buf, err := gocv.IMEncode(gocv.JPEGFileExt, cartoon)
	if err != nil {
		return nil, fmt.Errorf("encoding image: %w", err)
	}
	defer buf.Close()

	// 10. Deep copy to avoid segfault after buf.Close()
	nativeData := buf.GetBytes()
	finalBytes := make([]byte, len(nativeData))
	copy(finalBytes, nativeData)

	return finalBytes, nil
}
