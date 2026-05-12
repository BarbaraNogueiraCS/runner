package br.ufg.hubsaude.assinador.service;

import br.ufg.hubsaude.assinador.dto.SignRequest;
import br.ufg.hubsaude.assinador.dto.SignResponse;
import br.ufg.hubsaude.assinador.dto.ValidateRequest;
import br.ufg.hubsaude.assinador.dto.ValidateResponse;
import br.ufg.hubsaude.assinador.error.InvalidParameterException;
import br.ufg.hubsaude.assinador.validation.ParameterValidator;

import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.time.Instant;
import java.util.HexFormat;

public class FakeSignatureService implements SignatureService {
    private final ParameterValidator validator;

    public FakeSignatureService(ParameterValidator validator) {
        this.validator = validator;
    }

    @Override
    public SignResponse sign(SignRequest request) {
        var result = validator.validateSign(request);
        if (!result.valid()) {
            throw new InvalidParameterException(result.message(), result.details());
        }
        String signature = "assinatura-simulada-" + digest(request.document() + ":" + request.certificate()).substring(0, 16);
        return new SignResponse(true, "sign", signature, "Assinatura simulada criada com sucesso em " + Instant.now() + ".");
    }

    @Override
    public ValidateResponse validate(ValidateRequest request) {
        var result = validator.validateValidate(request);
        if (!result.valid()) {
            throw new InvalidParameterException(result.message(), result.details());
        }
        boolean valid = request.signature() != null && request.signature().startsWith("assinatura-simulada-");
        String message = valid ? "Assinatura simulada considerada válida." : "Assinatura simulada considerada inválida.";
        return new ValidateResponse(true, "validate", valid, message);
    }

    private String digest(String input) {
        try {
            MessageDigest md = MessageDigest.getInstance("SHA-256");
            byte[] data = md.digest(input.getBytes(StandardCharsets.UTF_8));
            return HexFormat.of().formatHex(data);
        } catch (Exception e) {
            throw new IllegalStateException("Não foi possível gerar assinatura simulada", e);
        }
    }
}
