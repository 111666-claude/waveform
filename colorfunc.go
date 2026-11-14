package waveform

import (
	"image/color"
	"math/rand"
	"time"
)

// ColorFunc is a function which accepts a variety of values which can be used
// to customize an output image.  These values include the current computed sample
// count (n), the current X coordinate (x), the current Y coordinate (y), and the
// maximum computed values for each of these. (maxN, maxX, maxY)
//
// A ColorFunc is applied during each image drawing iteration, and will
// return the appropriate color which should be drawn at the specified value
// for n, x, and y; possibly taking into account their maximum values.
type ColorFunc func(n int, x int, y int, maxN int, maxX int, maxY int) color.Color

// CheckerColor generates a ColorFunc which produces a checkerboard pattern,
// using the two input colors.  Each square is drawn to the size specified by
// the size parameter.
func CheckerColor(colorA color.Color, colorB color.Color, size uint) ColorFunc {
	return func(n int, x int, y int, maxN int, maxX int, maxY int) color.Color {
		if ((uint(x)/size)+(uint(y)/size))%2 == 0 {
			return colorA
		}

		return colorB
	}
}

// FuzzColor generates a ColorFunc which applies a random color on each call,
// selected from an input, variadic slice of colors.  This can be used to create
// a random fuzz or "static" effect in the resulting waveform image.
func FuzzColor(colors ...color.Color) ColorFunc {
	// Filter any nil values
	colors = filterNilColors(colors)

	// Seed RNG
	rand.Seed(time.Now().UnixNano())

	// Select a color at random on each call
	return func(n int, x int, y int, maxN int, maxX int, maxY int) color.Color {
		return colors[rand.Intn(len(colors))]
	}
}

// GradientColor generates a ColorFunc which produces a linear color gradient
// between two RGBA input colors, interpolating all four channels (including
// alpha).
//
// The gradient is a pure function of n and maxN: the color at n == 0 is the
// start color, the color at n == maxN (when maxN > 0) is the end color, and
// every channel moves monotonically between its two endpoints.  The function
// keeps no state between calls, so drawing callbacks that invoke it more than
// once for the same n, or invoke it with n out of order, always receive the
// same color.  Inputs outside the [0, maxN] range are clamped to the nearest
// endpoint, and no call can produce NaN or panic.
func GradientColor(start color.RGBA, end color.RGBA) ColorFunc {
	return func(n int, x int, y int, maxN int, maxX int, maxY int) color.Color {
		// A single sample (or a non-positive maxN) cannot describe a range:
		// return the start color without dividing by maxN.
		if maxN <= 0 {
			return start
		}

		// Clamp n into [0, maxN] so out-of-order or out-of-range callers,
		// such as a channel drawn in reverse, cannot overshoot the gradient.
		if n <= 0 {
			return start
		}
		if n >= maxN {
			return end
		}

		return color.RGBA{
			R: lerpUint8(start.R, end.R, n, maxN),
			G: lerpUint8(start.G, end.G, n, maxN),
			B: lerpUint8(start.B, end.B, n, maxN),
			A: lerpUint8(start.A, end.A, n, maxN),
		}
	}
}

// SolidColor generates a ColorFunc which simply returns the input color
// as the color which should be drawn at all coordinates.
//
// This is the default behavior of the waveform package.
func SolidColor(inColor color.Color) ColorFunc {
	return func(n int, x int, y int, maxN int, maxX int, maxY int) color.Color {
		return inColor
	}
}

// StripeColor generates a ColorFunc which applies one color from the input,
// variadic slice at each computed value.  Each color is used in order, and
// the rotation will repeat until the image is complete. This creates a stripe
// effect in the resulting waveform image.
func StripeColor(colors ...color.Color) ColorFunc {
	// Filter any nil values
	colors = filterNilColors(colors)

	var lastN int
	return func(n int, x int, y int, maxN int, maxX int, maxY int) color.Color {
		// For each new n value, use the next color in the slice
		if n > lastN {
			lastN = n
		}

		return colors[lastN%len(colors)]
	}
}

// lerpUint8 performs integer linear interpolation between two uint8 channel
// values at position n of maxN.  It uses only integer arithmetic, so it never
// produces NaN, and does constant work per call.
//
// The caller is responsible for ensuring maxN > 0 and 0 < n < maxN; under
// those conditions the result is exact at the endpoints (returning from and
// to) and bounded between from and to for every intermediate n, so repeated
// values stay flat and the sequence remains monotonic.
func lerpUint8(from uint8, to uint8, n int, maxN int) uint8 {
	// widen to int64 so (to-from)*n cannot overflow
	return uint8((int64(from)*int64(maxN-n) + int64(to)*int64(n)) / int64(maxN))
}

// filterNilColors strips any nil color.Color values from the input slice.
func filterNilColors(colors []color.Color) []color.Color {
	var cleanColors []color.Color
	for _, c := range colors {
		if c != nil {
			cleanColors = append(cleanColors, c)
		}
	}

	return cleanColors
}
