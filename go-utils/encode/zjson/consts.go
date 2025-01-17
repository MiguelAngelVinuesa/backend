package zjson

const (
	trueString     = "true"
	falseString    = "false"
	timestampMilli = "20060102150405.999Z0700"
)

// FixBufferSize adjust the given value to a more convenient multiple of a power of 2.
func FixBufferSize(capacity int) int {
	var m int
	if capacity < 1024 {
		m = 128 // small
	} else if capacity < 4096 {
		m = 1024 // medium
	} else {
		m = 4096 // large
	}

	c := (capacity - 1) / m
	c++
	return c * m
}
