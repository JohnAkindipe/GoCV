package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/qeesung/image2ascii/convert"
	"gocv.io/x/gocv"
)

const (
    KB = 1024
    IMAGE_PATH = "image_file.jpg"
    EDITED_IMAGE_PATH = "edited_image_file.jpg"
    // Use full path for Linux server
	MODEL_PATH = "/root/models/face-detect/res10_300x300_ssd_iter_140000.caffemodel"
	CONFIG_PATH= "/root/models/face-detect/deploy.prototxt"
	YUNET_MODEL_PATH = "/root/models/face-detect/face_detection_yunet_2023mar.onnx"
)

// ASCII set from dark to light
const asciiChars = "@%#*+=-:. "

// func tcp() {
//     // Listen for incoming connections
//     listener, err := net.Listen("tcp", ":8080")
//     if err != nil {
//         fmt.Println("Error:", err)
//         return
//     }
//     defer listener.Close()

//     fmt.Println("Server is listening on port 8080")

//     for {
//         // Accept incoming connections
//         conn, err := listener.Accept()
//         if err != nil {
//             fmt.Println("Error:", err)
//             continue
//         }

//         go handleClient(conn)
//     }
// }

type app struct {
	net gocv.Net
}

func httpServer() error {
	srvrPtr := &http.Server{
		Addr:    ":4000",
		Handler: routes(),
	}

	log.Printf("Starting server on %s", srvrPtr.Addr)
	return srvrPtr.ListenAndServe()
}

var appPtr *app

func main() {
    // Create/open log file
    logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatal("Failed to open log file:", err)
    }
    defer logFile.Close()

    // Set log output to both file and stdout
    multiWriter := io.MultiWriter(os.Stdout, logFile)
    log.SetOutput(multiWriter)
    
    // Add timestamps and file/line info to logs
    log.SetFlags(log.LstdFlags | log.Lshortfile)

	//Initialize gocv.net on app
	appPtr = &app{}
	appPtr.net = gocv.ReadNet(MODEL_PATH, CONFIG_PATH)
    if appPtr.net.Empty() {
        log.Fatal("failed to load DNN model")
    }
    defer appPtr.net.Close()

    err = httpServer() 
    if err != nil {
        log.Fatal(err)
    }
}

// func handleClient(conn net.Conn) {
//     defer conn.Close()

//     //create a file to hold the incoming stream
//     imageFile, closeFile := createFile(IMAGE_PATH) 
//     if imageFile == nil {
//         fmt.Println("unable to create image file")
//         return
//     }

//     //file created; defer closing the file
//     defer closeFile()
//     //get size of incoming file
//     sizeHeader := make([]byte, 20)
//     bytesRead, _ := conn.Read(sizeHeader)
//     sizeHeaderStr := string(sizeHeader[:bytesRead])
//     sizeHeaderStr = strings.TrimPrefix(sizeHeaderStr, "size: ")
//     sizeHeaderInt, err := strconv.Atoi(sizeHeaderStr)
//     if err != nil {
//         fmt.Println("invalid size header:", err)
//         return
//     }
//     fmt.Println(sizeHeaderInt)
//     //copy data from connection to image_file
//     written, err := io.CopyN(imageFile, conn, int64(sizeHeaderInt))
//     fmt.Println("Received image")
//     if err != nil {
//         fmt.Println("copy file:", err)
//         return
//     }
//     fmt.Println(written/KB)

//     //create a file for the blurred image
//     edited_image, closeEdited := createFile(EDITED_IMAGE_PATH)
//     if edited_image == nil {
//         fmt.Println("unable to create image file")
//         return
//     }
//     defer closeEdited()

//     // blur the image and save the blurred image to the specified path
//     // success := blurImage(IMAGE_PATH, EDITED_IMAGE_PATH)

//     success := glitchImage(IMAGE_PATH, EDITED_IMAGE_PATH)
//     if !success {
//         fmt.Println("Error: Could not save edited image")
//         return
//     }
//     fmt.Println("edited image saved to:", EDITED_IMAGE_PATH)

//     //copy the edited image back to the client
//     written, err = io.Copy(conn, edited_image)
//     if err != nil {
//         fmt.Println("copy edited:", err)
//     }
//     fmt.Printf("sent edited image with size: %d kb to client\n", written/KB)
// }

