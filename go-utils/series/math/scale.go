package math

import (
	"math"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/object"
)

// AutoScaleInt64 calculates the minimum & maximum value and step size for the scale of a graph.
func AutoScaleInt64(min, max, steps int64) (int64, int64, int64) {
	diff := (max - min) / steps
	mag := math.Floor(math.Log10(float64(diff)))
	pow := math.Pow(10, mag)
	frac := float64(diff) / pow

	var step int64
	switch {
	case frac <= 2:
		step = int64(2 * pow)
	case frac <= 5:
		step = int64(5 * pow)
	default:
		step = int64(10 * pow)
	}

	if step <= 0 {
		step = 1
	}

	min2 := (min / step) * step
	max2 := (max/step)*step + step

	return min2, max2, step
}

// AutoScaleValuesInt64 calculates the minimum & maximum value and step size for the scale of a graph and returns the scale values.
func AutoScaleValuesInt64(min, max, steps int64) []int64 {
	min2, max2, steps2 := AutoScaleInt64(min, max, steps)
	l := (max2 - min2) / steps2
	out := make(object.Int64s, 0, l)

	for v := min2; v <= max2; v += steps2 {
		out = append(out, v)
	}
	return out
}

// AutoScaleUTC calculates the minimum & maximum value and step size for the scale of a graph.
func AutoScaleUTC(min, max time.Time, steps int64) (time.Time, time.Time, time.Duration) {
	diff := max.Sub(min) / time.Duration(steps)

	days := int64(diff / (24 * time.Hour))
	hours := int64(diff / time.Hour)
	minutes := int64(diff / time.Minute)
	seconds := int64(diff / time.Second)

	var step int64
	switch {
	case days > 1:
		mag := math.Floor(math.Log10(float64(days)))
		pow := math.Pow(10, mag)
		frac := float64(days) / pow

		switch {
		case frac <= 2:
			step = int64(2 * pow)
		case frac <= 5:
			step = int64(5 * pow)
		default:
			step = int64(10 * pow)
		}
		step *= int64(24 * time.Hour / time.Millisecond)

	case hours > 1:
		switch {
		case minutes <= 2:
			step = 2
		case minutes <= 4:
			step = 4
		case minutes <= 6:
			step = 6
		case minutes <= 12:
			step = 12
		case minutes <= 24:
			step = 24
		default:
			step = 48
		}
		step *= int64(time.Hour / time.Millisecond)

	case minutes > 1:
		switch {
		case minutes <= 2:
			step = 2
		case minutes <= 5:
			step = 5
		case minutes <= 15:
			step = 15
		case minutes <= 30:
			step = 30
		case minutes <= 60:
			step = 60
		default:
			step = 120
		}
		step *= int64(time.Minute / time.Millisecond)

	default:
		switch {
		case seconds <= 2:
			step = 2
		case seconds <= 5:
			step = 5
		case seconds <= 15:
			step = 15
		case seconds <= 30:
			step = 30
		case seconds <= 60:
			step = 60
		default:
			step = 120
		}
		step *= int64(time.Second / time.Millisecond)

	}

	d := time.Duration(step) * time.Millisecond

	min2 := min.Round(d)
	if min2.After(min) {
		min2 = min2.Add(-d)
	}

	max2 := max.Round(d)
	if max2.Before(max) {
		max2 = max2.Add(d)
	}

	return min2, max2, d
}

// AutoScaleValuesUTC calculates the minimum & maximum value and step size for the scale of a graph and returns the scale values.
func AutoScaleValuesUTC(min, max time.Time, steps int64) []time.Time {
	min2, max2, steps2 := AutoScaleUTC(min, max, steps)
	l := max2.Sub(min2) / steps2
	out := make([]time.Time, 0, l)

	for v := min2; !v.After(max2); v = v.Add(steps2) {
		out = append(out, v)
	}
	return out
}
