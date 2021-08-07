package main

import (
	"echo-server-smet/controller"
	gorillamux "github.com/gorilla/mux"
	"gopkg.in/resty.v1"
	"log"
	"net/http"
	"os"
)

func init() {
	resty.SetDisableWarn(true)
}

func main() {
	controllers := controller.ApiControllers{
		MyController: &controller.MyController{},
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	serverBind := ":" + port
	log.Printf("Start server on %v", serverBind)
	router := controller.CreateApi(gorillamux.NewRouter(), &controllers)

	if err := http.ListenAndServe(serverBind, router); err != nil {
		log.Printf("server start error: %v", err.Error())
	}
	log.Printf("Server finished")
}
