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
	String() string
}

type FileWriter interface {
	io.Writer
	Fd() uintptr
	String() string
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
		In:  &file{File: os.Stdin},
		Out: &file{File: os.Stdout},
		Err: &file{File: os.Stderr},
		// TODO: check isTTY
		colorEnabled: EnvColorForced() || (!EnvColorDisabled()),
	}
}

// Returns an IOStreams with mock in/out/err values.
func Test() *IOStreams {
	return &IOStreams{
		In:           &mockFile{Buffer: &bytes.Buffer{}, fd: 0},
		Out:          &mockFile{Buffer: &bytes.Buffer{}, fd: 1},
		Err:          &mockFile{Buffer: &bytes.Buffer{}, fd: 2},
		colorEnabled: false,
	}
}

type file struct {
	*os.File
}

func (f *file) String() string {
	buf, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return string(buf)
}

type mockFile struct {
	*bytes.Buffer
	fd uintptr
}

func (m *mockFile) Fd() uintptr {
	return m.fd
}
