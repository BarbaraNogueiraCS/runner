package br.ufg.hubsaude.assinador;

import java.io.ByteArrayInputStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;
import java.time.Instant;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Set;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

class GuideSignatureService implements SignatureService {
    static final long MIN_TIMESTAMP = 1751328000L; // 2025-07-01T00:00:00Z
    static final long MAX_TIMESTAMP = 4102444800L; // 2100-01-01T00:00:00Z
    static final String POLICY_BASE = "https://fhir.saude.go.gov.br/r4/seguranca/ImplementationGuide/br.go.ses.seguranca";
    static final String DIGEST_ALG_SHA512 = "http://www.w3.org/2001/04/xmlenc#sha512";

    @Override
    public String sign(SignRequest request) throws Exception {
        SignInput input = readSignInput(request);
        validateTimestampForSigning(input.timestamp());
        validatePolicy(input.policy());
        validateStrategy(input.strategy(), input.config());
        validateCryptoMaterial(input.cryptoMaterial(), input.config());
        validateBundleAndProvenance(input.bundle(), input.provenance());
        List<String> chain = validateCertificateChain(input.certificateChain(), false);

        String payload = contentDigestPayload(input.bundle(), input.provenance());
        String protectedHeaderJson = protectedHeader(input.policy(), input.timestamp(), chain);
        String protectedB64 = CryptoUtil.base64Url(CryptoUtil.utf8(protectedHeaderJson));
        String signingInput = protectedB64 + "." + payload;
        String signature = simulatedSignature(signingInput);
        String digestValue = CryptoUtil.base64(CryptoUtil.sha512(CryptoUtil.utf8(String.join("", chain))));
        String jws = jwsJson(payload, protectedB64, signature, digestValue, "iat".equals(input.strategy()));
        String signatureData = CryptoUtil.base64(CryptoUtil.utf8(jws));
        String when = Instant.ofEpochSecond(input.timestamp()).toString();
        String signer = input.signer().isBlank() ? "Signatário simulado" : input.signer();

        return "{"
                + Json.prop("resourceType", "Signature") + ","
                + Json.rawProp("type", "[{" + Json.prop("system", "urn:iso-astm:E1762-95:2013") + "," + Json.prop("code", "1.2.840.10065.1.12.1.1") + "," + Json.prop("display", "Author's Signature") + "}]") + ","
                + Json.prop("when", when) + ","
                + Json.rawProp("who", "{" + Json.prop("display", signer) + "}") + ","
                + Json.prop("sigFormat", "application/jose+json") + ","
                + Json.prop("data", signatureData)
                + "}";
    }

