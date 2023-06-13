package command

import (
	"fmt"

	"github.com/flowshot-io/navigator-client-go/commandservice/v1"
	"github.com/spf13/cobra"
)

func (c *Command) NewDeleteCommand() *cobra.Command {
	var ID string

	cc := &cobra.Command{
		Use:   "delete",
		Short: "Delete asset",
		Long:  `Delete asset`,
		RunE: func(cmd *cobra.Command, args []string) error {
			request := &commandservice.DeleteAssetRequest{
				AssetID: ID,
			}

			client, err := c.driver.clientFactory.CommandClient()
			if err != nil {
				return err
			}

			resp, err := client.DeleteAsset(cmd.Context(), request)
			if err != nil {
				return fmt.Errorf("unable to delete asset: %w", err)
			}

			cmd.Println("Scheduled deletion of asset: ", ID, " (message) ", resp.Message)

			return nil
		},
	}

	cc.Flags().StringVarP(&ID, "id", "i", "", "Asset ID to delete")

	cc.MarkFlagRequired("id")

	return cc
}
