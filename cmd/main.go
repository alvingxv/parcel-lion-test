package main

import (
	"context"
	"lion-parcel-test/internal/app"
	"lion-parcel-test/internal/delivery/http"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, _ := context.WithCancel(context.Background())
	defer func() {
		r := recover()
		if r != nil {
			// Do something
		}
	}()

	app, err := app.NewApp(ctx)
	if err != nil {
		log.Fatal(err)
	}

	httpServer, err := http.NewHttpServer(app)
	if err != nil {
		log.Fatal(err)
	}

	wait := gracefulShutdown(context.Background(), 30*time.Second, []operationNew{
		{
			name: "server",
			op: func(ctx context.Context) error {
				return httpServer.Stop(ctx)
			},
		},
		{
			name: "app",
			op: func(ctx context.Context) error {
				return app.Close(ctx)
			},
		},
	})

	if err := httpServer.Run(); err != nil {
		log.Println(err)
	}

	<-wait

	os.Exit(0)
}

type operationNew struct {
	name string
	op   func(ctx context.Context) error
}

func gracefulShutdown(ctx context.Context, timeout time.Duration, ops []operationNew) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		// add any other syscalls that you want to be notified with
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		signal.Notify(s, os.Interrupt)
		<-s

		log.Println("shutting down")

		// set timeout for the ops to be done to prevent system hang
		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Printf("timeout %d ms has been elapsed, force exit", timeout.Milliseconds())
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		for _, op := range ops {
			innerOp := op.op
			innerKey := op.name

			log.Printf("cleaning up: %s", innerKey)
			if err := innerOp(ctx); err != nil {
				log.Printf("%s: clean up failed: %s", innerKey, err.Error())
				return
			}

			log.Printf("%s was shutdown gracefully", innerKey)
		}

		close(wait)
	}()

	return wait
}
