// See LICENSE.txt for licensing information.

package main

import (
	"fmt"
	"os"
	"os/signal"
)

func _sigwait(sigs ...os.Signal) {
	sig := make(chan os.Signal)
	signal.Notify(sig, sigs...)
	s := <-sig
	if s == sigs[0] {
		fmt.Println()
	}
	webuiLog.Printf("Signal '%s' received, stopping", s)
	return
}
