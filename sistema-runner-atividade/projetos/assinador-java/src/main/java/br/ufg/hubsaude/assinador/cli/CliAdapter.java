package br.ufg.hubsaude.assinador.cli;

import br.ufg.hubsaude.assinador.dto.SignRequest;
import br.ufg.hubsaude.assinador.dto.ValidateRequest;
import br.ufg.hubsaude.assinador.error.InvalidParameterException;
import br.ufg.hubsaude.assinador.service.SignatureService;
import br.ufg.hubsaude.assinador.util.JsonUtil;

import java.util.HashMap;
import java.util.Map;

public class CliAdapter {
    private final SignatureService service;

    public CliAdapter(SignatureService service) {
        this.service = service;
    }

    public void execute(String[] args) {
        if (args.length == 0 || "help".equals(args[0]) || "--help".equals(args[0])) {
            printHelp();
            return;
        }
        try {
            String operation = args[0];
            Map<String, String> values = parse(args);
            if ("sign".equalsIgnoreCase(operation)) {
                var response = service.sign(new SignRequest(values.get("documento"), values.get("certificado"), Map.of()));
                System.out.println(JsonUtil.toJson(response));
                return;
            }
            if ("validate".equalsIgnoreCase(operation)) {
                var response = service.validate(new ValidateRequest(values.get("documento"), values.get("assinatura"), values.get("certificado"), Map.of()));
                System.out.println(JsonUtil.toJson(response));
                return;
            }
            System.err.println("Operação desconhecida: " + operation);
            System.exit(2);
        } catch (InvalidParameterException e) {
            System.err.println("Erro: " + e.getMessage());
            System.err.println("Motivo: " + e.details());
            System.exit(2);
        }
    }

    private Map<String, String> parse(String[] args) {
        Map<String, String> values = new HashMap<>();
        for (int i = 1; i < args.length - 1; i++) {
            if (args[i].startsWith("--")) {
                String key = args[i].substring(2);
                values.put(key, args[i + 1]);
                i++;
            }
        }
        return values;
    }

    private void printHelp() {
        System.out.println("Uso: java -jar assinador.jar <sign|validate> [parâmetros]");
        System.out.println("  sign --documento <arquivo> --certificado <arquivo>");
        System.out.println("  validate --documento <arquivo> --assinatura <valor>");
        System.out.println("  server --port <porta>");
    }
}
