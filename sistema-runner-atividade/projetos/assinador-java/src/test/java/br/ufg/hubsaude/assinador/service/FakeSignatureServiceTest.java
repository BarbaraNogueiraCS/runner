package br.ufg.hubsaude.assinador.service;

import br.ufg.hubsaude.assinador.dto.SignRequest;
import br.ufg.hubsaude.assinador.dto.ValidateRequest;
import br.ufg.hubsaude.assinador.validation.ParameterValidator;
import org.junit.jupiter.api.Test;

import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;

class FakeSignatureServiceTest {
    private final FakeSignatureService service = new FakeSignatureService(new ParameterValidator());

    @Test
    void deveCriarAssinaturaSimulada() {
        var response = service.sign(new SignRequest("documento.json", "certificado.pem", Map.of()));
        assertTrue(response.success());
        assertTrue(response.signature().startsWith("assinatura-simulada-"));
    }

    @Test
    void deveValidarAssinaturaSimulada() {
        var response = service.validate(new ValidateRequest("documento.json", "assinatura-simulada-abc", "", Map.of()));
        assertTrue(response.success());
        assertTrue(response.valid());
    }
}