    @Override
    public String validate(ValidateRequest request) throws Exception {
        ValidateInput input = readValidateInput(request);
        validateTimestamp(input.timestamp());
        validatePolicy(input.policy());

        String signatureText = readFile(input.signature(), "signature");
        String data = extractSignatureData(signatureText);
        String jws = CryptoUtil.utf8(CryptoUtil.base64Decode(data, "FORMAT.SIGNATURE-DATA-INVALID"));
        String payload = requiredField(jws, "payload", "FORMAT.JWS-MALFORMED");
        String protectedB64 = requiredField(jws, "protected", "FORMAT.JWS-MALFORMED");
        String signature = requiredField(jws, "signature", "FORMAT.JWS-MALFORMED");
        if (Json.objectFieldRaw(jws, "header") == null) {
            throw new UserInputException("VALIDATION.LTV-EVIDENCE-INVALID: unprotected header com rRefs é obrigatório");
        }

        String headerJson = CryptoUtil.utf8(CryptoUtil.base64UrlDecode(protectedB64, "FORMAT.JWS-MALFORMED"));
        String alg = requiredField(headerJson, "alg", "FORMAT.JWS-MALFORMED");
        if (!alg.equals("RS256") && !alg.equals("ES256")) {
            throw new UserInputException("VALIDATION.UNSUPPORTED-ALGORITHM: algoritmo deve ser RS256 ou ES256");
        }
        List<String> chain = Json.stringArrayField(headerJson, "x5c");
        if (chain.isEmpty()) {
            throw new UserInputException("CERT.INVALID-FORMAT: protected header deve conter x5c não vazio");
        }
        validateCertificateChain("[\"" + String.join("\",\"", chain) + "\"]", true);

        String sigPId = Json.objectFieldRaw(headerJson, "sigPId");
        if (sigPId == null) {
            throw new UserInputException("POLICY.VERSION-UNSUPPORTED: protected header deve conter sigPId");
        }
        String policyInHeader = requiredField(sigPId, "id", "POLICY.VERSION-UNSUPPORTED");
        if (!policyInHeader.equals(input.policy())) {
            throw new UserInputException("POLICY.VERSION-UNSUPPORTED: política da assinatura não corresponde à política informada");
        }
        long iat = Json.longField(headerJson, "iat", -1);
        List<String> warnings = new ArrayList<>();
        if (iat < 0) {
            throw new UserInputException("TEMPORAL.IAT-INVALID: protected header deve conter iat na estratégia simulada");
        }
        validateTimestamp(iat);
        if (iat > input.timestamp()) {
            throw new UserInputException("TEMPORAL.IAT-INVALID: iat não pode ser posterior ao timestamp de referência");
        }
        if (Math.abs(iat - input.timestamp()) > 300) {
            warnings.add("TEMPORAL.CLOCK-SKEW-DETECTED");
        }

        if (!signature.equals(simulatedSignature(protectedB64 + "." + payload))) {
            throw new UserInputException("VALIDATION.SIGNATURE-VERIFICATION-FAILED: assinatura simulada não confere");
        }

        if (!input.bundle().isBlank() || !input.provenance().isBlank()) {
            if (input.bundle().isBlank() || input.provenance().isBlank()) {
                throw new UserInputException("FORMAT.CONTENT-INTEGRITY-INPUTS-INCOMPLETE: informe Bundle e Provenance para verificação de integridade");
            }
            String bundle = readFile(input.bundle(), "bundle");
            String provenance = readFile(input.provenance(), "provenance");
            validateBundleAndProvenance(bundle, provenance);
            String recomputed = contentDigestPayload(bundle, provenance);
            if (!constantTimeEquals(recomputed, payload)) {
                throw new UserInputException("VALIDATION.PAYLOAD-DIGEST-MISMATCH: integridade do conteúdo assinado não confirmada");
            }
        }

        return Json.operationOutcomeSuccess("Assinatura simulada validada conforme estrutura JWS JSON Serialization e FHIR Signature", warnings);
    }

    private SignInput readSignInput(SignRequest request) throws Exception {
        if (notBlank(request.input()) && isBlank(request.bundle())) {
            return legacyInput(request);
        }
        String bundle = readFile(request.bundle(), "bundle");
        String provenance = readFile(request.provenance(), "provenance");
        String crypto = readFile(request.cryptoMaterial(), "crypto-material");
        String chain = readFile(request.certificateChain(), "cert-chain");
        long timestamp = parseLong(required(request.timestamp(), "timestamp"), "TEMPORAL.REFERENCE-TIMESTAMP-INVALID");
        String strategy = defaultText(request.strategy(), "iat");
        String policy = required(request.policy(), "policy");
        String config = readOptionalFile(request.config());
        String signer = defaultText(request.signer(), "");
        return new SignInput(bundle, provenance, crypto, chain, timestamp, strategy, policy, config, signer);
    }

