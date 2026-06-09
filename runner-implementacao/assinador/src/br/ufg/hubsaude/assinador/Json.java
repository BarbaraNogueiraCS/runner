package br.ufg.hubsaude.assinador;

import java.util.LinkedHashMap;
import java.util.Map;

final class Json {
    private Json() {}

    static String prop(String key, String value) {
        return quote(key) + ":" + quote(value == null ? "" : value);
    }

    static String error(String code, String message) {
        return "{" + prop("status", "ERROR") + "," + prop("code", code) + "," + prop("message", message) + "}";
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
                default -> out.append(c);
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
        s = s.substring(1, s.length() - 1).trim();
        if (s.isEmpty()) {
            return result;
        }
        int i = 0;
        while (i < s.length()) {
            i = skipWhitespace(s, i);
            Parsed key = readString(s, i);
            i = skipWhitespace(s, key.next());
            if (i >= s.length() || s.charAt(i) != ':') {
                throw new UserInputException("JSON inválido: esperado ':'");
            }
            i++;
            i = skipWhitespace(s, i);
            Parsed value = readString(s, i);
            result.put(key.value(), value.value());
            i = skipWhitespace(s, value.next());
            if (i < s.length()) {
                if (s.charAt(i) != ',') {
                    throw new UserInputException("JSON inválido: esperado ','");
                }
                i++;
            }
        }
        return result;
    }

    private static int skipWhitespace(String s, int i) {
        while (i < s.length() && Character.isWhitespace(s.charAt(i))) {
            i++;
        }
        return i;
    }

    private static Parsed readString(String s, int start) {
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

    private record Parsed(String value, int next) {}
}
