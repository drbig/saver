// See LICENSE.txt for licensing information.

package main

import (
	"fmt"
	"time"
)

type Save struct {
	Stamp time.Time // save datetime stamp, also used for file name
	Path  string    // absolute path to the save
	Note  string    // optional user note
}

func (s *Save) PrintHeader() {
	fmt.Printf("%24s %s\n", "Last backup", "Note")
}

func (s *Save) Print() {
	fmt.Printf("%24s %s\n", s.Stamp.Format(timeFmt), s.Note)
}
