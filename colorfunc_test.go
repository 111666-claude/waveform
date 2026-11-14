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

// TestGradientColorEndpoints verifies that GradientColor returns the exact
// start color at n == 0 and the exact end color at n == maxN, including the
// alpha channel.
func TestGradientColorEndpoints(t *testing.T) {
	const maxN = 1000
	start := color.RGBA{R: 255, G: 0, B: 0, A: 128}
	end := color.RGBA{R: 0, G: 255, B: 0, A: 64}
	fn := GradientColor(start, end)

	if c := fn(0, 0, 0, maxN, 0, 0); c != start {
		t.Fatalf("n == 0: unexpected color: %v != %v", c, start)
	}

	if c := fn(maxN, 0, 0, maxN, 0, 0); c != end {
		t.Fatalf("n == maxN: unexpected color: %v != %v", c, end)
	}
}

// TestGradientColorMonotonic verifies that every interpolated channel stays
// within its two endpoints and changes monotonically as n advances, and that
// the alpha channel is interpolated rather than dropped.
func TestGradientColorMonotonic(t *testing.T) {
	const maxN = 1000
	start := color.RGBA{R: 255, G: 0, B: 0, A: 128}
	end := color.RGBA{R: 0, G: 255, B: 0, A: 64}
	fn := GradientColor(start, end)

	prev := start
	channels := []struct {
		name       string
		from, to   uint8
		get        func(color.RGBA) uint8
		descending bool
	}{
		{"R", start.R, end.R, func(c color.RGBA) uint8 { return c.R }, start.R > end.R},
		{"G", start.G, end.G, func(c color.RGBA) uint8 { return c.G }, start.G > end.G},
		{"B", start.B, end.B, func(c color.RGBA) uint8 { return c.B }, start.B > end.B},
		{"A", start.A, end.A, func(c color.RGBA) uint8 { return c.A }, start.A > end.A},
	}

	for n := 0; n <= maxN; n++ {
		c := fn(n, 0, 0, maxN, 0, 0).(color.RGBA)

		for _, ch := range channels {
			v := ch.get(c)

			// Every sample must stay within the inclusive endpoint range.
			lo, hi := ch.from, ch.to
			if lo > hi {
				lo, hi = hi, lo
			}
			if v < lo || v > hi {
				t.Fatalf("n == %d, channel %s: %d outside endpoint range [%d, %d]",
					n, ch.name, v, lo, hi)
			}

			// Each step must move only toward the end color.
			pv := ch.get(prev)
			if ch.descending && v > pv {
				t.Fatalf("n == %d, channel %s: not monotonic: %d -> %d",
					n, ch.name, pv, v)
			}
			if !ch.descending && v < pv {
				t.Fatalf("n == %d, channel %s: not monotonic: %d -> %d",
					n, ch.name, pv, v)
			}
		}

		prev = c
	}

	// Alpha must actually be interpolated: an interior sample cannot keep
	// either endpoint alpha value.
	mid := fn(maxN/2, 0, 0, maxN, 0, 0).(color.RGBA)
	if mid.A == start.A && mid.A == end.A {
		t.Fatalf("alpha channel was not interpolated, got %d", mid.A)
	}
}

// TestGradientColorSameStartEnd verifies that equal endpoints produce a
// constant color across the whole sequence.
func TestGradientColorSameStartEnd(t *testing.T) {
	const maxN = 1000
	c := color.RGBA{R: 42, G: 99, B: 7, A: 200}
	fn := GradientColor(c, c)

	for n := 0; n <= maxN; n++ {
		if got := fn(n, 0, 0, maxN, 0, 0); got != c {
			t.Fatalf("n == %d: unexpected color: %v != %v", n, got, c)
		}
	}
}

// TestGradientColorZeroMaxN verifies that a zero or negative maxN, where the
// gradient range is undefined, returns the start color without panicking or
// dividing by zero.
func TestGradientColorZeroMaxN(t *testing.T) {
	start := color.RGBA{R: 255, G: 0, B: 0, A: 128}
	end := color.RGBA{R: 0, G: 255, B: 0, A: 64}
	fn := GradientColor(start, end)

	for _, maxN := range []int{0, -1} {
		for _, n := range []int{-1, 0, 1} {
			if got := fn(n, 0, 0, maxN, 0, 0); got != start {
				t.Fatalf("maxN == %d, n == %d: unexpected color: %v != %v",
					maxN, n, got, start)
			}
		}
	}
}

// TestGradientColorClamping verifies that n values outside [0, maxN] are
// clamped to the nearest endpoint instead of overshooting it.
func TestGradientColorClamping(t *testing.T) {
	const maxN = 1000
	start := color.RGBA{R: 255, G: 0, B: 0, A: 128}
	end := color.RGBA{R: 0, G: 255, B: 0, A: 64}
	fn := GradientColor(start, end)

	if got := fn(-5, 0, 0, maxN, 0, 0); got != start {
		t.Fatalf("negative n: unexpected color: %v != %v", got, start)
	}
	if got := fn(maxN+5, 0, 0, maxN, 0, 0); got != end {
		t.Fatalf("n > maxN: unexpected color: %v != %v", got, end)
	}
}

