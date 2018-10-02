// See LICENSE.txt for licensing information.

package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

import (
	"code.cloudfoundry.org/bytefmt"
)

type Game struct {
	Name  string    // user-given mnemonic, also used for making save storage directory
	Path  string    // absolute path to the source save file/directory
	Root  string    // absolute path to the saved saves directory
	Stamp time.Time // last modification datetime stamp
	Saves []*Save   // slice of saves
	Size  uint64    `json:Size,omitempty` // total size of all saves, MINVER:1
}

func (g *Game) PrintHeader() {
	fmt.Printf("%-32s %-9s %24s %8s\n", "Name", "# Backups", "Last backup", "Size")
}

func (g *Game) Print() {
	fmt.Printf("%-32s %-9d %24s %8s\n", g.Name, len(g.Saves), g.Stamp.Format(timeFmt), bytefmt.ByteSize(g.Size))
}

func (g *Game) PrintInfo() {
	lsp := g.Saves[len(g.Saves)-1].Path
	if flagShort {
		fmt.Println(lsp)
		return
	}
	fmt.Printf(` Game files in: %s
Saved stuff in: %s
Latest save at: %s
`, g.Path, g.Root, lsp)
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

func (g *Game) Delete(from, to int) (n int, err error) {
	if from < 1 {
		return n, fmt.Errorf("Index from %d out of range", from)
	}
	if to > len(g.Saves) {
		return n, fmt.Errorf("Index to %d out of range", to)
	}
	spinner.Msg("Deleting saves...")
	from--
	i := from
	for ; i < to; i++ {
		s := g.Saves[i]
		if flagVerbose {
			spinner.Msg(fmt.Sprintf("removing save %d from %s", i+1, s.Stamp.Format(timeFmt)))
		}
		spinner.Tick()
		err = os.Remove(s.Path)
		if err != nil {
			break
		}
		g.Size -= s.Size
		n++
	}
	spinner.Finish()
	copy(g.Saves[from:], g.Saves[i:])
	for k, n := len(g.Saves)-i+from, len(g.Saves); k < n; k++ {
		g.Saves[k] = nil
	}
	g.Saves = g.Saves[:len(g.Saves)-i+from]
	return n, err
}

// Backup copies a save file, or zips a save directory.
// Note that currently proper file closing is depending on total
// exit of the program on any file operation failure.
func (g *Game) Backup() (sv *Save, err error) {
	fi, err := os.Stat(g.Path)
	if err != nil {
		return nil, err
	}
	s := time.Now()
	p := filepath.Join(g.Root, s.Format(fileFmt))
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
		spinner.Msg("Backing up save directory...")
		if flagVerbose {
			spinner.Msg(fmt.Sprintf("zipping %s to %s", g.Path, p))
		}
		z := zip.NewWriter(o)
		err = filepath.Walk(g.Path, func(path string, info os.FileInfo, ierr error) (err error) {
			if ierr != nil {
				return ierr
			}
			if info.IsDir() {
				if flagVerbose {
					spinner.Msg(fmt.Sprintf("skipping %s", path))
				}
				return nil
			}
			rp, err := filepath.Rel(g.Path, path)
			if err != nil {
				return err
			}
			rp = filepath.ToSlash(rp)
			i, err := os.Open(path)
			if err != nil {
				return err
			}
			o, err := z.Create(rp)
			if err != nil {
				return err
			}
			if flagVerbose {
				spinner.Msg(fmt.Sprintf("compressing %s", rp))
			}
			spinner.Tick()
			_, err = io.Copy(o, i)
			i.Close()
			return err
		})
		spinner.Finish()
		if err != nil {
			return nil, err
		}
		if err = z.Close(); err != nil {
			return nil, err
		}
	}
	si, err := o.Stat()
	if err != nil {
		return nil, err
	}
	size := uint64(si.Size())
	sv = &Save{
		Stamp: s,
		Path:  p,
		Size:  size,
	}
	g.Saves = append(g.Saves, sv)
	g.Size += size
	return sv, nil
}

// Restore removes the current save file or directory, and repopulates
// it with the data from given backup save file.
// Note that currently proper file closing is depending on total
// exit of the program on any file operation failure.
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
		fmt.Println("Restoring save from", sv.Path, "...")
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
	spinner.Msg(fmt.Sprint("Restoring save directory from ", sv.Stamp.Format(timeFmt), "..."))
	if flagVerbose {
		spinner.Msg(fmt.Sprintf("removing all from %s", g.Path))
	}
	err = os.RemoveAll(g.Path)
	if err != nil {
		return nil, err
	}
	if flagVerbose {
		spinner.Msg(fmt.Sprintf("unzipping %s to %s", sv.Path, g.Path))
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
		tp := filepath.Join(g.Path, filepath.FromSlash(f.Name))
		dp := filepath.Dir(tp)
		err = os.MkdirAll(dp, 0777)
		if err != nil {
			return nil, err
		}
		o, err := os.Create(tp)
		if err != nil {
			return nil, err
		}
		if flagVerbose {
			spinner.Msg(fmt.Sprintf("decompressing %s", f.Name))
		}
		spinner.Tick()
		n, err := io.Copy(o, i)
		if err != nil {
			if flagVerbose {
				fmt.Fprintf(os.Stderr, "error, copied %d\n", n)
			}
			return nil, err
		}
		o.Close()
		i.Close()
		if flagVerbose {
			spinner.Msg(fmt.Sprintf("ok, copied %d", n))
		}
	}
	spinner.Finish()
	return sv, nil
}

func (g *Game) ChecksumAll() {
	for i, s := range g.Saves {
		if flagShort {
			fmt.Print(s.Path, " ")
		} else {
			fmt.Print(i+1, " ", s.Stamp.Format(timeFmt), " ")
		}

		sum, err := s.Checksum()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(sum)
		}
	}
}
