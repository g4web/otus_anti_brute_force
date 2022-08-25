package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	app "github.com/g4web/otus_anti_brute_force/internal"
	"github.com/g4web/otus_anti_brute_force/internal/config"
	"github.com/g4web/otus_anti_brute_force/internal/server"
	memorystorage "github.com/g4web/otus_anti_brute_force/internal/storage/memory"
	sqlstorage "github.com/g4web/otus_anti_brute_force/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.env", "Path to configuration file")
}

func main() {
	flag.Parse()
	configs, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatalf("error reading configs: %v", err)
	}
	networkPersistentStorage, err := sqlstorage.NewSQLStorage(configs)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	networkFastStorage := memorystorage.NewMemoryStorage()

	application := app.NewApp(ctx, configs, networkFastStorage, networkPersistentStorage)

	grpc := server.NewABFServer(application, configs)
	defer func() {
		_ = grpc.Stop(ctx)
	}()

	err = grpc.Start(ctx)
	if err != nil {
		log.Println(err)
		return
	}
}
