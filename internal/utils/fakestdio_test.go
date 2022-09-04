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

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestFakeOut(t *testing.T) {
	var tests = []struct {
		wantOut string
	}{
		{"nope"},
		{"joe\n"},
		{"line1\nline2"},
		{"line1\nline2\n"},
		{strings.Repeat("joe ", 100)},
		{strings.Repeat("xyz\n", 300)},
	}

	for _, tt := range tests {
		testName := tt.wantOut
		if len(testName) > 30 {
			testName = testName[:30]
		}

		t.Run(testName, func(t *testing.T) {
			fs, err := NewFakeStdio()
			if err != nil {
				t.Fatal(err)
			}

			fmt.Print(tt.wantOut)

			s, _, err := fs.ReadAndRestore()
			if err != nil {
				t.Fatal(err)
			}

			if s != tt.wantOut {
				t.Errorf("got %q, want %q", s, tt.wantOut)
			}
		})
	}
}

func TestFakeOutLarge(t *testing.T) {
	fs, err := NewFakeStdio()
	if err != nil {
		t.Fatal(err)
	}

	var want strings.Builder
	for i := 0; i < 500000; i++ {
		snippet := strconv.Itoa(i)
		fmt.Print(snippet)
		want.WriteString(snippet)
	}

	s, _, err := fs.ReadAndRestore()
	if err != nil {
		t.Fatal(err)
	}

	if want.String() != s {
		t.Errorf("got %v, want %v", s, want)
	}
}

func TestFakeErr(t *testing.T) {
	var tests = []struct {
		wantErr string
	}{
		{"nope"},
		{"joe\n"},
		{"line1\nline2"},
		{"line1\nline2\n"},
		{strings.Repeat("joe ", 100)},
		{strings.Repeat("xyz\n", 300)},
	}

	for _, tt := range tests {
		testName := tt.wantErr
		if len(testName) > 30 {
			testName = testName[:30]
		}

		t.Run(testName, func(t *testing.T) {
			fs, err := NewFakeStdio()
			if err != nil {
				t.Fatal(err)
			}

			fmt.Fprint(os.Stderr, tt.wantErr)

			_, s, err := fs.ReadAndRestore()
			if err != nil {
				t.Fatal(err)
			}

			if s != tt.wantErr {
				t.Errorf("got %q, want %q", s, tt.wantErr)
			}
		})
	}
}

func TestFakeErrLarge(t *testing.T) {
	fs, err := NewFakeStdio()
	if err != nil {
		t.Fatal(err)
	}

	var want strings.Builder
	for i := 0; i < 500000; i++ {
		snippet := strconv.Itoa(i)
		fmt.Fprint(os.Stderr, snippet)
		want.WriteString(snippet)
	}

	_, s, err := fs.ReadAndRestore()
	if err != nil {
		t.Fatal(err)
	}

	if want.String() != s {
		t.Errorf("got %v, want %v", s, want)
	}
}