    private SignInput legacyInput(SignRequest request) throws Exception {
        String text = readFile(request.input(), "input");
        long now = Instant.now().getEpochSecond();
        String uuid = "11111111-1111-4111-8111-111111111111";
        String bundle = "{" + Json.prop("resourceType", "Bundle") + "," + Json.prop("type", "collection") + ","
                + Json.rawProp("entry", "[{" + Json.prop("fullUrl", "urn:uuid:" + uuid) + "," + Json.rawProp("resource", "{" + Json.prop("resourceType", "Binary") + "," + Json.prop("contentType", "text/plain") + "," + Json.prop("data", CryptoUtil.base64(text.getBytes(StandardCharsets.UTF_8))) + "}") + "}]") + "}";
        String provenance = "{" + Json.prop("resourceType", "Provenance") + "," + Json.rawProp("target", "[{" + Json.prop("reference", "urn:uuid:" + uuid) + "}]") + "}";
        String crypto = "{" + Json.prop("type", "PEM") + "," + Json.prop("simulation", "true") + "}";
        String chain = "[\"" + sampleCertificateBase64() + "\"]";
        return new SignInput(bundle, provenance, crypto, chain, now, "iat", POLICY_BASE + "|0.1.2", "{}", defaultText(request.signer(), "Signatário simulado"));
    }

    private ValidateInput readValidateInput(ValidateRequest request) {
        long timestamp = request.timestamp() == null || request.timestamp().isBlank()
                ? Instant.now().getEpochSecond()
                : parseLong(request.timestamp(), "TEMPORAL.REFERENCE-TIMESTAMP-INVALID");
        String policy = defaultText(request.policy(), POLICY_BASE + "|0.1.2");
        return new ValidateInput(required(request.signature(), "signature"), timestamp, policy, defaultText(request.config(), ""), defaultText(request.bundle(), defaultText(request.input(), "")), defaultText(request.provenance(), ""));
    }

    private void validateBundleAndProvenance(String bundle, String provenance) {
        if (!hasStringField(bundle, "resourceType", "Bundle")) {
            throw new UserInputException("FORMAT.BUNDLE-INVALID: recurso deve ser Bundle FHIR R4 em JSON");
        }
        if (!hasStringField(provenance, "resourceType", "Provenance")) {
            throw new UserInputException("FORMAT.PROVENANCE-INVALID: recurso deve ser Provenance FHIR R4 em JSON");
        }
        List<String> fullUrls = extractStringFields(bundle, "fullUrl");
        if (fullUrls.isEmpty()) {
            throw new UserInputException("FORMAT.BUNDLE-MALFORMED: Bundle.entry deve conter fullUrl");
        }
        requireUnique(fullUrls, "FORMAT.BUNDLE-MALFORMED: Bundle não deve conter fullUrl duplicado");
        List<String> targets = extractReferences(provenance);
        if (targets.isEmpty()) {
            throw new UserInputException("FORMAT.PROVENANCE-INVALID: Provenance.target deve conter pelo menos uma referência");
        }
        requireUnique(targets, "FORMAT.PROVENANCE-TARGET-DUPLICATE: Provenance.target não deve conter referências duplicadas");
        Set<String> available = new HashSet<>(fullUrls);
        Pattern uuid = Pattern.compile("^urn:uuid:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$");
        for (String target : targets) {
            if (!uuid.matcher(target).matches()) {
                throw new UserInputException("FORMAT.PROVENANCE-TARGET-INVALID: referência deve usar urn:uuid:UUID RFC 4122");
            }
            if (!available.contains(target)) {
                throw new UserInputException("FORMAT.PROVENANCE-TARGET-NOT-FOUND: referência não existe em Bundle.entry.fullUrl");
            }
        }
    }

    private void validateCryptoMaterial(String crypto, String config) throws Exception {
        String type = Json.stringField(crypto, "type");
        if (type == null || type.isBlank()) {
            throw new UserInputException("CRYPTO.MATERIAL-INVALID: material criptográfico deve informar type");
        }
        switch (type) {
            case "PEM", "PKCS12", "SMARTCARD", "TOKEN", "REMOTE" -> { }
            default -> throw new UserInputException("CRYPTO.MATERIAL-INVALID: type deve ser PEM, PKCS12, SMARTCARD, TOKEN ou REMOTE");
        }
        if ((type.equals("SMARTCARD") || type.equals("TOKEN")) && isBlank(Json.stringField(crypto, "pin"))) {
            throw new UserInputException("CRYPTO.PIN-REQUIRED: PIN obrigatório para SMARTCARD/TOKEN");
        }
        if (type.equals("SMARTCARD") || type.equals("TOKEN")) {
            Pkcs11Support.verify(crypto, config);
        }
    }

