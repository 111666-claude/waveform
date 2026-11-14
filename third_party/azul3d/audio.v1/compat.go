// Compatibility aliases restoring the type names used by the
// azul3d.org/audio.v1 era of the API.  The engine rewrite renamed the
// float sample slices (F64Samples -> Float64, F32Samples -> Float32);
// these aliases keep existing callers and the waveform package's public
// SampleReduceFunc signature source-compatible.

package audio

// F64 is a 64-bit floating-point audio sample.
type F64 = float64

// F32 is a 32-bit floating-point audio sample.
type F32 = float32

// F64Samples represents a slice of F64 encoded audio samples.
type F64Samples = Float64

// F32Samples represents a slice of F32 encoded audio samples.
type F32Samples = Float32
