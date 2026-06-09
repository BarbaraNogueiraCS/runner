package br.ufg.hubsaude.assinador;

/** Mantido apenas para compatibilidade de pacote. A saída atual é OperationOutcome em JSON. */
record ValidationResult(String operation, String status, boolean valid, String reason, String validatedAt) {}
