package rng

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRNG(t *testing.T) {
	t.Run("new RNG", func(t *testing.T) {
		r := NewRNG()
		require.NotNil(t, r)
		r.ReturnToPool()
	})
}

func TestNewRNGUniqueness(t *testing.T) {
	t.Run("new RNG uniqueness", func(t *testing.T) {
		max := 100000
		rpt := 100
		counts := make(map[int]int, max)

		for ix := 0; ix < rpt; ix++ {
			r := NewRNG()
			require.NotNil(t, r)

			for iy := 0; iy < max/rpt; iy++ {
				n := int(r.Uint32())
				counts[n] = counts[n] + 1
			}

			r.ReturnToPool()
		}

		for _, c := range counts {
			assert.LessOrEqual(t, c, 3, c)
		}
	})
}

func TestNewRNGWithRoundsUniqueness(t *testing.T) {
	t.Run("new RNG with 20 rounds uniqueness", func(t *testing.T) {
		max := 100000
		rpt := 100
		counts := make(map[int]int, max)

		for ix := 0; ix < rpt; ix++ {
			r := NewRNGWithRounds(20)
			require.NotNil(t, r)

			for iy := 0; iy < max/rpt; iy++ {
				n := int(r.Uint32())
				counts[n] = counts[n] + 1
			}

			r.ReturnToPool()
		}

		for _, c := range counts {
			assert.LessOrEqual(t, c, 3, c)
		}
	})
}

func TestRNG_Read(t *testing.T) {
	testCases := []struct {
		name string
		size int
	}{
		{"1 byte", 1},
		{"2 bytes", 2},
		{"3 bytes", 3},
		{"4 bytes", 4},
		{"5 bytes", 5},
		{"8 bytes", 8},
		{"10 bytes", 10},
		{"20 bytes", 20},
		{"50 bytes", 50},
		{"100 bytes", 100},
		{"200 bytes", 200},
		{"1000 bytes", 1000},
		{"2000 bytes", 2000},
	}

	r := NewCustom(1024, 12)
	require.NotNil(t, r)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := make([]byte, tc.size)
			r.Read(b)
			require.Equal(t, tc.size, len(b))
		})
	}
}

func TestRNG_Uint32(t *testing.T) {
	t.Run("RNG Uint32", func(t *testing.T) {
		r := NewRNG()
		require.NotNil(t, r)
		defer r.ReturnToPool()

		rpt := 500000
		cnt := make(map[uint32]int, rpt)

		for ix := 0; ix < rpt; ix++ {
			s := r.Uint32()
			cnt[s] = cnt[s] + 1
		}

		for _, c := range cnt {
			assert.LessOrEqual(t, c, 2)
		}
	})
}

func TestRNG_Uint64(t *testing.T) {
	t.Run("RNG Uint64", func(t *testing.T) {
		r := NewRNG()
		require.NotNil(t, r)
		defer r.ReturnToPool()

		rpt := 500000
		cnt := make(map[uint64]int, rpt)

		for ix := 0; ix < rpt; ix++ {
			s := r.Uint64()
			cnt[s] = cnt[s] + 1
		}

		for _, c := range cnt {
			assert.LessOrEqual(t, c, 2)
		}
	})
}

func TestRNG_IntN_large(t *testing.T) {
	t.Run("RNG IntN large", func(t *testing.T) {
		r := NewRNG()
		require.NotNil(t, r)
		defer r.ReturnToPool()

		rpt := 500000
		max := 1<<62 + 1
		cnt := make(map[int]int, rpt)

		for ix := 0; ix < rpt; ix++ {
			s := r.IntN(max)
			cnt[s] = cnt[s] + 1
		}

		assert.Greater(t, len(cnt), rpt*4/5)

		for _, c := range cnt {
			assert.LessOrEqual(t, c, 2)
		}
	})
}

