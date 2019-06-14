/*

This is a basic application that just has a /health API

*/

package main

import (
	"encoding/json"
	"fmt"
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
	LogInit()

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
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		versionMap := map[string]string{"time_now": time.Now().UTC().Format(time.UnixDate)}
		versionJSON, _ := json.Marshal(versionMap)
		fmt.Fprintf(w, string(versionJSON))
	})

	// Start HTTP server listener
	bindAddress := ":9999"
	if len(os.Getenv("BINDADDRESS")) > 0 {
		bindAddress = strings.ToLower(os.Getenv("BINDADDRESS"))
	}

	Info.Printf("HTTP server binding to: %s", bindAddress)

	http.ListenAndServe(bindAddress, nil)

}

func LogInit() {

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
