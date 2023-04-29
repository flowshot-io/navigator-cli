package command

import (
	"io"
	"os"
	"path/filepath"

	"github.com/flowshot-io/navigator-client-go/navigatorservice/v1"
	"github.com/spf13/cobra"
)

func (c *Command) NewDownloadCommand() *cobra.Command {
	var assetID string
	var destinationPath string
	var chunkSize int32

	cc := &cobra.Command{
		Use:   "download",
		Short: "Download file",
		Long:  `Download file`,
		Run: func(cmd *cobra.Command, args []string) {
			req := &navigatorservice.DownloadAssetRequest{
				Id:        assetID,
				ChunkSize: chunkSize,
			}

			client, err := c.driver.clientFactory.NavigatorClient(cmd)
			if err != nil {
				cmd.PrintErrln("Unable to create navigator client", err)
				return
			}

			stream, err := client.DownloadAsset(cmd.Context(), req)
			if err != nil {
				cmd.PrintErrln("Unable to create download stream", err)
				return
			}

			file, err := os.CreateTemp(destinationPath, "*")
			if err != nil {
				cmd.PrintErrln("Unable to create local file", err)
				return
			}
			defer file.Close()

			var fileName string

			for {
				res, err := stream.Recv()
				if err != nil {
					if err == io.EOF {
						break
					}
					cmd.PrintErrln("Unable to receive file chunk", err)
					return
				}

				if fileName == "" {
					fileName = filepath.Join(destinationPath, res.FileName)
				}

				_, err = file.Write(res.FileChunk)
				if err != nil {
					cmd.PrintErrln("Unable to write file chunk", err)
					return
				}

				if res.IsLastChunk {
					break
				}
			}

			err = os.Rename(file.Name(), fileName)
			if err != nil {
				cmd.PrintErrln("Unable to rename file", err)
				return
			}

			cmd.Println("Downloaded file to", fileName)
		},
	}

	cc.Flags().StringVarP(&assetID, "id", "i", "", "Asset ID to download")
	cc.Flags().Int32VarP(&chunkSize, "chunk", "c", 1024*1024, "Chunk size (in bytes) to use with upload")
	cc.Flags().StringVarP(&destinationPath, "destination", "d", ".", "Destination path to download file to")

	cc.MarkFlagRequired("id")

	return cc
}
