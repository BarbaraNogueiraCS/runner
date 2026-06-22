package br.ufg.hubsaude.assinador;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;

import java.io.IOException;
import java.io.OutputStream;
import java.nio.charset.StandardCharsets;
import java.time.Instant;
import java.util.Map;
import java.util.concurrent.atomic.AtomicLong;

/**
 * Controller HTTP do assinador.jar.
 *
 * Expõe os endpoints exigidos na Sprint 3:
 *   POST /sign
 *   POST /validate
 *
 * A classe não implementa regra de negócio própria. Ela apenas transforma HTTP/JSON
 * em SignRequest/ValidateRequest e delega para SignatureService, garantindo que o
 * modo servidor reutilize exatamente a mesma validação e simulação do modo CLI.
 */
final class SignatureController {
    private final SignatureService service;
    private final AtomicLong lastRequestMillis;

    SignatureController(SignatureService service, AtomicLong lastRequestMillis) {
        this.service = service;
        this.lastRequestMillis = lastRequestMillis;
    }

    void register(HttpServer server) {
        server.createContext("/health", this::health);
        server.createContext("/sign", this::sign);
        server.createContext("/validate", this::validate);
    }

    private void health(HttpExchange exchange) throws IOException {
        touch();
        if (!"GET".equalsIgnoreCase(exchange.getRequestMethod())) {
            send(exchange, 405, Json.error("METHOD_NOT_ALLOWED", "Use GET /health"));
            return;
        }
        send(exchange, 200, "{" + Json.prop("status", "UP") + ","
                + Json.prop("service", "assinador.jar") + ","
                + Json.prop("timestamp", Instant.now().toString()) + "}");
    }

    private void sign(HttpExchange exchange) throws IOException {
        touch();
        if (!"POST".equalsIgnoreCase(exchange.getRequestMethod())) {
            send(exchange, 405, Json.error("METHOD_NOT_ALLOWED", "Use POST /sign"));
            return;
        }
        try {
            Map<String, String> body = Json.parseFlat(readBody(exchange));
            SignRequest req = new SignRequest(
                    body.get("bundle"),
                    body.get("provenance"),
                    body.get("cryptoMaterial"),
                    body.get("certificateChain"),
                    body.get("timestamp"),
                    body.get("strategy"),
                    body.get("policy"),
                    body.get("config"),
                    body.get("signer"),
                    body.get("input")
            );
            send(exchange, 200, service.sign(req));
        } catch (UserInputException e) {
            send(exchange, 400, Json.error("USER_ERROR", e.getMessage()));
        } catch (Exception e) {
            send(exchange, 500, Json.error("SYSTEM_ERROR", e.getMessage()));
        }
    }

    private void validate(HttpExchange exchange) throws IOException {
        touch();
        if (!"POST".equalsIgnoreCase(exchange.getRequestMethod())) {
            send(exchange, 405, Json.error("METHOD_NOT_ALLOWED", "Use POST /validate"));
            return;
        }
        try {
            Map<String, String> body = Json.parseFlat(readBody(exchange));
            ValidateRequest req = new ValidateRequest(
                    body.get("signature"),
                    body.get("timestamp"),
                    body.get("policy"),
                    body.get("config"),
                    body.get("bundle"),
                    body.get("provenance"),
                    body.get("input")
            );
            send(exchange, 200, service.validate(req));
        } catch (UserInputException e) {
            send(exchange, 400, Json.error("USER_ERROR", e.getMessage()));
        } catch (Exception e) {
            send(exchange, 500, Json.error("SYSTEM_ERROR", e.getMessage()));
        }
    }

    private void touch() {
        lastRequestMillis.set(System.currentTimeMillis());
    }

    private static String readBody(HttpExchange exchange) throws IOException {
        return new String(exchange.getRequestBody().readAllBytes(), StandardCharsets.UTF_8);
    }

    static void send(HttpExchange exchange, int status, String body) throws IOException {
        byte[] b = body.getBytes(StandardCharsets.UTF_8);
        exchange.getResponseHeaders().set("Content-Type", "application/json; charset=utf-8");
        exchange.sendResponseHeaders(status, b.length);
        try (OutputStream os = exchange.getResponseBody()) {
            os.write(b);
        }
    }
}
