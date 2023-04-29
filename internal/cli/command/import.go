package command

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/flowshot-io/navigator-client-go/navigatorservice/v1"
	"github.com/flowshot-io/x/pkg/storager"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (c *Command) NewImportCommand() *cobra.Command {
	var path string
	var waitTime int

	cc := &cobra.Command{
		Use:   "import",
		Short: "Import assets",
		Long:  `Import assets`,
		Run: func(cmd *cobra.Command, args []string) {
			client, err := c.driver.clientFactory.NavigatorClient(cmd)
			if err != nil {
				cmd.PrintErrln("Unable to create navigator client", err)
			}

			err = c.importFilesByTypes(cmd.Context(), client, path, []string{".jpg", ".jpeg", ".png"}, time.Duration(waitTime)*time.Millisecond)
			if err != nil {
				cmd.PrintErrln("Unable to import assets", err)
				return
			}
		},
	}

	cc.Flags().StringVarP(&path, "path", "p", "", "Local path to import from (required)")
	cc.Flags().IntVarP(&waitTime, "wait", "w", 2000, "Time to wait between creating a asset and uploading a file (in ms)")

	cc.MarkFlagRequired("path")

	return cc
}

func (c *Command) importFilesByTypes(ctx context.Context, client navigatorservice.NavigatorServiceClient, dirPath string, fileTypes []string, waitTime time.Duration) error {
	// Check if the given path exists and is a directory
	fileInfo, err := os.Stat(dirPath)
	if err != nil {
		return fmt.Errorf("error accessing path: %v", err)
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("path is not a directory")
	}

	// Function to determine if a file has the desired type
	isDesiredFileType := func(filePath string) bool {
		for _, fileType := range fileTypes {
			if strings.HasSuffix(strings.ToLower(filePath), strings.ToLower(fileType)) {
				return true
			}
		}
		return false
	}

	err = filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type().IsRegular() && isDesiredFileType(path) {
			err = c.importFile(ctx, client, path, waitTime)
			if err != nil {
				return fmt.Errorf("error importing file: %v", err)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking the path: %v", err)
	}
	return nil
}

func (c *Command) importFile(ctx context.Context, client navigatorservice.NavigatorServiceClient, path string, waitTime time.Duration) error {
	name := c.humanizeFileName(filepath.Base(path))

	request := &navigatorservice.CreateAssetRequest{
		Name: name,
	}

	fmt.Println("Creating asset: ", name)

	resp, err := client.CreateAsset(ctx, request)
	if err != nil {
		return fmt.Errorf("unable to create asset: %v", err)
	}

	fmt.Println("Created asset: ", resp.Id)

	// Wait for the asset to be created
	time.Sleep(waitTime)

	fmt.Println("Uploading file: ", path)

	err = c.uploadFile(ctx, client, path, resp.Id, storager.DefaultChunkSize)
	if err != nil {
		return fmt.Errorf("unable to upload file: %v", err)
	}

	fmt.Println("Uploaded file: ", path)

	return nil
}

func (c *Command) humanizeFileName(filename string) string {
	if len(filename) == 0 {
		return filename
	}

	extension := filepath.Ext(filename)
	nameWithoutExtension := strings.TrimSuffix(filename, extension)

	return cases.Title(language.English).String(nameWithoutExtension)
}
