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
	cli.AddCommand(base, &VersionCommand{})
	cli.AddCommand(base, &DockerRepositoryCmd{})
	cli.AddCommand(base, &RepoTagsCmd{})
	// // add generate doc command
	// cli.AddCommand(base, &GenDocCommand{})
	err := cli.Run()
	utils.CheckAndExit(err)
}
