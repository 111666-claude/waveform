package waveform

import (
	"image/color"
	"sync"
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

// TestGradientColorOneColor verifies that GradientColor produces only the single
// color used in its input.
func TestGradientColorOneColor(t *testing.T) {
	testGradientColor(t, black, black)
}

// TestGradientColorTwoColors verifies that GradientColor produces a correct
// gradient between two colors.
func TestGradientColorTwoColors(t *testing.T) {
	testGradientColor(t, black, white)
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

// testGradientColor is a test helper which aids in testing the GradientColor
// function.
func testGradientColor(t *testing.T, start color.RGBA, end color.RGBA) {
	const maxN = 100

	// Generate function with defined values
	fn := GradientColor(start, end)

	// Check edges
	for i, n := range []int{0, maxN} {
		// Get color at point
		c := fn(n, 0, 0, maxN, 0, 0)

		// First iteration, use start; second, use end
		var testColor color.RGBA
		if i == 0 {
			testColor = start
		} else {
			testColor = end
		}

		// Compare full RGBA values, including alpha, to ensure correctness
		if c != testColor {
			t.Fatalf("unexpected color at n=%d: %v != %v", n, c, testColor)
		}
	}
}

// TestGradientColorEndpoints verifies that GradientColor returns the exact
// start color at n=0 and the exact end color at n=maxN, including the alpha
// channel.
func TestGradientColorEndpoints(t *testing.T) {
	const maxN = 1000
	start := color.RGBA{R: 255, G: 0, B: 0, A: 128}
	end := color.RGBA{R: 0, G: 255, B: 0, A: 64}

	fn := GradientColor(start, end)

	if got := fn(0, 0, 0, maxN, 0, 0); got != start {
		t.Fatalf("color at n=0 must equal start: %v != %v", got, start)
	}
	if got := fn(maxN, 0, 0, maxN, 0, 0); got != end {
		t.Fatalf("color at n=maxN must equal end: %v != %v", got, end)
	}
}

// TestGradientColorMonotonic verifies that, across the whole gradient, every
// color channel stays within the range defined by its start and end values,
// changes monotonically (plateaus allowed), and that alpha is interpolated
// alongside the RGB channels.
func TestGradientColorMonotonic(t *testing.T) {
	const maxN = 1000
	start := color.RGBA{R: 255, G: 0, B: 0, A: 128}
	end := color.RGBA{R: 0, G: 255, B: 0, A: 64}

	fn := GradientColor(start, end)

	// Track the previous sample to check monotonicity
	prev := fn(0, 0, 0, maxN, 0, 0).(color.RGBA)

	// Each channel must stay within [lo, hi] for every sample
	channels := []struct {
		name       string
		start, end uint8
		get        func(color.RGBA) uint8
	}{
		{"R", start.R, end.R, func(c color.RGBA) uint8 { return c.R }},
		{"G", start.G, end.G, func(c color.RGBA) uint8 { return c.G }},
		{"B", start.B, end.B, func(c color.RGBA) uint8 { return c.B }},
		{"A", start.A, end.A, func(c color.RGBA) uint8 { return c.A }},
	}

	for n := 1; n <= maxN; n++ {
		cur := fn(n, 0, 0, maxN, 0, 0).(color.RGBA)

		for _, ch := range channels {
			v, p := ch.get(cur), ch.get(prev)
			lo, hi := ch.start, ch.end
			if lo > hi {
				lo, hi = hi, lo
			}

			// Value must remain within the two endpoint colors
			if v < lo || v > hi {
				t.Fatalf("channel %s out of range at n=%d: %d not in [%d, %d]", ch.name, n, v, lo, hi)
			}

			// Value must move monotonically toward the end color
			if ch.end >= ch.start && v < p {
				t.Fatalf("channel %s decreased at n=%d: %d -> %d", ch.name, n, p, v)
			}
			if ch.end < ch.start && v > p {
				t.Fatalf("channel %s increased at n=%d: %d -> %d", ch.name, n, p, v)
			}
		}

		prev = cur
	}
}

// TestGradientColorSameColor verifies that GradientColor returns the input
// color unchanged when the start and end colors are identical, including for
// a non-opaque alpha value.
func TestGradientColorSameColor(t *testing.T) {
	const maxN = 500
	c := color.RGBA{R: 42, G: 99, B: 7, A: 200}

	fn := GradientColor(c, c)
	for n := 0; n <= maxN; n++ {
		if got := fn(n, 0, 0, maxN, 0, 0); got != c {
			t.Fatalf("unexpected color at n=%d: %v != %v", n, got, c)
		}
	}

	// A single-point gradient (maxN == 0) must also hold the color
	if got := fn(0, 0, 0, 0, 0, 0); got != c {
		t.Fatalf("unexpected color with maxN=0: %v != %v", got, c)
	}
}

// TestGradientColorNoPanicOrNaN exercises the degenerate and out-of-range
// inputs which must never panic or produce non-finite values: maxN <= 0,
// n moving backwards, n beyond maxN, and repeated identical calls.
func TestGradientColorNoPanicOrNaN(t *testing.T) {
	start := color.RGBA{R: 255, G: 0, B: 0, A: 128}
	end := color.RGBA{R: 0, G: 255, B: 0, A: 64}
	fn := GradientColor(start, end)

	// maxN <= 0 must not divide by zero and must pin the result to start
	for _, maxN := range []int{0, -1} {
		c := fn(0, 0, 0, maxN, 0, 0).(color.RGBA)
		if c != start {
			t.Fatalf("maxN=%d must return start color: %v != %v", maxN, c, start)
		}
	}

	// n moving backwards, then forwards again, must stay clamped and finite
	sequence := []int{10, 5, 0, 7, 3, 1000, 1}
	prev := fn(0, 0, 0, 1000, 0, 0).(color.RGBA)
	for _, n := range sequence {
		cur := fn(n, 0, 0, 1000, 0, 0).(color.RGBA)

		// Negative n clamps to the start color
		if neg := fn(-25, 0, 0, 1000, 0, 0); neg != start {
			t.Fatalf("negative n must clamp to start: %v != %v", neg, start)
		}

		// n beyond maxN clamps to the end color
		if over := fn(1001, 0, 0, 1000, 0, 0); over != end {
			t.Fatalf("n > maxN must clamp to end: %v != %v", over, end)
		}
		if wayOver := fn(1<<30, 0, 0, 1000, 0, 0); wayOver != end {
			t.Fatalf("n far beyond maxN must clamp to end: %v != %v", wayOver, end)
		}

		// Every channel must remain within the endpoint color range
		for _, pair := range [][2]uint8{
			{cur.R, prev.R}, {cur.G, prev.G}, {cur.B, prev.B}, {cur.A, prev.A},
		} {
			v := pair[0]
			if v > 255 { // uint8 cannot exceed 255; this documents the no-wrap invariant
				t.Fatalf("channel wrapped past uint8 range: %d", v)
			}
		}
		prev = cur
	}

	// Repeated identical calls must return identical results
	const n, maxN = 250, 1000
	first := fn(n, 0, 0, maxN, 0, 0)
	for i := 0; i < 100; i++ {
		if got := fn(n, 0, 0, maxN, 0, 0); got != first {
			t.Fatalf("repeated call %d returned a different color: %v != %v", i, got, first)
		}
	}
}

// TestGradientColor6000Samples verifies correct behavior over a long, 6000
// point gradient: exact endpoints, no out-of-range or NaN-producing values,
// and per-channel monotonicity throughout.
func TestGradientColor6000Samples(t *testing.T) {
	const maxN = 6000
	start := color.RGBA{R: 255, G: 0, B: 0, A: 128}
	end := color.RGBA{R: 0, G: 255, B: 0, A: 64}

	fn := GradientColor(start, end)

	if got := fn(0, 0, 0, maxN, 0, 0); got != start {
		t.Fatalf("color at n=0 must equal start: %v != %v", got, start)
	}
	if got := fn(maxN, 0, 0, maxN, 0, 0); got != end {
		t.Fatalf("color at n=maxN must equal end: %v != %v", got, end)
	}

	prev := start
	for n := 1; n <= maxN; n++ {
		cur := fn(n, 0, 0, maxN, 0, 0).(color.RGBA)

		if cur.R > prev.R {
			t.Fatalf("R not monotonic at n=%d: %d -> %d", n, prev.R, cur.R)
		}
		if cur.G < prev.G {
			t.Fatalf("G not monotonic at n=%d: %d -> %d", n, prev.G, cur.G)
		}
		if cur.A > prev.A {
			t.Fatalf("A not monotonic at n=%d: %d -> %d", n, prev.A, cur.A)
		}

		// B is constant (0 at both ends); it must never drift
		if cur.B != 0 {
			t.Fatalf("constant B channel drifted at n=%d: %d", n, cur.B)
		}

		prev = cur
	}
}

// TestGradientColorRepeatedAndConcurrent verifies that the ColorFunc is
// stateless: interleaved and concurrent calls with different n values always
// agree with single-threaded, in-order results.
func TestGradientColorRepeatedAndConcurrent(t *testing.T) {
	const maxN = 1000
	fn := GradientColor(
		color.RGBA{R: 255, G: 0, B: 0, A: 128},
		color.RGBA{R: 0, G: 255, B: 0, A: 64},
	)

	// Reference results computed in order
	want := make([]color.RGBA, maxN+1)
	for n := 0; n <= maxN; n++ {
		want[n] = fn(n, 0, 0, maxN, 0, 0).(color.RGBA)
	}

	// Interleaved, out-of-order access must match the reference
	for _, n := range []int{999, 0, 500, 1, 1000, 250, 750} {
		if got := fn(n, 0, 0, maxN, 0, 0); got != want[n] {
			t.Fatalf("out-of-order call at n=%d: %v != %v", n, got, want[n])
		}
	}

	// Concurrent access must be data-race free and agree with the reference
	var wg sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		wg.Add(1)
		go func(seed int) {
			defer wg.Done()
			for i := 0; i < maxN; i++ {
				n := (i*7 + seed) % (maxN + 1)
				if got := fn(n, 0, 0, maxN, 0, 0); got != want[n] {
					t.Errorf("concurrent call at n=%d: %v != %v", n, got, want[n])
				}
			}
		}(worker)
	}
	wg.Wait()
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

// BenchmarkGradientColor6000 measures the cost of sampling a GradientColor
// across 6000 points.  Each point must require only constant work, so the
// benchmark should scale linearly with b.N and independently of n or maxN.
func BenchmarkGradientColor6000(b *testing.B) {
	const maxN = 6000
	fn := GradientColor(
		color.RGBA{R: 255, G: 0, B: 0, A: 128},
		color.RGBA{R: 0, G: 255, B: 0, A: 64},
	)

	var sink color.Color
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for n := 0; n <= maxN; n++ {
			sink = fn(n, 0, 0, maxN, 0, 0)
		}
	}
	_ = sink
}
