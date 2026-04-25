package main

import (
	"gocv_project/effects"
	"io"
	"log"
	"net/http"
)

func (appPtr *app) testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Test endpoint reached\n"))
}

func (appPtr *app) glitchImageHandler(w http.ResponseWriter, r *http.Request) {
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

	log.Printf("Received frame: %d bytes, mode: glitch", len(body))
	//glitch image
	processed, err := effects.GlitchImage(body)
	if err != nil {
		log.Printf("Error processing glitch image: %v", err)
		http.Error(w, "Failed to process image", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(processed)
	log.Printf("Sent glitched image: %d bytes", len(processed))
}

func (appPtr *app) blurImageHandler(w http.ResponseWriter, r *http.Request) {
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

	log.Printf("Received frame: %d bytes, mode: blur", len(body))

	processed, err := effects.BlurImage(body)
	if err != nil {
		log.Printf("Error processing blur image: %v", err)
		http.Error(w, "Failed to process image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(processed)
	log.Printf("Sent blurred image: %d bytes", len(processed))
}

func (appPtr *app) sketchImageHandler(w http.ResponseWriter, r *http.Request) {
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

	log.Printf("Received frame: %d bytes, mode: sketch", len(body))

	processed, err := effects.PencilSketchImage(body)
	if err != nil {
		log.Printf("Error processing sketch image: %v", err)
		http.Error(w, "Failed to process image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(processed)
	log.Printf("Sent sketch image: %d bytes", len(processed))
}

func (appPtr *app) embossImageHandler(w http.ResponseWriter, r *http.Request) {
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

	log.Printf("Received frame: %d bytes, mode: emboss", len(body))

	processed, err := effects.EmbossImage(body)
	if err != nil {
		log.Printf("Error processing emboss image: %v", err)
		http.Error(w, "Failed to process image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(processed)
	log.Printf("Sent embossed image: %d bytes", len(processed))
}

func (appPtr *app) waveRippleImageHandler(w http.ResponseWriter, r *http.Request) {
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

	log.Printf("Received frame: %d bytes, mode: wave-ripple", len(body))

	processed, err := effects.WaveRippleImage(body)
	if err != nil {
		log.Printf("Error processing wave ripple image: %v", err)
		http.Error(w, "Failed to process image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(processed)
	log.Printf("Sent wave ripple image: %d bytes", len(processed))
}

func (appPtr *app) pixelateImageHandler(w http.ResponseWriter, r *http.Request) {
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

	log.Printf("Received frame: %d bytes, mode: pixelate", len(body))

	processed, err := effects.PixelateImage(body)
	if err != nil {
		log.Printf("Error processing pixelate image: %v", err)
		http.Error(w, "Failed to process image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(processed)
	log.Printf("Sent pixelated image: %d bytes", len(processed))
}

	// processed, err := appPtr.applyEffect(r, effects.HalftoneImage, HALFTONE)
	// if err != nil {
	// 	switch {
	// 	case errors.Is(err, ErrReadingBody) || errors.Is(err, ErrEmptyReqBody):
	// 		http.Error(w, err.Error(), http.StatusBadRequest)
	// 	case errors.Is(err, ErrApplyingEffect):
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	default:
	// 		http.Error(w, "an error occurred", http.StatusInternalServerError)
	// 	}
	// 	return
	// }

	// w.Header().Set("Content-Type", "image/jpeg")
	// w.Write(processed)	
	// log.Printf("Sent halftone image: %d bytes", len(processed))