func faceToASCIIImage(inputImage []byte) ([]byte, error) {
	// 1. Decode Input
	img, err := gocv.IMDecode(inputImage, gocv.IMReadColor)
	if err != nil || img.Empty() {
		return nil, fmt.Errorf("decoding image: %w", err)
	}
	defer img.Close()

	// 2. Create Blob for DNN Face Detection
	blob := gocv.BlobFromImage(img, 1.0, image.Pt(300, 300), gocv.NewScalar(104, 177, 123, 0), false, false)
	defer blob.Close()

	// 3. Run DNN Face Detection
	appPtr.net.SetInput(blob, "")
	detections := appPtr.net.Forward("")
	defer detections.Close()

	// 4. Parse Detections and Build Rectangle List
	confidenceThreshold := 0.5
	var rects []image.Rectangle

	for i := 0; i < detections.Total(); i += 7 {
		confidence := detections.GetFloatAt(0, i+2)

		if confidence > float32(confidenceThreshold) {
			// Get normalized coordinates
			left := detections.GetFloatAt(0, i+3)
			top := detections.GetFloatAt(0, i+4)
			right := detections.GetFloatAt(0, i+5)
			bottom := detections.GetFloatAt(0, i+6)

			// Convert to pixel coordinates
			x1 := int(left * float32(img.Cols()))
			y1 := int(top * float32(img.Rows()))
			x2 := int(right * float32(img.Cols()))
			y2 := int(bottom * float32(img.Rows()))

			rects = append(rects, image.Rect(x1, y1, x2, y2))
		}
	}

	if len(rects) == 0 {
		return nil, fmt.Errorf("no face detected")
	}

	// Pick the largest face
	faceRect := rects[0]
	for _, r := range rects {
		if r.Dx()*r.Dy() > faceRect.Dx()*faceRect.Dy() {
			faceRect = r
		}
	}

	// Clone to avoid memory issues
	region := img.Region(faceRect)
	faceROI := region.Clone()
	region.Close()
	defer faceROI.Close()

	// 4. Pre-process Face (Grayscale + Downscale)
	grayFace := gocv.NewMat()
	defer grayFace.Close()
	gocv.CvtColor(faceROI, &grayFace, gocv.ColorBGRToGray)

	// Grid Dimensions
	gridWidth := 80
	aspectRatio := float64(faceRect.Dy()) / float64(faceRect.Dx())
	gridHeight := int(float64(gridWidth) * aspectRatio * 0.55)

	smallFace := gocv.NewMat()
	defer smallFace.Close()
	gocv.Resize(grayFace, &smallFace, image.Point{X: gridWidth, Y: gridHeight}, 0, 0, gocv.InterpolationLinear)

	// 5. Create the Virtual Canvas
	charPixelWidth := 10
	charPixelHeight := 12

	canvasWidth := gridWidth * charPixelWidth
	canvasHeight := gridHeight * charPixelHeight

	// Create a blank black image
	canvas := gocv.NewMatWithSize(canvasHeight, canvasWidth, gocv.MatTypeCV8UC3)
	defer canvas.Close()
	black := gocv.NewScalar(0, 0, 0, 0)
	canvas.SetTo(black)

	// 6. Rendering Loop
	pixels, err := smallFace.DataPtrUint8()
	if err != nil {
		return nil, fmt.Errorf("getting pixel data: %w", err)
	}
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	for row := 0; row < gridHeight; row++ {
		for col := 0; col < gridWidth; col++ {
			brightness := pixels[row*gridWidth+col]
			charIndex := int(brightness) * (len(asciiChars) - 1) / 255
			char := string(asciiChars[charIndex])

			x := col * charPixelWidth
			y := (row + 1) * charPixelHeight

			gocv.PutText(&canvas, char, image.Point{X: x, Y: y},
				gocv.FontHersheyPlain, 0.8, white, 1)
		}
	}

	// 7. Encode to JPEG
	buf, err := gocv.IMEncode(gocv.JPEGFileExt, canvas)
	if err != nil {
		return nil, fmt.Errorf("encoding image: %w", err)
	}
	defer buf.Close()

	// 8. Deep Copy (Segfault Prevention)
	nativeData := buf.GetBytes()
	finalBytes := make([]byte, len(nativeData))
	copy(finalBytes, nativeData)

	return finalBytes, nil
}