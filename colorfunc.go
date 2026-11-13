package waveform

import (
	"image/color"
	"math"
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

// GradientColor generates a ColorFunc which produces a color gradient between two
// RGBA input colors.  The gradient gradually interpolates every channel, including
// alpha, from the start color at n=0 to the end color at n=maxN, creating a
// sweeping color change effect in the resulting waveform image.
//
// The returned ColorFunc is stateless: it holds no mutable data between calls,
// so repeated, out-of-order, or concurrent invocations always produce the same
// color for the same input values.  The interpolation factor derived from n and
// maxN is clamped to [0, 1], so n values outside the [0, maxN] range do not
// wrap around, and maxN <= 0 returns the start color instead of dividing by zero.
// Each invocation performs a constant amount of work regardless of n or maxN.
func GradientColor(start color.RGBA, end color.RGBA) ColorFunc {
	// Float equivalents of color values, and the per-channel distance to
	// travel from the start color to the end color.  These are computed once,
	// so the closure itself only does constant work on each call.
	startR, startG, startB, startA := float64(start.R), float64(start.G), float64(start.B), float64(start.A)
	deltaR := float64(end.R) - startR
	deltaG := float64(end.G) - startG
	deltaB := float64(end.B) - startB
	deltaA := float64(end.A) - startA

	return func(n int, x int, y int, maxN int, maxX int, maxY int) color.Color {
		// Interpolation factor across the waveform, clamped to [0, 1].
		// A non-positive maxN would divide by zero, so treat it as a
		// single-point gradient pinned at the start color.
		t := 0.0
		if maxN > 0 {
			t = float64(n) / float64(maxN)
			if t < 0.0 {
				t = 0.0
			} else if t > 1.0 {
				t = 1.0
			}
		}

		// Linearly interpolate each channel, including alpha
		return color.RGBA{
			R: uint8(math.Round(startR + deltaR*t)),
			G: uint8(math.Round(startG + deltaG*t)),
			B: uint8(math.Round(startB + deltaB*t)),
			A: uint8(math.Round(startA + deltaA*t)),
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
