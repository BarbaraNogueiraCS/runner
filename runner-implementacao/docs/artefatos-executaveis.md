# Artefatos executáveis de release

Este documento registra como o projeto Runner atende aos critérios de geração e distribuição dos artefatos executáveis.

## Artefatos esperados por release SemVer

Para uma tag `v1.0.0`, o workflow de release deve publicar no GitHub Releases:

| Aplicação | Plataforma | Artefato |
|---|---|---|
| assinatura | Windows amd64 | `assinatura-1.0.0-windows-amd64.exe` |
| assinatura | Linux amd64 | `assinatura-1.0.0-linux-amd64.AppImage` |
| assinatura | macOS amd64 | `assinatura-1.0.0-macos-amd64.dmg` |
| simulador | Windows amd64 | `simulador-1.0.0-windows-amd64.exe` |
| simulador | Linux amd64 | `simulador-1.0.0-linux-amd64.AppImage` |
| simulador | macOS amd64 | `simulador-1.0.0-macos-amd64.dmg` |

A versão `1.0.0` é apenas exemplo. A versão real é extraída automaticamente da tag SemVer, removendo o prefixo `v`. Por exemplo, a tag `v1.2.3` gera `assinatura-1.2.3-linux-amd64.AppImage`.

## Distribuição

A distribuição é automatizada por `.github/workflows/release.yml`.

O workflow é disparado por tags `v*` e valida que a tag segue SemVer no formato:

```text
vMAJOR.MINOR.PATCH
```

Exemplo:

```bash
git tag v1.0.0
git push origin v1.0.0
```

## Integridade

O workflow gera `checksums.txt` com SHA256 de todos os artefatos publicados.

## Assinatura dos artefatos

Cada artefato publicado, incluindo `checksums.txt`, é assinado com Cosign em modo keyless/OIDC. Para cada arquivo assinado, o release inclui:

- `<artefato>.sig`
- `<artefato>.pem`
- `<artefato>.bundle`

O workflow declara as permissões necessárias:

```yaml
permissions:
  contents: write
  id-token: write
```

O campo `id-token: write` permite a assinatura com identidade OIDC e o registro da assinatura no transparency log do Sigstore.

## Limite operacional

A criação e publicação real dos releases depende do projeto estar hospedado no GitHub e de uma tag SemVer ser enviada ao repositório. O código já contém o workflow necessário, mas a publicação efetiva ocorre no GitHub Actions.
