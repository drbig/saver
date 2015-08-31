# saver

Saver is a self-contained cross-platform save-scumming/backup utility.

Features:

- Handle single file and directory saves
- Directories are backed up as zip files
- Option to annotate saves
- No dependencies, single binary
- Tested on Linux and Windows (i386/amd64)

## Showcase

```bash
$ ./saver
Usage: ./saver (options...) <command>
saver v0.5 by Piotr S. Staszewski, see LICENSE.txt

Options:
  -c="saver.json": path to config file
  -v=false: be very verbose

Commands:
  [list]                       - list games
  <name> [add] <path>          - add new game
  <name> [kill]                - delete game and all saves
  <name> [l]ist                - list saves
  <name> [b]ackup (note)       - backup current save
  <name> [r]estore <id|ref>    - restore given save
  <name> [del]ete <id|from-to> - delete given save(s)

Where:
  name     - arbitrary name used to identify a game/character/world etc.
  path     - absolute path to save file/directory
  note     - optional quoted note, e.g. "haven't died yet"
  id       - particular save id
  from-to  - inclusive range of save ids
  ref      - nonpositive offset from the latest save, e.g. -1 is the save before the latest
```

- - -

```bash
$ ./saver list
  # Name                             Len                 Last mod
  1 dcss-mifi                        15       2014-11-23 20:00:58
  2 nh-pri                           1        2014-11-23 14:00:15
  3 cata-stateline                   1        2014-11-23 14:40:04
```

- - -

```bash
$ ./saver dcss-mifi l
Name                             Len                 Last mod
dcss-mifi                        15       2014-11-23 20:00:58

 ID                 Last mod Note
  1      2014-11-23 11:56:57 
  2      2014-11-23 12:16:38 
  3      2014-11-23 12:32:56 
  4      2014-11-23 12:45:36 
  5      2014-11-23 13:05:57 
  6      2014-11-23 13:10:58 
  7      2014-11-23 13:25:28 
  8      2014-11-23 13:39:21 
  9      2014-11-23 13:43:58 
 10      2014-11-23 13:46:32 
 11      2014-11-23 16:21:15 
 12      2014-11-23 16:23:14 Lair:7 clean
 13      2014-11-23 16:41:18 
 14      2014-11-23 16:46:04 Lair:8 clear
 15      2014-11-23 20:00:57 
```

- - -

```bash
$ ./saver dcss-mifi r -1
Restoring save directory from 2014-11-23 16:46:04 ...
Restored save from 2014-11-23 16:46:04
```

## Binaries

Binaries are back! This time for more platform-arch combinations thanks to the awesomeness of Go 1.5. The binaries are distributed outside of the repo, but md5 checksums are kept here for some notion of "authenticity".

Building on your platform is as simple as `go build`, and cross-compiling is as simple as `make` (assuming you have Go 1.5 or newer).

## Usage notes

Copy the binary wherever you want and run it from command line (for Windows folks that's via `cmd.exe`). By default the directory where the binary resides will also be the directory where the config/db file is saved and where the game directories will be made.

Note that restoring a save won't stash the current save, i.e. it will overwrite it without any prompts.

## Games

Some great games where `saver` may be useful:

- [NetHack](http://www.nethack.org/)
- [Dungeon Crawl Stone Soup](http://crawl.develz.org/wordpress/)
- [Cataclysm: Dark Days Ahead](http://en.cataclysmdda.com/)
- [Dwarf Fortress](http://www.bay12games.com/dwarves/)

## Contributing

Follow the usual GitHub development model:

1. Clone the repository
2. Make your changes on a separate branch
3. Make sure you run `gofmt` and `go test` before committing
4. Make a pull request

See licensing for legalese.

## Licensing

Standard two-clause BSD license, see LICENSE.txt for details.

Any contributions will be licensed under the same conditions.

Copyright (c) 2014 - 2015 Piotr S. Staszewski
