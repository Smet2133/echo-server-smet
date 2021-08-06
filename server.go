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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h1>Welcome to my web server!</h1>"))
	})

	serverBind := ":8080"
	log.Printf("Start server on %v", serverBind)
	router := controller.CreateApi(mux.NewRouter(), &controllers)

	if err := http.ListenAndServe(serverBind, router); err != nil {
		log.Printf("server start error: %v", err.Error())
	}
	log.Printf("Server finished")
}
