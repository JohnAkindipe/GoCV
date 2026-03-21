package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func routes() http.Handler {
	routerPtr := httprouter.New()

	routerPtr.HandlerFunc(http.MethodGet, "/test", testHandler)
	routerPtr.HandlerFunc(http.MethodPost, "/process-frame", processFrameHandler)
	routerPtr.HandlerFunc(http.MethodPost, "/glitch", glitchImageHandler)
	routerPtr.HandlerFunc(http.MethodGet, "/blur", blurImageHandler)

	
	return recoverPanic(enableCORS(routerPtr))
}