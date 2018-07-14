// See LICENSE.txt for licensing information.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

const (
	VERSION = `0.6`
	timeFmt = `2006-01-02 15:04:05`
	fileFmt = `2006-01-02_150405`
)

var (
	flagConfig  string
	flagVerbose bool
	cfg         *Config
	idRange     = regexp.MustCompile(`(\d+)-(\d+)`)
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s (options...) <command>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "saver v%s by Piotr S. Staszewski, see LICENSE.txt\n\n", VERSION)
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nCommands:")
		fmt.Fprintln(os.Stderr, "  [list]                       - list games")
		fmt.Fprintln(os.Stderr, "  <name> [add] <path>          - add new game")
		fmt.Fprintln(os.Stderr, "  <name> [kill]                - delete game and all saves")
		fmt.Fprintln(os.Stderr, "  <name> [l]ist                - list saves")
		fmt.Fprintln(os.Stderr, "  <name> [b]ackup (note)       - backup current save")
		fmt.Fprintln(os.Stderr, "  <name> [r]estore <id|ref>    - restore given save")
		fmt.Fprintln(os.Stderr, "  <name> [del]ete <id|from-to> - delete given save(s)")
		fmt.Fprintln(os.Stderr, "  [migrate]                    - migrate config, if needed")
		fmt.Fprintln(os.Stderr, "\nWhere:")
		fmt.Fprintln(os.Stderr, "  name     - arbitrary name used to identify a game/character/world etc.")
		fmt.Fprintln(os.Stderr, "  path     - absolute path to save file/directory")
		fmt.Fprintln(os.Stderr, "  note     - optional quoted note, e.g. \"haven't died yet\"")
		fmt.Fprintln(os.Stderr, "  id       - particular save id")
		fmt.Fprintln(os.Stderr, "  from-to  - inclusive range of save ids")
		fmt.Fprintln(os.Stderr, "  ref      - nonpositive offset from the latest save, e.g. -1 is the save before the latest")
	}
	flag.StringVar(&flagConfig, "c", "saver.json", "path to config file")
	flag.BoolVar(&flagVerbose, "v", false, "be very verbose")
}

func main() {
	flag.Parse()
	checkArgs(false, 1)

	if flagConfig == "" {
		fmt.Fprintln(os.Stderr, "Config file path can't be empty")
		os.Exit(2)
	}
	flagConfig, err := filepath.Abs(flagConfig)
	dieOnErr("Can't resolve config path", err)

	save := false
	if _, err := os.Stat(flagConfig); err != nil {
		r, err := os.Getwd()
		dieOnErr("Can't get current path", err)
		cfg = &Config{
			Root:  r,
			Games: make([]*Game, 0),
		}
		fmt.Println("Will save fresh config to", flagConfig)
		save = true
	} else {
		c, err := loadConfig(flagConfig)
		dieOnErr("Config load error", err)
		cfg = c
	}

	switch flag.Arg(0) {
	case "list":
		// list games
		cfg.PrintWhole()
	case "migrate":
		cfg.Migrate()
		save = true
	default:
		// per-game commands
		checkArgs(false, 2)
		game := flag.Arg(0)
		if flag.Arg(1) == "add" {
			// add a new game
			checkArgs(true, 3)
			if g := cfg.GetGame(game); g != nil {
				fmt.Fprintf(os.Stderr, "Game \"%s\" already exist\n", game)
				os.Exit(1)
			}
			p, err := filepath.Abs(flag.Arg(2))
			dieOnErr("Can't resolve path", err)
			gm, err := cfg.AddGame(game, p)
			dieOnErr("ERROR", err)
			gm.PrintHeader()
			gm.Print()
			save = true
		} else {
			g := cfg.GetGame(game)
			if g == nil {
				fmt.Fprintf(os.Stderr, "Game \"%s\" not found\n", game)
				os.Exit(1)
			}
			switch flag.Arg(1) {
			case "l", "list":
				// list game saves
				g.PrintWhole()
			case "kill":
				// remove game and all saves
				err := cfg.DelGame(game)
				dieOnErr("ERROR", err)
				fmt.Printf("Deleted game \"%s\" and all backed up saves\n", game)
				save = true
			case "b", "backup":
				// backup current save
				sv, err := g.Backup()
				dieOnErr("ERROR", err)
				g.Stamp = time.Now()
				if flag.NArg() > 2 {
					sv.Note = flag.Arg(2)
				}
				fmt.Printf("%3d ", len(g.Saves))
				sv.Print()
				save = true
			case "r", "restore":
				// restore selected save
				checkArgs(true, 3)
				if len(g.Saves) < 1 {
					fmt.Fprintf(os.Stderr, "Game \"%s\" has no saves backed up\n", g.Name)
					os.Exit(1)
				}
				i, err := strconv.Atoi(flag.Arg(2))
				dieOnErr(fmt.Sprintf("Malformed index \"%s\"", flag.Arg(2)), err)
				sv, err := g.Restore(i)
				dieOnErr("ERROR", err)
				g.Stamp = time.Now()
				fmt.Println("Restored save from", sv.Stamp.Format(timeFmt))
				if sv.Note != "" {
					fmt.Println("Save note:", sv.Note)
				}
				save = true
			case "del", "delete":
				// delete save(s)
				checkArgs(true, 3)
				if len(g.Saves) < 1 {
					fmt.Fprintf(os.Stderr, "Game \"%s\" has no saves backed up\n", g.Name)
					os.Exit(1)
				}
				var err error
				var f, t int
				f, err = strconv.Atoi(flag.Arg(2))
				t = f
				if err != nil {
					m := idRange.FindStringSubmatch(flag.Arg(2))
					f, err = strconv.Atoi(m[1])
					dieOnErr(fmt.Sprintf("Malformed index/range \"%s\"", flag.Arg(2)), err)
					t, err = strconv.Atoi(m[2])
					dieOnErr(fmt.Sprintf("Malformed index/range \"%s\"", flag.Arg(2)), err)
				}
				n, err := g.Delete(f, t)
				if err != nil {
					fmt.Fprintln(os.Stderr, "ERROR:", err)
				}
				fmt.Printf("Deleted %d save(s)\n", n)
				save = true
			default:
				flag.Usage()
				os.Exit(1)
			}
		}
	}

	if save {
		err := cfg.Save(flagConfig)
		dieOnErr("Can't save config", err)
	}
}

func dieOnErr(msg string, err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, msg+":", err)
		os.Exit(3)
	}
}

func checkArgs(exact bool, num int) {
	if (exact && flag.NArg() != num) ||
		(!exact && flag.NArg() < num) {
		flag.Usage()
		os.Exit(1)
	}
}