func TestRNG_IntN_small(t *testing.T) {
	t.Run("RNG IntN small", func(t *testing.T) {
		r := NewRNG()
		require.NotNil(t, r)
		defer r.ReturnToPool()

		rpt := 500000
		max := 9
		cnt := make(map[int]int, max)

		for ix := 0; ix < rpt; ix++ {
			s := r.IntN(max)
			cnt[s] = cnt[s] + 1
		}

		assert.Equal(t, max, len(cnt))

		low := rpt * 9 / (10 * max)
		high := rpt * 11 / (10 * max)

		for _, c := range cnt {
			assert.GreaterOrEqual(t, c, low)
			assert.LessOrEqual(t, c, high)
		}
	})
}

func TestRNG_IntsN_large(t *testing.T) {
	t.Run("RNG IntsN large", func(t *testing.T) {
		r := NewRNG()
		require.NotNil(t, r)
		defer r.ReturnToPool()

		rpt := 50000
		max := 1<<62 + 1
		buf := make([]int, 10)
		cnt := make(map[int]int, rpt)

		for ix := 0; ix < rpt; ix++ {
			r.IntsN(max, buf)
			for _, s := range buf {
				cnt[s] = cnt[s] + 1
			}
		}

		assert.Greater(t, len(cnt), rpt*4/5)

		for _, c := range cnt {
			assert.LessOrEqual(t, c, 2)
		}
	})
}

func TestRNG_IntsN_small(t *testing.T) {
	t.Run("RNG IntsN small", func(t *testing.T) {
		r := NewRNG()
		require.NotNil(t, r)
		defer r.ReturnToPool()

		rpt := 10000
		max := 13
		buf := make([]int, 50)
		cnt := make(map[int]int, max)

		for ix := 0; ix < rpt; ix++ {
			r.IntsN(max, buf)
			for _, s := range buf {
				cnt[s] = cnt[s] + 1
			}
		}

		assert.Equal(t, max, len(cnt))

		for _, c := range cnt {
			assert.GreaterOrEqual(t, c, rpt*9*50/(10*max))
			assert.LessOrEqual(t, c, rpt*11*50/(10*max))
		}
	})
}

func TestRNG_Int31N_exact(t *testing.T) {
	t.Run("RNG Int31N exact", func(t *testing.T) {
		r := NewRNG()
		require.NotNil(t, r)
		defer r.ReturnToPool()

		rpt := 10000
		max := 1<<31 - 1
		buf := make([]int, 50)
		cnt := make(map[int]int, rpt)

		for ix := 0; ix < rpt; ix++ {
			r.IntsN(max, buf)
			for _, s := range buf {
				cnt[s] = cnt[s] + 1
			}
		}

		assert.Greater(t, len(cnt), rpt*4/5)

		for _, c := range cnt {
			assert.LessOrEqual(t, c, 2)
		}
	})
}

func TestRNG_RandomN(t *testing.T) {
	t.Run("RNG random N", func(t *testing.T) {
		r := NewRNG()
		require.NotNil(t, r)
		defer r.ReturnToPool()

		rpt := 250000
		iy := 0

		for ix := 0; ix < rpt; ix++ {
			n := int(rand.Int31())

			iy++
			switch iy {
			case 1:
				n = n & 0xf
			case 2:
				n = n & 0xff
			case 3:
				n = n & 0xfff
			case 4:
				n = n & 0xffff
			case 5:
				n = n & 0xfffff
			case 6:
				n = n & 0xffffff
			case 7:
				n = n & 0xfffffff
			default:
				iy = 0
			}

			if n < 2 {
				n = 2
			}

			s := r.IntN(n)
			assert.GreaterOrEqual(t, s, 0)
			assert.Less(t, s, n)
		}
	})
}

func TestRNG_Error(t *testing.T) {
	t.Run("RNG IntN error", func(t *testing.T) {
		r := NewRNG()
		require.NotNil(t, r)
		defer r.ReturnToPool()

		var done bool
		defer func() {
			e := recover()
			require.NotNil(t, e)
			done = true
		}()

		r.IntN(0)
		assert.True(t, done)
	})

	t.Run("RNG IntsN error", func(t *testing.T) {
		r := NewRNG()
		require.NotNil(t, r)
		defer r.ReturnToPool()

		var done bool
		defer func() {
			e := recover()
			require.NotNil(t, e)
			done = true
		}()

		r.IntsN(0, nil)
		assert.True(t, done)
	})
}

