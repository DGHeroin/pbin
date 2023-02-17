package cmds

import (
    "github.com/spf13/cobra"
    "pbin/cmd/pbc/cmds/get"
    "pbin/cmd/pbc/cmds/put"
    "pbin/common/logger"
)

var (
    rootCmd = &cobra.Command{}
)

func Execute() {
    rootCmd.AddCommand(put.Cmd)
    rootCmd.AddCommand(get.Cmd)
    if err := rootCmd.Execute(); err != nil {
        logger.Info(err)
    }
}
