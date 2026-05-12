package br.ufg.hubsaude.assinador.validation;

import br.ufg.hubsaude.assinador.dto.SignRequest;
import br.ufg.hubsaude.assinador.dto.ValidateRequest;

public class ParameterValidator {
    public ValidationResult validateSign(SignRequest request) {
        if (request == null) {
            return ValidationResult.invalid("Requisição ausente", "O corpo da requisição não foi informado.");
        }
        if (isBlank(request.document())) {
            return ValidationResult.invalid("Parâmetro obrigatório ausente", "Informe o documento a ser assinado.");
        }
        if (isBlank(request.certificate())) {
            return ValidationResult.invalid("Parâmetro obrigatório ausente", "Informe o certificado usado na assinatura simulada.");
        }
        return ValidationResult.ok();
    }

    public ValidationResult validateValidate(ValidateRequest request) {
        if (request == null) {
            return ValidationResult.invalid("Requisição ausente", "O corpo da requisição não foi informado.");
        }
        if (isBlank(request.document())) {
            return ValidationResult.invalid("Parâmetro obrigatório ausente", "Informe o documento associado à assinatura.");
        }
        if (isBlank(request.signature())) {
            return ValidationResult.invalid("Parâmetro obrigatório ausente", "Informe a assinatura a ser validada.");
        }
        return ValidationResult.ok();
    }

    private boolean isBlank(String value) {
        return value == null || value.trim().isEmpty();
    }
}
