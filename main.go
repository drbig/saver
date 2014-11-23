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
	VERSION = `0.1`
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
		fmt.Fprintln(os.Stderr, "  [a]dd <name> <path>                 - add new game")
		fmt.Fprintln(os.Stderr, "  [d]el <name>                        - delete game and all saves")
		fmt.Fprintln(os.Stderr, "  [l]ist                              - list games")
		fmt.Fprintln(os.Stderr, "  [g]ame <name> [l]ist                - list saves")
		fmt.Fprintln(os.Stderr, "  [g]ame <name> [b]backup (note)      - backup current save")
		fmt.Fprintln(os.Stderr, "  [g]ame <name> [r]estore <id>        - restore given save")
		fmt.Fprintln(os.Stderr, "  [g]ame <name> [d]elete <id|from-to> - delete given save(s)")
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

	switch flag.Arg(0)[0] {
	case 'g': // games command
		checkArgs(false, 3)
		g := cfg.GetGame(flag.Arg(1))
		if g == nil {
			fmt.Fprintf(os.Stderr, "Game \"%s\" not found\n", flag.Arg(1))
			os.Exit(1)
		}
		switch flag.Arg(2)[0] {
		case 'l': // game list
			g.PrintWhole()
		case 'b': // game backup
			sv, err := g.Backup()
			dieOnErr("ERROR", err)
			g.Stamp = time.Now()
			if flag.NArg() > 3 {
				sv.Note = flag.Arg(3)
			}
			fmt.Println("Backed up at", sv.Stamp.Format(timeFmt))
			save = true
		case 'r': // game restore
			checkArgs(true, 4)
			if len(g.Saves) < 1 {
				fmt.Fprintf(os.Stderr, "Game \"%s\" has no saves backed up\n", g.Name)
				os.Exit(1)
			}
			i, err := strconv.Atoi(flag.Arg(3))
			dieOnErr(fmt.Sprintf("Malformed index \"%s\"", flag.Arg(3)), err)
			sv, err := g.Restore(i)
			dieOnErr("ERROR", err)
			g.Stamp = time.Now()
			fmt.Println("Restored save from", sv.Stamp.Format(timeFmt))
			save = true
		case 'd': // game delete saves
			checkArgs(true, 4)
			if len(g.Saves) < 1 {
				fmt.Fprintf(os.Stderr, "Game \"%s\" has no saves backed up\n", g.Name)
				os.Exit(1)
			}
			var err error
			var f, t int
			f, err = strconv.Atoi(flag.Arg(3))
			t = f
			if err != nil {
				m := idRange.FindStringSubmatch(flag.Arg(3))
				f, err = strconv.Atoi(m[1])
				dieOnErr(fmt.Sprintf("Malformed index/range \"%s\"", flag.Arg(3)), err)
				t, err = strconv.Atoi(m[2])
				dieOnErr(fmt.Sprintf("Malformed index/range \"%s\"", flag.Arg(3)), err)
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
	case 'a': // add game
		checkArgs(true, 3)
		if g := cfg.GetGame(flag.Arg(1)); g != nil {
			fmt.Fprintf(os.Stderr, "Game \"%s\" already exist\n", flag.Arg(1))
			os.Exit(1)
		}
		p, err := filepath.Abs(flag.Arg(2))
		dieOnErr("Can't resolve path", err)
		gm, err := cfg.AddGame(flag.Arg(1), p)
		dieOnErr("ERROR", err)
		gm.PrintHeader()
		gm.Print()
		save = true
	case 'd': // delete game
		checkArgs(true, 2)
		g := cfg.GetGame(flag.Arg(1))
		if g == nil {
			fmt.Fprintf(os.Stderr, "Game \"%s\" not found\n", flag.Arg(1))
			os.Exit(1)
		}
		err := cfg.DelGame(flag.Arg(1))
		dieOnErr("ERROR", err)
		fmt.Printf("Deleted game \"%s\" and all backed up saves\n", flag.Arg(1))
		save = true
	case 'l': // list games
		cfg.PrintWhole()
	default:
		flag.Usage()
		os.Exit(1)
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
