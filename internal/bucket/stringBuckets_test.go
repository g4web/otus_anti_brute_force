package bucket

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStringBucket(t *testing.T) {
	t.Run("Test count limit", func(t *testing.T) {
		ipBuckets := NewStringBuckets(1*time.Millisecond, 3)

		isBanned, _ := ipBuckets.IsBanned("-=@wesomeNikN@me=-")
		require.Equal(t, false, isBanned)

		isBanned, _ = ipBuckets.IsBanned("-=@wesomeNikN@me=-")
		require.Equal(t, false, isBanned)

		isBanned, _ = ipBuckets.IsBanned("-=@wesomeNikN@me=-")
		require.Equal(t, false, isBanned)

		isBanned, _ = ipBuckets.IsBanned("-=@wesomeNikN@me=-")
		require.Equal(t, true, isBanned)

		isBanned, _ = ipBuckets.IsBanned("-=@wesomeNikN@me=-2")
		require.Equal(t, false, isBanned)
	})

	t.Run("Test time limit", func(t *testing.T) {
		ipBuckets := NewStringBuckets(3*time.Millisecond, 3)

		isBanned, _ := ipBuckets.IsBanned("-=@wesomeNikN@me=-")
		require.Equal(t, false, isBanned)

		isBanned, _ = ipBuckets.IsBanned("-=@wesomeNikN@me=-")
		require.Equal(t, false, isBanned)

		time.Sleep(time.Millisecond * 2)
		isBanned, _ = ipBuckets.IsBanned("-=@wesomeNikN@me=-")
		require.Equal(t, false, isBanned)

		time.Sleep(time.Millisecond * 2)
		isBanned, _ = ipBuckets.IsBanned("-=@wesomeNikN@me=-")
		require.Equal(t, false, isBanned)

		isBanned, _ = ipBuckets.IsBanned("-=@wesomeNikN@me=-")
		require.Equal(t, true, isBanned)

		isBanned, _ = ipBuckets.IsBanned("-=@wesomeNikN@me=-")
		require.Equal(t, true, isBanned)

		isBanned, _ = ipBuckets.IsBanned("-=@wesomeNikN@me=-2")
		require.Equal(t, false, isBanned)
	})
}
