package command

import (
	"github.com/flowshot-io/navigator-client-go/navigatorservice/v1"
	"github.com/spf13/cobra"
)

func (c *Command) NewDeleteCommand() *cobra.Command {
	var ID string

	cc := &cobra.Command{
		Use:   "delete",
		Short: "Delete asset",
		Long:  `Delete asset`,
		Run: func(cmd *cobra.Command, args []string) {
			request := &navigatorservice.DeleteAssetRequest{
				Id: ID,
			}

			client, err := c.driver.clientFactory.NavigatorClient(cmd)
			if err != nil {
				cmd.PrintErrln("Unable to create navigator client", err)
				return
			}

			resp, err := client.DeleteAsset(cmd.Context(), request)
			if err != nil {
				cmd.PrintErrln("Unable to delete asset", err)
				return
			}

			cmd.Println("Scheduled deletion of asset: ", ID, "See", resp.Id)
		},
	}

	cc.Flags().StringVarP(&ID, "id", "i", "", "Asset ID to delete")

	cc.MarkFlagRequired("id")

	return cc
}
