// See LICENSE.txt for licensing information.

package main

import (
	"fmt"
	"os"
)

const (
	CFG_VER_1 = iota + 1
)

func (c *Config) Migrate() error {
	if c.Version < CFG_VER_1 {
		fmt.Println("Migrating -> 1")
		for _, g := range c.Games {
			fmt.Print("g")
			for _, s := range g.Saves {
				if s.Size > 0 {
					fmt.Print("S")
					continue
				}
				i, err := os.Stat(s.Path)
				if err != nil {
					return err
				}
				fmt.Print("s")
				size := uint64(i.Size())
				s.Size = size
				g.Size += size
			}
		}
		c.Version = CFG_VER_1
		fmt.Println("")
	}
	fmt.Println("All done")
	return nil
}
