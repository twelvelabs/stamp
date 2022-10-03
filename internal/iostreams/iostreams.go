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

// IOStream represents an input or output stream.
type IOStream interface {
	io.Reader
	io.Writer
	Fd() uintptr
	String() string
}

// Container for the three main CLI I/O streams.
type IOStreams struct {
	// os.Stdin (or mock when unit testing)
	In IOStream
	// os.Stdout (or mock when unit testing)
	Out IOStream
	// os.Stderr (or mock when unit testing)
	Err IOStream

	colorEnabled bool
}

func (s *IOStreams) ColorEnabled() bool {
	return s.colorEnabled
}

// Formatter returns a ANSI string formatter.
func (s *IOStreams) Formatter() *Formatter {
	return NewFormatter(s.ColorEnabled())
}

// Returns an IOStreams containing os.Stdin, os.Stdout, and os.Stderr.
func System() *IOStreams {
	return &IOStreams{
		In:  &systemIOStream{File: os.Stdin},
		Out: &systemIOStream{File: os.Stdout},
		Err: &systemIOStream{File: os.Stderr},
		// TODO: check isTTY
		colorEnabled: IsColorEnabled(),
	}
}

// Returns an IOStreams with mock in/out/err values.
func Test() *IOStreams {
	return &IOStreams{
		In:           &mockIOStream{Buffer: &bytes.Buffer{}, fd: 0},
		Out:          &mockIOStream{Buffer: &bytes.Buffer{}, fd: 1},
		Err:          &mockIOStream{Buffer: &bytes.Buffer{}, fd: 2},
		colorEnabled: false,
	}
}

var (
	_ IOStream = &systemIOStream{}
	_ IOStream = &mockIOStream{}
)

// Wrapper so we can make os.Stdin and friends fulfill IOStream.
type systemIOStream struct {
	*os.File
}

func (f *systemIOStream) String() string {
	buf, _ := io.ReadAll(f)
	return string(buf)
}

// Wrapper so we can make bytes.Buffer fulfill IOStream.
type mockIOStream struct {
	*bytes.Buffer
	fd uintptr
}

func (m *mockIOStream) Fd() uintptr {
	return m.fd
}
