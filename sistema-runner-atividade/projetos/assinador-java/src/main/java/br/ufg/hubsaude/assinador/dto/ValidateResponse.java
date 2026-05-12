package br.ufg.hubsaude.assinador.dto;

public record ValidateResponse(boolean success, String operation, boolean valid, String message) {
}
