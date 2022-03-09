package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alifpay/qs/app"
)

func main() {
	//nats server urls
	addr := os.Getenv("natsurl")
	if addr == "" {
		addr = *flag.String("natsurl", "nats://127.0.0.1:4222", "nats url")
	}

	//nats token
	token := os.Getenv("natstoken")
	if token == "" {
		token = *flag.String("natstoken", "123wexsx2asekcijyc", "nats token")
	}

	//queue service subject name
	subj := os.Getenv("subject")
	if subj == "" {
		subj = *flag.String("subject", "EnQueue", "qs subject name")
	}

	//queue service stream name
	str := os.Getenv("stream")
	if str == "" {
		str = *flag.String("stream", "qsStream21", "qs subject name")
	}

	appId := os.Getenv("HOSTNAME")
	if appId == "" {
		appId = *flag.String("appid", "qs23", "appid is unique id of app")
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
	j := app.New(addr, token, subj, str, appId)
	j.Run(ctx)
}
