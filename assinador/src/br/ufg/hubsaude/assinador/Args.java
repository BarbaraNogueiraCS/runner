package br.ufg.hubsaude.assinador;

import java.util.LinkedHashMap;
import java.util.Map;

final class Args {
    private Args() {}

    static Map<String, String> parse(String[] args, int start) {
        Map<String, String> flags = new LinkedHashMap<>();
        for (int i = start; i < args.length; i++) {
            String current = args[i];
            if (!current.startsWith("--")) {
                throw new UserInputException("Argumento inesperado: " + current);
            }
            String key = current.substring(2);
            if (i + 1 >= args.length || args[i + 1].startsWith("--")) {
                throw new UserInputException("Parâmetro sem valor: --" + key);
            }
            flags.put(key, args[++i]);
        }
        return flags;
    }
}
