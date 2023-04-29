package cli

import (
	"github.com/flowshot-io/navigator-cli/internal/cli/command"
	"github.com/spf13/cobra"

	cc "github.com/ivanpirog/coloredcobra"
)

var (
	FlagAddress string = "address"
)

func New() (*cobra.Command, error) {
	driver := command.NewDriver()

	rootCMD := command.NewCommand(driver)

	// Configure help template colours
	cc.Init(&cc.Config{
		RootCmd:         rootCMD,
		Headings:        cc.Cyan + cc.Bold + cc.Underline,
		Commands:        cc.Bold,
		ExecName:        cc.Bold,
		Flags:           cc.Bold,
		Aliases:         cc.Bold,
		Example:         cc.Italic,
		NoExtraNewlines: true,
		NoBottomNewline: true,
	})

	return rootCMD, nil
}
