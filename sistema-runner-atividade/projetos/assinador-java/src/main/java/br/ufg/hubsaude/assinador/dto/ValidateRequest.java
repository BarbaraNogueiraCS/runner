package br.ufg.hubsaude.assinador.dto;

import java.util.Map;

public record ValidateRequest(String document, String signature, String certificate, Map<String, String> parameters) {
}
