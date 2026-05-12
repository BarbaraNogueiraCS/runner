package br.ufg.hubsaude.assinador;

import br.ufg.hubsaude.assinador.cli.CliAdapter;
import br.ufg.hubsaude.assinador.controller.SignatureHttpServer;
import br.ufg.hubsaude.assinador.service.FakeSignatureService;
import br.ufg.hubsaude.assinador.validation.ParameterValidator;

public class Main {
    public static void main(String[] args) throws Exception {
        var service = new FakeSignatureService(new ParameterValidator());
        if (args.length > 0 && "server".equalsIgnoreCase(args[0])) {
            int port = readPort(args, 8080);
            new SignatureHttpServer(service).start(port);
            return;
        }
        new CliAdapter(service).execute(args);
    }

    private static int readPort(String[] args, int defaultPort) {
        for (int i = 0; i < args.length - 1; i++) {
            if ("--port".equals(args[i])) {
                return Integer.parseInt(args[i + 1]);
            }
        }
        return defaultPort;
    }
}
