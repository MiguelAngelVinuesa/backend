package object

// NormalizeSize normalizes the given size to a multiple of minSize.
func NormalizeSize(size, minSize int) int {
	minSize = max(4, minSize)
	if size > 0 {
		size--
	}
	size /= minSize
	size++
	return size * minSize
}
