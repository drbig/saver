// See LICENSE.txt for licensing information.

package main

import (
	"fmt"
	"os"
)

func (c *Config) Migrate() error {
	if c.Version < 1 {
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
		fmt.Println("")
	}
	return nil
}
