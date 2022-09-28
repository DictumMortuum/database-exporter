package main

import (
	"context"
	"fmt"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
	"github.com/heetch/confita/backend/flags"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Config struct {
	Databases map[string]string `config:"databases"`
	Exporter  struct {
		Database struct {
			Port string `config:"port"`
		}
	}
}

var (
	Cfg Config
)

func main() {
	loader := confita.NewLoader(
		file.NewBackend("/etc/conf.d/servusrc.yml"),
		flags.NewBackend(),
	)

	err := loader.Load(context.Background(), &Cfg)
	if err != nil {
		log.Fatal(err)
	}

	Init()

	serverDead := make(chan struct{})
	s := NewServer(NewClient())

	go func() {
		s.ListenAndServe()
		close(serverDead)
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		s.Stop()
	}()

	select {
	case <-ctx.Done():
	case <-serverDead:
	}

	version := "0.0.3"
	fmt.Printf("database-exporter v%s HTTP server stopped\n", version)
}
