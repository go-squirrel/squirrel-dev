package app

import (
	"github.com/spf13/cobra"

	"squirrel-dev/cmd/squ-agent/app/options"
	"squirrel-dev/internal/pkg/response"
)

var Version string

func NewServerCommand() *cobra.Command {
	o := options.NewAppOptions()
	cmd := &cobra.Command{
		Use:  "app",
		Long: `Long describe.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(o)
		},
	}
	cmd.Flags().StringVarP(&o.ConfFile, "config", "c", "config/agent.yaml", "Config file path.")
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of squirrel-dev.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("version:", Version)
		},
	}

	cmd.AddCommand(versionCmd)

	return cmd
}

func run(o *options.AppOptions) (err error) {
	// 初始化返回值
	response.Init()
	server, err := o.NewServer()
	server.Run()
	return
}
