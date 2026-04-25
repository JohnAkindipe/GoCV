package effects

import (
	"fmt"
	"image"

	"gocv.io/x/gocv"
)

// BlurImage applies a Gaussian blur to an in-memory JPEG image.
func BlurImage(inputImage []byte) ([]byte, error) {
	// 1. Decode
	img, err := gocv.IMDecode(inputImage, gocv.IMReadColor)
	if img.Empty() || err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}
	defer img.Close()

	// 2. Apply Gaussian blur (kernel 201x201)
	blurred := gocv.NewMat()
	defer blurred.Close()
	gocv.GaussianBlur(img, &blurred, image.Point{X: 51, Y: 51}, 0, 0, gocv.BorderDefault)

	// 3. Encode
	buf, err := gocv.IMEncode(gocv.JPEGFileExt, blurred)
	if err != nil {
		return nil, fmt.Errorf("encoding image: %w", err)
	}
	defer buf.Close()

	// 4. Deep copy to avoid segfault after buf.Close()
	nativeData := buf.GetBytes()
	finalBytes := make([]byte, len(nativeData))
	copy(finalBytes, nativeData)

	return finalBytes, nil
}