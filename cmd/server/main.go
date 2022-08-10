package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/g4web/otus_anti_brute_force/configs"
	app "github.com/g4web/otus_anti_brute_force/internal"
	"github.com/g4web/otus_anti_brute_force/server"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.env", "Path to configuration file")
}

func main() {
	flag.Parse()
	config, err := configs.NewConfig(configFile)
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	application := app.NewApp(ctx, config)

	grpc := server.NewServer(application, config)
	defer func() {
		_ = grpc.Stop(ctx)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	err = grpc.Start(ctx)
	if err != nil {
		log.Println(err)
		return
	}
}
