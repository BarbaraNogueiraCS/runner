package br.ufg.hubsaude.assinador;

import java.nio.file.Files;
import java.nio.file.Path;
import java.security.KeyStore;
import java.security.Provider;
import java.security.Security;
import java.util.Enumeration;

/**
 * Integração mínima com PKCS#11 para SMARTCARD/TOKEN.
 *
 * O guia SES-GO trata PIN, identificador, slot e token label como entrada do
 * material criptográfico, enquanto biblioteca/middleware PKCS#11 são detalhes
 * operacionais. Por isso, esta classe aceita simulação explícita para testes
 * acadêmicos e, quando uma biblioteca PKCS#11 é informada, inicializa o provider
 * SunPKCS11, abre o token com PIN e verifica a existência do alias/identificador.
 */
final class Pkcs11Support {
    private Pkcs11Support() {}

    static void verify(String cryptoJson, String configJson) throws Exception {
        if (isSimulation(cryptoJson) || isSimulation(configJson)) {
            return;
        }
        String library = firstNonBlank(
                Json.stringField(cryptoJson, "pkcs11Library"),
                Json.stringField(cryptoJson, "library"),
                Json.stringField(configJson, "pkcs11Library"),
                Json.stringField(configJson, "library")
        );
        if (isBlank(library)) {
            throw new UserInputException("PKCS11.LIBRARY-REQUIRED: SMARTCARD/TOKEN exige pkcs11Library no material criptográfico ou nas configurações operacionais, exceto quando simulation=true");
        }
        Path libraryPath = Path.of(library);
        if (!Files.isRegularFile(libraryPath)) {
            throw new UserInputException("PKCS11.LIBRARY-NOT-FOUND: biblioteca PKCS#11 não encontrada: " + library);
        }
        String pin = firstNonBlank(Json.stringField(cryptoJson, "pin"), Json.stringField(cryptoJson, "PIN"));
        if (isBlank(pin)) {
            throw new UserInputException("PKCS11.PIN-REQUIRED: PIN obrigatório para abrir token/smartcard PKCS#11");
        }
        String identifier = firstNonBlank(
                Json.stringField(cryptoJson, "identifier"),
                Json.stringField(cryptoJson, "identificador"),
                Json.stringField(cryptoJson, "alias")
        );
        String slot = firstNonBlank(Json.stringField(cryptoJson, "slotId"), Json.stringField(cryptoJson, "slot"));
        String tokenLabel = firstNonBlank(Json.stringField(cryptoJson, "tokenLabel"), Json.stringField(cryptoJson, "token"));

        Path cfg = Files.createTempFile("runner-pkcs11-", ".cfg");
        try {
            StringBuilder text = new StringBuilder();
            text.append("name=RunnerPKCS11\n");
            text.append("library=").append(libraryPath.toAbsolutePath()).append('\n');
            if (!isBlank(slot)) {
                text.append("slot=").append(slot).append('\n');
            }
            if (!isBlank(tokenLabel)) {
                text.append("attributes=compatibility\n");
            }
            Files.writeString(cfg, text.toString());

            Provider base = Security.getProvider("SunPKCS11");
            if (base == null) {
                throw new UserInputException("PKCS11.PROVIDER-NOT-AVAILABLE: provider SunPKCS11 indisponível; use um JDK/JRE com módulo jdk.crypto.cryptoki");
            }
            Provider configured = base.configure(cfg.toString());
            Security.addProvider(configured);
            KeyStore ks = KeyStore.getInstance("PKCS11", configured);
            ks.load(null, pin.toCharArray());
            if (!isBlank(identifier)) {
                if (!ks.containsAlias(identifier)) {
                    throw new UserInputException("PKCS11.KEY-NOT-FOUND: identificador/alias não encontrado no token: " + identifier);
                }
                return;
            }
            Enumeration<String> aliases = ks.aliases();
            if (!aliases.hasMoreElements()) {
                throw new UserInputException("PKCS11.KEY-NOT-FOUND: token aberto, mas nenhuma chave/certificado foi encontrado");
            }
        } finally {
            Files.deleteIfExists(cfg);
        }
    }

    private static boolean isSimulation(String json) {
        String v = firstNonBlank(Json.stringField(json, "simulation"), Json.stringField(json, "simulated"));
        return "true".equalsIgnoreCase(v) || "simulado".equalsIgnoreCase(v);
    }

    private static String firstNonBlank(String... values) {
        if (values == null) return "";
        for (String v : values) {
            if (!isBlank(v)) return v;
        }
        return "";
    }

    private static boolean isBlank(String value) {
        return value == null || value.isBlank();
    }
}
