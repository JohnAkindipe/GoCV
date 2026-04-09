package effects

import (
	"fmt"
	"image/color"
	"math"

	"gocv.io/x/gocv"
)





func WaveRippleImage(inputImage []byte) ([]byte, error) {
	// 1. Decode
	img, err := gocv.IMDecode(inputImage, gocv.IMReadColor)
	if img.Empty() || err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}
	defer img.Close()

	rows := img.Rows()
	cols := img.Cols()

	// Gentle water ripple parameters
	const (
		amplitudeX  = 5.0
		amplitudeY  = 5.0
		wavelengthX = 120.0
		wavelengthY = 120.0
	)

	// 2. Build map matrices for remapping
	mapX := gocv.NewMatWithSize(rows, cols, gocv.MatTypeCV32F)
	defer mapX.Close()
	mapY := gocv.NewMatWithSize(rows, cols, gocv.MatTypeCV32F)
	defer mapY.Close()

	// 3. Populate the remap matrices with sine wave offsets
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			srcX := float32(x) + float32(amplitudeX)*float32(math.Sin(2*math.Pi*float64(y)/wavelengthX))
			srcY := float32(y) + float32(amplitudeY)*float32(math.Sin(2*math.Pi*float64(x)/wavelengthY))

			mapX.SetFloatAt(y, x, srcX)
			mapY.SetFloatAt(y, x, srcY)
		}
	}

	// 4. Remap
	rippled := gocv.NewMat()
	defer rippled.Close()
	gocv.Remap(img, &rippled, &mapX, &mapY, gocv.InterpolationLinear, gocv.BorderReflect, color.RGBA{})

	// 5. Encode
	buf, err := gocv.IMEncode(gocv.JPEGFileExt, rippled)
	if err != nil {
		return nil, fmt.Errorf("encoding image: %w", err)
	}
	defer buf.Close()

	// 6. Deep copy to avoid segfault after buf.Close()
	nativeData := buf.GetBytes()
	finalBytes := make([]byte, len(nativeData))
	copy(finalBytes, nativeData)

	return finalBytes, nil
}