package app

import (
	"github.com/spf13/cobra"

	"squirrel-dev/cmd/squctl/app/options"
)

func NewServerCommand() *cobra.Command {
	o := options.NewAppOptions()
	cmd := &cobra.Command{
		Use:  "app",
		Long: `Long describe.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(o)
		},
	}
	cmd.Flags().StringVarP(&o.ConfFile, "config", "c", "config/squctl.yaml", "Config file path.")
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of squirrel-dev.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("abc")
		},
	}

	cmd.AddCommand(versionCmd)

	return cmd
}

func run(o *options.AppOptions) (err error) {
	server, err := o.NewServer()
	server.Run()
	return
}