func TestNewCustom_Error(t *testing.T) {
	t.Run("custom error buf size", func(t *testing.T) {
		var done bool
		defer func() {
			e := recover()
			require.NotNil(t, e)
			done = true
		}()

		NewCustom(10, 0)
		assert.True(t, done)
	})

	t.Run("custom error rounds", func(t *testing.T) {
		var done bool
		defer func() {
			e := recover()
			require.NotNil(t, e)
			done = true
		}()

		NewCustom(100, 0)
		assert.True(t, done)
	})
}

type badReader struct{}

func (r *badReader) Read(_ []byte) (int, error) {
	return 0, fmt.Errorf("bad reader, sharedlib whatcha expect")
}

func TestNewSeeder_Error(t *testing.T) {
	t.Run("RNG new seeder error", func(t *testing.T) {
		seederMutex.Lock()
		save := reader
		seederMutex.Unlock()

		var done bool
		defer func() {
			seederMutex.Lock()
			reader = save
			seederMutex.Unlock()

			e := recover()
			require.NotNil(t, e)
			done = true
		}()

		seederMutex.Lock()
		reader = &badReader{}
		seederMutex.Unlock()

		newSeeder()
		assert.True(t, done)
	})
}

func TestNewSeeder_Reseeding(t *testing.T) {
	t.Run("RNG new seeder re-seeding", func(t *testing.T) {
		seederMutex.Lock()
		saveMin, saveMax := minimumReseed, maximumReseed
		minimumReseed, maximumReseed = time.Millisecond, 2*time.Millisecond
		seederMutex.Unlock()

		newSeeder()

		counts := make(map[uint32]int, 100000)
		for ix := 0; ix < 50; ix++ {
			time.Sleep(time.Millisecond)
			r := NewRNG()
			for iy := 0; iy < 2000; iy++ {
				n := r.Uint32()
				counts[n] = counts[n] + 1
			}
			r.ReturnToPool()
		}

		seederMutex.Lock()
		minimumReseed, maximumReseed = saveMin, saveMax
		seederMutex.Unlock()

		newSeeder()

		seederMutex.Lock()
		require.Greater(t, reseeds, 25)
		seederMutex.Unlock()

		for ix, c := range counts {
			assert.LessOrEqual(t, c, 3, ix)
		}
	})
}

func TestNewSeeder_ReseedingFail(t *testing.T) {
	t.Run("RNG new seeder re-seeding fail", func(t *testing.T) {
		seederMutex.Lock()
		save := reader
		saveMin, saveMax := minimumReseed, maximumReseed
		saveRetry, saveFail := entropyRetry, entropyFail
		seederMutex.Unlock()

		var fail bool
		defer func() {
			seederMutex.Lock()
			reader = save
			minimumReseed, maximumReseed = saveMin, saveMax
			entropyRetry, entropyFail = saveRetry, saveFail
			seederMutex.Unlock()

			require.True(t, fail)
		}()

		seederMutex.Lock()
		minimumReseed, maximumReseed = time.Millisecond, 2*time.Millisecond
		entropyRetry = time.Millisecond
		entropyFail = func() { fail = true }
		seederMutex.Unlock()

		newSeeder()

		seederMutex.Lock()
		reader = &badReader{}
		seederMutex.Unlock()

		time.Sleep(200 * time.Millisecond)
	})
}

func BenchmarkNewRNG8Rounds(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := NewRNGWithRounds(8)
		r.ReturnToPool()
	}
}

func BenchmarkNewRNG12Rounds(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := NewRNGWithRounds(12)
		r.ReturnToPool()
	}
}

func BenchmarkNewRNG20Rounds(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := NewRNGWithRounds(20)
		r.ReturnToPool()
	}
}

