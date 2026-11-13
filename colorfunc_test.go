package waveform

import (
	"image/color"
	"testing"
)

// Named colors for easy testing
var (
	black = color.RGBA{0, 0, 0, 255}
	white = color.RGBA{255, 255, 255, 255}
	red   = color.RGBA{255, 0, 0, 255}
	green = color.RGBA{0, 255, 0, 255}
	blue  = color.RGBA{0, 0, 255, 255}
)

// TestCheckerColorOneColor verifies that CheckerColor produces only the single
// color used in its input.
func TestCheckerColorOneColor(t *testing.T) {
	testCheckerColor(t, black, black)
}

// TestCheckerColorTwoColors verifies that CheckerColor produces only colors
// which are used in its input.
func TestCheckerColorTwoColors(t *testing.T) {
	testCheckerColor(t, black, white)
}

// TestFuzzColorOneColor verifies that FuzzColor produces only the single
// color used in its input.
func TestFuzzColorOneColor(t *testing.T) {
	testFuzzColor(t, []color.Color{black})
}

// TestFuzzColorMultipleColors verifies that FuzzColor produces only colors
// which are used in its input.
func TestFuzzColorMultipleColors(t *testing.T) {
	testFuzzColor(t, []color.Color{black, white, red, green, blue})
}

// TestSolidColor verifies that SolidColor always returns the same input
// color, for all input values.
func TestSolidColor(t *testing.T) {
	colors := []color.Color{
		black,
		white,
	}

	for i, c := range colors {
		fn := SolidColor(c)
		if out := fn(i, i, i, i, i, i); out != c {
			t.Fatalf("unexpected SolidColor color: %v != %v", out, c)
		}
	}
}

// TestStripeColorOneColor verifies that StripeColor produces a correct
// color sequence with a single input color.
func TestStripeColorOneColor(t *testing.T) {
	testStripeColor(t, []color.Color{black}, []color.Color{
		black, black, black, black,
	})
}

// TestStripeColorMultipleColors verifies that StripeColor produces a correct
// color sequence with multiple input colors.
func TestStripeColorMultipleColors(t *testing.T) {
	testStripeColor(t, []color.Color{
		black, white, white, red, green, green, green, blue,
	}, []color.Color{
		black, white, white, red, green, green, green, blue,
		black, white, white, red, green, green, green, blue,
	})
}

// testCheckerColor is a test helper which aids in testing the CheckerColor function.
func testCheckerColor(t *testing.T, colorA color.Color, colorB color.Color) {
	// Predefined values for test
	const maxX, maxY, size = 1000, 1000, 10

	// Generate checker function with input values
	fn := CheckerColor(colorA, colorB, size)

	// Iterate all coordinates and check color at each
	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {
			// Check color at specified coordinate
			c := fn(0, x, y, 0, maxX, maxY)

			// Apply checker algorithm to determine if color A or B should be used
			if ((uint(x)/size)+(uint(y)/size))%2 == 0 {
				if c != colorA {
					t.Fatalf("unexpected color: %v != %v", c, colorA)
				}
			} else {
				if c != colorB {
					t.Fatalf("unexpected color: %v != %v", c, colorB)
				}
			}
		}
	}
}

// testFuzzColor is a test helper which aids in testing the FuzzColor function.
func testFuzzColor(t *testing.T, in []color.Color) {
	// Make a set of colors from input slice
	set := make(map[color.RGBA]struct{})
	for _, c := range in {
		set[c.(color.RGBA)] = struct{}{}
	}

	// Validate that FuzzColor only produces colors which are present in
	// the input slice.
	fn := FuzzColor(in...)
	for i := 0; i < 10000; i++ {
		if out, ok := set[fn(i, i, i, i, i, i).(color.RGBA)]; !ok {
			t.Fatalf("color not in set: %v", out)
		}
	}
}

// TestGradientColorEndpoints verifies that GradientColor returns exactly the
// start color at n=0 and exactly the end color at n=maxN.
func TestGradientColorEndpoints(t *testing.T) {
	start := color.RGBA{255, 0, 0, 128}
	end := color.RGBA{0, 255, 0, 64}
	const maxN = 1000

	fn := GradientColor(start, end)

	if out := fn(0, 0, 0, maxN, 0, 0); out != start {
		t.Fatalf("unexpected color at n=0: %v != %v", out, start)
	}
	if out := fn(maxN, 0, 0, maxN, 0, 0); out != end {
		t.Fatalf("unexpected color at n=maxN: %v != %v", out, end)
	}
}

// TestGradientColorMonotonic verifies that every channel (including alpha)
// stays within the endpoints' range and changes monotonically as n advances.
func TestGradientColorMonotonic(t *testing.T) {
	start := color.RGBA{255, 0, 0, 128}
	end := color.RGBA{0, 255, 0, 64}
	const maxN = 1000

	fn := GradientColor(start, end)

	prev := fn(0, 0, 0, maxN, 0, 0).(color.RGBA)
	for n := 1; n <= maxN; n++ {
		cur := fn(n, 0, 0, maxN, 0, 0).(color.RGBA)

		// R descends 255 -> 0, G ascends 0 -> 255, B stays 0, A descends 128 -> 64
		checkChannel(t, n, "R", prev.R, cur.R, start.R, end.R)
		checkChannel(t, n, "G", prev.G, cur.G, start.G, end.G)
		checkChannel(t, n, "B", prev.B, cur.B, start.B, end.B)
		checkChannel(t, n, "A", prev.A, cur.A, start.A, end.A)

		prev = cur
	}
}

