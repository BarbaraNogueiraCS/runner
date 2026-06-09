package br.ufg.hubsaude.assinador;

import java.nio.file.Files;
import java.nio.file.Path;
import java.security.MessageDigest;
import java.time.Instant;
import java.util.HexFormat;

final class FakeSignatureService implements SignatureService {
    @Override
    public SignResult sign(SignRequest request) throws Exception {
        requireText(request.input(), "input");
        requireText(request.signer(), "signer");
        Path input = Path.of(request.input());
        if (!Files.isRegularFile(input)) {
            throw new UserInputException("Arquivo de entrada não encontrado: " + request.input());
        }
        byte[] content = Files.readAllBytes(input);
        String hash = sha256(content);
        String value = "SIMULATED-SIGNATURE-" + hash.substring(0, 24);
        return new SignResult("SIGN", "SUCCESS", request.signer(), hash, value, Instant.now().toString());
    }

    @Override
    public ValidationResult validate(ValidateRequest request) throws Exception {
        requireText(request.signature(), "signature");
        Path signature = Path.of(request.signature());
        if (!Files.isRegularFile(signature)) {
            throw new UserInputException("Arquivo de assinatura não encontrado: " + request.signature());
        }
        String text = Files.readString(signature);
        boolean hasSimulatedSignature = text.contains("SIMULATED-SIGNATURE-");
        String reason = hasSimulatedSignature
                ? "Assinatura simulada reconhecida pelo assinador.jar"
                : "O arquivo informado não contém uma assinatura simulada reconhecida";
        return new ValidationResult("VALIDATE", "SUCCESS", hasSimulatedSignature, reason, Instant.now().toString());
    }

    private static void requireText(String value, String field) {
        if (value == null || value.isBlank()) {
            throw new UserInputException("Parâmetro obrigatório ausente ou vazio: --" + field);
        }
    }

    private static String sha256(byte[] content) throws Exception {
        MessageDigest digest = MessageDigest.getInstance("SHA-256");
        return HexFormat.of().formatHex(digest.digest(content));
    }
}
