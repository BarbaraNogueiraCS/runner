package br.ufg.hubsaude.assinador;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

final class Json {
    private Json() {}

    static String prop(String key, String value) {
        return quote(key) + ":" + quote(value == null ? "" : value);
    }

    static String rawProp(String key, String rawValue) {
        return quote(key) + ":" + rawValue;
    }

    static String error(String code, String message) {
        return "{" + prop("resourceType", "OperationOutcome") + ","
                + prop("status", "ERROR") + ","
                + "\"valid\":false,"
                + prop("code", code) + ","
                + prop("message", message) + ","
                + rawProp("issue", "[{" + prop("severity", "error") + "," + prop("code", code) + "," + prop("diagnostics", message) + "}]")
                + "}";
    }

    static String operationOutcomeSuccess(String diagnostics, List<String> warnings) {
        StringBuilder issues = new StringBuilder("[");
        issues.append("{").append(prop("severity", "information")).append(",")
                .append(prop("code", "VALIDATION.OK")).append(",")
                .append(prop("diagnostics", diagnostics)).append("}");
        for (String warning : warnings) {
            issues.append(",{").append(prop("severity", "warning")).append(",")
                    .append(prop("code", warning)).append(",")
                    .append(prop("diagnostics", warning)).append("}");
        }
        issues.append("]");
        return "{" + prop("resourceType", "OperationOutcome") + ","
                + prop("status", "SUCCESS") + ","
                + "\"valid\":true,"
                + rawProp("issue", issues.toString()) + "}";
    }

    static String quote(String s) {
        StringBuilder out = new StringBuilder("\"");
        for (int i = 0; i < s.length(); i++) {
            char c = s.charAt(i);
            switch (c) {
                case '\\' -> out.append("\\\\");
                case '"' -> out.append("\\\"");
                case '\n' -> out.append("\\n");
                case '\r' -> out.append("\\r");
                case '\t' -> out.append("\\t");
                default -> {
                    if (c < 0x20) {
                        out.append(String.format("\\u%04x", (int) c));
                    } else {
                        out.append(c);
                    }
                }
            }
        }
        return out.append('"').toString();
    }

    static Map<String, String> parseFlat(String json) {
        Map<String, String> result = new LinkedHashMap<>();
        if (json == null || json.isBlank()) {
            return result;
        }
        String s = json.trim();
        if (!s.startsWith("{") || !s.endsWith("}")) {
            throw new UserInputException("Corpo da requisição deve ser um objeto JSON simples");
        }
        int i = 1;
        while (i < s.length() - 1) {
            while (i < s.length() - 1 && (Character.isWhitespace(s.charAt(i)) || s.charAt(i) == ',')) {
                i++;
            }
            if (i >= s.length() - 1) break;
            Parsed key = readJsonStringValue(s, i);
            i = key.next();
            while (i < s.length() && Character.isWhitespace(s.charAt(i))) i++;
            if (i >= s.length() || s.charAt(i) != ':') {
                throw new UserInputException("JSON inválido: esperado ':'");
            }
            i++;
            while (i < s.length() && Character.isWhitespace(s.charAt(i))) i++;
            if (i >= s.length() || s.charAt(i) != '"') {
                throw new UserInputException("JSON inválido: apenas valores string são aceitos neste endpoint");
            }
            Parsed value = readJsonStringValue(s, i);
            result.put(key.value(), value.value());
            i = value.next();
        }
        return result;
    }

    static String stringField(String json, String field) {
        if (json == null) return null;
        int pos = json.indexOf("\"" + field + "\"");
        if (pos < 0) return null;
        int colon = json.indexOf(':', pos);
        if (colon < 0) return null;
        int i = colon + 1;
        while (i < json.length() && Character.isWhitespace(json.charAt(i))) i++;
        if (i >= json.length() || json.charAt(i) != '"') return null;
        return readJsonStringValue(json, i).value();
    }

