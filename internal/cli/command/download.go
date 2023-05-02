package command

import (
	"fmt"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			req := &navigatorservice.DownloadAssetRequest{
				Id:        assetID,
				ChunkSize: chunkSize,
			}

			client, err := c.driver.clientFactory.NavigatorClient(cmd)
			if err != nil {
				return fmt.Errorf("unable to create navigator client: %w", err)
			}

			stream, err := client.DownloadAsset(cmd.Context(), req)
			if err != nil {
				return fmt.Errorf("unable to create download stream: %w", err)
			}

			file, err := os.CreateTemp(destinationPath, "*")
			if err != nil {
				return fmt.Errorf("unable to create local file: %w", err)
			}
			defer file.Close()

			var fileName string

			for {
				res, err := stream.Recv()
				if err != nil {
					if err == io.EOF {
						break
					}
					return fmt.Errorf("unable to receive file chunk: %w", err)
				}

				if fileName == "" {
					fileName = filepath.Join(destinationPath, res.FileName)
				}

				_, err = file.Write(res.FileChunk)
				if err != nil {
					return fmt.Errorf("unable to write file chunk: %w", err)
				}

				if res.IsLastChunk {
					break
				}
			}

			err = os.Rename(file.Name(), fileName)
			if err != nil {
				return fmt.Errorf("unable to rename file: %w", err)
			}

			cmd.Println("Downloaded file to", fileName)

			return nil
		},
	}

	cc.Flags().StringVarP(&assetID, "id", "i", "", "Asset ID to download")
	cc.Flags().Int32VarP(&chunkSize, "chunk", "c", 1024*1024, "Chunk size (in bytes) to use with upload")
	cc.Flags().StringVarP(&destinationPath, "destination", "d", ".", "Destination path to download file to")

	cc.MarkFlagRequired("id")

	return cc
}
