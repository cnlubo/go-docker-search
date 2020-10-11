package main

import (
	"github.com/cnlubo/go-docker-search/registry"
	"github.com/cnlubo/go-docker-search/utils"
	"github.com/gookit/color"
	"github.com/mitchellh/go-homedir"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	dockerSearchDesc = "Docker Search toolkit. " +
		"Flags and arguments can be input to do what actually you wish. "
)

// Option uses to define the global options.
type Option struct {
	noColor bool
	Env     registry.Environment
}

type Cli struct {
	Option
	rootCmd *cobra.Command
}

// NewCli creates an instance of 'Cli'.
func NewCli() *Cli {
	cc := &Cli{
		rootCmd: &cobra.Command{
			Use:   "docker-search",
			Short: "Docker-Search toolkit",
			Long:  dockerSearchDesc,
			// disable displaying auto generation tag in cli docs
			DisableAutoGenTag: true,
		},
	}
	cobra.EnableCommandSorting = false

	// hide help subCommand
	cc.rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	cc.rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {

		if cc.noColor {
			color.Enable = false
		}

		if len(args) == 0 {
			err := Displaylogo()
			utils.CheckAndExit(err)
		}

		err := cc.Initialize()
		if err != nil {
			return err
		}
		return nil
	}

	cc.rootCmd.Run = func(cmd *cobra.Command, args []string) {

		_ = cc.rootCmd.Help()
	}

	return cc
}

// SetFlags sets all global options.
func (c *Cli) SetFlags() *Cli {

	var defaultStorePath, defaultDockerID, defaultDockerPass = os.Getenv("DOCKER_CONFIG_HOME"), os.Getenv("DOCKER_ID"), os.Getenv("DOCKER_PASS")

	home, _ := homedir.Dir()

	if len(defaultStorePath) == 0 {
		defaultStorePath = filepath.Join(home, ".dockerhub")
	} else {
		if fp, _ := filepath.Rel("~", defaultStorePath); fp != "" {
			defaultStorePath = filepath.Join(home, fp)
		}
	}
	if len(defaultDockerID) == 0 {
		defaultDockerID = ""
	}

	if len(defaultDockerPass) == 0 {
		defaultDockerPass = ""
	}

	flags := c.rootCmd.PersistentFlags()
	flags.BoolVar(&c.Option.noColor, "no-color", false, "Disable color when outputting message.")
	flags.StringVar(&c.Option.Env.StorePath, "configPath", defaultStorePath, "Path where store profiles.\ncan also be set by the DOCKER_CONFIG_HOME environment variable.")
	// flags.StringVar(&c.Option.Env.RegistryUrl, "registry", "https://index.docker.io", "Registry url.")
	flags.StringVar(&c.Option.Env.DockerID, "dockerID", defaultDockerID, "docker hub login docker ID.")
	flags.StringVar(&c.Option.Env.DockerPass, "dockerPass", defaultDockerPass, "docker hub login Password.")
	return c
}

func (c *Cli) Run() error {

	basename := filepath.Base(os.Args[0])
	co, _, _ := c.rootCmd.Find([]string{basename})
	if co != nil {
		for _, ca := range co.Aliases {
			if basename == ca {
				c.rootCmd.SetArgs(os.Args)
				break
			}
		}
	}

	return c.rootCmd.Execute()

}

// AddCommand add a subCommand.
func (c *Cli) AddCommand(parent, child Command) {

	child.Init(c)
	parentCmd := parent.Cmd()
	childCmd := child.Cmd()
	// make command error not return command usage and error
	childCmd.SilenceUsage = true
	childCmd.SilenceErrors = true
	childCmd.DisableFlagsInUseLine = true
	childCmd.PreRun = func(cmd *cobra.Command, args []string) {

	}
	parentCmd.AddCommand(childCmd)
}

// initialize
func (c *Cli) Initialize() error {

	return nil

}

// NewTableAsciiDisplay creates a display instance, and uses to format output with AsciiTable.
func (c *Cli) NewAsciiTableDisplay() *DisplayTable {

	w := tablewriter.NewWriter(os.Stdout)
	return &DisplayTable{w}
}

// Print outputs the obj's fields.
func (c *Cli) PrintTable(tableHead []string, headerColors []tablewriter.Colors, columnColors []tablewriter.Colors, rowData [][]string,defaultColWidth int) {

	display := c.NewAsciiTableDisplay()
	if tableHead != nil {
		display.SetHeader(tableHead)
		if !c.noColor {
			display.SetHeaderColor(headerColors...)
		}
	}
	display.SetAutoFormatHeaders(true)
	display.SetColWidth(defaultColWidth)
	display.SetCenterSeparator(" ")
	display.SetColumnSeparator(" ")
	display.SetRowLine(false)
	display.SetBorder(false)

	display.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})

	if !c.noColor {
		if columnColors != nil {
			display.SetColumnColor(columnColors...)
		}
	}
	if rowData != nil {
		display.AppendBulk(rowData)
	}
	display.Render()
}
