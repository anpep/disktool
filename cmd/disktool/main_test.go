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

package main

import (
	"disktool/internal/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestExecuteMain(t *testing.T) {
	stdio, err := utils.NewFakeStdio()
	assert.NoError(t, err)

	oldArgs := os.Args
	os.Args = []string{"disktool"}
	defer func() { os.Args = oldArgs }()

	exitCode := executeMain()
	assert.Equal(t, 0, exitCode)

	stdout, stderr, err := stdio.ReadAndRestore()
	assert.NoError(t, err)
	assert.Regexp(t, `^Usage:`, stdout)
	assert.Empty(t, stderr)
}

func TestExecuteMainUnrecognizedCommand(t *testing.T) {
	stdio, err := utils.NewFakeStdio()
	assert.NoError(t, err)

	oldArgs := os.Args
	os.Args = []string{"disktool", "sfdeljknesv"}
	defer func() { os.Args = oldArgs }()

	exitCode := executeMain()
	assert.Equal(t, 1, exitCode)

	stdout, stderr, err := stdio.ReadAndRestore()
	assert.NoError(t, err)
	assert.Empty(t, stdout)
	assert.Regexp(t, "^error: Unknown command `sfdeljknesv'", stderr)
}
