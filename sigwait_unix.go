// See LICENSE.txt for licensing information.
// +build !windows

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// sigwait processes signals such as a CTRL-C hit.
func sigwait() {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	s := <-sig
	if s == syscall.SIGINT {
		fmt.Println()
	}
	webuiLog.Printf("Signal '%s' received, stopping", s)
	return
}
