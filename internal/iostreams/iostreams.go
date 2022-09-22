package iostreams

import (
	"bytes"
	"io"
	"os"
)

type FileReader interface {
	io.Reader
	Fd() uintptr
}

type FileWriter interface {
	io.Writer
	Fd() uintptr
}

// Container for the three main CLI I/O streams.
type IOStreams struct {
	// os.Stdin (or mock when unit testing)
	In FileReader
	// os.Stdout (or mock when unit testing)
	Out FileWriter
	// os.Stderr (or mock when unit testing)
	Err FileWriter

	colorEnabled bool
}

func (s *IOStreams) ColorEnabled() bool {
	return s.colorEnabled
}

// ColorScheme returns a configured ColorScheme struct.
func (s *IOStreams) ColorScheme() *ColorScheme {
	return NewColorScheme(s.ColorEnabled())
}

// Returns an IOStreams containing os.Stdin, os.Stdout, and os.Stderr
func System() *IOStreams {
	return &IOStreams{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
		// TODO: check isTTY
		colorEnabled: EnvColorForced() || (!EnvColorDisabled()),
	}
}

// Returns an IOStreams with mock bytes.Buffer values.
// Also returns the raw in/out/err buffers for ease of use in tests.
func Test() (*IOStreams, *bytes.Buffer, *bytes.Buffer, *bytes.Buffer) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}

	return &IOStreams{
		In:           &fdReader{Reader: in, fd: 0},
		Out:          &fdWriter{Writer: out, fd: 1},
		Err:          &fdWriter{Writer: err, fd: 2},
		colorEnabled: false,
	}, in, out, err
}

type fdReader struct {
	io.Reader
	fd uintptr
}

func (r *fdReader) Fd() uintptr {
	return r.fd
}

type fdWriter struct {
	io.Writer
	fd uintptr
}

func (w *fdWriter) Fd() uintptr {
	return w.fd
}
