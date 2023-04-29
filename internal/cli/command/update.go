package command

import (
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
		Run: func(cmd *cobra.Command, args []string) {
			request := &navigatorservice.UpdateAssetRequest{
				Id:   ID,
				Name: name,
			}

			client, err := c.driver.clientFactory.NavigatorClient(cmd)
			if err != nil {
				cmd.PrintErrln("Unable to create navigator client", err)
				return
			}

			resp, err := client.UpdateAsset(cmd.Context(), request)
			if err != nil {
				cmd.PrintErrln("Unable to update asset", err)
				return
			}

			cmd.Println("Updated asset: ", resp.Asset.Id)
		},
	}

	cc.Flags().StringVarP(&ID, "id", "i", "", "Asset ID to delete")
	cc.Flags().StringVarP(&name, "name", "n", "", "New name for asset")

	cc.MarkFlagRequired("id")

	return cc
}
