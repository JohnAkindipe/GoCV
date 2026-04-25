package effects

import (
	"fmt"
	"image"
	"math/rand"
	"time"

	"gocv.io/x/gocv"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GlitchImage(inputImage []byte) ([]byte, error) {
	fmt.Println("Applying glitch effect")
	// 1. Decode
	mat, err := gocv.IMDecode(inputImage, gocv.IMReadColor)
	if mat.Empty() || err != nil {
		return nil, fmt.Errorf("error reading image: %w", err)
	}
	defer mat.Close()

	rows, cols := mat.Rows(), mat.Cols()

	// 2. Channel Split & Chromatic Aberration (Vectorized)
	// We use WarpAffine to shift whole channels instantly rather than pixel loops
	channels := gocv.Split(mat)
	defer channels[0].Close()
	defer channels[1].Close()
	defer channels[2].Close()

	// Random shifts for RGB
	// Blue: Shift Left
	shiftedBlue := shiftChannel(channels[0], float64(-rand.Intn(20)-10), 0)
	defer shiftedBlue.Close()
	// Green: Shift vertical slightly
	shiftedGreen := shiftChannel(channels[1], float64(rand.Intn(10)-5), float64(rand.Intn(10)-5))
	defer shiftedGreen.Close()
	// Red: Shift Right
	shiftedRed := shiftChannel(channels[2], float64(rand.Intn(20)+10), 0)
	defer shiftedRed.Close()

	// Merge back
	result := gocv.NewMat()
	defer result.Close()
	gocv.Merge([]gocv.Mat{shiftedBlue, shiftedGreen, shiftedRed}, &result)

	// 3. Displacement Bands (Using ROI - Region of Interest)
	// We cut strips and move them. Much faster than pixel loops.
	numBands := rand.Intn(20) + 10
	for i := 0; i < numBands; i++ {
		y := rand.Intn(rows - 10)
		height := rand.Intn(30) + 5
		if y+height >= rows {
			continue
		}

		shift := rand.Intn(50) - 25 // Shift left or right
		if shift == 0 {
			continue
		}

		// Define the source region (strip)
		rect := image.Rect(0, y, cols, y+height)
		roi := result.Region(rect)

		// Create a displaced copy
		shiftedRoi := gocv.NewMat()
		
		// Use WarpAffine again for the strip shift to handle borders gracefully
		shiftMat := gocv.NewMatWithSize(2, 3, gocv.MatTypeCV32F)
		shiftMat.SetFloatAt(0, 0, 1)
		shiftMat.SetFloatAt(0, 1, 0)
		shiftMat.SetFloatAt(0, 2, float32(shift))
		shiftMat.SetFloatAt(1, 0, 0)
		shiftMat.SetFloatAt(1, 1, 1)
		shiftMat.SetFloatAt(1, 2, 0)
		
		gocv.WarpAffine(roi, &shiftedRoi, shiftMat, image.Point{X: cols, Y: height})
		
		// Copy shifted strip back onto the main image
		shiftedRoi.CopyTo(&roi)
		
		roi.Close()
		shiftedRoi.Close()
		shiftMat.Close()
	}

	// 4. Pixel Level Effects (Smear, Noise, Scanlines)
	// OPTIMIZATION: Instead of calling SetUCharAt 2 million times,
	// we get the raw byte slice once. This is 100x faster.
	applyRawPixelEffects(&result)

	// 5. Encode
	buf, err := gocv.IMEncode(gocv.JPEGFileExt, result)
	if err != nil {
		return nil, fmt.Errorf("encoding image: %w", err)
	}
	defer buf.Close()

	nativeBytes := buf.GetBytes()
	imageBytes := make([]byte, len(nativeBytes))
	copy(imageBytes, nativeBytes)

	return imageBytes, nil
}

// shiftChannel creates a translation matrix and warps the image (fast)
func shiftChannel(src gocv.Mat, xShift, yShift float64) gocv.Mat {
	dst := gocv.NewMat()
	// Create a 2x3 affine transform matrix
	m := gocv.NewMatWithSize(2, 3, gocv.MatTypeCV32F)
	defer m.Close()

	// [ 1, 0, xShift ]
	// [ 0, 1, yShift ]
	m.SetFloatAt(0, 0, 1)
	m.SetFloatAt(0, 1, 0)
	m.SetFloatAt(0, 2, float32(xShift))
	m.SetFloatAt(1, 0, 0)
	m.SetFloatAt(1, 1, 1)
	m.SetFloatAt(1, 2, float32(yShift))

	// Apply affine wrap
	gocv.WarpAffine(src, &dst, m, image.Point{X: src.Cols(), Y: src.Rows()})
	return dst
}

// applyRawPixelEffects handles smearing, noise, and scanlines via direct memory access
func applyRawPixelEffects(img *gocv.Mat) {
	// Unsafe access to underlying C++ data for max speed
	// Note: In GoCV, DataPtrUint8 returns a slice header pointing to C memory.
	// Modification is immediate.
	data, _ := img.DataPtrUint8()
	rows := img.Rows()
	cols := img.Cols()
	channels := 3 // BGR

	stride := cols * channels

	// Pre-calculate random thresholds to avoid math in the loop
	noiseThreshold := 15 // 15% chance
	
	for y := 0; y < rows; y++ {
		rowStart := y * stride
		
		// Determine if this row should smear
		isSmearRow := rand.Intn(100) < 5 // 5% of rows have heavy smear
		smearCount := 0
		var smearB, smearG, smearR uint8
		
		// Determine Scanline intensity for this row
		scanlineFactor := 1.0
		if y%2 == 0 {
			scanlineFactor = 0.7 // Darken alternating rows
		}

		for x := 0; x < cols; x++ {
			pixelIndex := rowStart + (x * channels)

			// 1. Smear Effect
			if isSmearRow {
				// Start a smear
				if smearCount == 0 && rand.Intn(20) == 0 {
					smearCount = rand.Intn(50) + 10
					smearB = data[pixelIndex]
					smearG = data[pixelIndex+1]
					smearR = data[pixelIndex+2]
				}
				
				// Apply smear
				if smearCount > 0 {
					data[pixelIndex] = smearB
					data[pixelIndex+1] = smearG
					data[pixelIndex+2] = smearR
					smearCount--
				}
			}

			// 2. Color Noise (Random static)
			if rand.Intn(100) < noiseThreshold {
				// Add simple noise
				noise := uint8(rand.Intn(100))
				// Using simple safe addition (clamping happens automatically on overflow/wrap for uint8)
				// Use addition to brighten, subtraction to darken, or just replace
				data[pixelIndex] = data[pixelIndex] + noise
				data[pixelIndex+1] = data[pixelIndex+1] + noise
				data[pixelIndex+2] = data[pixelIndex+2] + noise
			}

			// 3. Scanlines & Darkening (Apply intensity)
			if scanlineFactor < 1.0 {
				data[pixelIndex] = uint8(float64(data[pixelIndex]) * scanlineFactor)
				data[pixelIndex+1] = uint8(float64(data[pixelIndex+1]) * scanlineFactor)
				data[pixelIndex+2] = uint8(float64(data[pixelIndex+2]) * scanlineFactor)
			}
		}
	}
}