package main

import (
	"io"
	"log"
	"net/http"
	"time"

	terminal "github.com/buildkite/terminal-to-html/v3"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Test endpoint reached\n"))
}

func processFrameHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if len(body) == 0 {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	mode := r.URL.Query().Get("mode")
	timestamp := time.Now().Format("15:04:05")
	log.Printf("[%s] Received frame: %d bytes, mode: %s", timestamp, len(body), mode)

	switch mode {
	case "glitch":
		processed, err := glitchImage(body)
		if err != nil {
			log.Printf("Error processing glitch image: %v", err)
			http.Error(w, "Failed to process image", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write(processed)
		log.Printf("[%s] Sent glitched image: %d bytes", timestamp, len(processed))

	default:
		// Default to ASCII mode
		asciiBytes, err := imageToASCII(body)
		if err != nil {
			log.Printf("Error converting to ASCII: %v", err)
			http.Error(w, "Failed to process image", http.StatusInternalServerError)
			return
		}
		htmlOutput := terminal.Render(asciiBytes)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(htmlOutput))
		log.Printf("[%s] Sent ASCII art: %d bytes", timestamp, len(asciiBytes))
	}
}

func glitchImageHandler(w http.ResponseWriter, r *http.Request) {
	//glitch image
}

func blurImageHandler(w http.ResponseWriter, r *http.Request) {
	//blur image.
}