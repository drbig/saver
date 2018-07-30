// See LICENSE.txt for licensing information.
// +build !windows

package main

import "syscall"

// sigwait processes signals such as a CTRL-C hit.
func sigwait() {
	_sigwait(syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
}
