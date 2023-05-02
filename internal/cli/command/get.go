package command

import (
	"fmt"

	"github.com/flowshot-io/navigator-client-go/navigatorservice/v1"
	"github.com/spf13/cobra"
)

func (c *Command) NewGetCommand() *cobra.Command {
	var ID string

	cc := &cobra.Command{
		Use:   "get",
		Short: "Get asset",
		Long:  `Get asset`,
		RunE: func(cmd *cobra.Command, args []string) error {
			request := &navigatorservice.GetAssetRequest{
				Id: ID,
			}

			client, err := c.driver.clientFactory.NavigatorClient(cmd)
			if err != nil {
				return fmt.Errorf("unable to create navigator client: %w", err)
			}

			resp, err := client.GetAsset(cmd.Context(), request)
			if err != nil {
				return fmt.Errorf("unable to get asset: %w", err)
			}

			assets := []*navigatorservice.Asset{
				resp.Asset,
			}

			cmd.Println(renderAssets(assets))

			return nil
		},
	}

	cc.Flags().StringVarP(&ID, "id", "i", "", "Asset ID")
	cc.MarkFlagRequired("id")

	return cc
}
