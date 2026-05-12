package br.ufg.hubsaude.assinador.util;

import br.ufg.hubsaude.assinador.dto.ErrorResponse;
import br.ufg.hubsaude.assinador.dto.SignRequest;
import br.ufg.hubsaude.assinador.dto.SignResponse;
import br.ufg.hubsaude.assinador.dto.ValidateRequest;
import br.ufg.hubsaude.assinador.dto.ValidateResponse;

import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public final class JsonUtil {
    private JsonUtil() {}

    public static String getString(String json, String field) {
        if (json == null) return "";
        Pattern pattern = Pattern.compile("\\\"" + Pattern.quote(field) + "\\\"\\s*:\\s*\\\"([^\\\"]*)\\\"");
        Matcher matcher = pattern.matcher(json);
        if (matcher.find()) {
            return matcher.group(1);
        }
        return "";
    }

    public static SignRequest signRequestFromJson(String json) {
        return new SignRequest(getString(json, "document"), getString(json, "certificate"), Map.of());
    }

    public static ValidateRequest validateRequestFromJson(String json) {
        return new ValidateRequest(getString(json, "document"), getString(json, "signature"), getString(json, "certificate"), Map.of());
    }

    public static String toJson(SignResponse response) {
        return "{"
                + quote("success") + ":" + response.success() + ","
                + quote("operation") + ":" + quote(response.operation()) + ","
                + quote("signature") + ":" + quote(response.signature()) + ","
                + quote("message") + ":" + quote(response.message())
                + "}";
    }

    public static String toJson(ValidateResponse response) {
        return "{"
                + quote("success") + ":" + response.success() + ","
                + quote("operation") + ":" + quote(response.operation()) + ","
                + quote("valid") + ":" + response.valid() + ","
                + quote("message") + ":" + quote(response.message())
                + "}";
    }

    public static String toJson(ErrorResponse response) {
        return "{"
                + quote("success") + ":false,"
                + quote("error") + ":{"
                + quote("code") + ":" + quote(response.error().code()) + ","
                + quote("message") + ":" + quote(response.error().message()) + ","
                + quote("details") + ":" + quote(response.error().details())
                + "}}";
    }

    public static String quote(String value) {
        return "\"" + escape(value) + "\"";
    }

    private static String escape(String value) {
        if (value == null) return "";
        return value.replace("\\", "\\\\").replace("\"", "\\\"").replace("\n", "\\n");
    }
}
