package br.ufg.hubsaude.assinador;

interface SignatureService {
    SignResult sign(SignRequest request) throws Exception;
    ValidationResult validate(ValidateRequest request) throws Exception;
}
