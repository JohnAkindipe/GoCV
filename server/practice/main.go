package main

import (
	"fmt"

	"github.com/qeesung/image2ascii/convert"
)

func main() {
	imageConverter := convert.NewImageConverter()
	convertOpts := convert.DefaultOptions
	convertOpts.Colored = false
	convertOpts.FixedWidth = 120
	convertOpts.FixedHeight = 45
	imageASCII := imageConverter.ImageFile2ASCIIString("image_file.jpg", &convertOpts)
	fmt.Println(imageASCII)
}