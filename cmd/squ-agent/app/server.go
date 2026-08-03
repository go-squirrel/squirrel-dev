package app

import (
	"fmt"

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

	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrate(o)
		},
	}

	rollbackCmd := &cobra.Command{
		Use:   "rollback [version]",
		Short: "Rollback one database migration.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRollback(o, args[0])
		},
	}

	cmd.AddCommand(versionCmd)
	cmd.AddCommand(migrateCmd)
	cmd.AddCommand(rollbackCmd)

	return cmd
}

func run(o *options.AppOptions) (err error) {
	// 初始化返回值
	response.Init()
	server, err := o.NewServer()
	if err != nil {
		return err
	}
	return server.Run()
}

func runMigrate(o *options.AppOptions) error {
	server, err := o.NewServer()
	if err != nil {
		return err
	}
	if err := server.Migrate(); err != nil {
		return err
	}
	fmt.Println("migrations completed")
	return nil
}

func runRollback(o *options.AppOptions, version string) error {
	server, err := o.NewServer()
	if err != nil {
		return err
	}
	if err := server.RollbackMigration(version); err != nil {
		return err
	}
	fmt.Printf("rollback completed: %s\n", version)
	return nil
}
