package command

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/flowshot-io/navigator-client-go/navigatorservice/v1"
	"github.com/flowshot-io/x/pkg/storager"
	"github.com/spf13/cobra"
)

func (c *Command) NewUploadCommand() *cobra.Command {
	var AssetID string
	var fileName string
	var chunkSizeBytes int

	cc := &cobra.Command{
		Use:   "upload",
		Short: "Upload a file to storage",
		Long:  `Upload a file to storage. A chunk size of 5mb or greater will use multipart upload into s3 storage.`,
		Run: func(cmd *cobra.Command, args []string) {
			client, err := c.driver.clientFactory.NavigatorClient(cmd)
			if err != nil {
				cmd.PrintErrln("Unable to create navigator client", err)
				return
			}

			err = c.uploadFile(cmd.Context(), client, fileName, AssetID, chunkSizeBytes)
			if err != nil {
				cmd.PrintErrln("Unable to upload file", err)
				return
			}

			cmd.Println("Upload complete")
		},
	}

	cc.Flags().StringVarP(&AssetID, "id", "i", "", "Asset ID to upload to (required)")
	cc.Flags().StringVarP(&fileName, "file", "f", "", "Local file to upload (required)")
	cc.Flags().IntVarP(&chunkSizeBytes, "chunk", "c", storager.DefaultChunkSize, "Chunk size (in bytes) to use with upload")

	cc.MarkFlagRequired("file")
	cc.MarkFlagRequired("id")

	return cc
}

func (c *Command) uploadFile(ctx context.Context, client navigatorservice.NavigatorServiceClient, filePath string, assetID string, chunkSizeBytes int) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("unable to open file: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("unable to stat file: %w", err)
	}

	streamRequest := &navigatorservice.CreateAssetStreamRequest{
		FileName: filepath.Base(filePath),
		Id:       assetID,
	}

	// Create a upload stream for the asset
	uploadStream, err := client.CreateAssetStream(ctx, streamRequest)
	if err != nil {
		return fmt.Errorf("unable to create upload stream: %w", err)
	}

	// Initialize the upload stream
	stream, err := client.UploadAssetStream(ctx)
	if err != nil {
		return fmt.Errorf("unable to create upload stream: %w", err)
	}

	complete := 0
	buf := make([]byte, chunkSizeBytes)

	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				// Send last chunk
				lastChunkReq := &navigatorservice.AssetUploadStream{
					Id:          uploadStream.Id,
					FileChunk:   nil,
					IsLastChunk: true,
				}

				if err := stream.Send(lastChunkReq); err != nil {
					return fmt.Errorf("unable to send last chunk: %w", err)
				}

				break
			}

			return fmt.Errorf("unable to read file: %w", err)
		}

		chunk := buf[:n]
		req := &navigatorservice.AssetUploadStream{
			Id:          uploadStream.Id,
			FileChunk:   chunk,
			IsLastChunk: false,
		}

		if err := stream.Send(req); err != nil {
			return fmt.Errorf("unable to send file chunk: %w", err)
		}

		complete += len(chunk)
		fmt.Printf("Uploaded %d/%d bytes\n", complete, stat.Size())
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("something went wrong: %w", err)
	}

	return nil
}
