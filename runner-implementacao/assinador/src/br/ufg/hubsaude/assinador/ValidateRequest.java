package br.ufg.hubsaude.assinador;

record ValidateRequest(
        String signature,
        String timestamp,
        String policy,
        String config,
        String bundle,
        String provenance,
        String input
) {}
