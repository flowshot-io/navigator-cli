package get

import (
	"fmt"

	"github.com/flowshot-io/navigator-cli/internal/cli/display"
	"github.com/flowshot-io/navigator-client-go/commandservice/v1"
	"github.com/spf13/cobra"
)

func (c *Command) NewAssetCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "asset",
		Short: "Get asset",
		Long:  `Get asset`,
		RunE: func(cmd *cobra.Command, args []string) error {
			request := &commandservice.GetAssetRequest{
				AssetID: c.id,
			}

			client, err := c.driver.clientFactory.ResourceClient()
			if err != nil {
				return err
			}

			resp, err := client.GetAsset(cmd.Context(), request)
			if err != nil {
				return fmt.Errorf("unable to get asset: %w", err)
			}

			cmd.Println(display.Assets(resp.Asset))

			return nil
		},
	}

	return cc
}
