package br.ufg.hubsaude.assinador.validation;

import br.ufg.hubsaude.assinador.dto.SignRequest;
import br.ufg.hubsaude.assinador.dto.ValidateRequest;
import org.junit.jupiter.api.Test;

import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;

class ParameterValidatorTest {
    private final ParameterValidator validator = new ParameterValidator();

    @Test
    void deveAceitarSignRequestValido() {
        var result = validator.validateSign(new SignRequest("documento.json", "certificado.pem", Map.of()));
        assertTrue(result.valid());
    }

    @Test
    void deveRejeitarSignRequestSemDocumento() {
        var result = validator.validateSign(new SignRequest("", "certificado.pem", Map.of()));
        assertFalse(result.valid());
    }

    @Test
    void deveRejeitarValidateRequestSemAssinatura() {
        var result = validator.validateValidate(new ValidateRequest("documento.json", "", "", Map.of()));
        assertFalse(result.valid());
    }
}
