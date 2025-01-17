package rng

import (
	crypto "crypto/rand"
	"encoding/binary"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/aead/chacha20/chacha"
)

var (
	defaultRounds = 12 // compromise higher security for higher speed.
)

// RNG is a cryptographically-strong PRNG utilizing the randomness of a ChaCha stream cipher.
// See https://en.wikipedia.org/wiki/Salsa20#ChaCha_variant for more details on ChaCha.
// Implementation is based on https://github.com/lukechampine/frand but with significant changes for speed and usability.
// The CPRNG uses a moderately sized buffer (~16kB) to speed up the random number generation.
// Various functions are provided to return random data as byte(s) or integers.
// Take note of the IntN() & IntsN() functions as they are optimized specifically for generating small numbers,
// as is needed for games like slot machines. These functions are also designed to prevent modulo bias.
// A global hidden "master" CPRNG is created at application startup (see init() function).
// Every new CPRNG generated for the application itself with NewPRG() is seeded from this "master" CPRNG.
// Note that the RNG.IntN() function is faster than the standard rand.IntN(), whilst being cryptographically strong and
// preventing modula bias.
// It is about twice faster than using the cryptographically strong PRNG from the standard rand package as it doesn't use a mutex.
// The functions also cannot fail due to entropy starvation as may happen with the system CPRNG.
// The design makes the CPRNG perfect for games of chance both for short-term use (e.g. during a single API call),
// and prolonged use in simulators requiring a high speed PRNG.
// However, because it doesn't use a mutex, RNG is not safe for concurrent use across multiple go-routines.
type RNG struct {
	buf     byteSlice // the byte buffer (including the seed).
	ptr     byteSlice // offset into the buffer for the next random byte to retrieve.
	rounds  int       // number of rounds for the ChaCha stream cipher.
	lastN   int       // last N used in int31n().
	lastMax uint32    // max used to prevent modulo bias for last N.
}

// NewRNG returns a new CPRNG from the memory pool, initialized with the default number of ChaCha rounds.
// Using a memory pool makes them cheap to initialize and tear down.
// As every new CPRNG is seeded from the "master" CPRNG it will produce totally unpredictable numbers.
func NewRNG() *RNG {
	r := rngPool.Get().(*RNG)
	r.rounds = defaultRounds
	readSeed(r.buf[:chacha.KeySize])
	r.fillBuffer()
	return r
}

// NewRNGWithRounds returns a new CPRNG from the memory pool, initialized with the given rounds for the ChaCha cipher.
// Using a memory pool makes them cheap to initialize and tear down.
// As every new CPRNG is seeded from the "master" CPRNG it will produce totally unpredictable numbers.
// The function will panic if rounds is not 8, 12 or 20.
func NewRNGWithRounds(rounds int) *RNG {
	r := rngPool.Get().(*RNG)
	r.rounds = rounds
	readSeed(r.buf[:chacha.KeySize])
	r.fillBuffer()
	return r
}

// ReturnToPool puts the CPRNG back in the memory pool.
// Never re-use the CPRNG after calling this function, just call NewRNG() if you need another one.
func (r *RNG) ReturnToPool() {
	if r != nil {
		// clear the buffer; it is re-populated when it comes back in use.
		for ix := range r.buf {
			r.buf[ix] = 0
		}
		r.lastN = 0
		rngPool.Put(r)
	}
}

// Read fills b with random data from the internal buffer.
// If the buffer is depleted it is re-populated with random data using the ChaCha stream cipher.
// If the size of b exceeds the size of the internal buffer, b will be filled using the ChaCha stream cipher directly.
// This is considered an exceptional use-case which is not common and doesn't warrant any optimization.
func (r *RNG) Read(b []uint8) {
	if len(b) <= len(r.ptr) {
		n := copy(b, r.ptr)
		r.ptr = r.ptr[n:]
	} else if len(b) <= len(r.ptr)+len(r.buf[chacha.KeySize:]) {
		n := copy(b, r.ptr)
		r.fillBuffer()
		n = copy(b[n:], r.ptr)
		r.ptr = r.ptr[n:]
	} else {
		// not using memory pool!
		buf := make(byteSlice, chacha.KeySize)
		r.Read(buf) // safe recursive call as buf is always smaller than the CPRNG buffer.
		chacha.XORKeyStream(b, b, dummyNonce, buf, r.rounds)
	}
}

