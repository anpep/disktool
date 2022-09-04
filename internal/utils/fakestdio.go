// Copyright (c) 2022 Ángel Pérez <ap@anpep.co>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package utils

// The source of this file is:
// <https://github.com/eliben/code-for-blog/tree/master/2020/go-fake-stdio>

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

// FakeStdio can be used to fake stdin and capture stdout.
// Between creating a new FakeStdio and calling ReadAndRestore on it,
// code reading os.Stdin will get the contents of stdinText passed to New.
// Output to os.Stdout will be captured and returned from ReadAndRestore.
// FakeStdio is not reusable; don't attempt to use it after calling
// ReadAndRestore, but it should be safe to create a new FakeStdio.
type FakeStdio struct {
	origStdout *os.File
	origStderr *os.File

	stdoutReader *os.File
	stderrReader *os.File

	outCh chan []byte
	errCh chan []byte
}

func NewFakeStdio() (*FakeStdio, error) {
	// Pipe for stdout.
	//
	//               ======
	//  w ----------->||||------> r
	// (os.Stdout)   ======      (stdoutReader)
	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	// Pipe for stderr.
	//
	//               ======
	//  w ----------->||||------> r
	// (os.Stdout)   ======      (stdoutReader)
	stderrReader, stderrWriter, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	origStdout := os.Stdout
	origStderr := os.Stderr
	os.Stdout = stdoutWriter
	os.Stderr = stderrWriter

	outCh := make(chan []byte)
	errCh := make(chan []byte)

	// This goroutine reads stdout into a buffer in the background.
	go func() {
		var b bytes.Buffer
		if _, err := io.Copy(&b, stdoutReader); err != nil {
			log.Println(err)
		}
		outCh <- b.Bytes()
	}()

	// This goroutine reads stderr into a buffer in the background.
	go func() {
		var b bytes.Buffer
		if _, err := io.Copy(&b, stderrReader); err != nil {
			log.Println(err)
		}
		errCh <- b.Bytes()
	}()

	return &FakeStdio{
		origStdout,
		origStderr,
		stdoutReader,
		stderrReader,
		outCh,
		errCh,
	}, nil
}

// ReadAndRestore collects all captured stdout and returns it; it also restores
// os.Stdin and os.Stdout to their original values.
func (sf *FakeStdio) ReadAndRestore() (string, string, error) {
	if sf.stdoutReader == nil || sf.stderrReader == nil {
		return "", "", fmt.Errorf("ReadAndRestore from closed FakeStdio")
	}

	// Close the writer side of the faked stdout pipe. This signals to the
	// background goroutine that it should exit.
	os.Stdout.Close()
	os.Stderr.Close()
	out := <-sf.outCh
	err := <-sf.errCh

	os.Stdout = sf.origStdout
	os.Stderr = sf.origStderr

	if sf.stdoutReader != nil {
		sf.stdoutReader.Close()
		sf.stdoutReader = nil
	}

	if sf.stderrReader != nil {
		sf.stderrReader.Close()
		sf.stderrReader = nil
	}

	return string(out), string(err), nil
}
