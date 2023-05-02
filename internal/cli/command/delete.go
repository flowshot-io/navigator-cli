package command

import (
	"fmt"

	"github.com/flowshot-io/navigator-client-go/navigatorservice/v1"
	"github.com/spf13/cobra"
)

func (c *Command) NewDeleteCommand() *cobra.Command {
	var ID string

	cc := &cobra.Command{
		Use:   "delete",
		Short: "Delete asset",
		Long:  `Delete asset`,
		RunE: func(cmd *cobra.Command, args []string) error {
			request := &navigatorservice.DeleteAssetRequest{
				Id: ID,
			}

			client, err := c.driver.clientFactory.NavigatorClient(cmd)
			if err != nil {
				return fmt.Errorf("unable to create navigator client: %w", err)
			}

			resp, err := client.DeleteAsset(cmd.Context(), request)
			if err != nil {
				return fmt.Errorf("unable to delete asset: %w", err)
			}

			cmd.Println("Scheduled deletion of asset: ", ID, "See", resp.Id)

			return nil
		},
	}

	cc.Flags().StringVarP(&ID, "id", "i", "", "Asset ID to delete")

	cc.MarkFlagRequired("id")

	return cc
}
