package br.ufg.hubsaude.assinador;

record SignResult(String operation, String status, String signer, String documentHash, String signatureValue, String signedAt) {
    String toJson() {
        return "{" +
                Json.prop("operation", operation) + "," +
                Json.prop("status", status) + "," +
                Json.prop("signer", signer) + "," +
                Json.prop("documentHash", documentHash) + "," +
                Json.prop("signatureValue", signatureValue) + "," +
                Json.prop("signedAt", signedAt) +
                "}";
    }
}
