package cmds

import (
    "github.com/spf13/cobra"
    "pbin/cmd/pbs/cmds/server"
    "pbin/common/logger"
)

var (
    rootCmd = &cobra.Command{}
)

func Execute() {
    rootCmd.AddCommand(server.Cmd)
    if err := rootCmd.Execute(); err != nil {
        logger.Info(err)
    }
}
