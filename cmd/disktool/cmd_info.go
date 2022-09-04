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
	"fmt"
	"io"
	"os"

	"github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/disk"
	"github.com/jessevdk/go-flags"
)

type CmdInfo struct {
	formatMixin

	Positional struct {
		Filename flags.Filename `positional-arg-name:"FILENAME"`
	} `positional-args:"yes" required:"yes"`
}

const cmdInfoShort = `Show information about a disk image`
const cmdInfoLong = `
The info command shows information about a disk image.
`

type gptDiskInfo struct {
	Type          string `json:"type" yaml:"type"`
	FileType      string `json:"file_type" yaml:"file_type"`
	Writable      bool   `json:"writable" yaml:"writable"`
	BlockSize     int64  `json:"block_size" yaml:"block_size"`
	PhysBlockSize int64  `json:"phys_block_size" yaml:"phys_block_size"`
	Size          int64  `json:"size" yaml:"size"`
}

func (i *gptDiskInfo) PrettyPrint(w io.Writer) error {
	fmt.Fprint(w, "GPT-formatted ")

	if i.FileType == "file" {
		fmt.Fprint(w, "disk image file ")
	} else {
		fmt.Fprint(w, "block device ")
	}

	if i.Writable {
		fmt.Fprintln(w, "(read-write)")
	} else {
		fmt.Fprintln(w, "(read-only)")
	}

	fmt.Fprintf(w, "Logical Block Size: %d\n", i.BlockSize)
	fmt.Fprintf(w, "Physical Block Size: %d\n", i.PhysBlockSize)
	fmt.Fprintf(w, "Total Size: %s\n", utils.FormatBytes(i.Size))

	return nil
}

type infoOutput struct {
	DiskInfo prettyPrintable `json:"disk_info" yaml:"disk-info"`
}

func (o *infoOutput) PrettyPrint(w io.Writer) error {
	return o.DiskInfo.PrettyPrint(w)
}

func (cmd *CmdInfo) Execute(args []string) error {
	if len(args) > 0 {
		return ErrExtraArgs
	}

	d, err := diskfs.OpenWithMode(string(cmd.Positional.Filename), diskfs.ReadOnly)
	if err != nil {
		return err
	}

	table, err := d.GetPartitionTable()
	if err != nil {
		return err
	}

	switch table.Type() {
	case "mbr":
		/*info := mbrInfo{
			Type:          table.Type(),
			Writable:      d.Writable,
			BlockSize:     d.LogicalBlocksize,
			PhysBlockSize: d.PhysicalBlocksize,
			Size:          d.Size,
		}

		if d.Type == disk.File {
			info.FileType = "file"
		} else if d.Type == disk.Device {
			info.FileType = "device"
		}

		if cmd.Format == "" {
			info.PrettyPrint(os.Stdout)
			return nil
		} else {
			return cmd.PrintStruct(info)
		}*/
	case "gpt":
		info := gptDiskInfo{
			Type:          table.Type(),
			Writable:      d.Writable,
			BlockSize:     d.LogicalBlocksize,
			PhysBlockSize: d.PhysicalBlocksize,
			Size:          d.Size,
		}

		if d.Type == disk.File {
			info.FileType = "file"
		} else if d.Type == disk.Device {
			info.FileType = "device"
		}

		if cmd.Format == "" {
			info.PrettyPrint(os.Stdout)
			return nil
		} else {
			return cmd.PrintStruct(&infoOutput{DiskInfo: &info})
		}
	}

	return nil
}

func init() {
	parser.AddCommand("info", cmdInfoShort, cmdInfoLong, &CmdInfo{})
}