func BenchmarkNewRNGDefault(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := NewRNG()
		r.ReturnToPool()
	}
}

func BenchmarkRNGDefault_NewBuf(b *testing.B) {
	r := NewRNG()
	defer r.ReturnToPool()

	for i := 0; i < b.N; i++ {
		r.fillBuffer()
	}
}

func BenchmarkRNG8Rounds_IntN_small(b *testing.B) {
	r := NewRNGWithRounds(8)
	defer r.ReturnToPool()

	n := 9
	for i := 0; i < b.N; i++ {
		r.IntN(n)
	}
}

func BenchmarkRNG8Rounds_IntN_medium(b *testing.B) {
	r := NewRNGWithRounds(8)
	defer r.ReturnToPool()

	n := 1000000
	for i := 0; i < b.N; i++ {
		r.IntN(n)
	}
}

func BenchmarkRNG8Rounds_IntN_large(b *testing.B) {
	r := NewRNGWithRounds(8)
	defer r.ReturnToPool()

	n := 1<<62 + 1
	for i := 0; i < b.N; i++ {
		r.IntN(n)
	}
}

func BenchmarkRNG12Rounds_IntN_small(b *testing.B) {
	r := NewRNGWithRounds(12)
	defer r.ReturnToPool()

	n := 9
	for i := 0; i < b.N; i++ {
		r.IntN(n)
	}
}

func BenchmarkRNG12Rounds_IntN_medium(b *testing.B) {
	r := NewRNGWithRounds(12)
	defer r.ReturnToPool()

	n := 1000000
	for i := 0; i < b.N; i++ {
		r.IntN(n)
	}
}

func BenchmarkRNG12Rounds_IntN_large(b *testing.B) {
	r := NewRNGWithRounds(12)
	defer r.ReturnToPool()

	n := 1<<62 + 1
	for i := 0; i < b.N; i++ {
		r.IntN(n)
	}
}

func BenchmarkRNG20Rounds_IntN_small(b *testing.B) {
	r := NewRNGWithRounds(20)
	defer r.ReturnToPool()

	n := 9
	for i := 0; i < b.N; i++ {
		r.IntN(n)
	}
}

func BenchmarkRNG20Rounds_IntN_medium(b *testing.B) {
	r := NewRNGWithRounds(20)
	defer r.ReturnToPool()

	n := 1000000
	for i := 0; i < b.N; i++ {
		r.IntN(n)
	}
}

func BenchmarkRNG20Rounds_IntN_large(b *testing.B) {
	r := NewRNGWithRounds(20)
	defer r.ReturnToPool()

	n := 1<<62 + 1
	for i := 0; i < b.N; i++ {
		r.IntN(n)
	}
}

func BenchmarkGoRand_IntN_small(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := 9
	for i := 0; i < b.N; i++ {
		r.Intn(n)
	}
}

func BenchmarkGoRand_IntN_medium(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := 1000000
	for i := 0; i < b.N; i++ {
		r.Intn(n)
	}
}

func BenchmarkGoRand_IntN_large(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := 1<<62 + 1
	for i := 0; i < b.N; i++ {
		r.Intn(n)
	}
}

func BenchmarkGoCryptoRand_IntN_small(b *testing.B) {
	buf := make([]byte, 4)
	n := uint32(9)
	for i := 0; i < b.N; i++ {
		rand.Read(buf)
		r := binary.LittleEndian.Uint32(buf)
		r = r % n
	}
}

func BenchmarkGoCryptoRand_IntN_medium(b *testing.B) {
	buf := make([]byte, 4)
	n := uint32(1000000)
	for i := 0; i < b.N; i++ {
		rand.Read(buf)
		r := binary.LittleEndian.Uint32(buf)
		r = r % n
	}
}

func BenchmarkGoCryptoRand_IntN_large(b *testing.B) {
	buf := make([]byte, 8)
	n := uint64(1<<62 + 1)
	for i := 0; i < b.N; i++ {
		rand.Read(buf)
		r := binary.LittleEndian.Uint64(buf)
		r = r % n
	}
}
