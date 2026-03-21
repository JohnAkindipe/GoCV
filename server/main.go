package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gocv.io/x/gocv"
)

const (
    KB = 1024
    // IMAGE_PATH = "image_file.jpg"
    // EDITED_IMAGE_PATH = "edited_image_file.jpg"
    // Use full path for Linux server
	// MODEL_PATH = "/root/models/face-detect/res10_300x300_ssd_iter_140000.caffemodel"
	// CONFIG_PATH= "/root/models/face-detect/deploy.prototxt"
	// YUNET_MODEL_PATH = "/root/models/face-detect/face_detection_yunet_2023mar.onnx"
)

// ASCII set from dark to light
const asciiChars = "@%#*+=-:. "

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
	//Initialize gocv.net on app
	appPtr = &app{}
	// appPtr.net = gocv.ReadNet(MODEL_PATH, CONFIG_PATH)
    // if appPtr.net.Empty() {
    //     log.Fatal("failed to load DNN model")
    // }
    // defer appPtr.net.Close()

    err := httpServer() 
    if err != nil {
        log.Fatal(err)
    }
}

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