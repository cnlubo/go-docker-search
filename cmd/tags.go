package main

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/cnlubo/go-docker-search/registry"
	"github.com/cnlubo/go-docker-search/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"sort"
)

var RepoTagsDesc = "Tool for image repository tags ..... "

// RepoTagsCmd use to implement 'RepoTags' command.
type RepoTagsCmd struct {
	baseCommand
}

// Init initialize command.
func (cc *RepoTagsCmd) Init(c *Cli) {
	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "tags",
		Aliases: []string{"mtags"},
		Short:   "Tool for image repository tags (alias: mtags)",
		Long:    RepoTagsDesc,
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
	c.AddCommand(cc, &repoTagsListCmd{})
	c.AddCommand(cc, &repoTagsListallCmd{})
}

var listRepoTagsDesc = "List all tags for repo ..... "

type repoTagsListCmd struct {
	baseCommand
	pagesize uint8
	filter   string
}

// Init initializes command.
func (cc *repoTagsListCmd) Init(c *Cli) {

	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all tags for repo (alias:ls)",
		Long:    listRepoTagsDesc,
		RunE: func(cmd *cobra.Command, args []string) error {

			if err := utils.ArgumentsCheck(len(args), 1, 1); err != nil {
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "args input failed")
			} else {
				return cc.runListRepoTags(args)
			}
		},
	}
	cc.addFlags()
}

// addFlags adds flags for specific command.
func (cc *repoTagsListCmd) addFlags() {
	flagSet := cc.cmd.Flags()
	flagSet.SetInterspersed(false)
	flagSet.Uint8VarP(&cc.pagesize, "pagesize", "p", 50, "Max number of search results (default 50)")
	flagSet.StringVarP(&cc.filter, "filter", "f", "", "example filter=5.1")
}

var tagsTableHead = []string{
	"TAG",
}

func (cc *repoTagsListCmd) runListRepoTags(args []string) error {

	var major, minor int64

	v, err := semver.NewVersion(cc.filter)
	if err != nil {
		major = 0
		minor = 0
	} else {
		major = v.Major()
		minor = v.Minor()
	}

	repository := registry.RepositoryRequest{
		RepositoryUrl: registry.HostURL + "/v2",
		User:          "cnak47",
		Password:      "inLg9wRN",
		Repo:          args[0],
		Endpoint:      "/tags/list",
		PageSize:      cc.pagesize,
		Major:         major,
		Minor:         minor,
	}
	semverTags, _, err := registry.ListRepoTags(&repository)
	if err != nil {
		return err
	}
	sort.Sort(sort.Reverse(semver.Collection(semverTags)))
	cc.cli.PrintTable(displayRepoTags(semverTags))
	return nil
}

func displayRepoTags(semVerTags []*semver.Version) (tableHeader []string, hColors []tablewriter.Colors, colColors []tablewriter.Colors, data [][]string, defaultColWidth int) {

	var rowData [][]string
	colWidth := 50
	for _, t := range semVerTags {
		var row []string

		row = append(row, t.String())

		rowData = append(rowData, row)
	}

	var headerColors, columnColors []tablewriter.Colors
	for i := 1; i <= len(tagsTableHead); i++ {
		headerColors = append(headerColors, tablewriter.Colors{tablewriter.FgHiMagentaColor, tablewriter.Bold})
		columnColors = append(columnColors, tablewriter.Colors{tablewriter.FgHiMagentaColor})
	}
	return tagsTableHead, headerColors, columnColors, rowData, colWidth
}

var listallRepoTagsDesc = "List all tags for repo ..... "

type repoTagsListallCmd struct {
	baseCommand
}

// Init initializes command.
func (cc *repoTagsListallCmd) Init(c *Cli) {

	cc.cli = c
	cc.cmd = &cobra.Command{
		Use:     "listall",
		Aliases: []string{"la"},
		Short:   "List all tags for repo (alias:ls)",
		Long:    listallRepoTagsDesc,
		RunE: func(cmd *cobra.Command, args []string) error {

			if err := utils.ArgumentsCheck(len(args), 1, 1); err != nil {
				_ = cc.Cmd().Help()
				fmt.Println()
				return errors.WithMessage(err, "args input failed")
			} else {
				return cc.runListallRepoTags(args)
			}
		},
	}
	cc.addFlags()
}

func (cc *repoTagsListallCmd) addFlags() {

}

func (cc *repoTagsListallCmd) runListallRepoTags(args []string) error {
	url      := "https://index.docker.io/"
	username := "cnak47" // anonymous
	password := "inLg9wRN" // anonymous
	hub, _ := registry.New(url, username, password)
	// repositories, _ := hub.Repositories()
	// fmt.Println(repositories)
	tags, _ := hub.Tags(args[0])
	fmt.Println(tags)
	return nil
}
