// See LICENSE.txt for licensing information.

package main

import (
	"fmt"
	"strings"
)

var (
	SPINNER_STRINGS = []string{"◢ ", "◣ ", "◤ ", "◥ "}
	SPINNER_LEN     = len(SPINNER_STRINGS)
)

type Spinner struct {
	running  bool   // indiacte we're actually printing
	msg      string // current message
	pos      int    // position in SPINNER_STRINGS
	last_len int    // total length of last status
}

func (s *Spinner) Msg(msg string) {
	s.msg = msg
	s.Tick()
}

func (s *Spinner) Tick() {
	if s.running {
		fmt.Print("\r", strings.Repeat(" ", s.last_len), "\r")
	} else {
		s.running = true
	}
	s.last_len, _ = fmt.Print(SPINNER_STRINGS[s.pos], s.msg)
	s.pos = (s.pos + 1) % SPINNER_LEN
}

func (s *Spinner) Finish() {
	s.Msg("All done")
	fmt.Println()
	s.running = false
	s.msg = ""
	s.pos = 0
	s.last_len = 0
}
