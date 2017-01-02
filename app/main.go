package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/kung-foo/freki"
	log "github.com/sirupsen/logrus"
)

func onErrorExit(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func onInterruptSignal(fn func()) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		<-sig
		fn()
	}()
}

func main() {
	processor := freki.New()

	err := processor.Init()
	onErrorExit(err)

	exit := func() {
		onErrorExit(processor.Cleanup())
		os.Exit(0)
	}

	defer exit()
	onInterruptSignal(exit)

	// TODO: move
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello on port %d", 0)
	})

	go http.ListenAndServe(":8080", nil)

	// TODO: pass in stop channel
	processor.Start()
}