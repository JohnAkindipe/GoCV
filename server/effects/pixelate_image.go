package effects

import (
	"fmt"
	"image"

	"gocv.io/x/gocv"
)


func PixelateImage(inputImage []byte) ([]byte, error) {
	// 1. Decode
	img, err := gocv.IMDecode(inputImage, gocv.IMReadColor)
	if img.Empty() || err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}
	defer img.Close()

	rows := img.Rows()
	cols := img.Cols()

	const blockSize = 8

	// 2. Shrink down — averaging naturally happens during downscale
	small := gocv.NewMat()
	defer small.Close()
	gocv.Resize(img, &small, image.Point{X: cols / blockSize, Y: rows / blockSize}, 0, 0, gocv.InterpolationLinear)

	// 3. Scale back up with nearest neighbor to produce hard block edges
	pixelated := gocv.NewMat()
	defer pixelated.Close()
	gocv.Resize(small, &pixelated, image.Point{X: cols, Y: rows}, 0, 0, gocv.InterpolationNearestNeighbor)

	// 4. Encode
	buf, err := gocv.IMEncode(gocv.JPEGFileExt, pixelated)
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