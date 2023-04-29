package command

import (
	"github.com/flowshot-io/navigator-client-go/navigatorservice/v1"
	"github.com/spf13/cobra"
)

func (c *Command) NewCreateCommand() *cobra.Command {
	var name string

	cc := &cobra.Command{
		Use:   "create",
		Short: "Create assets",
		Long:  `Create assets`,
		Run: func(cmd *cobra.Command, args []string) {
			request := &navigatorservice.CreateAssetRequest{
				Name: name,
			}

			client, err := c.driver.clientFactory.NavigatorClient(cmd)
			if err != nil {
				cmd.PrintErrln("Unable to create navigator client", err)
				return
			}

			resp, err := client.CreateAsset(cmd.Context(), request)
			if err != nil {
				cmd.PrintErrln("Unable to create asset", err)
				return
			}

			cmd.Println("Created asset: ", resp.Id)
		},
	}

	cc.Flags().StringVarP(&name, "name", "n", "", "Name for asset")

	cc.MarkFlagRequired("name")

	return cc
}
