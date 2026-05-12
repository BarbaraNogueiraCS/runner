package br.ufg.hubsaude.assinador.validation;

public record ValidationResult(boolean valid, String code, String message, String details) {
    public static ValidationResult ok() {
        return new ValidationResult(true, "OK", "Requisição válida", "");
    }

    public static ValidationResult invalid(String message, String details) {
        return new ValidationResult(false, "INVALID_PARAMETER", message, details);
    }
}
