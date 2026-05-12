package br.ufg.hubsaude.assinador.dto;

import java.util.Map;

public record SignRequest(String document, String certificate, Map<String, String> parameters) {
}
