package main

import (
	"fmt"
	"github.com/cnlubo/go-docker-search/registry"
	"github.com/cnlubo/go-docker-search/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"sort"
	"strconv"
)

var DockerRepositoryDesc = "Tool for image repository ..... "

type DockerRepositoryCmd struct {
	baseCommand
}

// Init initialize command.
func (cc *DockerRepositoryCmd) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "repo",
		Aliases: []string{"mrepo"},
		Short:   "Tool for image repository (alias: mrepo)",
		Long:    DockerRepositoryDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ArgumentsCheck(len(args), -1, 0); err != nil {
				utils.Displaylogo()
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "args input failed")
			} else {
				return cc.cmd.Help()
			}
		},
	}
	c.AddCommand(cc, &containerSearchCmd{})
}

var SearchContainerDesc = "Search the Docker Hub for image Repository ...."

type containerSearchCmd struct {
	baseCommand
	pagesize uint8
}

// Init initializes command.
func (cc *containerSearchCmd) Init(c *Cli) {

	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Search the Docker Hub for image Repository (alias:ls)",
		Long:    SearchContainerDesc,
		RunE: func(cmd *cobra.Command, args []string) error {

			if err := utils.ArgumentsCheck(len(args), 1, 1); err != nil {
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "args input failed")
			} else {
				return cc.runSearchContainer(args)
			}
		},
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *containerSearchCmd) addFlags() {
	flagSet := cc.cmd.Flags()
	flagSet.SetInterspersed(false)
	flagSet.Uint8VarP(&cc.pagesize, "pagesize", "p", 25, "Max number of search results (default 25)")
}

var tableHead = []string{
	"NAME",
	"STARS",
	"OFFICIAL",
	"DESCRIPTION",
}

func (cc *containerSearchCmd) runSearchContainer(args []string) error {

	err, rs := registry.SearchContainer(args[0], cc.pagesize)
	if err != nil {
		return errors.WithMessage(err, "search Container failure")
	}
	cc.cli.PrintTable(displayDockerImage(rs))
	return nil
}

func displayDockerImage(containers registry.SearchResults) (tableHeader []string, hColors []tablewriter.Colors, colColors []tablewriter.Colors, data [][]string, defaultColWidth int) {

	var rowData [][]string
	colWidth := 50
	images := containers.Results
	sort.Sort(images)
	for _, r := range images {
		var row []string
		if r.IsOfficial {
			row = append(row, r.Name, strconv.Itoa(r.StarCount), "[OK]", r.Description)
		} else {
			row = append(row, r.Name, strconv.Itoa(r.StarCount), "", r.Description)

		}
		rowData = append(rowData, row)
	}

	var headerColors, columnColors []tablewriter.Colors
	for i := 1; i <= len(tableHead); i++ {
		headerColors = append(headerColors, tablewriter.Colors{tablewriter.FgHiMagentaColor, tablewriter.Bold})
		columnColors = append(columnColors, tablewriter.Colors{tablewriter.FgHiMagentaColor})
	}
	return tableHead, headerColors, columnColors, rowData, colWidth
}

