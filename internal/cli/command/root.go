package command

import (
	"github.com/flowshot-io/navigator-cli/internal/cli/command/create"
	"github.com/flowshot-io/navigator-cli/internal/cli/command/delete"
	"github.com/flowshot-io/navigator-cli/internal/cli/command/get"
	"github.com/flowshot-io/navigator-cli/internal/cli/command/list"
	"github.com/flowshot-io/navigator-cli/internal/cli/factory"
	"github.com/spf13/cobra"
)

type Driver struct {
	clientFactory factory.ClientFactory
}

type Command struct {
	driver *Driver
}

func NewDriver() *Driver {
	clientFactory := factory.NewClientFactory()

	return &Driver{
		clientFactory: clientFactory,
	}
}

func NewCommand(d *Driver) *cobra.Command {
	c := &cobra.Command{
		Use:           "navigator",
		Short:         "A command-line tool for Navigator.",
		Long:          `A command-line tool for managing stored assets via Navigator.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	c.SetVersionTemplate("{{.Version}}\n")

	cmd := &Command{
		driver: d,
	}

	c.AddCommand(
		cmd.NewCreateCommand(),
		cmd.NewGetCommand(),
		cmd.NewListCommand(),
		cmd.NewDeleteCommand(),
	)

	return c
}

func (c *Command) NewCreateCommand() *cobra.Command {
	driver := create.NewDriver(c.driver.clientFactory)
	return driver.NewCreateCommand()
}

func (c *Command) NewGetCommand() *cobra.Command {
	driver := get.NewDriver(c.driver.clientFactory)
	return driver.NewGetCommand()
}

func (c *Command) NewListCommand() *cobra.Command {
	driver := list.NewDriver(c.driver.clientFactory)
	return driver.NewListCommand()
}

func (c *Command) NewDeleteCommand() *cobra.Command {
	driver := delete.NewDriver(c.driver.clientFactory)
	return driver.NewDeleteCommand()
}
