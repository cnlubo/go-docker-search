package main

import (
	"github.com/cnlubo/go-docker-search/utils"
)

func main() {

	cli := NewCli()
	// set global flags for rootCmd in cli.

	cli.SetFlags()

	base := &baseCommand{cmd: cli.rootCmd, cli: cli}
	base.Cmd().SilenceErrors = true
	// disable sort flags
	base.Cmd().Flags().SortFlags = false
	base.Cmd().PersistentFlags().SortFlags = false

	// Add all subCommands.
	// cli.AddCommand(base, &KeyCommand{})
	// cli.AddCommand(base, &CfgCommand{})
	// cli.AddCommand(base, &ClusterCommand{})
	// cli.AddCommand(base, &HostAliasCommand{})
	// cli.AddCommand(base, &InstallCommand{})
	// cli.AddCommand(base, &UninstallCommand{})
	// cli.AddCommand(base, &BackupCommand{})
	cli.AddCommand(base, &VersionCommand{})

	// // add generate doc command
	// cli.AddCommand(base, &GenDocCommand{})
	err := cli.Run()
	utils.CheckAndExit(err)
}
