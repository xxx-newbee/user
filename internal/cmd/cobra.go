package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xxx-newbee/user/internal/cmd/api"
	"github.com/xxx-newbee/user/internal/cmd/migrate"
)

var rootCmd = &cobra.Command{
	Use:          "user",
	Short:        "user",
	SilenceUsage: true,
	Long:         "user-srv",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			tip()
			return errors.New("requires at least 1 arg")
		}
		return nil
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		tip()
	},
}

func tip() {
	usageStr := `欢迎使用服务，请使用 -h 查看命令`
	fmt.Printf("%s\n", usageStr)
}

func init() {
	rootCmd.AddCommand(migrate.StartCmd)
	rootCmd.AddCommand(api.StartCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
