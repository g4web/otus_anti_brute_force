package bucket

import (
	"testing"
	"time"

	memorystorage "github.com/g4web/otus_anti_brute_force/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestIpBucket(t *testing.T) {
	t.Run("test whitelist", func(t *testing.T) {
		networkPersistentStorage := memorystorage.NewMemoryStorage()
		networkFastStorage := memorystorage.NewMemoryStorage()

		ipBuckets := NewIPBuckets(1*time.Millisecond, 1, networkFastStorage, networkPersistentStorage)

		err := ipBuckets.AddWhiteListNetwork("192.168.1.00/23")
		require.Error(t, err)

		err = ipBuckets.AddWhiteListNetwork("192.168.0.0/24")
		require.NoError(t, err)

		isBanned, _ := ipBuckets.IsBanned("192.168.0.1")
		require.Equal(t, false, isBanned)

		isBanned, _ = ipBuckets.IsBanned("192.168.0.1")
		require.Equal(t, false, isBanned)

		isBanned, _ = ipBuckets.IsBanned("192.168.0.255")
		require.Equal(t, false, isBanned)

		isBanned, _ = ipBuckets.IsBanned("192.168.0.255")
		require.Equal(t, false, isBanned)

		isBanned, _ = ipBuckets.IsBanned("192.168.1.1")
		require.Equal(t, false, isBanned)

		isBanned, _ = ipBuckets.IsBanned("192.168.1.1")
		require.Equal(t, true, isBanned)
	})

	t.Run("Test blacklist", func(t *testing.T) {
		networkPersistentStorage := memorystorage.NewMemoryStorage()
		networkFastStorage := memorystorage.NewMemoryStorage()

		ipBuckets := NewIPBuckets(1*time.Millisecond, 1, networkFastStorage, networkPersistentStorage)

		err := ipBuckets.AddWhiteListNetwork("192.168.0000.1/23")
		require.Error(t, err)

		err = ipBuckets.AddBlackListNetwork("192.168.0.0/24")
		require.NoError(t, err)

		isBanned, _ := ipBuckets.IsBanned("192.168.0.1")
		require.Equal(t, true, isBanned)

		isBanned, _ = ipBuckets.IsBanned("192.168.1.1")
		require.Equal(t, false, isBanned)
	})

	t.Run("Test persistent blacklist", func(t *testing.T) {
		networkFastStorage := memorystorage.NewMemoryStorage()

		networkPersistentStorage := memorystorage.NewMemoryStorage()
		err := networkPersistentStorage.AddToBlackList("192.168.0.0/24")
		require.NoError(t, err)
		err = networkPersistentStorage.AddToWhiteList("192.168.1.0/24")
		require.NoError(t, err)

		ipBuckets := NewIPBuckets(1*time.Millisecond, 1, networkFastStorage, networkPersistentStorage)

		isBanned, _ := ipBuckets.IsBanned("192.168.0.1")
		require.Equal(t, true, isBanned)

		isBanned, _ = ipBuckets.IsBanned("192.168.1.1")
		require.Equal(t, false, isBanned)
		isBanned, _ = ipBuckets.IsBanned("192.168.1.1")
		require.Equal(t, false, isBanned)
	})
}
