package commands

import (
	"fmt"

	"github.com/kyriosdata/assinatura/internal/assinatura/usecase"
	"github.com/kyriosdata/assinatura/internal/cli/formatter"
	"github.com/kyriosdata/assinatura/internal/config"
	"github.com/spf13/cobra"
)

var startJar string
var startPort int
var stopPort int

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia o assinador.jar em modo servidor",
	Run: func(cmd *cobra.Command, args []string) {
		metadata, err := usecase.NewStartAssinadorUseCase().Execute(startJar, startPort)
		if err != nil {
			fmt.Fprintln(cmd.ErrOrStderr(), formatter.FormatError(err))
			return
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Assinador iniciado. Porta: %d | PID: %d | Health: %s\n", metadata.Port, metadata.PID, metadata.HealthEndpoint)
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Interrompe o assinador.jar em modo servidor",
	Run: func(cmd *cobra.Command, args []string) {
		if err := usecase.NewStopAssinadorUseCase().Execute(stopPort); err != nil {
			fmt.Fprintln(cmd.ErrOrStderr(), formatter.FormatError(err))
			return
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Assinador encerrado na porta %d\n", stopPort)
	},
}

func init() {
	startCmd.Flags().StringVar(&startJar, "jar", "", "Caminho para o assinador.jar")
	startCmd.Flags().IntVar(&startPort, "port", config.DefaultAssinadorPort, "Porta do assinador")
	stopCmd.Flags().IntVar(&stopPort, "port", config.DefaultAssinadorPort, "Porta do assinador")
}
