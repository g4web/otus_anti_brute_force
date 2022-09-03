package app

import (
	"context"
	"log"
	"testing"

	"github.com/g4web/otus_anti_brute_force/internal/config"
	memorystorage "github.com/g4web/otus_anti_brute_force/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	t.Run("Test request", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		configs, err := config.NewConfig("../configs/config_test.env")
		if err != nil {
			log.Fatalf("error reading configs: %v", err)
		}
		networkPersistentStorage := memorystorage.NewMemoryStorage()
		networkFastStorage := memorystorage.NewMemoryStorage()
		app := NewApp(ctx, configs, networkPersistentStorage, networkFastStorage)

		isOk, err := app.IsOk("192.168.0.1", "-=@wesomeNikN@me=-", "gfhjkm")
		require.NoError(t, err)
		require.Equal(t, true, isOk)

		isOk, err = app.IsOk("192.168.0.1", "-=@wesomeNikN@me=-", "gfhjkm")
		require.NoError(t, err)
		require.Equal(t, true, isOk)

		isOk, err = app.IsOk("192.168.0.1", "-=@wesomeNikN@me=-", "gfhjkm")
		require.NoError(t, err)
		require.Equal(t, true, isOk)

		isOk, err = app.IsOk("192.168.0.12", "-=@wesomeNikN@me=-", "gfhjkm")
		require.NoError(t, err)
		require.Equal(t, false, isOk)

		isOk, err = app.IsOk("192.168.0.1", "-=@wesomeNikN@me=-2", "gfhjkm")
		require.NoError(t, err)
		require.Equal(t, false, isOk)

		isOk, err = app.IsOk("192.168.0.1", "-=@wesomeNikN@me=-", "gfhjkm2")
		require.NoError(t, err)
		require.Equal(t, false, isOk)

		isOk, err = app.IsOk("192.168.0.12", "-=@wesomeNikN@me=-2", "gfhjkm2")
		require.NoError(t, err)
		require.Equal(t, true, isOk)
	})
}
