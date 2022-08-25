package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/g4web/otus_anti_brute_force/internal/config"
	"github.com/g4web/otus_anti_brute_force/internal/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	configFile string
	network    string
	ctx        context.Context
	cancel     context.CancelFunc
	client     proto.AntiBruteForceClient
	conn       *grpc.ClientConn
)

var rootCmd = &cobra.Command{
	Use:   "abf-cli",
	Short: "The CLI for \"anti brute force\" application",
	Long:  `The command line client for "anti brute force" application`,
}

func Execute() {
	startGrpcClient()
	err := rootCmd.Execute()
	stopGrpcClient()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVar(
		&configFile,
		"c",
		"./configs/config.env",
		"A path to config file",
	)
}

func startGrpcClient() {
	configs, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatal("error reading configs:", err)
	}

	conn, err = grpc.Dial(
		net.JoinHostPort(configs.GrpcHost, configs.GrpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		fmt.Println("error connect to GRPC server:", err)
	}

	ctx, cancel = context.WithCancel(context.Background())

	client = proto.NewAntiBruteForceClient(conn)
}

func stopGrpcClient() {
	_ = conn.Close()
	cancel()
}
