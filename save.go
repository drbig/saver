// See LICENSE.txt for licensing information.

package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"time"
)

import (
	"code.cloudfoundry.org/bytefmt"
)

type Save struct {
	Stamp time.Time // save datetime stamp, also used for file name
	Path  string    // absolute path to the save
	Note  string    // optional user note
	Size  uint64    `json:Size,omitempty` // save file size, MINVER:1
}

func (s *Save) PrintHeader() {
	if flagVerbose {
		fmt.Printf("%24s %8s %s\n", "Last backup", "Size", "Note")
	} else {
		fmt.Printf("%24s %s\n", "Last backup", "Note")
	}
}

func (s *Save) Print() {
	if flagVerbose {
		fmt.Printf("%24s %8s %s\n", s.Stamp.Format(timeFmt), bytefmt.ByteSize(s.Size), s.Note)
	} else {
		fmt.Printf("%24s %s\n", s.Stamp.Format(timeFmt), s.Note)
	}
}

func (s *Save) Checksum() (string, error) {
	i, err := os.Open(s.Path)
	if err != nil {
		return "", err
	}
	defer i.Close()

	h := md5.New()
	if _, err := io.Copy(h, i); err != nil {
		return "", nil
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
