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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v3"
)

type formatMixin struct {
	Format string `short:"f" long:"format" description:"Output format. Valid values are: 'json', 'yaml', 'pretty' (default)"`
}

type prettyPrintable interface {
	PrettyPrint(w io.Writer) error
}

func (f *formatMixin) PrintStruct(v prettyPrintable) error {
	switch f.Format {
	case "", "pretty":
		if err := v.PrettyPrint(os.Stdout); err != nil {
			return fmt.Errorf("couldn't display output: %w", err)
		}
		return nil
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "    ")
		return encoder.Encode(v)
	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		return encoder.Encode(v)
	default:
		return fmt.Errorf("unrecognized output format: '%s'", f.Format)
	}
}

type Options struct {
	Verbose bool `short:"v" description:"Show verbose debug information"`
}

var ErrExtraArgs = errors.New("unrecognized extra arguments")

var options Options
var parser = flags.NewParser(&options, flags.HelpFlag|flags.PassDoubleDash)

func executeMain() int {
	_, err := parser.Parse()
	if err, ok := err.(*flags.Error); ok {
		if err.Type == flags.ErrHelp || err.Type == flags.ErrCommandRequired {
			parser.WriteHelp(os.Stdout)
			return 0
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(executeMain())
}
