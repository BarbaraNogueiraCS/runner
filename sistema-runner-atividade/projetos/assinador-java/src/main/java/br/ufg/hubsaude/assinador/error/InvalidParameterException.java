package br.ufg.hubsaude.assinador.error;

public class InvalidParameterException extends RuntimeException {
    private final String details;

    public InvalidParameterException(String message, String details) {
        super(message);
        this.details = details;
    }

    public String details() {
        return details;
    }
}
