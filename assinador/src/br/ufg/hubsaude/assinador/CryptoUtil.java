package br.ufg.hubsaude.assinador;

import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.util.Base64;

final class CryptoUtil {
    private CryptoUtil() {}

    static byte[] sha256(byte[] content) throws Exception {
        return MessageDigest.getInstance("SHA-256").digest(content);
    }

    static byte[] sha512(byte[] content) throws Exception {
        return MessageDigest.getInstance("SHA-512").digest(content);
    }

    static String base64(byte[] content) {
        return Base64.getEncoder().encodeToString(content);
    }

    static byte[] base64Decode(String content, String code) {
        try {
            return Base64.getDecoder().decode(content);
        } catch (IllegalArgumentException e) {
            throw new UserInputException(code + ": valor não está codificado em base64 padrão");
        }
    }

    static String base64Url(byte[] content) {
        return Base64.getUrlEncoder().withoutPadding().encodeToString(content);
    }

    static byte[] base64UrlDecode(String content, String code) {
        try {
            return Base64.getUrlDecoder().decode(content);
        } catch (IllegalArgumentException e) {
            throw new UserInputException(code + ": valor não está codificado em base64Url");
        }
    }

    static String utf8(byte[] bytes) {
        return new String(bytes, StandardCharsets.UTF_8);
    }

    static byte[] utf8(String text) {
        return text.getBytes(StandardCharsets.UTF_8);
    }
}
