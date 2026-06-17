# ADR 0004 — Provisionamento por release.json e assinatura de artefatos com Cosign

## Contexto

O Runner precisa executar aplicações Java sem exigir que o usuário conheça detalhes de instalação do Java ou do simulador/validador do HubSaúde. Também precisa publicar binários verificáveis para Windows, Linux e macOS.

## Decisão

1. O Runner usa `https://raw.githubusercontent.com/kyriosdata/runner/main/release.json` como manifesto remoto padrão.
2. O CLI compara a versão do artefato remoto com a versão local em `~/.hubsaude/<artefato>/*.version`.
3. O download é feito apenas quando o artefato não existe, está desatualizado ou possui checksum SHA256 divergente.
4. O Java é procurado em `RUNNER_JAVA`, `~/.hubsaude/jdk/bin/java` e no `PATH`; se Java 21+ não estiver disponível, o Runner baixa o runtime Temurin indicado no manifesto.
5. A release do projeto gera artefatos para Windows, Linux e macOS e assina cada arquivo com Cosign em modo keyless/OIDC, incluindo transparency log.

## Consequências

- A instalação local fica reprodutível e rastreável.
- O usuário não precisa baixar manualmente o simulador/validador quando usa `simulador start` sem `--jar`.
- O pipeline de release passa a depender de permissões `id-token: write` para assinatura keyless com Cosign.
