# saver

Saver is a self-contained cross-platform save-scumming/backup utility.

Features:

- Handle single file and directory saves
- Directories are backed up as zip files
- Option to annotate saves
- No dependencies, single binary
- Tested on Linux and Windows (i386/amd64)

## Showcase

    $ ./saver
    Usage: ./saver (options...) <command>
    saver v0.2 by Piotr S. Staszewski, see LICENSE.txt
    
    Options:
      -c="saver.json": path to config file
      -v=false: be very verbose
    
    Commands:
      [a]dd <name> <path>                 - add new game
      [d]el <name>                        - delete game and all saves
      [l]ist                              - list games
      [g]ame <name> [l]ist                - list saves
      [g]ame <name> [b]backup (note)      - backup current save
      [g]ame <name> [r]estore <id>        - restore given save
      [g]ame <name> [d]elete <id|from-to> - delete given save(s)

- - -

    $ ./saver l
      # Name                             Len                 Last mod
      1 dcss-mifi                        15       2014-11-23 20:00:58
      2 nh-pri                           1        2014-11-23 14:00:15
      3 cata-stateline                   1        2014-11-23 14:40:04

- - -

    $ ./saver g dcss-mifi l
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

- - -

    $ ./saver g dcss-mifi r -1
    Restoring save directory from 2014-11-23 16:46:04 ...
    Restored save from 2014-11-23 16:46:04

## Builds

Building from source is just one `go build` away, but for those not inclined here are 2014-11-23 builds:

SHA-1                                    | File
-----------------------------------------|---------------------------
4905daad7235eff0ed6e9529352f3f37d0ca9688 | [saver-linux-386](http://insomniac.pl/~drbig/saver/saver-linux-386)
4e0d55129a27ae0cf39d355c1aa5b3f93c635674 | [saver-linux-amd64](http://insomniac.pl/~drbig/saver/saver-linux-amd64)
3bae864974b1e6a7b63d03972a0282d86a58c5ac | [saver-windows-386.exe](http://insomniac.pl/~drbig/saver/saver-windows-386.exe)
6a48f705044225acbbd79f4869a297054308e3ac | [saver-windows-amd64.exe](http://insomniac.pl/~drbig/saver/saver-windows-amd64.exe)

## Usage notes

Copy the binary wherever you want and run it from command line (for Windows folks that's via `cmd.exe`). By default the directory where the binary resides will also be the directory where the config/db file is saved and where the game directories will be made.

The only thing not mentioned in the help is that the game restore sub-command can take either a positive index which corresponds to the save ID as printed by list, or a 0-or-negative index which is relative to the last save in the list (i.e. `g dcss-mifi r 0` will restore the last save for that particular game; likewise -1 will restore the one before the last save).

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

Copyright (c) 2014 Piotr S. Staszewski
