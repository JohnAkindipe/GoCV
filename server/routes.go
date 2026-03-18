package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func routes() http.Handler {
	routerPtr := httprouter.New()

	routerPtr.HandlerFunc(http.MethodGet, "/test", testHandler)
	routerPtr.HandlerFunc(http.MethodPost, "/webrtc/offer", offerHandler)
	routerPtr.HandlerFunc(http.MethodPost, "/webrtc/candidate", candidateHandler)
	
	return recoverPanic(enableCORS(routerPtr))
}