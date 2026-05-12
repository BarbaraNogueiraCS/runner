package br.ufg.hubsaude.assinador.dto;

public record SignResponse(boolean success, String operation, String signature, String message) {
}
