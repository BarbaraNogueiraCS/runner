package br.ufg.hubsaude.assinador.controller;

import br.ufg.hubsaude.assinador.dto.ErrorResponse;
import br.ufg.hubsaude.assinador.error.InvalidParameterException;
import br.ufg.hubsaude.assinador.service.SignatureService;
import br.ufg.hubsaude.assinador.util.JsonUtil;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.nio.charset.StandardCharsets;
import java.util.concurrent.Executors;

public class SignatureHttpServer {
    private final SignatureService service;
    private HttpServer server;

    public SignatureHttpServer(SignatureService service) {
        this.service = service;
    }

    public void start(int port) throws IOException {
        server = HttpServer.create(new InetSocketAddress(port), 0);
        server.createContext("/health", this::health);
        server.createContext("/sign", this::sign);
        server.createContext("/validate", this::validate);
        server.createContext("/shutdown", this::shutdown);
        server.setExecutor(Executors.newVirtualThreadPerTaskExecutor());
        server.start();
        System.out.println("assinador.jar em execução na porta " + port);
    }

    private void health(HttpExchange exchange) throws IOException {
        write(exchange, 200, "{\"status\":\"UP\"}");
    }

    private void sign(HttpExchange exchange) throws IOException {
        if (!"POST".equalsIgnoreCase(exchange.getRequestMethod())) {
            write(exchange, 405, "{\"success\":false,\"error\":\"Método não permitido\"}");
            return;
        }
        try {
            String body = new String(exchange.getRequestBody().readAllBytes(), StandardCharsets.UTF_8);
            var request = JsonUtil.signRequestFromJson(body);
            var response = service.sign(request);
            write(exchange, 200, JsonUtil.toJson(response));
        } catch (InvalidParameterException e) {
            write(exchange, 400, JsonUtil.toJson(new ErrorResponse(false, new ErrorResponse.ErrorBody("INVALID_PARAMETER", e.getMessage(), e.details()))));
        } catch (Exception e) {
            write(exchange, 500, JsonUtil.toJson(new ErrorResponse(false, new ErrorResponse.ErrorBody("INTERNAL_ERROR", "Erro interno", e.getMessage()))));
        }
    }

    private void validate(HttpExchange exchange) throws IOException {
        if (!"POST".equalsIgnoreCase(exchange.getRequestMethod())) {
            write(exchange, 405, "{\"success\":false,\"error\":\"Método não permitido\"}");
            return;
        }
        try {
            String body = new String(exchange.getRequestBody().readAllBytes(), StandardCharsets.UTF_8);
            var request = JsonUtil.validateRequestFromJson(body);
            var response = service.validate(request);
            write(exchange, 200, JsonUtil.toJson(response));
        } catch (InvalidParameterException e) {
            write(exchange, 400, JsonUtil.toJson(new ErrorResponse(false, new ErrorResponse.ErrorBody("INVALID_PARAMETER", e.getMessage(), e.details()))));
        } catch (Exception e) {
            write(exchange, 500, JsonUtil.toJson(new ErrorResponse(false, new ErrorResponse.ErrorBody("INTERNAL_ERROR", "Erro interno", e.getMessage()))));
        }
    }

    private void shutdown(HttpExchange exchange) throws IOException {
        write(exchange, 200, "{\"success\":true,\"message\":\"Encerrando assinador\"}");
        new Thread(() -> {
            try { Thread.sleep(200); } catch (InterruptedException ignored) { }
            server.stop(0);
        }).start();
    }

    private void write(HttpExchange exchange, int statusCode, String body) throws IOException {
        byte[] bytes = body.getBytes(StandardCharsets.UTF_8);
        exchange.getResponseHeaders().set("Content-Type", "application/json; charset=utf-8");
        exchange.sendResponseHeaders(statusCode, bytes.length);
        exchange.getResponseBody().write(bytes);
        exchange.close();
    }
}
