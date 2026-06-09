package br.ufg.hubsaude.assinador;

/** Mantido apenas para compatibilidade de pacote. A saída atual é FHIR Signature em JSON. */
record SignResult(String operation, String status, String signer, String documentHash, String signatureValue, String signedAt) {}
