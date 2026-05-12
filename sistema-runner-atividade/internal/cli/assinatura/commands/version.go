package commands

import (
	"fmt"

	"github.com/kyriosdata/assinatura/internal/appinfo"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Exibe a versão do CLI assinatura",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "assinatura %s (%s, %s)\n", appinfo.Version, appinfo.Commit, appinfo.Date)
	},
}
