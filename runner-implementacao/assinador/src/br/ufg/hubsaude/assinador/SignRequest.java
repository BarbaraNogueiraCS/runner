package br.ufg.hubsaude.assinador;

record SignRequest(
        String bundle,
        String provenance,
        String cryptoMaterial,
        String certificateChain,
        String timestamp,
        String strategy,
        String policy,
        String config,
        String signer,
        String input
) {}
