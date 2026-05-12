package commands

import (
	"fmt"

	"github.com/kyriosdata/assinatura/internal/assinatura/dto"
	"github.com/kyriosdata/assinatura/internal/assinatura/usecase"
	"github.com/kyriosdata/assinatura/internal/cli/formatter"
	"github.com/kyriosdata/assinatura/internal/config"
	"github.com/spf13/cobra"
)

var validateOpts dto.ValidateCommand

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Valida uma assinatura digital simulada",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := usecase.NewValidateUseCase().Execute(validateOpts)
		if err != nil {
			fmt.Fprintln(cmd.ErrOrStderr(), formatter.FormatError(err))
			return
		}
		fmt.Fprintln(cmd.OutOrStdout(), formatter.FormatResult(formatter.Result(result)))
	},
}

func init() {
	validateCmd.Flags().StringVar(&validateOpts.Documento, "documento", "", "Documento relacionado à assinatura")
	validateCmd.Flags().StringVar(&validateOpts.Assinatura, "assinatura", "", "Assinatura a ser validada")
	validateCmd.Flags().StringVar(&validateOpts.Certificado, "certificado", "", "Certificado usado na validação simulada")
	validateCmd.Flags().BoolVar(&validateOpts.Local, "local", false, "Força execução local via java -jar")
	validateCmd.Flags().IntVar(&validateOpts.Port, "port", config.DefaultAssinadorPort, "Porta do assinador em modo servidor")
	validateCmd.Flags().StringVar(&validateOpts.JarPath, "jar", "", "Caminho para o assinador.jar")
}
