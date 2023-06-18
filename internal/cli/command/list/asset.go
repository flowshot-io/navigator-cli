package list

import (
	"fmt"

	"github.com/flowshot-io/navigator-cli/internal/cli/display"
	"github.com/flowshot-io/navigator-client-go/commandservice/v1"
	"github.com/spf13/cobra"
)

func (c *Command) NewAssetCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "asset",
		Short: "List assets.",
		Long:  `List assets.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			request := &commandservice.ListAssetsRequest{}

			client, err := c.driver.clientFactory.ResourceClient()
			if err != nil {
				return err
			}

			resp, err := client.ListAssets(cmd.Context(), request)
			if err != nil {
				return fmt.Errorf("unable to list files: %w", err)
			}

			cmd.Println(display.Assets(resp.Assets...))

			return nil
		},
	}

	return cc
}
