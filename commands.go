package main

import "errors"

type command struct {
	Name string
	Args []string
}

type commands struct {
	table map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.table[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.table[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	return f(s, cmd)
}
