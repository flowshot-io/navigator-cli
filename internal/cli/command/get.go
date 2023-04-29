package command

import (
	"github.com/flowshot-io/navigator-client-go/navigatorservice/v1"
	"github.com/spf13/cobra"
)

func (c *Command) NewGetCommand() *cobra.Command {
	var ID string

	cc := &cobra.Command{
		Use:   "get",
		Short: "Get asset",
		Long:  `Get asset`,
		Run: func(cmd *cobra.Command, args []string) {
			request := &navigatorservice.GetAssetRequest{
				Id: ID,
			}

			client, err := c.driver.clientFactory.NavigatorClient(cmd)
			if err != nil {
				cmd.PrintErrln("Unable to create navigator client", err)
				return
			}

			resp, err := client.GetAsset(cmd.Context(), request)
			if err != nil {
				cmd.PrintErrln("Unable to get asset", err)
				return
			}

			assets := []*navigatorservice.Asset{
				resp.Asset,
			}

			cmd.Println(renderAssets(assets))
		},
	}

	cc.Flags().StringVarP(&ID, "id", "i", "", "Asset ID")
	cc.MarkFlagRequired("id")

	return cc
}
