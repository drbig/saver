// See LICENSE.txt for licensing information.
// +build windows

package main

import "os"

// sigwait processes signals such as a CTRL-C hit.
func sigwait() {
	_sigwait(os.Interrupt, os.Kill)
}