//create a file located at the specified path, return nil
//if there was an error creating the file.
func createFile(filePath string) (file *os.File, closeFile func()) {
    file, err := os.Create(filePath)
    if err != nil {
        fmt.Println("create file:", err)
        return nil, nil
    }
    //function to close the file
    closeFile = func () {
		err := file.Close()
		if err != nil {
			fmt.Println("close file:", err)
		}
    }
    return file, closeFile
}

// blur image function
func blurImage(inputPath, outputPath string) bool {
    // Read the image
    img := gocv.IMRead(inputPath, gocv.IMReadColor)
    if img.Empty() {
        fmt.Println("Error: Could not read image")
        return false
    }
    defer img.Close()

    // Create a Mat to hold the blurred image
    blurred := gocv.NewMat()
    defer blurred.Close()

    // Apply Gaussian blur
    // Parameters:
    // - img: source image
    // - &blurred: destination image
    // - image.Point{X: 201, Y: 201}: kernel size (must be odd numbers)
    // - 0: sigmaX (0 means calculate from kernel size)
    // - 0: sigmaY (0 means calculate from kernel size)
    gocv.GaussianBlur(img, &blurred, image.Point{X: 201, Y: 201}, 0, 0, gocv.BorderDefault)

    // Save the blurred image
    return gocv.IMWrite(outputPath, blurred)
}

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

