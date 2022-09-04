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
	"github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/partition/gpt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jessevdk/go-flags"
	"io"
	"strings"
)

type CmdLs struct {
	formatMixin

	ShowAll bool `short:"a" long:"--show-all" description:"Show all entries in the partition table, even the empty ones"`

	Positional struct {
		Filename flags.Filename `positional-arg-name:"FILENAME"`
	} `positional-args:"yes" required:"yes"`
}

const cmdLsShort = `List partitions on a disk image`
const cmdLsLong = `
The ls command displays a list of partitions present on a disk image.
`

type partInfo struct {
	Index   int   `json:"index" yaml:"index"`
	Size    int64 `json:"size" yaml:"size"`
	Sectors int64 `json:"sectors" yaml:"sectors"`
	Start   int64 `json:"start" yaml:"start"`
	End     int64 `json:"end" yaml:"end"`

	GPT *gptPartInfo `json:"gpt_info" yaml:"gpt-info"`
}

type gptPartInfo struct {
	GUID     string   `json:"guid" yaml:"guid"`
	Type     string   `json:"type" yaml:"type"`
	TypeGUID string   `json:"type_guid" yaml:"type-guid"`
	Name     string   `json:"name" yaml:"name"`
	Attrs    []string `json:"attrs" yaml:"attrs"`
}

var partTypes = map[gpt.Type]string{
	gpt.Unused:                   "unused",
	gpt.MbrBoot:                  "mbr",
	gpt.EFISystemPartition:       "efi",
	gpt.BiosBoot:                 "bios",
	gpt.MicrosoftReserved:        "ms_reserved",
	gpt.MicrosoftBasicData:       "ms_basic_data",
	gpt.MicrosoftLDMMetadata:     "ms_ldm_meta",
	gpt.MicrosoftLDMData:         "ms_ldm_data",
	gpt.MicrosoftWindowsRecovery: "ms_winrecovery",
	gpt.LinuxFilesystem:          "linux",
	gpt.LinuxRaid:                "linux_raid",
	gpt.LinuxRootX86:             "linux_root_x86",
	gpt.LinuxRootX86_64:          "linux_root_x86_64",
	gpt.LinuxRootArm32:           "linux_root_arm32",
	gpt.LinuxRootArm64:           "linux_root_arm64",
	gpt.LinuxSwap:                "linux_swap",
	gpt.LinuxLVM:                 "linux_lvm",
	gpt.LinuxDMCrypt:             "linux_dmcrypt",
	gpt.LinuxLUKS:                "linux_luks",
	gpt.VMWareFilesystem:         "vmware",
}

type lsOutput struct {
	Partitions []partInfo `json:"partitions" yaml:"partitions"`
}

func (o *lsOutput) PrettyPrint(w io.Writer) error {
	t := table.NewWriter()
	t.SetOutputMirror(w)

	t.AppendHeader(table.Row{"#", "GUID", "Name", "Size", "Sectors", "Start", "End", "Type", "Type GUID", "Attributes"})
	for i, p := range o.Partitions {
		if p.GPT == nil || i == 0 {
			t.AppendRow(table.Row{p.Index, "", "", utils.FormatBytes(p.Size), p.Sectors, p.Start, p.End})
		} else {
			t.AppendRow(table.Row{
				p.Index,
				p.GPT.GUID,
				p.GPT.Name,
				utils.FormatBytes(p.Size),
				p.Sectors,
				p.Start,
				p.End,
				p.GPT.Type,
				p.GPT.TypeGUID,
				strings.Join(p.GPT.Attrs, ", "),
			})
		}
	}
	t.Render()

	return nil
}

func (cmd *CmdLs) Execute(args []string) error {
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

	var gptTable []*gpt.Partition = nil
	gptInfo, err := gpt.Read(d.File, int(d.LogicalBlocksize), int(d.PhysicalBlocksize))
	isGPT := err == nil
	if isGPT {
		gptTable = gptInfo.Partitions
	}

	parts := make([]partInfo, 0)
	for i, part := range table.GetPartitions() {
		start, size := part.GetStart(), part.GetSize()
		isEmptyEntry := start == 0 || size == 0
		if isEmptyEntry && !cmd.ShowAll {
			continue
		}

		info := partInfo{
			Index:   i,
			Size:    size,
			Sectors: size / d.LogicalBlocksize,
			Start:   start / d.LogicalBlocksize,
			End:     (start / d.LogicalBlocksize) + (size / d.LogicalBlocksize) - 1,
		}

		if isGPT {
			info.GPT = &gptPartInfo{
				GUID:     gptTable[i].GUID,
				TypeGUID: string(gptTable[i].Type),
				Type:     partTypes[gptTable[i].Type],
				Name:     gptTable[i].Name,
				Attrs:    make([]string, 0),
			}

			knownAttributes := map[int]string{
				0:  "platform_required",
				1:  "efi_ignore",
				2:  "bios_bootable",
				56: "cros_boot_ok",
				60: "ms_basic_read_only",
				61: "ms_basic_shadow_copy",
				62: "ms_basic_hidden",
				63: "ms_basic_no_automount",
			}

			for b, s := range knownAttributes {
				if gptTable[i].Attributes&(1<<b) != 0 {
					info.GPT.Attrs = append(info.GPT.Attrs, s)
				}
			}
		}

		parts = append(parts, info)
	}

	return cmd.PrintStruct(&lsOutput{Partitions: parts})
}

func init() {
	parser.AddCommand("ls", cmdLsShort, cmdLsLong, &CmdLs{})
}
