package test

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/g4web/otus_anti_brute_force/internal/config"
	"github.com/g4web/otus_anti_brute_force/internal/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	cancel context.CancelFunc
	ctx    context.Context
	client proto.AntiBruteForceClient
	conn   *grpc.ClientConn
)

func init() {
	configs, err := config.NewConfig("/abf/configs/config_test.env")
	if err != nil {
		log.Fatalf("error reading configs: %v", err)
	}

	ctx, cancel = context.WithCancel(context.Background())

	conn, err = grpc.Dial(
		net.JoinHostPort(configs.GrpcHost, configs.GrpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		fmt.Println("error connect to GRPC server:", err)
	}
	client = proto.NewAntiBruteForceClient(conn)
}

func TestApp(t *testing.T) {
	t.Helper()
	defer func() {
		_ = conn.Close()
		cancel()
	}()

	t.Run("Test isOk", func(t *testing.T) {
		dontBlockBecauseLimitHasNotEnded(t)

		blockByIPBecauseLimitHasBeenReached(t)
		blockByLoginBecauseLimitHasBeenReached(t)
		blockByPasswordBecauseLimitHasBeenReached(t)
	})

	t.Run("Test whitelist", func(t *testing.T) {
		time.Sleep(time.Millisecond * 50)
		dontBlockBecauseLimitHasNotEnded(t)

		blockByIPBecauseLimitHasBeenReached(t)
		blockByLoginBecauseLimitHasBeenReached(t)
		blockByPasswordBecauseLimitHasBeenReached(t)

		addNetworkToWhiteList(t)
		dontBlockBecauseIPInWhiteList(t)
		removeNetworkFromWhiteList(t)

		blockByIPBecauseLimitHasBeenReached(t)
		blockByLoginBecauseLimitHasBeenReached(t)
		blockByPasswordBecauseLimitHasBeenReached(t)
	})

	t.Run("Test blacklist", func(t *testing.T) {
		time.Sleep(time.Millisecond * 50)
		AddNetworkToBlackList(t)

		blockBecauseIPInBlackList(t)

		removeNetworkFromBlackList(t)

		dontBlockBecauseLimitHasNotEnded(t)
	})
}

func dontBlockBecauseLimitHasNotEnded(t *testing.T) {
	t.Helper()
	r, err := client.IsOk(ctx, &proto.UserRequest{IP: "192.168.0.1", Login: "-=@wesomeNikN@me=-", Password: "gfhjkm"})
	require.NoError(t, err)
	require.Equal(t, true, r.IsOk)

	r, err = client.IsOk(ctx, &proto.UserRequest{IP: "192.168.0.1", Login: "-=@wesomeNikN@me=-", Password: "gfhjkm"})
	require.NoError(t, err)
	require.Equal(t, true, r.IsOk)

	r, err = client.IsOk(ctx, &proto.UserRequest{IP: "192.168.0.1", Login: "-=@wesomeNikN@me=-", Password: "gfhjkm"})
	require.NoError(t, err)
	require.Equal(t, true, r.IsOk)
}

func blockByIPBecauseLimitHasBeenReached(t *testing.T) {
	t.Helper()
	r, err := client.IsOk(ctx, &proto.UserRequest{IP: "192.168.0.1", Login: "-=@wesomeNikN@me=-1", Password: "gfhjkm1"})
	require.NoError(t, err)
	require.Equal(t, false, r.IsOk)
}

func blockByLoginBecauseLimitHasBeenReached(t *testing.T) {
	t.Helper()
	r, err := client.IsOk(ctx, &proto.UserRequest{IP: "192.168.0.12", Login: "-=@wesomeNikN@me=-", Password: "gfhjkm1"})
	require.NoError(t, err)
	require.Equal(t, false, r.IsOk)
}

func blockByPasswordBecauseLimitHasBeenReached(t *testing.T) {
	t.Helper()
	r, err := client.IsOk(ctx, &proto.UserRequest{IP: "192.168.0.12", Login: "-=@wesomeNikN@me=-1", Password: "gfhjkm"})
	require.NoError(t, err)
	require.Equal(t, false, r.IsOk)
}

func addNetworkToWhiteList(t *testing.T) {
	t.Helper()
	r, err := client.AddNetworkToWhiteList(ctx, &proto.AddNetworkToWhiteListRequest{Network: "192.168.0.0/24"})
	require.NoError(t, err)
	require.Equal(t, true, r.IsSuccess)
}

func dontBlockBecauseIPInWhiteList(t *testing.T) {
	t.Helper()
	r, err := client.IsOk(ctx, &proto.UserRequest{IP: "192.168.0.2", Login: "-=@wesomeNikN@me=-2", Password: "gfhjkm2"})
	require.NoError(t, err)
	require.Equal(t, true, r.IsOk)
}

func removeNetworkFromWhiteList(t *testing.T) {
	t.Helper()
	r, err := client.RemoveNetworkFromWhiteList(ctx, &proto.RemoveNetworkFromWhiteListRequest{Network: "192.168.0.0/24"})
	require.NoError(t, err)
	require.Equal(t, true, r.IsSuccess)
}

func AddNetworkToBlackList(t *testing.T) {
	t.Helper()
	r, err := client.AddNetworkToBlackList(ctx, &proto.AddNetworkToBlackListRequest{Network: "192.168.0.0/24"})
	require.NoError(t, err)
	require.Equal(t, true, r.IsSuccess)
}

func blockBecauseIPInBlackList(t *testing.T) {
	t.Helper()
	r, err := client.IsOk(ctx, &proto.UserRequest{IP: "192.168.0.3", Login: "-=@wesomeNikN@me=-33", Password: "gfhjkm33"})
	require.NoError(t, err)
	require.Equal(t, false, r.IsOk)
}

func removeNetworkFromBlackList(t *testing.T) {
	t.Helper()
	r, err := client.RemoveNetworkFromBlackList(ctx, &proto.RemoveNetworkFromBlackListRequest{Network: "192.168.0.0/24"})
	require.NoError(t, err)
	require.Equal(t, true, r.IsSuccess)
}
