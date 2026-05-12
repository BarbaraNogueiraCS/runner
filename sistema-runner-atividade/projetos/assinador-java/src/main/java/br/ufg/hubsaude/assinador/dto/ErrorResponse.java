package br.ufg.hubsaude.assinador.dto;

public record ErrorResponse(boolean success, ErrorBody error) {
    public record ErrorBody(String code, String message, String details) {}
}