// TestGradientColorInterpolatesAlpha verifies that alpha is interpolated along
// with the color channels, landing on the expected midpoint value.
func TestGradientColorInterpolatesAlpha(t *testing.T) {
	start := color.RGBA{255, 0, 0, 128}
	end := color.RGBA{0, 255, 0, 64}
	const maxN = 1000

	fn := GradientColor(start, end)
	mid := fn(maxN/2, 0, 0, maxN, 0, 0).(color.RGBA)

	if mid.A != 96 {
		t.Fatalf("unexpected midpoint alpha: %v != 96", mid.A)
	}
	if mid.R != 128 || mid.G != 128 {
		t.Fatalf("unexpected midpoint color: %v", mid)
	}
}

// TestGradientColorDegenerate verifies that GradientColor never panics and
// always returns a sane endpoint color for degenerate inputs: maxN=0,
// negative n, and n beyond maxN.
func TestGradientColorDegenerate(t *testing.T) {
	start := color.RGBA{255, 0, 0, 128}
	end := color.RGBA{0, 255, 0, 64}

	fn := GradientColor(start, end)

	// maxN=0 must not divide by zero or panic
	if out := fn(0, 0, 0, 0, 0, 0); out != start {
		t.Fatalf("unexpected color for maxN=0: %v != %v", out, start)
	}
	if out := fn(5, 0, 0, 0, 0, 0); out != start {
		t.Fatalf("unexpected color for n>maxN=0: %v != %v", out, start)
	}

	// Negative n clamps to start
	if out := fn(-1, 0, 0, 1000, 0, 0); out != start {
		t.Fatalf("unexpected color for negative n: %v != %v", out, start)
	}

	// n beyond maxN clamps to end
	if out := fn(1001, 0, 0, 1000, 0, 0); out != end {
		t.Fatalf("unexpected color for n>maxN: %v != %v", out, end)
	}
}

// TestGradientColorPure verifies that GradientColor is a pure function of its
// inputs: repeated calls, out-of-order n values, and backwards n sequences all
// produce identical results for identical inputs.
func TestGradientColorPure(t *testing.T) {
	start := color.RGBA{255, 0, 0, 128}
	end := color.RGBA{0, 255, 0, 64}
	const maxN = 1000

	fn := GradientColor(start, end)

	// Walk forward, then backward, then repeat individual values; the color
	// for a given n must always be the same.
	forward := make([]color.Color, maxN+1)
	for n := 0; n <= maxN; n++ {
		forward[n] = fn(n, 0, 0, maxN, 0, 0)
	}
	for n := maxN; n >= 0; n-- {
		if out := fn(n, 0, 0, maxN, 0, 0); out != forward[n] {
			t.Fatalf("n=%d: color changed between calls: %v != %v", n, out, forward[n])
		}
	}
	for n := 0; n <= maxN; n++ {
		if out := fn(n, 0, 0, maxN, 0, 0); out != forward[n] {
			t.Fatalf("n=%d: repeated call returned different color: %v != %v", n, out, forward[n])
		}
	}
}

// TestGradientColorSameEndpoints verifies that identical start and end colors
// produce a constant color at every sample.
func TestGradientColorSameEndpoints(t *testing.T) {
	c := color.RGBA{10, 20, 30, 40}
	fn := GradientColor(c, c)

	for n := 0; n <= 100; n++ {
		if out := fn(n, 0, 0, 100, 0, 0); out != c {
			t.Fatalf("n=%d: unexpected color: %v != %v", n, out, c)
		}
	}
}

// TestGradientColorManySamples verifies that GradientColor produces valid,
// in-range colors across a long waveform's worth of samples.
func TestGradientColorManySamples(t *testing.T) {
	start := color.RGBA{255, 0, 0, 128}
	end := color.RGBA{0, 255, 0, 64}
	const maxN = 6000

	fn := GradientColor(start, end)

	prev := fn(0, 0, 0, maxN, 0, 0).(color.RGBA)
	for n := 1; n <= maxN; n++ {
		cur := fn(n, 0, 0, maxN, 0, 0).(color.RGBA)

		checkChannel(t, n, "R", prev.R, cur.R, start.R, end.R)
		checkChannel(t, n, "G", prev.G, cur.G, start.G, end.G)
		checkChannel(t, n, "B", prev.B, cur.B, start.B, end.B)
		checkChannel(t, n, "A", prev.A, cur.A, start.A, end.A)

		prev = cur
	}
}

// checkChannel verifies that a channel value stays within the inclusive range
// of its endpoints and moves monotonically toward the end value.
func checkChannel(t *testing.T, n int, name string, prev, cur, start, end uint8) {
	lo, hi := start, end
	if lo > hi {
		lo, hi = hi, lo
	}

	if cur < lo || cur > hi {
		t.Fatalf("n=%d: channel %s out of range: %d not in [%d, %d]", n, name, cur, lo, hi)
	}

	if start <= end && cur < prev {
		t.Fatalf("n=%d: channel %s not monotonically increasing: %d < %d", n, name, cur, prev)
	}
	if start >= end && cur > prev {
		t.Fatalf("n=%d: channel %s not monotonically decreasing: %d > %d", n, name, cur, prev)
	}
}

// testStripeColor is a test helper which aids in testing the StripeColor function.
func testStripeColor(t *testing.T, in []color.Color, out []color.Color) {
	// Validate that StripeColor produces expected output at each index
	fn := StripeColor(in...)
	for i := 0; i < len(out); i++ {
		if c := fn(i, 0, 0, 0, 0, 0); c != out[i] {
			t.Fatalf("[%02d] unexpected output color: %v != %v", i, c, out[i])
		}
	}
}
