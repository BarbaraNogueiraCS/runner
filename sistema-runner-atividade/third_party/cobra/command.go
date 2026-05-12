// Package cobra provides a very small compatibility subset used only so this
// academic project can be compiled in offline environments. In production,
// remove the replace directive in go.mod and use the official Cobra module.
package cobra

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Command struct {
	Use          string
	Short        string
	SilenceUsage bool
	Run          func(cmd *Command, args []string)
	RunE         func(cmd *Command, args []string) error
	commands     []*Command
	flags        *flag.FlagSet
	out          io.Writer
	err          io.Writer
}

func (c *Command) AddCommand(commands ...*Command) {
	c.commands = append(c.commands, commands...)
}

func (c *Command) Flags() *flag.FlagSet {
	if c.flags == nil {
		c.flags = flag.NewFlagSet(c.name(), flag.ContinueOnError)
		c.flags.SetOutput(c.ErrOrStderr())
	}
	return c.flags
}

func (c *Command) Execute() error {
	return c.execute(os.Args[1:])
}

func (c *Command) execute(args []string) error {
	if len(args) == 0 || args[0] == "--help" || args[0] == "help" {
		c.printHelp()
		return nil
	}
	for _, child := range c.commands {
		if args[0] == child.name() {
			return child.run(args[1:])
		}
	}
	return c.run(args)
}

func (c *Command) run(args []string) error {
	if c.flags != nil {
		if err := c.flags.Parse(args); err != nil {
			return err
		}
		args = c.flags.Args()
	}
	if c.RunE != nil {
		return c.RunE(c, args)
	}
	if c.Run != nil {
		c.Run(c, args)
		return nil
	}
	c.printHelp()
	return nil
}

func (c *Command) OutOrStdout() io.Writer {
	if c.out != nil {
		return c.out
	}
	return os.Stdout
}

func (c *Command) ErrOrStderr() io.Writer {
	if c.err != nil {
		return c.err
	}
	return os.Stderr
}

func (c *Command) name() string {
	parts := strings.Fields(c.Use)
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

func (c *Command) printHelp() {
	fmt.Fprintf(c.OutOrStdout(), "%s\n\n", c.Short)
	fmt.Fprintf(c.OutOrStdout(), "Uso: %s [comando] [opções]\n", c.name())
	if len(c.commands) > 0 {
		fmt.Fprintln(c.OutOrStdout(), "\nComandos disponíveis:")
		for _, child := range c.commands {
			fmt.Fprintf(c.OutOrStdout(), "  %-12s %s\n", child.name(), child.Short)
		}
	}
	if c.flags != nil {
		fmt.Fprintln(c.OutOrStdout(), "\nOpções:")
		c.flags.SetOutput(c.OutOrStdout())
		c.flags.PrintDefaults()
	}
}