    private List<String> validateCertificateChain(String certificateChainJson, boolean fullValidation) throws Exception {
        List<String> chain = Json.stringArray(certificateChainJson);
        if (chain.isEmpty()) {
            throw new UserInputException("CERT.CHAIN-INCOMPLETE: cadeia de certificados vazia");
        }
        CertificateFactory factory = CertificateFactory.getInstance("X.509");
        for (String certBase64 : chain) {
            byte[] der = CryptoUtil.base64Decode(certBase64, "CERT.BASE64-INVALID");
            try {
                factory.generateCertificate(new ByteArrayInputStream(der));
            } catch (Exception e) {
                if (fullValidation) {
                    throw new UserInputException("CERT.INVALID-FORMAT: certificado não é X.509 DER válido");
                }
                // Em modo simulado, aceita certificado base64 DER sintaticamente inválido apenas para permitir testes sem AC real.
            }
        }
        return chain;
    }

    private String contentDigestPayload(String bundle, String provenance) throws Exception {
        String canonical = Json.canonicalizeLoose(bundle) + Json.canonicalizeLoose(provenance);
        return CryptoUtil.base64Url(CryptoUtil.sha256(CryptoUtil.utf8(canonical)));
    }

    private String protectedHeader(String policy, long timestamp, List<String> chain) {
        StringBuilder x5c = new StringBuilder("[");
        for (int i = 0; i < chain.size(); i++) {
            if (i > 0) x5c.append(',');
            x5c.append(Json.quote(chain.get(i)));
        }
        x5c.append(']');
        return "{" + Json.prop("alg", "RS256") + ","
                + Json.rawProp("x5c", x5c.toString()) + ","
                + Json.rawProp("sigPId", "{" + Json.prop("id", policy) + "}") + ","
                + Json.rawProp("iat", Long.toString(timestamp))
                + "}";
    }

    private String jwsJson(String payload, String protectedB64, String signature, String digestValue, boolean iatStrategy) {
        String rRefs = "{" + Json.rawProp("ocspRefs", "[]") + "," + Json.rawProp("crlRefs", "[{" + Json.prop("digestAlg", DIGEST_ALG_SHA512) + "," + Json.prop("digestValue", digestValue) + "}]") + "}";
        String header = "{" + Json.rawProp("rRefs", rRefs) + "," + Json.rawProp("simulation", "true") + "}";
        return "{" + Json.prop("payload", payload) + ","
                + Json.rawProp("signatures", "[{" + Json.prop("protected", protectedB64) + "," + Json.rawProp("header", header) + "," + Json.prop("signature", signature) + "}]")
                + "}";
    }

    private String simulatedSignature(String signingInput) throws Exception {
        return CryptoUtil.base64Url(CryptoUtil.sha256(CryptoUtil.utf8(signingInput + ".RUNNER-SIMULATED-JADES")));
    }

    private String extractSignatureData(String signatureText) {
        String data = Json.stringField(signatureText, "data");
        if (data != null && !data.isBlank()) {
            return data;
        }
        String trimmed = signatureText.trim();
        if (trimmed.startsWith("{") && trimmed.contains("OperationOutcome")) {
            throw new UserInputException("FORMAT.SIGNATURE-DATA-INVALID: arquivo contém OperationOutcome, não Signature.data");
        }
        return trimmed;
    }

    private void validateTimestampForSigning(long timestamp) {
        validateTimestamp(timestamp);
        long now = Instant.now().getEpochSecond();
        if (Math.abs(now - timestamp) > 300) {
            throw new UserInputException("TEMPORAL.REFERENCE-TIMESTAMP-OUT-OF-CURRENT-WINDOW: timestamp deve estar a ±5 minutos do UTC atual para criação");
        }
    }

    private void validateTimestamp(long timestamp) {
        if (timestamp < MIN_TIMESTAMP || timestamp > MAX_TIMESTAMP) {
            throw new UserInputException("TEMPORAL.REFERENCE-TIMESTAMP-OUT-OF-RANGE: timestamp fora do intervalo permitido");
        }
    }

