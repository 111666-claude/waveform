module github.com/mdlayher/waveform

go 1.27.4

require (
	azul3d.org/audio.v1 v1.0.0
	azul3d.org/audio/wav.v1 v1.0.0
	github.com/azul3d/audio-flac v1.0.0
)

require (
	github.com/icza/bitio v1.1.0 // indirect
	github.com/mewkiz/flac v1.0.14 // indirect
	github.com/mewkiz/pkg v0.0.0-20250417130911-3f050ff8c56d // indirect
	github.com/mewpkg/term v0.0.0-20241026122259-37a80af23985 // indirect
)

replace (
	azul3d.org/audio.v1 => ./third_party/azul3d/audio.v1
	azul3d.org/audio/wav.v1 => ./third_party/azul3d/audio-wav.v1
	github.com/azul3d/audio-flac => ./third_party/azul3d/audio-flac
)
