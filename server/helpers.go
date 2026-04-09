package main

import (
	"errors"
	"gocv_project/effects"
	"io"
	"log"
	"net/http"
)

var (
	ErrReadingBody = errors.New("error reading body")
	ErrEmptyReqBody = errors.New("empty request body")
	ErrApplyingEffect = errors.New("error applying effect")

)


func (appPtr *app) applyEffect(r *http.Request, effectFunc effects.EffectFunc, effect string) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading body: %v", err)
		return nil, ErrReadingBody
	}
	defer r.Body.Close()

	if len(body) == 0 {
		log.Printf("empty request body: %v", err)
		return nil, ErrEmptyReqBody
	}

	log.Printf("Received frame: %d bytes, mode: halftone", len(body))

	processed, err := effectFunc(body)
	if err != nil {
		log.Printf("Error applying %s effect: %v", effect, err)
		return nil, ErrApplyingEffect
	}

	return processed, nil
}