package main

import (
	"archive/zip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	timeFmt = `2006-01-02 15:04:05`
)

var (
	flagConfig  string
	flagVerbose bool
	cfg         *Config
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options...] <command>\n\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nCommands:")
		fmt.Fprintln(os.Stderr, "(a)dd <name> <path>          - add new game")
		fmt.Fprintln(os.Stderr, "(d)el <name>                 - delete game and all saves")
		fmt.Fprintln(os.Stderr, "(l)ist                       - list games")
		fmt.Fprintln(os.Stderr, "(l)ist <name>                - list game saves")
		fmt.Fprintln(os.Stderr, "(g)ame <name> (b)backup      - backup current save")
		fmt.Fprintln(os.Stderr, "(g)ame <name> (r)estore <id> - restore given save")
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
		case 'b':
			sv, err := g.Backup()
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR:", err)
				os.Exit(3)
			}
			g.Stamp = time.Now()
			fmt.Println("Backed up", sv.Stamp.Format(timeFmt))
			save = true
		case 'r':
			if flag.NArg() != 4 {
				flag.Usage()
				os.Exit(1)
			}
			if len(g.Saves) < 1 {
				fmt.Fprintf(os.Stderr, "Game \"%s\" has not saves backed up, yet\n", g.Name)
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
			fmt.Println("Restored save", sv.Stamp.Format(timeFmt))
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
		p, err := filepath.Abs(flag.Arg(2))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Can't resolve path:", err)
			os.Exit(3)
		}
		gm, err := cfg.AddGame(flag.Arg(1), p)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error adding game:", err)
			os.Exit(3)
		}
		gm.PrintHeader()
		gm.Print()
		save = true
	case 'l':
		if flag.NArg() == 1 {
			cfg.PrintWhole()
			break
		}
		g := cfg.GetGame(flag.Arg(1))
		if g == nil {
			fmt.Fprintf(os.Stderr, "Game \"%s\" not found\n", flag.Arg(1))
			os.Exit(1)
		}
		g.PrintWhole()
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

type Save struct {
	Stamp time.Time // save datetime stamp, also used for file name
	Path  string    // absolute path to the save
	Note  string    // optional user note
}

type Game struct {
	Name  string    // user-given mnemonic, also used for making save storage directory
	Path  string    // absolute path to the source save file/directory
	Root  string    // absolute path to the saved saves directory
	Stamp time.Time // last modification datetime stamp
	Saves []*Save   // slice of saves
}

type Config struct {
	Root  string  // absolute path for game directories
	Games []*Game // slice of games
}

