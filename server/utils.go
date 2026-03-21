package main

import (
	"fmt"
	"image"

	"gocv.io/x/gocv"
)

// preprocessForASCII enhances image contrast using CLAHE histogram equalization
// to produce sharper, more defined ASCII art output.
func preprocessForASCII(inputImage []byte) ([]byte, error) {
    // Decode with gocv
    mat, err := gocv.IMDecode(inputImage, gocv.IMReadColor)
    if err != nil || mat.Empty() {
        return nil, fmt.Errorf("decoding image for preprocessing: %w", err)
    }
    defer mat.Close()

    // Convert to grayscale for CLAHE
    gray := gocv.NewMat()
    defer gray.Close()
    gocv.CvtColor(mat, &gray, gocv.ColorBGRToGray)

    // Apply CLAHE (Contrast Limited Adaptive Histogram Equalization)
    // clipLimit=2.0, tileGridSize=8x8 for balanced contrast enhancement
    clahe := gocv.NewCLAHEWithParams(2.0, image.Pt(8, 8))
    defer clahe.Close()

    enhanced := gocv.NewMat()
    defer enhanced.Close()
    clahe.Apply(gray, &enhanced)

    // Convert back to BGR for JPEG encoding
    result := gocv.NewMat()
    defer result.Close()
    gocv.CvtColor(enhanced, &result, gocv.ColorGrayToBGR)

    // Encode back to JPEG with high quality
    buf, err := gocv.IMEncodeWithParams(gocv.JPEGFileExt, result, []int{gocv.IMWriteJpegQuality, 95})
    if err != nil {
        return nil, fmt.Errorf("encoding preprocessed image: %w", err)
    }
    defer buf.Close()

    // Deep copy to avoid segfault
    nativeData := buf.GetBytes()
    finalBytes := make([]byte, len(nativeData))
    copy(finalBytes, nativeData)

    return finalBytes, nil
}

// faceDetect detects a face using YuNet and draws a green rectangle around it
// func faceDetect(inputImage []byte) ([]byte, error) {
//     // 1. Decode Input Image
//     img, err := gocv.IMDecode(inputImage, gocv.IMReadColor)
//     if err != nil || img.Empty() {
//         return nil, fmt.Errorf("decoding image: %w", err)
//     }
//     defer img.Close()

//     // 2. Create YuNet Face Detector
//     // Parameters: model path, config (empty for ONNX), input size, score threshold, NMS threshold, top_k
//     faceDetector := gocv.NewFaceDetectorYN(YUNET_MODEL_PATH, "", image.Pt(img.Cols(), img.Rows()))
//     defer faceDetector.Close()

//     // Set detection parameters
//     faceDetector.SetScoreThreshold(0.5)
//     faceDetector.SetNMSThreshold(0.3)
//     faceDetector.SetTopK(5000)

//     // 3. Detect Faces
//     faces := gocv.NewMat()
//     defer faces.Close()
    
//     faceDetector.Detect(img, &faces)

//     if faces.Empty() || faces.Rows() == 0 {
//         return nil, fmt.Errorf("no face detected")
//     }

//     // 4. Find face with highest confidence
//     var bestConfidence float32 = 0.0
//     var bestRect image.Rectangle

//     // YuNet returns faces as rows with 15 columns:
//     // [x, y, w, h, x_re, y_re, x_le, y_le, x_nt, y_nt, x_rcm, y_rcm, x_lcm, y_lcm, score]
//     for i := 0; i < faces.Rows(); i++ {
//         confidence := faces.GetFloatAt(i, 14)
        
//         if confidence > bestConfidence {
//             bestConfidence = confidence
            
//             x := int(faces.GetFloatAt(i, 0))
//             y := int(faces.GetFloatAt(i, 1))
//             w := int(faces.GetFloatAt(i, 2))
//             h := int(faces.GetFloatAt(i, 3))
            
//             bestRect = image.Rect(x, y, x+w, y+h)
//         }
//     }

//     // 5. Draw rectangle on the face with highest confidence
//     green := color.RGBA{R: 0, G: 255, B: 0, A: 255}
//     gocv.Rectangle(&img, bestRect, green, 3)

//     // 6. Encode to JPEG
//     buf, err := gocv.IMEncode(gocv.JPEGFileExt, img)
//     if err != nil {
//         return nil, fmt.Errorf("encoding image: %w", err)
//     }
//     defer buf.Close()

//     // 7. Deep Copy (Segfault Prevention)
//     nativeData := buf.GetBytes()
//     finalBytes := make([]byte, len(nativeData))
//     copy(finalBytes, nativeData)

//     return finalBytes, nil
// }