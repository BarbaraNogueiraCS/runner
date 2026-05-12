package commands

import (
	"fmt"

	"github.com/kyriosdata/assinatura/internal/assinatura/dto"
	"github.com/kyriosdata/assinatura/internal/assinatura/usecase"
	"github.com/kyriosdata/assinatura/internal/cli/formatter"
	"github.com/kyriosdata/assinatura/internal/config"
	"github.com/spf13/cobra"
)

var signOpts dto.SignCommand

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Cria uma assinatura digital simulada",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := usecase.NewSignUseCase().Execute(signOpts)
		if err != nil {
			fmt.Fprintln(cmd.ErrOrStderr(), formatter.FormatError(err))
			return
		}
		fmt.Fprintln(cmd.OutOrStdout(), formatter.FormatResult(formatter.Result(result)))
	},
}

func init() {
	signCmd.Flags().StringVar(&signOpts.Documento, "documento", "", "Documento a ser assinado")
	signCmd.Flags().StringVar(&signOpts.Certificado, "certificado", "", "Certificado usado na assinatura simulada")
	signCmd.Flags().StringVar(&signOpts.Algoritmo, "algoritmo", "simulado", "Algoritmo conceitual da simulação")
	signCmd.Flags().BoolVar(&signOpts.Local, "local", false, "Força execução local via java -jar")
	signCmd.Flags().IntVar(&signOpts.Port, "port", config.DefaultAssinadorPort, "Porta do assinador em modo servidor")
	signCmd.Flags().IntVar(&signOpts.TimeoutMinutes, "timeout", 0, "Tempo de inatividade em minutos")
	signCmd.Flags().StringVar(&signOpts.JarPath, "jar", "", "Caminho para o assinador.jar")
}
