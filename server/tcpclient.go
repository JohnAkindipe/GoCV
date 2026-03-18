package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

const (
	_EDITED_IMAGE_PATH = "edited_image_file.jpg"
	_IMAGE_PATH        = `C:\Users\hp\Pictures\HAIR DO.jpg`
	_KB                = 1024
)

func tcpclient() {
	// Connect to the server
	conn, err := net.Dial("tcp", "161.35.36.3:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	imageFile, err := os.Open(_IMAGE_PATH)
	if err != nil {
		fmt.Println("open file:", err)
		return
	}
	defer func() {
		err = imageFile.Close()
		if err != nil {
			fmt.Println("close file:", err)
		}
	}()
	fileInfo, err := imageFile.Stat()
	if err != nil {
		fmt.Println("get file info:", err)
		return
	}
	fmt.Println(fileInfo.Name(), fileInfo.Size()/_KB)

	//tell the server the size of the file i'm sending
	sizeHeader := fmt.Sprintf("size: %d", fileInfo.Size())
	conn.Write([]byte(sizeHeader))
	// stream file to server
	written, err := io.Copy(conn, imageFile)
	if err != nil {
		fmt.Println("stream file:", err)
		return
	}
	fmt.Println(written)

	//create file to store edited image
	editedImageFile, closeFile := _createFile(_EDITED_IMAGE_PATH)
	if editedImageFile == nil {
		return
	}
	defer closeFile()

	//receive edited image from server
	written, err = io.Copy(editedImageFile, conn)
	if err != nil {
		fmt.Println("receive edited image:", err)
	}
	fmt.Printf("received edited image with size %d bytes located at %s", written, _EDITED_IMAGE_PATH)
}

// create a file located at the specified path, return nil
// if there was an error creating the file.
func _createFile(filePath string) (file *os.File, closeFile func()) {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("create file:", err)
		return nil, nil
	}
	//function to close the file
	closeFile = func() {
		err := file.Close()
		if err != nil {
			fmt.Println("close file:", err)
		}
	}
	return file, closeFile
}