package display

import (
	"fmt"
	"image"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/flowshot-io/navigator-client-go/commandservice/v1"
	"github.com/flowshot-io/navigator-client-go/fileservice/v1"
	"github.com/nfnt/resize"
	"github.com/olekukonko/tablewriter"
)

const (
	ImageHost = "localhost:50054/files"
)

type (
	DisplayFile struct {
		Data  *fileservice.File
		Image *image.Image
		URL   string
	}

	Options struct {
		ImageHost    string
		DisplayImage bool
	}

	Service struct {
		fileClient   fileservice.FileServiceClient
		imageHost    string
		displayImage bool
	}
)

func NewService(fileClient fileservice.FileServiceClient, opts *Options) *Service {
	if opts == nil {
		opts = &Options{}
	}

	if opts.ImageHost == "" {
		opts.ImageHost = ImageHost
	}

	return &Service{
		fileClient:   fileClient,
		imageHost:    opts.ImageHost,
		displayImage: opts.DisplayImage,
	}
}

func (s *Service) Assets(assets ...*commandservice.Asset) string {
	return Assets(assets...)
}

func (s *Service) Files(files ...*fileservice.File) string {
	var displayFiles []*DisplayFile
	for _, file := range files {
		url := fmt.Sprintf("http://%s/%s.%s", s.imageHost, file.FileID, file.Extension)

		var image *image.Image
		if strings.HasPrefix(file.Mime, "image/") && s.displayImage {
			var err error
			image, err = getImageFromURL(url)
			if err != nil {
				fmt.Printf("Failed to get image from url: %s\n", err)
			}
		}

		displayFiles = append(displayFiles, &DisplayFile{
			Data:  file,
			Image: image,
			URL:   url,
		})
	}

	return s.RenderFiles(displayFiles...)
}

func Assets(assets ...*commandservice.Asset) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	table.SetCaption(true, "Query Results.")
	table.SetHeader([]string{"ID", "Name", "CreatedAt", "UpdatedAt"})

	var rows [][]string
	for _, result := range assets {
		rows = append(rows, []string{
			result.AssetID,
			result.Name,
			result.CreatedAt,
			result.UpdatedAt,
		})
	}

	table.AppendBulk(rows)
	table.Render()

	return tableString.String()
}

func (s *Service) RenderFiles(files ...*DisplayFile) string {
	if len(files) == 0 {
		return ""
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	headers := []string{"ID", "AssetID", "Status", "URL"}

	if s.displayImage {
		headers = append(headers, "Image")
	}

	table.SetCaption(true, "Query Results.")
	table.SetHeader(headers)
	table.SetRowLine(true)
	table.SetRowSeparator("-")

	var imageWidth uint = 20
	var fileID string
	var assetID string
	var status string
	var rows [][]string

	if len(files) == 1 {
		imageWidth = 32
	}

	for _, result := range files {
		if result.Data != nil {
			fileID = result.Data.FileID
			assetID = result.Data.AssetID
			status = result.Data.Status.String()
		}

		data := []string{
			fileID,
			assetID,
			status,
			result.URL,
		}

		if s.displayImage {
			imgStr := "Image not found"
			if result.Image != nil {
				imgStr = renderImageToBlocks(*result.Image, imageWidth)
			}

			data = append(data, imgStr)
		}

		rows = append(rows, data)
	}

	table.AppendBulk(rows)
	table.Render()

	return tableString.String()
}

func renderImageToBlocks(img image.Image, maxWidth uint) string {
	// Resize image
	img = resize.Resize(maxWidth, maxWidth*uint(img.Bounds().Dy())/uint(img.Bounds().Dx())*2, img, resize.Bilinear)

	// Create a Builder and use it to write our image to a string
	var builder strings.Builder

	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 2 {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			upper := img.At(x, y)
			r1, g1, b1, _ := upper.RGBA()

			var r2, g2, b2 uint32
			if y+1 < bounds.Max.Y {
				lower := img.At(x, y+1)
				r2, g2, b2, _ = lower.RGBA()
			}

			// Add an ANSI escape code for color before each character
			builder.WriteString(fmt.Sprintf("\033[48;2;%d;%d;%dm\033[38;2;%d;%d;%dm▄", r2>>8, g2>>8, b2>>8, r1>>8, g1>>8, b1>>8))
		}

		builder.WriteString("\033[0m\n")
	}

	return builder.String()
}

func getImageFromURL(url string) (*image.Image, error) {
	// Perform the HTTP request to get the image data
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get image from url: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	// Decode the image
	img, err := imaging.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	return &img, nil
}
