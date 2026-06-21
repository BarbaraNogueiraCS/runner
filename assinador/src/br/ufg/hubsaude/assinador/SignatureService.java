package br.ufg.hubsaude.assinador;

interface SignatureService {
    String sign(SignRequest request) throws Exception;
    String validate(ValidateRequest request) throws Exception;
}
