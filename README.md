# saver [![Build Status](https://travis-ci.org/drbig/saver.svg?branch=master)](https://travis-ci.org/drbig/saver)

Saver is a self-contained cross-platform save-scumming/backup utility.

**Now at version 0.9.3** I'll be updating the examples here only on user-significant changes.

Features:

- Handle single file and directory saves
- Directories are backed up as zip files
- Option to annotate saves
- No dependencies, single binary
- Tested on Linux and Windows (i386/amd64)
- Keeps track of space usage (and you will use space)
- Has a spinner thingy so you can see it ain't dead
- Minor bugfixes and improvements since the version from 2014
- Can also do MD5 checksums of backups
- Starting to add shorter/scriptable output modes

## Showcase

```bash
$ ./saver
Usage: ./saver (options...) <command>
saver v0.9.3 by Piotr S. Staszewski, see LICENSE.txt
binary build by drbig@swordfish on Sat 6 Oct 18:49:30 CEST 2018

Options:
  -c string
        path to config file (default "saver.json")
  -s    be short, be scriptful
  -v    be very verbose

Commands:
  [list]                       - list games
  <name> [add] <path>          - add new game
  <name> [b]ackup (note)       - backup current save
  <name> [l]ist                - list saves
  <name> [r]estore <id|ref>    - restore given save
  <name> [del]ete <id|from-to> - delete given save(s)
  <name> [i]nfo                - game info, mostly paths to stuff
  <name> check[sum]s           - checksum (MD5) all backups
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
$ ./saver list
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
$ ./saver cata-gruver info
 Game files in: /home/drbig/Projects/cdda-dev/save/Gruver
Saved stuff in: /home/drbig/Projects/go/src/github.com/drbig/saver/cata-gruver
Latest save at: /home/drbig/Projects/go/src/github.com/drbig/saver/cata-gruver/2014-11-29_122018
```

- - -

```bash
$ ./saver cata-gruver l
Name                             # Backups              Last backup     Size
cata-gruver                      1              2014-11-29 12:37:39    25.6M

 ID              Last backup Note
  1      2014-11-29 12:20:18 before sleep, for debug

```

Or if you want to see how much each save takes:

```bash
$ ./saver -v cata-gruver l
loading config from /home/drbig/Projects/go/src/github.com/drbig/saver/saver.json
Name                             # Backups              Last backup     Size
cata-gruver                      1              2014-11-29 12:37:39    25.6M

 ID              Last backup     Size Note
  1      2014-11-29 12:20:18    25.6M before sleep, for debug

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

Copyright (c) 2014 - 2019 Piotr S. Staszewski
