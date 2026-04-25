package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func routes() http.Handler {
	routerPtr := httprouter.New()

	routerPtr.HandlerFunc(http.MethodGet, "/test", appPtr.testHandler)
	routerPtr.HandlerFunc(http.MethodPost, "/glitch", appPtr.glitchImageHandler)
	routerPtr.HandlerFunc(http.MethodPost, "/blur", appPtr.blurImageHandler)
	routerPtr.HandlerFunc(http.MethodPost, "/sketch", appPtr.sketchImageHandler)
	routerPtr.HandlerFunc(http.MethodPost, "/emboss", appPtr.embossImageHandler)
	routerPtr.HandlerFunc(http.MethodPost, "/wave-ripple", appPtr.waveRippleImageHandler)
	routerPtr.HandlerFunc(http.MethodPost, "/pixelate", appPtr.pixelateImageHandler)

	return recoverPanic(enableCORS(routerPtr))
}