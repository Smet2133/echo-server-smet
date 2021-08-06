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

	serverBind := "localhost:8080"
	log.Printf("Start server on %v", serverBind)
	router := controller.CreateApi(mux.NewRouter(), &controllers)

	if err := http.ListenAndServe(serverBind, router); err != nil {
		log.Printf(err.Error())
	}
	log.Printf("Successfully finished")
}


