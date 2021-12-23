package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alifpay/qs/app"
)

func main() {

	//nats server urls
	addr := os.Getenv("NATSURLS")
	if addr == "" {
		addr = "nats://127.0.0.1:4222"
	}

	//nats token
	token := os.Getenv("NATSTOKEN")
	if token == "" {
		token = "123wexsx2asekcijyc"
	}

	//queue service subject name
	subj := os.Getenv("SUBJECT")
	if subj == "" {
		subj = "EnQueue"
	}

	ctx, cancel := context.WithCancel(context.Background())

	j := app.New(addr, token, subj)
	go j.Run(ctx)

	log.Println("Queue Service is running")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	log.Println("Shutting down")
	signal.Stop(sig)
	close(sig)
	cancel()

}
