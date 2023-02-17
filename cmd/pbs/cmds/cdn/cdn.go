package cdn

import (
    "github.com/mailgun/groupcache"
    "github.com/spf13/cobra"
)

var (
    Cmd = &cobra.Command{
        RunE: func(cmd *cobra.Command, args []string) error {
            return nil
        },
    }
)
var (
    gc groupcache.Group
)
