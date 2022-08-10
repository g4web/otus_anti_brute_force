package bucket

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIpBucket(t *testing.T) {
	t.Run("Test count limit", func(t *testing.T) {
		ipBuckets := NewIPBuckets(1*time.Millisecond, 3)

		isBanned, _ := ipBuckets.IsBanned("192.168.0.1")
		require.Equal(t, false, isBanned)

		isBanned, _ = ipBuckets.IsBanned("192.168.0.1")
		require.Equal(t, false, isBanned)

		isBanned, _ = ipBuckets.IsBanned("192.168.0.1")
		require.Equal(t, false, isBanned)

		isBanned, _ = ipBuckets.IsBanned("192.168.0.1")
		require.Equal(t, true, isBanned)

		isBanned, _ = ipBuckets.IsBanned("192.168.0.2")
		require.Equal(t, false, isBanned)
	})

	t.Run("Test time limit", func(t *testing.T) {
		ipBuckets := NewIPBuckets(3*time.Millisecond, 3)

		isBanned, _ := ipBuckets.IsBanned("192.168.0.1")
		require.Equal(t, false, isBanned)

		isBanned, _ = ipBuckets.IsBanned("192.168.0.1")
		require.Equal(t, false, isBanned)

		time.Sleep(time.Millisecond * 2)
		isBanned, _ = ipBuckets.IsBanned("192.168.0.1")
		require.Equal(t, false, isBanned)

		time.Sleep(time.Millisecond * 2)
		isBanned, _ = ipBuckets.IsBanned("192.168.0.1")
		require.Equal(t, false, isBanned)

		isBanned, _ = ipBuckets.IsBanned("192.168.0.1")
		require.Equal(t, true, isBanned)

		isBanned, _ = ipBuckets.IsBanned("192.168.0.1")
		require.Equal(t, true, isBanned)

		isBanned, _ = ipBuckets.IsBanned("192.168.0.2")
		require.Equal(t, false, isBanned)
	})

	t.Run("test whitelist", func(t *testing.T) {
		ipBuckets := NewIPBuckets(1*time.Millisecond, 1)

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
		ipBuckets := NewIPBuckets(1*time.Millisecond, 1)

		err := ipBuckets.AddWhiteListNetwork("192.168.0000.1/23")
		require.Error(t, err)

		err = ipBuckets.AddBlackListNetwork("192.168.0.0/24")
		require.NoError(t, err)

		isBanned, _ := ipBuckets.IsBanned("192.168.0.1")
		require.Equal(t, true, isBanned)

		isBanned, _ = ipBuckets.IsBanned("192.168.1.1")
		require.Equal(t, false, isBanned)
	})
}
