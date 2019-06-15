/*

This is a basic application that just has a /health API

*/

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	// "unicode/utf8"
)

var (
	Debug *log.Logger
	Info  *log.Logger
	Error *log.Logger
)

func init() {
}

func main() {

	// Setup Logging To StdOut
	logInit()

	// Handle ^C and SIGTERM gracefully
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		Debug.Printf("Caught signal: %+v", sig)
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()

	// Prepare to handle /health requests to HTTP server
	http.HandleFunc("/health", HealthCheckHandler)

	// Start HTTP server listener
	bindAddress := ":9999"
	if len(os.Getenv("BINDADDRESS")) > 0 {
		bindAddress = strings.ToLower(os.Getenv("BINDADDRESS"))
	}

	Info.Printf("HTTP server binding to: %s", bindAddress)

	http.ListenAndServe(bindAddress, nil)

}

func logInit() {

	errorHandle := os.Stderr
	infoHandle := os.Stdout

	debugHandle := ioutil.Discard
	if strings.ToLower(os.Getenv("VERBOSE")) == "true" {
		debugHandle = os.Stderr
	}

	Debug = log.New(debugHandle,
		"DEB: ",
		log.Ldate|log.Lmicroseconds|log.LUTC)

	Info = log.New(infoHandle,
		"INF: ",
		log.Ldate|log.Lmicroseconds|log.LUTC)

	Error = log.New(errorHandle,
		"ERR: ",
		log.Ldate|log.Lmicroseconds|log.LUTC)

	// no condition here, as you'll only see the message if
	// Verbose logging really is enabled!
	Debug.Printf("Verbose logging enabled")

}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	statusMap := map[string]string{"status": "ok"}
	status, err := json.Marshal(statusMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(status)
}
