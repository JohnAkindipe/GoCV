package effects

import (
	"fmt"
	"image"

	"gocv.io/x/gocv"
)





func EmbossImage(inputImage []byte) ([]byte, error) {
	// 1. Decode
	img, err := gocv.IMDecode(inputImage, gocv.IMReadColor)
	if img.Empty() || err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}
	defer img.Close()

	// 2. Convert to grayscale
	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

	// 3. Apply emboss kernel via filter2D
	embossed := gocv.NewMat()
	defer embossed.Close()

	// Emboss kernel — simulates light hitting from top-left
	kernel := gocv.NewMatWithSize(3, 3, gocv.MatTypeCV32F)
	defer kernel.Close()
	kernel.SetFloatAt(0, 0, -2)
	kernel.SetFloatAt(0, 1, -1)
	kernel.SetFloatAt(0, 2, 0)
	kernel.SetFloatAt(1, 0, -1)
	kernel.SetFloatAt(1, 1, 1)
	kernel.SetFloatAt(1, 2, 1)
	kernel.SetFloatAt(2, 0, 0)
	kernel.SetFloatAt(2, 1, 1)
	kernel.SetFloatAt(2, 2, 2)

	// Apply the kernel — delta of 128 shifts result to mid-gray baseline
	gocv.Filter2D(gray, &embossed, -1, kernel, image.Point{X: -1, Y: -1}, 128, gocv.BorderDefault)

	// 4. Encode
	buf, err := gocv.IMEncode(gocv.JPEGFileExt, embossed)
	if err != nil {
		return nil, fmt.Errorf("encoding image: %w", err)
	}
	defer buf.Close()

	// 5. Deep copy to avoid segfault after buf.Close()
	nativeData := buf.GetBytes()
	finalBytes := make([]byte, len(nativeData))
	copy(finalBytes, nativeData)

	return finalBytes, nil
}