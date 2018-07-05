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
$ ./bin/saver-linux-amd64-0.6
Usage: ./saver (options...) <command>
saver v0.6 by Piotr S. Staszewski, see LICENSE.txt

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
$ ./bin/saver-linux-amd64-0.6 list
  # Name                             # Backups              Last backup
  1 cata-gruver                      1              2014-11-29 12:37:39
  2 cata-marty                       1              2014-12-20 22:56:42
  3 urw-tut                          1              2015-02-24 19:12:21
  4 df-2015                          1              2015-06-15 17:27:48
  5 urw-legacy                       1              2015-04-24 15:23:58
  6 urw-helena                       1              2015-05-04 16:07:48
  7 urw-kamputuuri                   1              2015-07-04 11:10:31
  8 urw-xena                         1              2015-06-27 00:13:12
  9 cata-richlawn                    1              2015-09-05 00:54:17
 10 cata-morrill                     1              2015-09-17 23:09:37
 11 ds                               1              2015-12-09 22:49:24
```

- - -

```bash
$ ./bin/saver-linux-amd64-0.6 cata-gruver l
Name                             # Backups              Last backup
cata-gruver                      1              2014-11-29 12:37:39

 ID              Last backup Note
  1      2014-11-29 12:20:18 before sleep, for debug

```

- - -

```bash
$ ./bin/saver-linux-amd64-0.6 dcss-mifi r -1
Restoring save directory from 2014-11-23 16:46:04 ...
Restored save from 2014-11-23 16:46:04
```

## Binaries

Binaries are back! This time for more platform-arch combinations thanks to the awesomeness of Go 1.5. The binaries are distributed outside of the repo, but md5 checksums are kept here for some notion of "authenticity".

Building on your platform is as simple as `go build`, and cross-compiling is as simple as `make` (assuming you have Go 1.5 or newer).

And now the [binaries link is public](https://insomniac.pl/~drbig/binaries/). Apply caution please!

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

Copyright (c) 2014 - 2018 Piotr S. Staszewski