// Uint32 returns a random uint32.
func (r *RNG) Uint32() uint32 {
	if len(r.ptr) < uint32size {
		r.fillBuffer()
	}
	out := binary.LittleEndian.Uint32(r.ptr)
	r.ptr = r.ptr[uint32size:]
	return out
}

// Uint64 returns a random uint64.
func (r *RNG) Uint64() uint64 {
	if len(r.ptr) < uint64size {
		r.fillBuffer()
	}
	out := binary.LittleEndian.Uint64(r.ptr)
	r.ptr = r.ptr[uint64size:]
	return out
}

// IntN returns a uniform random int in the half open interval [0,n].
// The function is optimized for generating small numbers and prevents modulo bias.
// It panics if n <= 0.
func (r *RNG) IntN(n int) int {
	if n <= 0 {
		panic("invalid n for IntN: " + strconv.Itoa(n))
	}

	if n <= math.MaxInt32 {
		return r.int31n(n)
	}

	// prevent modulo bias; unoptimized as this part of the code should occur infrequently.
	max := math.MaxUint64 - math.MaxUint64%uint64(n)
	for {
		if i := r.Uint64(); i < max {
			return int(i % uint64(n))
		}
	}
}

// IntsN fills the given slice with uniform random ints in the half open interval [0,n].
// The function is optimized for generating small numbers and prevents modulo bias.
// It panics if n <= 0.
func (r *RNG) IntsN(n int, out []int) {
	if n <= 0 {
		panic("invalid n for IntsN: " + strconv.Itoa(n))
	}

	if n <= math.MaxInt32 {
		for ix := range out {
			out[ix] = r.int31n(n)
		}
		return
	}

	// prevent modulo bias; unoptimized as this part of the code should occur infrequently.
	max := math.MaxUint64 - math.MaxUint64%uint64(n)
	for ix := range out {
		for {
			if i := r.Uint64(); i < max {
				out[ix] = int(i % uint64(n))
				break
			}
		}
	}
}

// int31n is significantly faster than the standard modulo operation commonly used.
// See the Golang standard source code for rand.int31n() for a more detailed explanation.
// This function also prevents modulo bias by skipping some high values that would cause the bias.
func (r *RNG) int31n(n int) int {
	// it's most likely that the RNG is used with the same N, so this saves us a modulo operation every invocation.
	if n != r.lastN {
		r.lastN = n
		r.lastMax = math.MaxUint32 - math.MaxUint32%uint32(n)
	}

	produce := func() uint32 {
		for {
			if i := r.Uint32(); i < r.lastMax {
				return i
			}
		}
	}

	prod := uint64(produce()) * uint64(n)
	low := uint32(prod)
	if low < uint32(n) {
		thresh := uint32(-n) % uint32(n)
		for low < thresh {
			prod = uint64(produce()) * uint64(n)
			low = uint32(prod)
		}
	}
	return int(prod >> 32)
}

// fillBuffer (re-)populates the CPRNG internal buffer using the ChaCha stream cipher.
func (r *RNG) fillBuffer() {
	chacha.XORKeyStream(r.buf, r.buf, dummyNonce, r.buf[:chacha.KeySize], r.rounds)
	r.ptr = r.buf[chacha.KeySize:]
}

// NewCustom returns a new RNG instance seeded with the "master" CPRNG and
// using the specified buffer size and number of ChaCha rounds.
// It panics if bufSize < 32, or rounds != 8, 12 or 20.
func NewCustom(bufSize, rounds int) *RNG {
	if bufSize < chacha.KeySize {
		panic("invalid bufSize; must be at least 32: " + strconv.Itoa(bufSize))
	}
	if !(rounds == 8 || rounds == 12 || rounds == 20) {
		panic(" invalid rounds; must be 8, 12, or 20: " + strconv.Itoa(rounds))
	}

	r := &RNG{buf: make(byteSlice, chacha.KeySize+bufSize), rounds: rounds}
	readSeed(r.buf[:chacha.KeySize])
	r.fillBuffer()
	return r
}

