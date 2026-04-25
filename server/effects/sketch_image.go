package effects

import (
	"fmt"
	"image"

	"gocv.io/x/gocv"
)

// pencil sketch effect function
func PencilSketchImage(inputImage []byte) ([]byte, error) {
	// 1. Decode Image
	img, err := gocv.IMDecode(inputImage, gocv.IMReadColor)
	if img.Empty() || err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}
	defer img.Close()

	// 2. Convert to Grayscale
	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

	// 3. Invert the Grayscale
	inverted := gocv.NewMat()
	defer inverted.Close()
	gocv.BitwiseNot(gray, &inverted)

	// 4. Gaussian Blur the Inverted Image
	blurred := gocv.NewMat()
	defer blurred.Close()
	gocv.GaussianBlur(inverted, &blurred, image.Point{X: 51, Y: 51}, 0, 0, gocv.BorderDefault)

	// 5. MATH FIX: Floating Point Division
	// We cannot divide uint8 directly (results in 0). 
	// Convert both to Float32 for precision.
	grayFloat := gocv.NewMat()
	defer grayFloat.Close()
	gray.ConvertTo(&grayFloat, gocv.MatTypeCV32F)

	blurredFloat := gocv.NewMat()
	defer blurredFloat.Close()
	blurred.ConvertTo(&blurredFloat, gocv.MatTypeCV32F)

	// Perform Division: Gray / Blur
	sketchFloat := gocv.NewMat()
	defer sketchFloat.Close()
	gocv.Divide(grayFloat, blurredFloat, &sketchFloat)

	// 6. Scale and Convert back to 8-bit
	// The Color Dodge formula requires multiplying the result by 256.0.
	// ConvertToWithParams handles the scaling (alpha=256.0) and type conversion efficiently.
	sketch := gocv.NewMat()
	defer sketch.Close()
	sketchFloat.ConvertToWithParams(&sketch, gocv.MatTypeCV8U, 256.0, 0)

	// 7. Output Formatting (Back to BGR)
	result := gocv.NewMat()
	defer result.Close()
	gocv.CvtColor(sketch, &result, gocv.ColorGrayToBGR)

	// 8. Encode
	buf, err := gocv.IMEncode(gocv.JPEGFileExt, result)
	if err != nil {
		return nil, fmt.Errorf("encoding image: %w", err)
	}
	defer buf.Close()

	// 9. Deep Copy (Segfault Prevention)
	nativeData := buf.GetBytes()
	finalBytes := make([]byte, len(nativeData))
	copy(finalBytes, nativeData)

	return finalBytes, nil
}