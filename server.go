package main

import (
	"echo-server-smet/controller"
	"github.com/gorilla/mux"
	"gopkg.in/resty.v1"
	"log"
	"net/http"
)

func init() {
	resty.SetDisableWarn(true)
}

func main() {
	controllers := controller.ApiControllers{
		MyController: &controller.MyController{},
	}

	serverBind := ":80"
	log.Printf("Start server on %v", serverBind)
	router := controller.CreateApi(mux.NewRouter(), &controllers)

	if err := http.ListenAndServe(serverBind, router); err != nil {
		log.Printf("server start error: %v", err.Error())
	}
	log.Printf("Successfully finished")
}
