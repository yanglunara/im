package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/yanglunara/im/internal/api"
)

// ApiCommand is a struct that represents the api command
type ApiCommand struct {
	*RootCommand
	ctx       context.Context
	configMap map[string]any
	apiConf   *api.Conf
}

func NewApiCommand() *ApiCommand {
	var (
		apiConf api.Conf
	)
	ret := &ApiCommand{
		apiConf: &apiConf,
		ctx:     context.WithValue(context.Background(), "version", "v1.0.0"),
	}
	ret.configMap = map[string]any{
		ImAPIConfigFileName:  &apiConf.Rpc,
		ConsulConfigFileName: &apiConf.Consul,
		ShareFileName:        &apiConf.Share,
	}
	ret.RootCommand = NewRootCommand("api", WithConfigs(ret.configMap))
	ret.Command.RunE = func(_ *cobra.Command, _ []string) error {
		return ret.start()
	}
	return ret
}

func (a *ApiCommand) start() error {
	return api.Start(a.ctx, a.processIndex, a.apiConf)
}

func (a *ApiCommand) Exec() error {
	return a.Command.Execute()
}