// pencil sketch effect function
func pencilSketchImage(inputImage []byte) ([]byte, error) {
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

//convert an image to ascii characters
func imageToASCII(inputImage []byte) ([]byte, error) {
    // 1. Preprocess with gocv: enhance contrast for better ASCII output
    // preprocessed, err := preprocessForASCII(inputImage)
    // if err != nil {
    //     // Fall back to original image if preprocessing fails
    //     log.Printf("preprocessing failed, using original: %v", err)
    //     preprocessed = inputImage
    // }

    readerPtr := bytes.NewReader(inputImage)
    //decode image from webrtc to image.Image using image/jpeg
    jpegImage, err := jpeg.Decode(readerPtr)
    if err != nil {
        return nil, fmt.Errorf("imageToASCII %w", err)
    }
    fmt.Printf("width: %d, height: %d\n", jpegImage.Bounds().Dx(), jpegImage.Bounds().Dy())
    //convert image.Image to ASCII using image2ascii library
    image2ASCIIConverter := convert.NewImageConverter()
    convertOptions := convert.DefaultOptions
    convertOptions.Colored = false
    convertOptions.FixedWidth = 60
    convertOptions.FixedHeight = 22
    asciiString := image2ASCIIConverter.Image2ASCIIString(jpegImage, &convertOptions)
    //send ASCII string as byte slice
    fmt.Println(asciiString)
    return []byte(asciiString), nil
}

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
func faceDetect(inputImage []byte) ([]byte, error) {
    // 1. Decode Input Image
    img, err := gocv.IMDecode(inputImage, gocv.IMReadColor)
    if err != nil || img.Empty() {
        return nil, fmt.Errorf("decoding image: %w", err)
    }
    defer img.Close()

    // 2. Create YuNet Face Detector
    // Parameters: model path, config (empty for ONNX), input size, score threshold, NMS threshold, top_k
    faceDetector := gocv.NewFaceDetectorYN(YUNET_MODEL_PATH, "", image.Pt(img.Cols(), img.Rows()))
    defer faceDetector.Close()

    // Set detection parameters
    faceDetector.SetScoreThreshold(0.5)
    faceDetector.SetNMSThreshold(0.3)
    faceDetector.SetTopK(5000)

    // 3. Detect Faces
    faces := gocv.NewMat()
    defer faces.Close()
    
    faceDetector.Detect(img, &faces)

    if faces.Empty() || faces.Rows() == 0 {
        return nil, fmt.Errorf("no face detected")
    }

    // 4. Find face with highest confidence
    var bestConfidence float32 = 0.0
    var bestRect image.Rectangle

    // YuNet returns faces as rows with 15 columns:
    // [x, y, w, h, x_re, y_re, x_le, y_le, x_nt, y_nt, x_rcm, y_rcm, x_lcm, y_lcm, score]
    for i := 0; i < faces.Rows(); i++ {
        confidence := faces.GetFloatAt(i, 14)
        
        if confidence > bestConfidence {
            bestConfidence = confidence
            
            x := int(faces.GetFloatAt(i, 0))
            y := int(faces.GetFloatAt(i, 1))
            w := int(faces.GetFloatAt(i, 2))
            h := int(faces.GetFloatAt(i, 3))
            
            bestRect = image.Rect(x, y, x+w, y+h)
        }
    }

    // 5. Draw rectangle on the face with highest confidence
    green := color.RGBA{R: 0, G: 255, B: 0, A: 255}
    gocv.Rectangle(&img, bestRect, green, 3)

    // 6. Encode to JPEG
    buf, err := gocv.IMEncode(gocv.JPEGFileExt, img)
    if err != nil {
        return nil, fmt.Errorf("encoding image: %w", err)
    }
    defer buf.Close()

    // 7. Deep Copy (Segfault Prevention)
    nativeData := buf.GetBytes()
    finalBytes := make([]byte, len(nativeData))
    copy(finalBytes, nativeData)

    return finalBytes, nil
}

// digital glitch effect function
// func glitchImage(inputImage []byte) ([]byte, error) {
//     // Read the image
//     mat, err := gocv.IMDecode(inputImage, gocv.IMReadColor)
//     if mat.Empty() || err != nil {
//         return nil, fmt.Errorf("error reading image: %w", err)
//     }
//     defer func () {
//        err = mat.Close()
//        if err != nil {
//            fmt.Println("Error closing mat:", err)
//        }
//     }()
//     // Get image dimensions
//     rows := mat.Rows()
//     cols := mat.Cols()

//     // Create output image
//     result := gocv.NewMatWithSize(rows, cols, mat.Type())
//     defer result.Close()
//     mat.CopyTo(&result)

//     // Step 1: Split into RGB channels
//     channels := gocv.Split(mat)
//     defer channels[0].Close()
//     defer channels[1].Close()
//     defer channels[2].Close()

//     // Step 2: Create EXTREME RGB channel separation with large varying shifts
//     redChannel := gocv.NewMatWithSize(rows, cols, gocv.MatTypeCV8U)
//     defer redChannel.Close()
//     greenChannel := gocv.NewMatWithSize(rows, cols, gocv.MatTypeCV8U)
//     defer greenChannel.Close()
//     blueChannel := gocv.NewMatWithSize(rows, cols, gocv.MatTypeCV8U)
//     defer blueChannel.Close()

//     // Apply MUCH stronger horizontal shifts to each channel
//     for y := 0; y < rows; y++ {
//         // Much larger shifts - 40-80 pixels for dramatic effect
//         redShift := 50 + (y%60) - 30     // Varies between 20-80 pixels
//         blueShift := -40 - (y%50) + 25   // Varies between -65 to -15 pixels
//         greenShift := (y % 30) - 15      // Green shift for extra chaos
        
//         for x := 0; x < cols; x++ {
//             // Red channel - shift right
//             srcX := x - redShift
//             if srcX >= 0 && srcX < cols {
//                 redChannel.SetUCharAt(y, x, channels[2].GetUCharAt(y, srcX))
//             }
            
//             // Blue channel - shift left
//             srcX = x - blueShift
//             if srcX >= 0 && srcX < cols {
//                 blueChannel.SetUCharAt(y, x, channels[0].GetUCharAt(y, srcX))
//             }
            
//             // Green channel - varying shift
//             srcX = x - greenShift
//             if srcX >= 0 && srcX < cols {
//                 greenChannel.SetUCharAt(y, x, channels[1].GetUCharAt(y, srcX))
//             }
//         }
//     }

//     // Merge shifted channels
//     shiftedChannels := []gocv.Mat{blueChannel, greenChannel, redChannel}
//     rgbShifted := gocv.NewMat()
//     defer rgbShifted.Close()
//     gocv.Merge(shiftedChannels, &rgbShifted)
//     rgbShifted.CopyTo(&result)

//     // Step 3: Add MANY heavy horizontal displacement bands
//     numGlitchBands := 60 // Many more glitch bands
//     for i := 0; i < numGlitchBands; i++ {
//         pos := (rows * i / numGlitchBands) + (i * 11) % 30
//         if pos >= rows {
//             continue
//         }
        
//         // Larger strip heights
//         stripHeight := 5 + (i % 25)
        
//         // MUCH stronger displacement - up to 150 pixels
//         displacement := ((i * 53) % 300) - 150
        
//         for y := pos; y < pos+stripHeight && y < rows; y++ {
//             for x := 0; x < cols; x++ {
//                 srcX := x - displacement
//                 if srcX >= 0 && srcX < cols {
//                     vec := rgbShifted.GetVecbAt(y, srcX)
//                     result.SetVecbAt(y, x, vec)
//                 }
//             }
//         }
//     }

//     // Step 4: Add LONG horizontal smear/stretch effects on many rows
//     for y := 0; y < rows; y++ {
//         // Apply smear effect to roughly 50% of rows
//         if (y*17)%10 < 5 {
//             smearLength := 30 + (y % 100) // Much longer smears
//             startX := (y * 13) % (cols / 2)
            
//             if startX < cols {
//                 vec := result.GetVecbAt(y, startX)
                
//                 for x := startX; x < startX+smearLength && x < cols; x++ {
//                     result.SetVecbAt(y, x, vec)
//                 }
//             }
            
//             // Add a second smear on some rows
//             if (y*23)%10 < 3 {
//                 startX2 := (y * 29) % cols
//                 smearLength2 := 50 + (y % 80)
//                 if startX2 < cols {
//                     vec := result.GetVecbAt(y, startX2)
                    
//                     for x := startX2; x < startX2+smearLength2 && x < cols; x++ {
//                         result.SetVecbAt(y, x, vec)
//                     }
//                 }
//             }
//         }
//     }

//     // Step 5: Add heavy color noise/static - 15% of pixels
//     for y := 0; y < rows; y++ {
//         for x := 0; x < cols; x++ {
//             noiseChance := ((x * 31) + (y * 17)) % 100
//             if noiseChance < 15 {
//                 noiseR := uint8(((x * 7 + y * 13) % 256))
//                 noiseG := uint8(((x * 11 + y * 23) % 256))
//                 noiseB := uint8(((x * 17 + y * 29) % 256))
//                 result.SetVecbAt(y, x, gocv.NewVecb(noiseB, noiseG, noiseR))
//             }
//         }
//     }

//     // Step 6: Add HEAVY scanlines - every row affected
//     for y := 0; y < rows; y++ {
//         // Very dark scanlines
//         scanlineIntensity := 0.4 + float64((y%4))*0.15 // Varies between 0.4-0.85
        
//         for x := 0; x < cols; x++ {
//             vec := result.GetVecbAt(y, x)
//             b := uint8(float64(vec[0]) * scanlineIntensity)
//             g := uint8(float64(vec[1]) * scanlineIntensity)
//             r := uint8(float64(vec[2]) * scanlineIntensity)
            
//             result.SetVecbAt(y, x, gocv.NewVecb(b, g, r))
//         }
//     }

//     // Step 7: Add horizontal line artifacts (bright/dark lines)
//     for i := 0; i < 40; i++ {
//         lineY := (i * 37) % rows
//         lineIntensity := 0.3 + float64((i%5))*0.2 // Some lines darker, some brighter
        
//         for x := 0; x < cols; x++ {
//             vec := result.GetVecbAt(lineY, x)
//             b := uint8(float64(vec[0]) * lineIntensity)
//             g := uint8(float64(vec[1]) * lineIntensity)
//             r := uint8(float64(vec[2]) * lineIntensity)
            
//             result.SetVecbAt(lineY, x, gocv.NewVecb(b, g, r))
//         }
//     }

//     // Step 8: Add vertical glitch columns (characteristic of digital glitches)
//     for i := 0; i < 15; i++ {
//         colX := (i * 43) % cols
//         colWidth := 2 + (i % 5)
//         displacement := ((i * 71) % 60) - 30
        
//         for y := 0; y < rows; y++ {
//             srcY := y - displacement
//             if srcY >= 0 && srcY < rows {
//                 for dx := 0; dx < colWidth && colX+dx < cols; dx++ {
//                     vec := result.GetVecbAt(srcY, colX+dx)
//                     result.SetVecbAt(y, colX+dx, vec)
//                 }
//             }
//         }
//     }

//     // Step 9: Darken overall image significantly
//     for y := 0; y < rows; y++ {
//         for x := 0; x < cols; x++ {
//             vec := result.GetVecbAt(y, x)
//             b := uint8(float64(vec[0]) * 0.9)
//             g := uint8(float64(vec[1]) * 0.9)
//             r := uint8(float64(vec[2]) * 0.9)
            
//             // Reduce overall brightness by 10% for slight moody effect
//             result.SetVecbAt(y, x, gocv.NewVecb(b, g, r))
//         }
//     }

//     // Save the glitched image
//     buf, err := gocv.IMEncode(gocv.JPEGFileExt, result)
//     if err != nil {
//         return nil, fmt.Errorf("encoding image: %w", err)
//     }
//     imageByte := buf.GetBytes()
//     defer buf.Close()

//     return imageByte, nil
// }