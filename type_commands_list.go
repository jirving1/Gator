package main

import "fmt"

type Commands struct {
	commandList map[string]func(*State, Command) error
}

func (c *Commands) run(s *State, cmd Command) error {
	handler, ok := c.commandList[cmd.name]
	if !ok {
		return fmt.Errorf("command does not exist")
	}
	err := handler(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *Commands) register(name string, f func(*State, Command) error) {
	if name == "" {
		fmt.Println("no name given to command")
	}
	c.commandList[name] = f
}
