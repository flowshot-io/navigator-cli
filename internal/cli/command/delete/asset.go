package delete

import (
	"fmt"

	"github.com/flowshot-io/navigator-client-go/commandservice/v1"
	"github.com/spf13/cobra"
)

func (c *Command) NewAssetCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "asset",
		Short: "Delete asset",
		Long:  `Delete asset`,
		RunE: func(cmd *cobra.Command, args []string) error {
			request := &commandservice.DeleteAssetRequest{
				AssetID: c.id,
			}

			client, err := c.driver.clientFactory.CommandClient()
			if err != nil {
				return err
			}

			_, err = client.DeleteAsset(cmd.Context(), request)
			if err != nil {
				return fmt.Errorf("unable to delete asset: %w", err)
			}

			cmd.Println("Scheduled deletion of asset: ", c.id)

			return nil
		},
	}

	return cc
}
