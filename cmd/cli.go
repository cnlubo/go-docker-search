package main

import (
	"fmt"
	"github.com/cnlubo/go-docker-search/dockerutils"
	"github.com/cnlubo/go-docker-search/utils"
	"github.com/gookit/color"
	"github.com/mitchellh/go-homedir"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
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
	Env     dockerutils.Environment
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
			Short: "My ssh toolkit",
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
			dockerutils.Displaylogo()
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

	var defaultStorePath, defaultSkmPath, defaultSSHPath = os.Getenv("MYSSH_CONFIG_HOME"), os.Getenv("MKM_PATH"), os.Getenv("SSH_PATH")

	home, _ := homedir.Dir()

	if len(defaultStorePath) == 0 {
		defaultStorePath = filepath.Join(home, ".myssh")
	} else {
		if fp, _ := filepath.Rel("~", defaultStorePath); fp != "" {
			defaultStorePath = filepath.Join(home, fp)
		}
	}
	if len(defaultSkmPath) == 0 {
		defaultSkmPath = filepath.Join(home, ".mkm")
	} else {
		if fp, _ := filepath.Rel("~", defaultSkmPath); fp != "" {
			defaultSkmPath = filepath.Join(home, fp)
		}
	}

	if len(defaultSSHPath) == 0 {
		defaultSSHPath = filepath.Join(home, ".ssh")
	} else {
		if fp, _ := filepath.Rel("~", defaultSSHPath); fp != "" {
			defaultSSHPath = filepath.Join(home, fp)
		}
	}

	flags := c.rootCmd.PersistentFlags()
	flags.BoolVar(&c.Option.noColor, "no-color", false, "Disable color when outputting message.")
	flags.StringVar(&c.Option.Env.StorePath, "configPath", defaultStorePath, "Path where store myssh profiles.\ncan also be set by the MYSSH_CONFIG_HOME environment variable.")

	// flags.StringVar(&c.Option.Env.SKMPath, "mkmPath", defaultSkmPath, "Path where myssh should store multi SSHKeys.\ncan also be set by the MKM_PATH environment variable.")
	// flags.StringVar(&c.Option.Env.SSHPath, "sshPath", defaultSSHPath, "Path to .ssh folder.\ncan also be set by the SSH_PATH environment variable.")
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

func (c *Cli) initConfigDir() error {

	if exists := utils.PathExist(c.Env.StorePath); !exists {

		// default context dir
		defCtxPath := filepath.Join(c.Env.StorePath, "contexts")
		if err := os.MkdirAll(filepath.Join(defCtxPath, "default"), 0755); err != nil {
			return errors.Wrap(err, "Create context dir error")
		}

		// default cluster configfile
		// if err := myssh.ClustersConfigExample(&c.Env).WriteTo(filepath.Join(defCtxPath, "default", "default.yaml")); err != nil {
		// 	return errors.Wrap(err, "create cluster default.yaml error")
		// }

		// default SSH config file
		// sshConfigFile := filepath.Join(defCtxPath, "default", "sshconfig")
		// // create file
		// if err := ioutil.WriteFile(sshConfigFile, []byte("Include include/*"+"\n"), 0777); err != nil {
		// 	return errors.Wrap(err, "Create SSh config file error")
		// }
		// if err := os.MkdirAll(filepath.Join(defCtxPath, "default", "include"), 0755); err != nil {
		// 	return errors.Wrap(err, "create default context include dir failed")
		// }

		// create symlink ~/.ssh/config
		// if err := myssh.CreateSSHlink("default", &c.Env); err != nil {
		// 	return errors.Wrap(err, "Create default symlink error")
		// }

		// default main config file
		// if err := myssh.MainConfigExample(&c.Env).WriteTo(filepath.Join(c.Env.StorePath, "main.yaml")); err != nil {
		// 	return errors.Wrap(err, "create main.yaml failed")
		// }

	}

	// skm key store
	// Remove existing empty key store
	// if ok, err := utils.IsEmpty(c.Env.SKMPath); ok {
	// 	err = os.Remove(c.Env.SKMPath)
	// 	if err != nil {
	// 		return errors.Wrap(err, "Remove empty dir failed")
	// 	}
	// }

	// // create key store
	//
	// if exits := utils.PathExist(c.Env.SKMPath); !exits {
	// 	err := os.Mkdir(c.Env.SKMPath, 0755)
	// 	if err != nil {
	// 		return errors.Wrap(err, "Create ssh key store failed")
	// 	}
	//
	// 	// Check & move existing keys into default folder
	// 	err = myssh.MoveDefaultSSKey(&c.Env)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

// initialize myssh
func (c *Cli) Initialize() error {

	op := &c.Option.Env
	if err := c.initConfigDir(); err != nil {
		return err
	}

	mainConfigFile := filepath.Join(op.StorePath, "main.yaml")
	// load main config
	if err := myssh.Main.LoadFrom(mainConfigFile); err != nil {
		return errors.Wrap(err, "load mainConfigFile failed")
	}

	// check context
	if len(myssh.Main.Contexts) == 0 {
		return errors.New("get context failed")
	}

	// get current context
	ctx, exists := myssh.Main.Contexts.FindContextByName(myssh.Main.Current)
	if !exists {
		return errors.New(fmt.Sprintf("current context: %s not found\n", myssh.Main.Current))
	}
	// load current context cluster
	err := myssh.ClustersCfg.LoadFrom(ctx.ClusterConfig)
	if err != nil {
		return errors.Wrap(err, "load current cluster failed")
	}

	return nil

}

// NewTableAsciiDisplay creates a display instance, and uses to format output with AsciiTable.
func (c *Cli) NewAsciiTableDisplay() *DisplayTable {

	w := tablewriter.NewWriter(os.Stdout)
	return &DisplayTable{w}
}

// Print outputs the obj's fields.
func (c *Cli) PrintTable(tableHead []string, headerColors []tablewriter.Colors, columnColors []tablewriter.Colors, rowData [][]string) {

	display := c.NewAsciiTableDisplay()
	if tableHead != nil {
		display.SetHeader(tableHead)
		if !c.noColor {
			display.SetHeaderColor(headerColors...)
		}
	}

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
