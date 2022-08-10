package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/g4web/otus_anti_brute_force/configs"
	app "github.com/g4web/otus_anti_brute_force/internal"
	"github.com/g4web/otus_anti_brute_force/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	grpcServer *ABFServer
	cancel     context.CancelFunc
	ctx        context.Context
	client     proto.AntiBruteForceClient
	conn       *grpc.ClientConn
)

func init() {
	config, err := configs.NewConfig("../configs/config_test.env")
	if err != nil {
		fmt.Println("error reading config:", err)
	}

	ctx, cancel = context.WithCancel(context.Background())
	application := app.NewApp(ctx, config)

	grpcServer = NewServer(application, config)
	go func() {
		err = grpcServer.Start(ctx)
		if err != nil {
			log.Fatalf("error reading config: %v", err)
		}
	}()

	conn, err = grpc.Dial(
		net.JoinHostPort(config.GrpcHost, config.GrpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		fmt.Println("error connect to GRPC server:", err)
	}
	client = proto.NewAntiBruteForceClient(conn)
}

func TestApp(t *testing.T) {
	defer func() {
		_ = grpcServer.Stop(ctx)
		_ = conn.Close()
		cancel()
	}()

	time.Sleep(time.Millisecond * 50)

	t.Run("Test isOk", func(t *testing.T) {
		r, err := client.IsOk(ctx, &proto.UserRequest{IP: "192.168.0.1", Login: "-=@wesomeNikN@me=-", Password: "gfhjkm"})
		require.NoError(t, err)
		require.Equal(t, true, r.IsOk)

		r, err = client.IsOk(ctx, &proto.UserRequest{IP: "192.168.0.1", Login: "-=@wesomeNikN@me=-", Password: "gfhjkm"})
		require.NoError(t, err)
		require.Equal(t, true, r.IsOk)

		r, err = client.IsOk(ctx, &proto.UserRequest{IP: "192.168.0.1", Login: "-=@wesomeNikN@me=-", Password: "gfhjkm"})
		require.NoError(t, err)
		require.Equal(t, true, r.IsOk)

		r, err = client.IsOk(ctx, &proto.UserRequest{IP: "192.168.0.1", Login: "-=@wesomeNikN@me=-1", Password: "gfhjkm1"})
		require.NoError(t, err)
		require.Equal(t, false, r.IsOk)

		r, err = client.IsOk(ctx, &proto.UserRequest{IP: "192.168.0.12", Login: "-=@wesomeNikN@me=-", Password: "gfhjkm1"})
		require.NoError(t, err)
		require.Equal(t, false, r.IsOk)

		r, err = client.IsOk(ctx, &proto.UserRequest{IP: "192.168.0.12", Login: "-=@wesomeNikN@me=-1", Password: "gfhjkm"})
		require.NoError(t, err)
		require.Equal(t, false, r.IsOk)

		r, err = client.IsOk(ctx, &proto.UserRequest{IP: "192.168.0.12", Login: "-=@wesomeNikN@me=-1", Password: "gfhjkm1"})
		require.NoError(t, err)
		require.Equal(t, true, r.IsOk)

		time.Sleep(time.Millisecond * 50)
	})

	t.Run("Test whitelist", func(t *testing.T) {
		r, err := client.IsOk(ctx, &proto.UserRequest{IP: "192.168.0.2", Login: "-=@wesomeNikN@me=-2", Password: "gfhjkm22"})
		require.NoError(t, err)
		require.Equal(t, true, r.IsOk)

		r, err = client.IsOk(ctx, &proto.UserRequest{IP: "192.168.0.2", Login: "-=@wesomeNikN@me=-22", Password: "gfhjkm22"})
		require.NoError(t, err)
		require.Equal(t, true, r.IsOk)

		r, err = client.IsOk(ctx, &proto.UserRequest{IP: "192.168.0.2", Login: "-=@wesomeNikN@me=-22", Password: "gfhjkm22"})
		require.NoError(t, err)
		require.Equal(t, true, r.IsOk)

		r, err = client.IsOk(ctx, &proto.UserRequest{IP: "192.168.0.2", Login: "-=@wesomeNikN@me=-23", Password: "gfhjkm23"})
		require.NoError(t, err)
		require.Equal(t, false, r.IsOk)

		r2, err := client.AddNetworkToWhiteList(ctx, &proto.AddNetworkToWhiteListRequest{Network: "192.168.0.0/24"})
		require.NoError(t, err)
		require.Equal(t, true, r2.IsSuccess)

		r, err = client.IsOk(ctx, &proto.UserRequest{IP: "192.168.0.2", Login: "-=@wesomeNikN@me=-2", Password: "gfhjkm2"})
		require.NoError(t, err)
		require.Equal(t, true, r.IsOk)

		r3, err := client.RemoveNetworkFromWhiteList(ctx, &proto.RemoveNetworkFromWhiteListRequest{Network: "192.168.0.0/24"})
		require.NoError(t, err)
		require.Equal(t, true, r3.IsSuccess)

		r, err = client.IsOk(ctx, &proto.UserRequest{IP: "192.168.0.2", Login: "-=@wesomeNikN@me=-2", Password: "gfhjkm2"})
		require.NoError(t, err)
		require.Equal(t, false, r.IsOk)

		time.Sleep(time.Millisecond * 50)
	})

	t.Run("Test blacklist", func(t *testing.T) {
		r2, err := client.AddNetworkToBlackList(ctx, &proto.AddNetworkToBlackListRequest{Network: "192.168.0.0/24"})
		require.NoError(t, err)
		require.Equal(t, true, r2.IsSuccess)

		r, err := client.IsOk(ctx, &proto.UserRequest{IP: "192.168.0.3", Login: "-=@wesomeNikN@me=-33", Password: "gfhjkm33"})
		require.NoError(t, err)
		require.Equal(t, false, r.IsOk)

		r3, err := client.RemoveNetworkFromBlackList(ctx, &proto.RemoveNetworkFromBlackListRequest{Network: "192.168.0.0/24"})
		require.NoError(t, err)
		require.Equal(t, true, r3.IsSuccess)

		r, err = client.IsOk(ctx, &proto.UserRequest{IP: "192.168.0.3", Login: "-=@wesomeNikN@me=-33", Password: "gfhjkm33"})
		require.NoError(t, err)
		require.Equal(t, true, r.IsOk)
	})
}
