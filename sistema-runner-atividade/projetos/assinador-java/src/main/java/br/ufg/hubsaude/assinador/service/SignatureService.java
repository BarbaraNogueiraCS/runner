package br.ufg.hubsaude.assinador.service;

import br.ufg.hubsaude.assinador.dto.SignRequest;
import br.ufg.hubsaude.assinador.dto.SignResponse;
import br.ufg.hubsaude.assinador.dto.ValidateRequest;
import br.ufg.hubsaude.assinador.dto.ValidateResponse;

public interface SignatureService {
    SignResponse sign(SignRequest request);
    ValidateResponse validate(ValidateRequest request);
}
