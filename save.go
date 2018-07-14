// See LICENSE.txt for licensing information.

package main

import (
	"fmt"
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
