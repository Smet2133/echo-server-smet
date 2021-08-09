package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"time"

	//"golang.org/x/net/websocket"
	"log"
	"net/http"
	"os"
	"strconv"
)

var myChan chan string
var chanMap map[string]chan string

const (
	// Time allowed to write a message to the peer.
	writeWait = 60 * time.Minute

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

func init() {
	myChan = make(chan string)
	chanMap = make(map[string]chan string)
}

type ApiControllers struct {
	MyController *MyController
}

func CreateApi(r *mux.Router, controllers *ApiControllers) http.Handler {
	register := func(method string, path string, handler http.HandlerFunc) {
		r.Handle(path,
			handlers.LoggingHandler(
				os.Stdout,
				http.HandlerFunc(handler),
			)).Methods(method)
	}

	// rabbit api
	register(http.MethodPost, "/echo", controllers.MyController.Echo)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h1>Welcome to my web server!</h1>"))
	})

	r.HandleFunc("/websocket", echoWebsocket)

	return errorHandlerWrapper(r)
}

var upgrader = websocket.Upgrader{} // use default options

func echoWebsocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	c.SetWriteDeadline(time.Now().Add(time.Minute * 60))

	id := uuid.New().String()
	socketEchoChan := make(chan string)
	chanMap[id] = socketEchoChan

	defer close(socketEchoChan)
	defer delete(chanMap, id)

	stdoutDone := make(chan struct{})
	go ping(c, stdoutDone)

	for {
		msg := <-socketEchoChan
		err = c.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func ping(ws *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Println("ping:", err)
			}
		case <-done:
			return
		}
	}
}

func errorHandlerWrapper(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() { // error handler, when error occurred it sends request with http status 400 and body with error message
			if err := recover(); err != nil {
				log.Printf("Error: %v\n", err)
				respondWithError(context.Background(), w, http.StatusInternalServerError, err.(string))
				return
			}
		}()
		handler.ServeHTTP(w, r)
	})
}

func respondWithError(ctx context.Context, w http.ResponseWriter, code int, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Printf(msg)
	respondWithJson(ctx, w, code, map[string]string{"error": msg})
}

func respondWithJson(ctx context.Context, w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Error marshall payload: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unknown error"))
		return
	}

	log.Printf("Send response code: %v, body: %v", code, string(response))
	w.Header().Set("Content-Type", "application/json")
	if payload != nil {
		w.Header().Set("Content-Length", strconv.Itoa(len(response)))
	}
	w.WriteHeader(code)
	if payload != nil {
		w.Write(response)
	}
}
