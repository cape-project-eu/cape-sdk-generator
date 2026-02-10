package main

import (
	"log"
	"net/http"

	s_v1 "cape-project.eu/sdk-generator/mockserver/foundation/storage/v1"
	ws_v1 "cape-project.eu/sdk-generator/mockserver/foundation/workspace/v1"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	ws_v1.RegisterServer(router)
	s_v1.RegisterServer(router)

	s := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:8080",
	}

	log.Fatal(s.ListenAndServe())
}