    private void validatePolicy(String policy) {
        if (!policy.matches(Pattern.quote(POLICY_BASE) + "\\|\\d+\\.\\d+\\.\\d+")) {
            throw new UserInputException("POLICY.ID-INVALID: política deve seguir {baseUri}|major.minor.patch");
        }
    }

    private void validateStrategy(String strategy, String config) {
        if (!strategy.equals("iat") && !strategy.equals("tsa")) {
            throw new UserInputException("CONFIG.INVALID-PARAMETER: estratégia deve ser iat ou tsa");
        }
        if (strategy.equals("tsa")) {
            String tsaUrl = Json.stringField(config, "tsaUrl");
            if (tsaUrl == null || !tsaUrl.startsWith("https://")) {
                throw new UserInputException("CONFIG.TSA-URL-INVALID: estratégia tsa exige tsaUrl HTTPS nas configurações");
            }
        }
    }

    private static boolean constantTimeEquals(String a, String b) {
        if (a == null || b == null) return false;
        int diff = a.length() ^ b.length();
        for (int i = 0; i < Math.min(a.length(), b.length()); i++) {
            diff |= a.charAt(i) ^ b.charAt(i);
        }
        return diff == 0;
    }

    private static String readFile(String path, String field) throws Exception {
        String p = required(path, field);
        Path file = Path.of(p);
        if (!Files.isRegularFile(file)) {
            throw new UserInputException("INPUT.FILE-NOT-FOUND: arquivo --" + field + " não encontrado: " + p);
        }
        return Files.readString(file, StandardCharsets.UTF_8);
    }

    private static String readOptionalFile(String path) throws Exception {
        if (isBlank(path)) return "{}";
        return readFile(path, "config");
    }

    private static String required(String value, String field) {
        if (isBlank(value)) {
            throw new UserInputException("INPUT.MISSING-PARAMETER: informe --" + field);
        }
        return value;
    }

    private static String requiredField(String json, String field, String code) {
        String value = Json.stringField(json, field);
        if (value == null || value.isBlank()) {
            throw new UserInputException(code + ": campo obrigatório ausente: " + field);
        }
        return value;
    }

    private static long parseLong(String value, String code) {
        try {
            return Long.parseLong(value);
        } catch (NumberFormatException e) {
            throw new UserInputException(code + ": timestamp deve ser inteiro Unix UTC");
        }
    }

    private static boolean hasStringField(String json, String field, String expected) {
        return expected.equals(Json.stringField(json, field));
    }

    private static List<String> extractStringFields(String json, String field) {
        List<String> values = new ArrayList<>();
        Matcher m = Pattern.compile("\\\"" + Pattern.quote(field) + "\\\"\\s*:\\s*\\\"((?:\\\\.|[^\\\"])*)\\\"").matcher(json);
        while (m.find()) {
            values.add(m.group(1));
        }
        return values;
    }

    private static List<String> extractReferences(String provenance) {
        return extractStringFields(provenance, "reference");
    }

    private static void requireUnique(List<String> values, String message) {
        Set<String> seen = new HashSet<>();
        for (String value : values) {
            if (!seen.add(value)) {
                throw new UserInputException(message);
            }
        }
    }

    private static boolean isBlank(String value) {
        return value == null || value.isBlank();
    }

    private static boolean notBlank(String value) {
        return !isBlank(value);
    }

    private static String defaultText(String value, String def) {
        return isBlank(value) ? def : value;
    }

    private static String sampleCertificateBase64() {
        // DER sintático simplificado em base64. No fluxo novo, informe --cert-chain com certificados reais.
        return "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsimulatedCertificateForRunnerOnly000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000IDAQAB";
    }

    private record SignInput(String bundle, String provenance, String cryptoMaterial, String certificateChain, long timestamp, String strategy, String policy, String config, String signer) {}
    private record ValidateInput(String signature, long timestamp, String policy, String config, String bundle, String provenance) {}
}