    static long longField(String json, String field, long def) {
        Matcher m = Pattern.compile("\\\"" + Pattern.quote(field) + "\\\"\\s*:\\s*(-?\\d+)").matcher(json == null ? "" : json);
        return m.find() ? Long.parseLong(m.group(1)) : def;
    }

    static List<String> stringArray(String json) {
        if (json == null || json.isBlank()) return List.of();
        String s = json.trim();
        if (!s.startsWith("[") || !s.endsWith("]")) {
            throw new UserInputException("CERT.CHAIN-MALFORMED: cadeia de certificados deve ser array JSON de strings base64 DER");
        }
        List<String> values = new ArrayList<>();
        int i = 1;
        while (i < s.length() - 1) {
            while (i < s.length() - 1 && (Character.isWhitespace(s.charAt(i)) || s.charAt(i) == ',')) i++;
            if (i >= s.length() - 1) break;
            if (s.charAt(i) != '"') {
                throw new UserInputException("CERT.CHAIN-MALFORMED: todos os certificados devem ser strings JSON");
            }
            Parsed parsed = readJsonStringValue(s, i);
            values.add(parsed.value());
            i = parsed.next();
        }
        return values;
    }

    static List<String> stringArrayField(String json, String field) {
        String array = arrayFieldRaw(json, field);
        if (array == null) return List.of();
        return stringArray(array);
    }

    static String arrayFieldRaw(String json, String field) {
        int fieldPos = json == null ? -1 : json.indexOf("\"" + field + "\"");
        if (fieldPos < 0) return null;
        int bracket = json.indexOf('[', fieldPos);
        if (bracket < 0) return null;
        return bracketed(json, bracket, '[', ']');
    }

    static String objectFieldRaw(String json, String field) {
        int fieldPos = json == null ? -1 : json.indexOf("\"" + field + "\"");
        if (fieldPos < 0) return null;
        int brace = json.indexOf('{', fieldPos);
        if (brace < 0) return null;
        return bracketed(json, brace, '{', '}');
    }

    private static String bracketed(String json, int start, char open, char close) {
        int depth = 0;
        boolean inString = false;
        boolean escape = false;
        for (int i = start; i < json.length(); i++) {
            char c = json.charAt(i);
            if (escape) { escape = false; continue; }
            if (c == '\\') { escape = inString; continue; }
            if (c == '"') { inString = !inString; continue; }
            if (inString) continue;
            if (c == open) depth++;
            if (c == close) {
                depth--;
                if (depth == 0) return json.substring(start, i + 1);
            }
        }
        return null;
    }

    static String canonicalizeLoose(String json) {
        if (json == null) return "";
        StringBuilder out = new StringBuilder();
        boolean inString = false;
        boolean escape = false;
        for (int i = 0; i < json.length(); i++) {
            char c = json.charAt(i);
            if (escape) { out.append(c); escape = false; continue; }
            if (c == '\\') { out.append(c); escape = inString; continue; }
            if (c == '"') { out.append(c); inString = !inString; continue; }
            if (!inString && Character.isWhitespace(c)) continue;
            out.append(c);
        }
        return out.toString();
    }

    private record Parsed(String value, int next) {}

    private static Parsed readJsonStringValue(String s, int start) {
        if (start >= s.length() || s.charAt(start) != '"') {
            throw new UserInputException("JSON inválido: esperado texto entre aspas");
        }
        StringBuilder out = new StringBuilder();
        int i = start + 1;
        while (i < s.length()) {
            char c = s.charAt(i++);
            if (c == '"') {
                return new Parsed(out.toString(), i);
            }
            if (c == '\\') {
                if (i >= s.length()) {
                    throw new UserInputException("JSON inválido: escape incompleto");
                }
                char e = s.charAt(i++);
                switch (e) {
                    case '"' -> out.append('"');
                    case '\\' -> out.append('\\');
                    case '/' -> out.append('/');
                    case 'n' -> out.append('\n');
                    case 'r' -> out.append('\r');
                    case 't' -> out.append('\t');
                    default -> out.append(e);
                }
            } else {
                out.append(c);
            }
        }
        throw new UserInputException("JSON inválido: texto não encerrado");
    }
}
