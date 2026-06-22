package br.ufg.hubsaude.assinador;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.nio.charset.StandardCharsets;
import java.time.Instant;
import java.util.Map;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicLong;

public final class Main {
    private static final SignatureService SERVICE = new FakeSignatureService();

    private Main() {}

    public static void main(String[] args) throws Exception {
        if (args.length == 0 || "help".equals(args[0]) || "--help".equals(args[0])) {
            usage();
            return;
        }
        try {
            switch (args[0]) {
                case "sign" -> localSign(args);
                case "validate" -> localValidate(args);
                case "server" -> startServer(args);
                default -> {
                    System.err.println("Comando desconhecido: " + args[0]);
                    usage();
                    System.exit(1);
                }
            }
        } catch (UserInputException e) {
            System.err.println(Json.error("USER_ERROR", e.getMessage()));
            System.exit(1);
        } catch (Exception e) {
            System.err.println(Json.error("SYSTEM_ERROR", "Falha inesperada: " + e.getMessage()));
            System.exit(4);
        }
    }

    private static void localSign(String[] args) throws Exception {
        Map<String, String> flags = Args.parse(args, 1);
        SignRequest req = new SignRequest(flags.get("bundle"), flags.get("provenance"), flags.get("crypto-material"), flags.get("cert-chain"), flags.get("timestamp"), flags.get("strategy"), flags.get("policy"), flags.get("config"), flags.get("signer"), flags.get("input"));
        System.out.println(SERVICE.sign(req));
    }

    private static void localValidate(String[] args) throws Exception {
        Map<String, String> flags = Args.parse(args, 1);
        ValidateRequest req = new ValidateRequest(flags.get("signature"), flags.get("timestamp"), flags.get("policy"), flags.get("config"), flags.get("bundle"), flags.get("provenance"), flags.get("input"));
        System.out.println(SERVICE.validate(req));
    }

    private static void startServer(String[] args) throws IOException {
        Map<String, String> flags = Args.parse(args, 1);
        int port = parseInt(flags.getOrDefault("port", "8080"), "port");
        int idleMinutes = parseInt(flags.getOrDefault("idle-timeout-minutes", "0"), "idle-timeout-minutes");
        HttpServer server = HttpServer.create(new InetSocketAddress("127.0.0.1", port), 0);
        server.setExecutor(Executors.newCachedThreadPool());
        AtomicLong lastRequest = new AtomicLong(System.currentTimeMillis());

        server.createContext("/health", exchange -> {
            lastRequest.set(System.currentTimeMillis());
            send(exchange, 200, "{\"status\":\"UP\",\"service\":\"assinador.jar\",\"timestamp\":\"" + Instant.now() + "\"}");
        });
        server.createContext("/sign", exchange -> {
            lastRequest.set(System.currentTimeMillis());
            if (!"POST".equalsIgnoreCase(exchange.getRequestMethod())) {
                send(exchange, 405, Json.error("METHOD_NOT_ALLOWED", "Use POST /sign"));
                return;
            }
            try {
                Map<String, String> body = Json.parseFlat(readBody(exchange));
                SignRequest req = new SignRequest(body.get("bundle"), body.get("provenance"), body.get("cryptoMaterial"), body.get("certificateChain"), body.get("timestamp"), body.get("strategy"), body.get("policy"), body.get("config"), body.get("signer"), body.get("input"));
                send(exchange, 200, SERVICE.sign(req));
            } catch (UserInputException e) {
                send(exchange, 400, Json.error("USER_ERROR", e.getMessage()));
            } catch (Exception e) {
                send(exchange, 500, Json.error("SYSTEM_ERROR", e.getMessage()));
            }
        });
        server.createContext("/validate", exchange -> {
            lastRequest.set(System.currentTimeMillis());
            if (!"POST".equalsIgnoreCase(exchange.getRequestMethod())) {
                send(exchange, 405, Json.error("METHOD_NOT_ALLOWED", "Use POST /validate"));
                return;
            }
            try {
                Map<String, String> body = Json.parseFlat(readBody(exchange));
                ValidateRequest req = new ValidateRequest(body.get("signature"), body.get("timestamp"), body.get("policy"), body.get("config"), body.get("bundle"), body.get("provenance"), body.get("input"));
                send(exchange, 200, SERVICE.validate(req));
            } catch (UserInputException e) {
                send(exchange, 400, Json.error("USER_ERROR", e.getMessage()));
            } catch (Exception e) {
                send(exchange, 500, Json.error("SYSTEM_ERROR", e.getMessage()));
            }
        });
        server.createContext("/shutdown", exchange -> {
            send(exchange, 200, "{\"status\":\"SHUTTING_DOWN\"}");
            new Thread(() -> {
                try { Thread.sleep(250); } catch (InterruptedException ignored) { Thread.currentThread().interrupt(); }
                server.stop(1);
                System.exit(0);
            }, "shutdown-thread").start();
        });

        if (idleMinutes > 0) {
            ScheduledExecutorService scheduler = Executors.newSingleThreadScheduledExecutor();
            scheduler.scheduleAtFixedRate(() -> {
                long idleMillis = System.currentTimeMillis() - lastRequest.get();
                if (idleMillis > TimeUnit.MINUTES.toMillis(idleMinutes)) {
                    server.stop(1);
                    System.exit(0);
                }
            }, 30, 30, TimeUnit.SECONDS);
        }
        server.start();
        System.err.printf("assinador.jar server started at http://localhost:%d%n", port);
    }

    private static String readBody(HttpExchange exchange) throws IOException {
        return new String(exchange.getRequestBody().readAllBytes(), StandardCharsets.UTF_8);
    }

    private static void send(HttpExchange exchange, int status, String body) throws IOException {
        byte[] b = body.getBytes(StandardCharsets.UTF_8);
        exchange.getResponseHeaders().set("Content-Type", "application/json; charset=utf-8");
        exchange.sendResponseHeaders(status, b.length);
        try (OutputStream os = exchange.getResponseBody()) {
            os.write(b);
        }
    }

    private static int parseInt(String value, String field) {
        try {
            return Integer.parseInt(value);
        } catch (NumberFormatException e) {
            throw new UserInputException("Parâmetro --" + field + " deve ser numérico");
        }
    }

    private static void usage() {
        System.out.println("""
                assinador.jar - simulação de assinatura digital

                Uso:
                  java -jar assinador.jar sign --bundle <bundle.json> --provenance <provenance.json> --crypto-material <crypto.json> --cert-chain <certs.json> --timestamp <unix> --policy <uri> [--strategy iat]
                  java -jar assinador.jar validate --signature <signature.json> --timestamp <unix> --policy <uri> [--bundle <bundle.json> --provenance <provenance.json>]
                  java -jar assinador.jar server [--port 8080] [--idle-timeout-minutes 10]
                """);
    }
}