// TestGradientColorInvariance verifies that the color is a pure function of
// (n, maxN): repeated calls, calls after n moves backwards, and shuffled call
// order all return the same color as a direct call.  This covers drawing
// callbacks invoked repeatedly or with channels drawn in reverse.
func TestGradientColorInvariance(t *testing.T) {
	const maxN = 1000
	fn := GradientColor(
		color.RGBA{R: 255, G: 0, B: 0, A: 128},
		color.RGBA{R: 0, G: 255, B: 0, A: 64},
	)

	// Reference colors computed in ascending order.
	ref := make(map[int]color.RGBA)
	for n := 0; n <= maxN; n++ {
		ref[n] = fn(n, 0, 0, maxN, 0, 0).(color.RGBA)
	}

	// Repeated calls at the same n must be identical.
	for _, n := range []int{0, 1, 2, maxN / 2, maxN - 1, maxN} {
		first := fn(n, 0, 0, maxN, 0, 0)
		for i := 0; i < 5; i++ {
			if got := fn(n, 0, 0, maxN, 0, 0); got != first {
				t.Fatalf("n == %d: repeated call changed color: %v != %v",
					n, got, first)
			}
		}
	}

	// Descending order must return the same colors as ascending order.
	for n := maxN; n >= 0; n-- {
		if got := fn(n, 0, 0, maxN, 0, 0); got != ref[n] {
			t.Fatalf("n == %d: reverse-order color mismatch: %v != %v",
				n, got, ref[n])
		}
	}

	// A shuffled sequence interleaving repeats and backwards jumps.
	shuffled := []int{3, 700, 3, 0, 999, 250, 250, 500, 1, 1000, 1000, 42}
	for _, n := range shuffled {
		if got := fn(n, 0, 0, maxN, 0, 0); got != ref[n] {
			t.Fatalf("n == %d: shuffled-order color mismatch: %v != %v",
				n, got, ref[n])
		}
	}
}

// TestGradientColorSingleStep verifies the maxN == 1 boundary, where the only
// valid positions are the two endpoints.
func TestGradientColorSingleStep(t *testing.T) {
	start := color.RGBA{R: 10, G: 20, B: 30, A: 40}
	end := color.RGBA{R: 200, G: 100, B: 50, A: 250}
	fn := GradientColor(start, end)

	if got := fn(0, 0, 0, 1, 0, 0); got != start {
		t.Fatalf("maxN == 1, n == 0: %v != %v", got, start)
	}
	if got := fn(1, 0, 0, 1, 0, 0); got != end {
		t.Fatalf("maxN == 1, n == 1: %v != %v", got, end)
	}
}

// TestGradientColorProperties exhaustively checks several arbitrary endpoint
// pairs, asserting endpoint exactness, range bounds, and monotonicity for
// every interior sample.
func TestGradientColorProperties(t *testing.T) {
	const maxN = 6000
	pairs := []struct {
		start, end color.RGBA
	}{
		{color.RGBA{0, 0, 0, 0}, color.RGBA{255, 255, 255, 255}},
		{color.RGBA{255, 255, 255, 255}, color.RGBA{0, 0, 0, 0}},
		{color.RGBA{1, 254, 128, 7}, color.RGBA{254, 1, 127, 250}},
		{color.RGBA{123, 45, 67, 89}, color.RGBA{123, 45, 67, 89}},
		{color.RGBA{0, 0, 0, 255}, color.RGBA{0, 255, 0, 32}},
	}

	for _, p := range pairs {
		fn := GradientColor(p.start, p.end)

		if got := fn(0, 0, 0, maxN, 0, 0); got != p.start {
			t.Fatalf("start mismatch: %v != %v", got, p.start)
		}
		if got := fn(maxN, 0, 0, maxN, 0, 0); got != p.end {
			t.Fatalf("end mismatch: %v != %v", got, p.end)
		}

		prev := p.start
		for n := 1; n < maxN; n++ {
			c := fn(n, 0, 0, maxN, 0, 0).(color.RGBA)
			got := [4]uint8{c.R, c.G, c.B, c.A}
			from := [4]uint8{p.start.R, p.start.G, p.start.B, p.start.A}
			to := [4]uint8{p.end.R, p.end.G, p.end.B, p.end.A}

			for i := 0; i < 4; i++ {
				lo, hi := from[i], to[i]
				if lo > hi {
					lo, hi = hi, lo
				}
				if got[i] < lo || got[i] > hi {
					t.Fatalf("pair %v->%v, n == %d, channel %d: %d out of [%d, %d]",
						p.start, p.end, n, i, got[i], lo, hi)
				}

				pv := [4]uint8{prev.R, prev.G, prev.B, prev.A}[i]
				if from[i] <= to[i] && got[i] < pv {
					t.Fatalf("pair %v->%v, n == %d, channel %d: non-monotonic increase",
						p.start, p.end, n, i)
				}
				if from[i] >= to[i] && got[i] > pv {
					t.Fatalf("pair %v->%v, n == %d, channel %d: non-monotonic decrease",
						p.start, p.end, n, i)
				}
			}
			prev = c
		}
	}
}

// BenchmarkGradientColor6000 measures the cost of sampling a long gradient;
// each sample must require only constant work.
func BenchmarkGradientColor6000(b *testing.B) {
	const maxN = 6000
	fn := GradientColor(
		color.RGBA{R: 255, G: 0, B: 0, A: 128},
		color.RGBA{R: 0, G: 255, B: 0, A: 64},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for n := 0; n < maxN; n++ {
			_ = fn(n, 0, 0, maxN, 0, 0)
		}
	}
}
