package delete

import (
	"github.com/flowshot-io/navigator-cli/internal/cli/factory"
	"github.com/spf13/cobra"
)

type Driver struct {
	clientFactory factory.ClientFactory
}

type Command struct {
	driver *Driver
	id     string
}

func NewDriver(clientFactory factory.ClientFactory) *Driver {
	return &Driver{
		clientFactory: clientFactory,
	}
}

func (d *Driver) NewDeleteCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "delete",
		Short: "Delete an object.",
		Long:  `Delete an object by ID.`,
	}

	cmd := &Command{
		driver: d,
	}

	cc.AddCommand(
		cmd.NewAssetCommand(),
	)

	cc.PersistentFlags().StringVarP(&cmd.id, "id", "i", "", "The ID of the object to get.")
	cc.MarkFlagRequired("id")

	return cc
}
