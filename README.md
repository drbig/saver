# saver

Saver is a self-contained cross-platform save-scumming/backup utility.

**Now at version 0.9.1** I'll be updating the examples here only on user-significant changes (so the 0.9.1 spinner thing doesn't count).

Features:

- Handle single file and directory saves
- Directories are backed up as zip files
- Option to annotate saves
- No dependencies, single binary
- Tested on Linux and Windows (i386/amd64)
- Keeps track of space usage (and you will use space)
- Has a spinner thingy so you can see it ain't dead
- Minor bugfixes and improvements since the version from 2014

## Showcase

```bash
$ ./bin/saver-linux-amd64-0.8
Usage: ./bin/saver-linux-amd64-0.8 (options...) <command>
saver v0.8 by Piotr S. Staszewski, see LICENSE.txt
binary build by drbig@swordfish on Sun 15 Jul 20:15:37 CEST 2018

Options:
  -c string
        path to config file (default "saver.json")
  -v    be very verbose

Commands:
  [list]                       - list games
  <name> [add] <path>          - add new game
  <name> [b]ackup (note)       - backup current save
  <name> [l]ist                - list saves
  <name> [r]estore <id|ref>    - restore given save
  <name> [del]ete <id|from-to> - delete given save(s)
  <name> [kill]                - delete game and all saves
  [migrate]                    - migrate config, if needed

Where:
  name     - arbitrary name used to identify a game/character/world etc.
  path     - absolute path to save file/directory
  note     - optional quoted note, e.g. "haven't died yet"
  id       - particular save id
  from-to  - inclusive range of save ids
  ref      - non-positive offset from the latest save, e.g. -1 is the save before the latest
```

- - -

```bash
$ ./bin/saver-linux-amd64-0.8 list
  # Name                             # Backups              Last backup     Size
  1 cata-gruver                      1              2014-11-29 12:37:39    25.6M
  2 cata-marty                       1              2014-12-20 22:56:42     2.2M
  3 urw-tut                          1              2015-02-24 19:12:21     9.9M
  4 df-2015                          1              2018-07-14 22:05:55     7.2M
  5 urw-legacy                       1              2015-04-24 15:23:58       4M
  6 urw-helena                       1              2015-05-04 16:07:48     6.3M
  7 urw-kamputuuri                   1              2015-07-04 11:10:31     5.2M
  8 urw-xena                         1              2015-06-27 00:13:12     5.6M
  9 cata-richlawn                    1              2015-09-05 00:54:17   851.7K
 10 cata-morrill                     1              2015-09-17 23:09:37     6.2M
 11 ds                               1              2015-12-09 22:49:24   262.2K

```

- - -

```bash
$ ./bin/saver-linux-amd64-0.8 cata-gruver l
Name                             # Backups              Last backup     Size
cata-gruver                      1              2014-11-29 12:37:39    25.6M

 ID              Last backup Note
  1      2014-11-29 12:20:18 before sleep, for debug

```

Or if you want to see how much each save takes:

```bash
$ ./bin/saver-linux-amd64-0.8 -v cata-gruver l
loading config from /home/drbig/Projects/go/src/github.com/drbig/saver/saver.json
Name                             # Backups              Last backup     Size
cata-gruver                      1              2014-11-29 12:37:39    25.6M

 ID              Last backup     Size Note
  1      2014-11-29 12:20:18    25.6M before sleep, for debug

```

- - -

```bash
$ ./bin/saver-linux-amd64-0.8 dcss-mifi r -1
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
