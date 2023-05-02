package command

import (
	"fmt"

	"github.com/flowshot-io/navigator-client-go/navigatorservice/v1"
	"github.com/spf13/cobra"
)

func (c *Command) NewUpdateCommand() *cobra.Command {
	var ID string
	var name string

	cc := &cobra.Command{
		Use:   "update",
		Short: "Update asset",
		Long:  `Update asset`,
		RunE: func(cmd *cobra.Command, args []string) error {
			request := &navigatorservice.UpdateAssetRequest{
				Id:   ID,
				Name: name,
			}

			client, err := c.driver.clientFactory.NavigatorClient(cmd)
			if err != nil {
				return fmt.Errorf("unable to create navigator client: %w", err)
			}

			resp, err := client.UpdateAsset(cmd.Context(), request)
			if err != nil {
				return fmt.Errorf("unable to update asset: %w", err)
			}

			cmd.Println("Updated asset: ", resp.Asset.Id)

			return nil
		},
	}

	cc.Flags().StringVarP(&ID, "id", "i", "", "Asset ID to delete")
	cc.Flags().StringVarP(&name, "name", "n", "", "New name for asset")

	cc.MarkFlagRequired("id")

	return cc
}
