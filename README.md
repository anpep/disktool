# disktool
> Tool to work with raw disk images and its partitions.

## 🚧 Disclaimer 🚧
This tool has been primarily conceived to work with disk images and not physical
drives. This software is only experimental and is not yet reliably tested, and
we can't guarantee this program's correctness. **Proceed at your own risk.**

## Sample usage
### Get information about a disk image
```
$ disktool info /dev/sda
GPT-formatted block device (read-only)
Logical Block Size: 512
Physical Block Size: 512
Total Size: 5.0 GiB
```

### Get YAML/JSON output in any subcommand
```
$ disktool info /dev/sda -f yaml
disk-info:
    type: gpt
    file_type: device
    writable: false
    block_size: 512
    phys_block_size: 512
    size: 5368709120
    
$ disktool info /dev/sda -f json
{
    "disk_info": {
        "type": "gpt",
        "file_type": "device",
        "writable": false,
        "block_size": 512,
        "phys_block_size": 512,
        "size": 5368709120
    }
}
```

### List partitions
```
$ disktool ls /dev/sda
+----+--------------------------------------+------+----------+----------+--------+----------+------+--------------------------------------+------------+
|  # | GUID                                 | NAME | SIZE     |  SECTORS |  START |      END | TYPE | TYPE GUID                            | ATTRIBUTES |
+----+--------------------------------------+------+----------+----------+--------+----------+------+--------------------------------------+------------+
|  0 |                                      |      | 4.9 GiB  | 10278879 | 206848 | 10485726 |      |                                      |            |
| 14 | 624E86F2-B792-4DD9-981A-3CF034269DC9 |      | 99.0 MiB |   202753 |   2048 |   204800 | efi  | C12A7328-F81F-11D2-BA4B-00A0C93EC93B |            |
+----+--------------------------------------+------+----------+----------+--------+----------+------+--------------------------------------+------------+

$ disktool ls /dev/sda -f yaml
partitions:
    - index: 0
      size: 5262786048
      sectors: 10278879
      start: 206848
      end: 10485726
      gpt-info:
        guid: 4B7414F6-AC72-45F9-BD47-F549FEFF9CF0
        type: linux
        type-guid: 0FC63DAF-8483-4772-8E79-3D69D8477DE4
        name: ""
        attrs: []
    - index: 14
      size: 103809536
      sectors: 202753
      start: 2048
      end: 204800
      gpt-info:
        guid: 624E86F2-B792-4DD9-981A-3CF034269DC9
        type: efi
        type-guid: C12A7328-F81F-11D2-BA4B-00A0C93EC93B
        name: ""
        attrs: []
```

## License

This software is licensed under the GNU General Public License v2.0 only
(`GPL-2.0`). 

```
Copyright (c) 2022 Ángel Pérez <ap@anpep.co>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, version 2.

This program is distributed in the hope that it will be useful, but
WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
```