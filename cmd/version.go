package main

import (
	"encoding/base64"
	"fmt"
	"github.com/cnlubo/go-docker-search/utils"
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var versionTpl = `%s

%s V%s
%s

OsArch:   %s_%s
GoVersion:%s
BuildTime:%s
GitCommit:%s
`

var versionDescription = "Display the version information of docker-search， " +
	"including GoVersion, KernelVersion, Os, Version, Arch, BuildTime and GitCommit."

// VersionCommand use to implement 'version' command.
type VersionCommand struct {
	baseCommand
}

// Init initialize version command.
func (v *VersionCommand) Init(c *Cli) {
	v.cli = c
	v.cmd = &cobra.Command{
		Use:   "version",
		Short: "Print versions about Docker-search",
		Args:  cobra.NoArgs,
		Long:  versionDescription,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return v.runVersion()
		},
	}
	v.addFlags()
}

// addFlags adds flags for specific command.
func (v *VersionCommand) addFlags() {
	// TODO: add flags here
}

// runVersion is the entry of version command.
func (v *VersionCommand) runVersion() error {

	result, err := utils.Version()
	if err != nil {
		return errors.Wrap(err, "failed to get system version")
	}
	banner, _ := base64.StdEncoding.DecodeString(result.Banner)
	fmt.Printf(color.FgLightGreen.Render(versionTpl), banner, result.Name, result.Version, result.GitHub, result.Os, result.Arch, result.GoVersion, result.BuildTime, result.GitCommit)
	return nil
}

var logo = `%s

%s V%s
%s

`

func Displaylogo() error {

	result, err := utils.Version()
	if err != nil {
		return errors.Wrap(err, "failed to get system version")
	}
	banner, _ := base64.StdEncoding.DecodeString(result.Banner)
	fmt.Printf(color.FgGreen.Render(logo), banner, result.Name, result.Version, color.FgMagenta.Render(result.GitHub))
	return nil
}
