package iostreams

/*
This file started out as a copy of https://github.com/cli/cli/blob/trunk/pkg/iostreams/iostreams.go
Original license:

MIT License

Copyright (c) 2019 GitHub Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
*/

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

// Formatter returns a ANSI string formatter.
func (s *IOStreams) Formatter() *Formatter {
	return NewFormatter(s.ColorEnabled())
}

func EnvColorDisabled() bool {
	// See: https://bixense.com/clicolors/
	return os.Getenv("NO_COLOR") != "" || os.Getenv("CLICOLOR") == "0"
}

func EnvColorForced() bool {
	// See: https://bixense.com/clicolors/
	return os.Getenv("CLICOLOR_FORCE") != "" && os.Getenv("CLICOLOR_FORCE") != "0"
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