const (
	defaultBufSize = 1<<14 - chacha.KeySize // ~16kB
	uint32size     = 4
	uint64size     = 8
	entropyStarved = "not enough entropy for master seeder"
)

var (
	// the global hidden "master" CPRNG for seeding every new application CPRNG.
	seeder *RNG
	// mutex to protect the "master" CPRNG from concurrent use across multiple go-routines.
	seederMutex sync.Mutex
	// rngPool is the memory pool for CPRNG instances.
	// The New() function instantiates a new RNG structure seeded from the "master" CPRNG with a populated buffer.
	rngPool = sync.Pool{
		New: func() interface{} { return NewCustom(defaultBufSize, defaultRounds) },
	}
	// dummyNonce is the default nonce used for the ChaCha stream cipher.
	// There is no need for a real nonce when using the cipher for generating random numbers.
	dummyNonce = make(byteSlice, chacha.NonceSize)
	// minimumReseed is the minimum delay before re-seeding of the master CPRNG is performed.
	minimumReseed = 40 * time.Second
	// maximumReseed is the maximum delay before re-seeding of the master CPRNG is performed.
	maximumReseed = 60 * time.Second
	// getDelay can be used to generate a random interval between minimumReseed and maximumReseed.
	// It can be overridden in unit-tests os the re-seeding mechanism can be triggered faster.
	getDelay = func() time.Duration {
		seederMutex.Lock()
		t := minimumReseed + time.Duration(seeder.Uint64()%uint64(maximumReseed-minimumReseed))
		seederMutex.Unlock()
		return t
	}
	// reseeds is a counter of the number of times re-seeding of the master CPRNG occurred.
	// It is only used for statistical purposes (e.g. unit-tests).
	reseeds = 0
	// following fields are only to allow override for unit-tests (e.g. to test panics).
	reader       = crypto.Reader
	entropyRetry = 250 * time.Millisecond
	entropyFail  = func() { panic(entropyStarved) }
)

// byteSlice is a convenience type for a slice of bytes.
type byteSlice []uint8

// newSeeder instantiates the "master" CPRNG.
// It is initialized with a 4k buffer and will utilize a ChaCha stream cipher with 20 rounds for enhanced security.
// Using a "master" seeder decouples the generation of every new CPRNG from the OS,
// so if the OS doesn't have enough entropy we're not stuck waiting for it.
// This initialization function will panic if the OS is already starved for entropy.
// The master CPRNG will also be re-seeded at random intervals between minimumReseed and maximumReseed.
func newSeeder() {
	seed := make(byteSlice, chacha.KeySize)
	if n, err := reader.Read(seed); err != nil || n != len(seed) {
		entropyFail()
	}

	seeder = &RNG{buf: make(byteSlice, 4096), rounds: 20}
	copy(seeder.buf, seed)
	seeder.fillBuffer()

	go func() {
		seed2 := make(byteSlice, chacha.KeySize)
		t := time.NewTimer(getDelay())
		var retries int

		for {
			select {
			case <-t.C:
				go func() {
					newDelay := getDelay()

					seederMutex.Lock()
					defer seederMutex.Unlock()

					if n, err := reader.Read(seed2); err != nil || n != len(seed) {
						// on Linux, it uses /dev/urandom which should never cause the following code to trigger.
						if retries++; retries >= 100 {
							entropyFail()
						}
						newDelay = entropyRetry
					} else {
						copy(seeder.buf, seed2)
						seeder.fillBuffer()
						reseeds++
						retries = 0
					}

					t.Reset(newDelay)
				}()
			}
		}
	}()
}

// readSeed is a blocking reader for the "master" CPRNG and used only to seed every new application CPRNG.
func readSeed(buf byteSlice) {
	seederMutex.Lock()
	seeder.Read(buf)
	seederMutex.Unlock()
}

// init creates the global hidden "master" CPRNG at application startup.
func init() {
	newSeeder()
}
