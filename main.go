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
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	if flagConfig == "" {
		fmt.Fprintln(os.Stderr, "Config file path can't be empty")
		os.Exit(2)
	}
	flagConfig, err := filepath.Abs(flagConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't resolve config path:", err)
		os.Exit(3)
	}

	save := false
	if _, err := os.Stat(flagConfig); err != nil {
		r, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Can't get current path:", err)
			os.Exit(3)
		}
		cfg = &Config{
			Root:  r,
			Games: make([]*Game, 0),
		}
		fmt.Println("Will save fresh config to", flagConfig)
		save = true
	} else {
		c, err := loadConfig(flagConfig)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Config load error:", err)
			os.Exit(3)
		}
		cfg = c
	}

	switch flag.Arg(0)[0] {
	case 'g':
		if flag.NArg() < 3 {
			flag.Usage()
			os.Exit(1)
		}
		g := cfg.GetGame(flag.Arg(1))
		if g == nil {
			fmt.Fprintf(os.Stderr, "Game \"%s\" not found\n", flag.Arg(1))
			os.Exit(1)
		}
		switch flag.Arg(2)[0] {
		case 'l':
			g.PrintWhole()
		case 'b':
			sv, err := g.Backup()
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR:", err)
				os.Exit(3)
			}
			g.Stamp = time.Now()
			if flag.NArg() > 3 {
				sv.Note = flag.Arg(3)
			}
			fmt.Println("Backed up at", sv.Stamp.Format(timeFmt))
			save = true
		case 'r':
			if flag.NArg() != 4 {
				flag.Usage()
				os.Exit(1)
			}
			if len(g.Saves) < 1 {
				fmt.Fprintf(os.Stderr, "Game \"%s\" has no saves backed up\n", g.Name)
				os.Exit(1)
			}
			i, err := strconv.Atoi(flag.Arg(3))
			if err != nil {
				fmt.Fprintln(os.Stderr, "Malormed index \"%s\": %s\n", flag.Arg(3), err)
				os.Exit(1)
			}
			sv, err := g.Restore(i)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR:", err)
				os.Exit(3)
			}
			g.Stamp = time.Now()
			fmt.Println("Restored save from", sv.Stamp.Format(timeFmt))
			save = true
		case 'd':
			if flag.NArg() != 4 {
				flag.Usage()
				os.Exit(1)
			}
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
				if err != nil {
					fmt.Fprintf(os.Stderr, "Malformed index/range \"%s\": %s\n", flag.Arg(3), err)
					os.Exit(3)
				}
				t, err = strconv.Atoi(m[2])
				if err != nil {
					fmt.Fprintf(os.Stderr, "Malformed index/range \"%s\": %s\n", flag.Arg(3), err)
					os.Exit(3)
				}
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
	case 'a':
		if flag.NArg() != 3 {
			flag.Usage()
			os.Exit(1)
		}
		if g := cfg.GetGame(flag.Arg(1)); g != nil {
			fmt.Fprintf(os.Stderr, "Game \"%s\" already exist\n", flag.Arg(1))
			os.Exit(1)
		}
		p, err := filepath.Abs(flag.Arg(2))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Can't resolve path:", err)
			os.Exit(3)
		}
		gm, err := cfg.AddGame(flag.Arg(1), p)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			os.Exit(3)
		}
		gm.PrintHeader()
		gm.Print()
		save = true
	case 'd':
		if flag.NArg() != 2 {
			flag.Usage()
			os.Exit(1)
		}
		g := cfg.GetGame(flag.Arg(1))
		if g == nil {
			fmt.Fprintf(os.Stderr, "Game \"%s\" not found\n", flag.Arg(1))
			os.Exit(1)
		}
		if err := cfg.DelGame(flag.Arg(1)); err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			os.Exit(3)
		}
		fmt.Printf("Deleted game \"%s\" and all backed up saves\n", flag.Arg(1))
		save = true
	case 'l':
		cfg.PrintWhole()
	default:
		flag.Usage()
		os.Exit(1)
	}

	if save {
		if err := cfg.Save(flagConfig); err != nil {
			fmt.Fprintln(os.Stderr, "Can't save config:", err)
			os.Exit(3)
		}
	}
}