func loadConfig(path string) (*Config, error) {
	if flagVerbose {
		fmt.Fprintln(os.Stderr, "loading config from", path)
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	d, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var c *Config
	if err := json.Unmarshal(d, &c); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) Save(path string) error {
	if flagVerbose {
		fmt.Fprintln(os.Stderr, "saving config to", path)
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	j := json.NewEncoder(f)
	return j.Encode(c)
}

func (s *Save) PrintHeader() {
	fmt.Printf("%24s %s\n", "Last mod", "Note")
}

func (s *Save) Print() {
	fmt.Printf("%24s %s\n", s.Stamp.Format(timeFmt), s.Note)
}

func (g *Game) PrintHeader() {
	fmt.Printf("%-32s %-3s %24s\n", "Name", "Len", "Last mod")
}

func (g *Game) Print() {
	fmt.Printf("%-32s %-3d %24s\n", g.Name, len(g.Saves), g.Stamp.Format(timeFmt))
}

func (g *Game) PrintWhole() {
	g.PrintHeader()
	g.Print()
	fmt.Println()
	if len(g.Saves) > 0 {
		fmt.Printf("%3s ", "ID")
		g.Saves[0].PrintHeader()
		for i, s := range g.Saves {
			fmt.Printf("%3d ", i+1)
			s.Print()
		}
		fmt.Println()
	}
}

func (g *Game) Backup() (sv *Save, err error) {
	fi, err := os.Stat(g.Path)
	if err != nil {
		return nil, err
	}
	s := time.Now()
	p := filepath.Join(g.Root, s.Format(timeFmt))
	o, err := os.Create(p)
	if err != nil {
		return nil, err
	}
	defer func() {
		o.Close()
		if err != nil {
			if lerr := os.Remove(p); lerr != nil {
				fmt.Fprintln(os.Stderr, "Error removing partial backup file:", lerr)
			}
		}
	}()
	// single file
	if fi.Mode().IsRegular() {
		i, err := os.Open(g.Path)
		if err != nil {
			return nil, err
		}
		defer i.Close()
		fmt.Println("Backing up", i.Name(), "...")
		if flagVerbose {
			fmt.Fprintf(os.Stderr, "copying %s to %s\n", g.Path, p)
		}
		n, err := io.Copy(o, i)
		if err != nil {
			if flagVerbose {
				fmt.Fprintf(os.Stderr, "error, copied %d\n", n)
			}
			return nil, err
		}
		if flagVerbose {
			fmt.Fprintf(os.Stderr, "ok, copied %d\n", n)
		}
	} else {
		// directory
		fmt.Println("Backing up save directory...")
		if flagVerbose {
			fmt.Fprintf(os.Stderr, "zipping %s to %s\n", g.Path, p)
		}
		z := zip.NewWriter(o)
		err = filepath.Walk(g.Path, func(path string, info os.FileInfo, ierr error) (err error) {
			if ierr != nil {
				return ierr
			}
			if info.IsDir() {
				if flagVerbose {
					fmt.Fprintf(os.Stderr, "skipping %s\n", path)
				}
				return nil
			}
			rp, err := filepath.Rel(g.Path, path)
			if err != nil {
				return err
			}
			i, err := os.Open(path)
			if err != nil {
				return err
			}
			defer i.Close()
			o, err := z.Create(rp)
			if err != nil {
				return err
			}
			if flagVerbose {
				fmt.Fprintf(os.Stderr, "compressing %s\n", rp)
			}
			_, err = io.Copy(o, i)
			return err
		})
		if err != nil {
			return nil, err
		}
		if err = z.Close(); err != nil {
			return nil, err
		}
	}
	sv = &Save{
		Stamp: s,
		Path:  p,
	}
	g.Saves = append(g.Saves, sv)
	return sv, nil
}

func (g *Game) Restore(index int) (sv *Save, err error) {
	var i int

	if index > 0 {
		i = index - 1
		if i > len(g.Saves) {
			return nil, fmt.Errorf("Save ID %d out of range (1 ~ %d)", index, len(g.Saves))
		}
	} else {
		i = len(g.Saves) - 1 + index
		if i < 0 {
			return nil, fmt.Errorf("Save offset %d out of range (%d ~ 0)", index, -len(g.Saves)+1)
		}
	}
	sv = g.Saves[i]
	fi, err := os.Stat(g.Path)
	if err != nil {
		return nil, err
	}
	// single file
	if fi.Mode().IsRegular() {
		i, err := os.Open(sv.Path)
		if err != nil {
			return nil, err
		}
		defer i.Close()
		o, err := os.Create(g.Path)
		if err != nil {
			return nil, err
		}
		defer o.Close()
		fmt.Println("Restoring save from", sv.Stamp.Format(timeFmt), "...")
		if flagVerbose {
			fmt.Fprintf(os.Stderr, "copying %s to %s\n", sv.Path, g.Path)
		}
		n, err := io.Copy(o, i)
		if flagVerbose {
			if err != nil {
				fmt.Fprintf(os.Stderr, "error, copied %d\n", n)
			} else {
				fmt.Fprintf(os.Stderr, "ok, copied %d\n", n)
			}
		}
		return sv, nil
	}
	// directory
	fmt.Println("Restoring save directory from", sv.Stamp.Format(timeFmt), "...")
	if flagVerbose {
		fmt.Fprintf(os.Stderr, "removing all from %s\n", g.Path)
	}
	err = os.RemoveAll(g.Path)
	if err != nil {
		return nil, err
	}
	if flagVerbose {
		fmt.Fprintf(os.Stderr, "unzipping %s to %s\n", sv.Path, g.Path)
	}
	z, err := zip.OpenReader(sv.Path)
	if err != nil {
		return nil, err
	}
	defer z.Close()
	for _, f := range z.File {
		i, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer i.Close()
		tp := filepath.Join(g.Path, f.Name)
		dp := filepath.Dir(tp)
		err = os.MkdirAll(dp, 0777)
		if err != nil {
			return nil, err
		}
		o, err := os.Create(tp)
		if err != nil {
			return nil, err
		}
		defer o.Close()
		fmt.Fprintf(os.Stderr, "decompressing %s\n", f.Name)
		n, err := io.Copy(o, i)
		if err != nil {
			if flagVerbose {
				fmt.Fprintf(os.Stderr, "error, copied %d\n", n)
			}
			return nil, err
		}
		if flagVerbose {
			fmt.Fprintf(os.Stderr, "ok, copied %d\n", n)
		}
	}
	return sv, nil
}

func (c *Config) PrintWhole() {
	if len(c.Games) < 1 {
		fmt.Println("No games defined")
		return
	}
	fmt.Printf("%3s ", "#")
	c.Games[0].PrintHeader()
	for i, g := range c.Games {
		fmt.Printf("%3d ", i+1)
		g.Print()
	}
	fmt.Println()
}

func (c *Config) GetGame(name string) *Game {
	for _, g := range c.Games {
		if g.Name == name {
			return g
		}
	}
	return nil
}

func (c *Config) AddGame(name, path string) (gm *Game, err error) {
	if g := c.GetGame(name); g != nil {
		return nil, fmt.Errorf(`Game "%s" already exist`, name)
	}
	r := filepath.Join(c.Root, name)
	if _, err = os.Stat(r); os.IsExist(err) {
		return
	}
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return
	}
	if err = os.MkdirAll(r, 0777); err != nil {
		return
	}
	gm = &Game{
		Name:  name,
		Path:  path,
		Root:  r,
		Stamp: time.Now(),
		Saves: make([]*Save, 0),
	}
	c.Games = append(c.Games, gm)
	return gm, nil
}
