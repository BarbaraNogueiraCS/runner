package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// version deve ser variável, e não constante, para permitir sobrescrita no build
// por meio de: -ldflags "-X main.version=<tag>".
var version = "dev"

func newRootCommand() *cobra.Command {
	var showVersion bool

	rootCmd := &cobra.Command{
		Use:           "assinatura",
		Short:         "CLI do Sistema Runner para operações de assinatura digital simulada",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if showVersion {
				fmt.Fprintln(cmd.OutOrStdout(), version)
				return nil
			}

			return cmd.Help()
		},
	}

	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "exibe a versão atual do CLI")

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Exibe a versão atual do CLI assinatura",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), version)
		},
	})

	return rootCmd
}

func main() {
	if err := newRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
