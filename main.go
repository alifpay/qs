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

	//queue service subject name
	str := os.Getenv("STREAM")
	if str == "" {
		str = "qsStream21"
	}

	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		signal.Stop(sig)
		close(sig)
		cancel()
	}()

	log.Println("Queue Service is running")
	j := app.New(addr, token, subj, str)
	j.Run(ctx)
}
