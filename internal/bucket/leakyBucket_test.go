package bucket

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLeakyBucket(t *testing.T) {
	t.Run("Test leaky bucket", func(t *testing.T) {
		leakyBucket := NewLeakyBucket(10*time.Millisecond, 3)

		require.Equal(t, false, leakyBucket.isBan())
		require.Equal(t, false, leakyBucket.isBan())
		require.Equal(t, false, leakyBucket.isBan())

		require.Equal(t, true, leakyBucket.isBan())
		require.Equal(t, false, leakyBucket.isGarbage())

		time.Sleep(9 * time.Millisecond)
		require.Equal(t, false, leakyBucket.isGarbage())

		time.Sleep(2 * time.Millisecond)
		require.Equal(t, true, leakyBucket.isGarbage())
		require.Equal(t, false, leakyBucket.isBan())
		require.Equal(t, false, leakyBucket.isGarbage())
	})
}
