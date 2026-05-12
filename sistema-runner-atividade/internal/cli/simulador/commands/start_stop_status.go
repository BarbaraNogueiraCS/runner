package commands

import (
	"fmt"

	"github.com/kyriosdata/assinatura/internal/cli/formatter"
	"github.com/kyriosdata/assinatura/internal/config"
	simdto "github.com/kyriosdata/assinatura/internal/simulador/dto"
	"github.com/kyriosdata/assinatura/internal/simulador/usecase"
	"github.com/spf13/cobra"
)

var startOpts simdto.SimulatorConfig
var stopPort int
var statusPort int

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia o Simulador do HubSaúde",
	Run: func(cmd *cobra.Command, args []string) {
		status, err := usecase.NewStartSimulatorUseCase().Execute(startOpts)
		if err != nil {
			fmt.Fprintln(cmd.ErrOrStderr(), formatter.FormatError(err))
			return
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Status: %s\nPorta: %d\nPID: %d\n", status.Message, status.Port, status.PID)
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Para o Simulador do HubSaúde",
	Run: func(cmd *cobra.Command, args []string) {
		if err := usecase.NewStopSimulatorUseCase().Execute(stopPort); err != nil {
			fmt.Fprintln(cmd.ErrOrStderr(), formatter.FormatError(err))
			return
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Simulador encerrado na porta %d\n", stopPort)
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Consulta o status do Simulador do HubSaúde",
	Run: func(cmd *cobra.Command, args []string) {
		status, err := usecase.NewStatusSimulatorUseCase().Execute(statusPort)
		if err != nil {
			fmt.Fprintln(cmd.ErrOrStderr(), formatter.FormatError(err))
			return
		}
		running := "não"
		if status.Running {
			running = "sim"
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Em execução: %s\nPorta: %d\nPID: %d\nMensagem: %s\n", running, status.Port, status.PID, status.Message)
	},
}

func init() {
	startCmd.Flags().IntVar(&startOpts.Port, "port", config.DefaultSimulatorPort, "Porta do simulador")
	startCmd.Flags().StringVar(&startOpts.JarPath, "jar", "", "Caminho para o simulador.jar")
	startCmd.Flags().StringVar(&startOpts.Source, "source", "", "URL alternativa para download do simulador.jar")
	startCmd.Flags().StringVar(&startOpts.SHA256, "sha256", "", "Checksum SHA256 esperado do simulador.jar")
	stopCmd.Flags().IntVar(&stopPort, "port", config.DefaultSimulatorPort, "Porta do simulador")
	statusCmd.Flags().IntVar(&statusPort, "port", config.DefaultSimulatorPort, "Porta do simulador")
}
