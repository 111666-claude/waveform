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

// GradientColor generates a ColorFunc which produces a color gradient between two
// RGBA input colors.  The gradient linearly interpolates every channel (including
// alpha) between the two colors as n advances from 0 to maxN, creating a sweeping
// color change effect in the resulting waveform image.
//
// The returned function is pure: its output depends only on n and maxN, so
// repeated, out-of-order, or decreasing n values always yield consistent colors.
func GradientColor(start color.RGBA, end color.RGBA) ColorFunc {
	return func(n int, x int, y int, maxN int, maxX int, maxY int) color.Color {
		// Degenerate or out-of-range inputs clamp to the nearest endpoint,
		// so a zero maxN or a backwards n can never panic or misbehave.
		if maxN <= 0 || n <= 0 {
			return start
		}
		if n >= maxN {
			return end
		}

		// Linearly interpolate a single channel between a and b using integer
		// math only, so results are exact, monotonic, and free of NaN.
		lerp := func(a, b uint8) uint8 {
			return uint8((int(a)*(maxN-n) + int(b)*n + maxN/2) / maxN)
		}

		return color.RGBA{
			R: lerp(start.R, end.R),
			G: lerp(start.G, end.G),
			B: lerp(start.B, end.B),
			A: lerp(start.A, end.A),
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
