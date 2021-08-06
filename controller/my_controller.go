package controller

import (
	"bytes"
	"context"
	"log"
	"net/http"
)
type MyController struct {

}

func (c *MyController) Echo(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	// read request body to string
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	requestStr := buf.String()
	log.Printf("Remote addr: %v", r.RemoteAddr)
	log.Printf("Request body: %v", requestStr)

	respondWithJson(context.Background(), w, http.StatusOK, "")
}