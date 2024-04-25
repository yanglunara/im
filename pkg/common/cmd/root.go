package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yunbaifan/pkg/utils/errs"
)

type CommandOptions func(*CmdOption)

func WithLoggerPreName(name string) CommandOptions {
	return func(o *CmdOption) {
		o.preName = name
	}
}

func WithConfigs(configMap map[string]any) CommandOptions {
	return func(o *CmdOption) {
		o.configMap = configMap
	}
}

type CmdOption struct {
	preName   string
	configMap map[string]any
}

type RootCommand struct {
	Command      cobra.Command
	processName  string
	processIndex int
}

func NewRootCommand(processName string, opts ...CommandOptions) *RootCommand {
	rootCmd := &RootCommand{
		processName: processName,
	}
	cmd := cobra.Command{
		Use:  "Start Im application",
		Long: fmt.Sprintf(`Start %s`, processName),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.persistentPreRunE(cmd, opts...)
		},
		SilenceErrors: false,
		SilenceUsage:  true,
	}
	cmd.Flags().StringP(FlagConf, "c", "", "config file path")
	cmd.Flags().IntP(FlagProcessIndex, "i", 0, "process startup sequence number")
	rootCmd.Command = cmd
	return rootCmd
}

func (r *RootCommand) persistentPreRunE(cmd *cobra.Command, opts ...CommandOptions) error {
	cmdOpts := r.applyCmdOpts(opts...)
	if err := r.initConfig(cmd, cmdOpts); err != nil {
		return errs.Warp(err, "init config failed")
	}
	if err := r.initLogger(cmdOpts); err != nil {
		return errs.Warp(err, "init logger failed")
	}
	return nil
}

// initConfig 初始化配置
func (r *RootCommand) initConfig(cmd *cobra.Command, cmdOpt *CmdOption) error {
	confDirectory, processIndex, err := r.getFlag(cmd)
	if err != nil {
		return err
	}
	r.processIndex = processIndex
	fmt.Printf("confDirectory: %s, processIndex: %d\n", confDirectory, processIndex)
	return nil
}

// initLogger 初始化日志
func (r *RootCommand) initLogger(cmdOpt *CmdOption) error {
	return nil
}

func (r *RootCommand) applyCmdOpts(opts ...CommandOptions) *CmdOption {
	cmdOpts := newDefaultCmdOpts()
	for _, opt := range opts {
		opt(cmdOpts)
	}
	return cmdOpts
}

func newDefaultCmdOpts() *CmdOption {
	return &CmdOption{
		preName: "im-serveric-log",
	}
}

func (r *RootCommand) Execute() error {
	return r.Command.Execute()
}

func (r *RootCommand) ProcessIndex() int {
	return r.processIndex
}

func (r *RootCommand) getFlag(cmd *cobra.Command) (string, int, error) {
	confDirectory, err := cmd.Flags().GetString(FlagConf)
	if err != nil {
		return "", 0, err
	}
	processIndex, err := cmd.Flags().GetInt(FlagProcessIndex)
	if err != nil {
		return "", 0, err
	}
	return confDirectory, processIndex, nil
}
