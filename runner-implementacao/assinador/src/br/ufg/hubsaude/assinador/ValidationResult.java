package br.ufg.hubsaude.assinador;

record ValidationResult(String operation, String status, boolean valid, String reason, String validatedAt) {
    String toJson() {
        return "{" +
                Json.prop("operation", operation) + "," +
                Json.prop("status", status) + "," +
                "\"valid\":" + valid + "," +
                Json.prop("reason", reason) + "," +
                Json.prop("validatedAt", validatedAt) +
                "}";
    }
}
