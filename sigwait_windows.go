// See LICENSE.txt for licensing information.
// +build windows

package main

import (
	"fmt"
	"os"
	"os/signal"
)

// sigwait processes signals such as a CTRL-C hit.
func sigwait() {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	s := <-sig
	if s == os.Interrupt {
		fmt.Println()
	}
	webuiLog.Printf("Signal '%s' received, stopping", s)
	return
}
