package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Root  string  // absolute path for game directories
	Games []*Game // slice of games
}

func (c *Config) PrintWhole() {
	if len(c.Games) < 1 {
		fmt.Println("No games defined")
		return
	}
	fmt.Printf("%3s ", "#")
	c.Games[0].PrintHeader()
	for i, g := range c.Games {
		fmt.Printf("%3d ", i+1)
		g.Print()
	}
	fmt.Println()
}

func (c *Config) GetGame(name string) *Game {
	for _, g := range c.Games {
		if g.Name == name {
			return g
		}
	}
	return nil
}

func (c *Config) AddGame(name, path string) (gm *Game, err error) {
	if g := c.GetGame(name); g != nil {
		return nil, fmt.Errorf(`Game "%s" already exist`, name)
	}
	r := filepath.Join(c.Root, name)
	if _, err = os.Stat(r); os.IsExist(err) {
		return
	}
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return
	}
	if err = os.MkdirAll(r, 0777); err != nil {
		return
	}
	gm = &Game{
		Name:  name,
		Path:  path,
		Root:  r,
		Stamp: time.Now(),
		Saves: make([]*Save, 0),
	}
	c.Games = append(c.Games, gm)
	return gm, nil
}

func loadConfig(path string) (*Config, error) {
	if flagVerbose {
		fmt.Fprintln(os.Stderr, "loading config from", path)
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	d, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var c *Config
	if err := json.Unmarshal(d, &c); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) Save(path string) error {
	if flagVerbose {
		fmt.Fprintln(os.Stderr, "saving config to", path)
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	j := json.NewEncoder(f)
	return j.Encode(c)
}
