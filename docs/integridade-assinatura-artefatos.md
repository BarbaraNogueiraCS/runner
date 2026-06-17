# Integridade e assinatura de artefatos

Este documento registra como o projeto Runner atende ao requisito de integridade e assinatura dos artefatos publicados em GitHub Releases.

## 1. Objetivo

Todos os binários e artefatos distribuídos em uma release devem permitir verificação independente de autenticidade e integridade. Para isso, o pipeline de release usa Cosign, do ecossistema Sigstore, com assinatura keyless baseada em OIDC.

## 2. Requisito implementado

O workflow `.github/workflows/release.yml`, localizado na raiz do repositório Git, executa automaticamente a assinatura durante a publicação da release.

A assinatura é feita no job `publish`, após a geração dos artefatos para Windows, Linux e macOS. O workflow declara:

```yaml
permissions:
  contents: write
  id-token: write
```

A permissão `id-token: write` permite que o GitHub Actions forneça uma identidade OIDC para o Cosign. A assinatura é enviada para o transparency log do Sigstore por meio da opção explícita:

```bash
--tlog-upload=true
```

## 3. Artefatos assinados

Para uma tag `v1.0.0`, os principais artefatos executáveis publicados são:

```text
assinatura-1.0.0-windows-amd64.exe
assinatura-1.0.0-linux-amd64.AppImage
assinatura-1.0.0-macos-amd64.dmg
simulador-1.0.0-windows-amd64.exe
simulador-1.0.0-linux-amd64.AppImage
simulador-1.0.0-macos-amd64.dmg
```

A versão `1.0.0` é apenas exemplo. A versão real é extraída da tag SemVer usada na release. Por exemplo, a tag `v1.2.3` gera arquivos com a versão `1.2.3`.

Além desses artefatos, o workflow também publica o `assinador-<versão>.jar` e `checksums.txt`. Esses arquivos também são assinados para manter a rastreabilidade completa da release.

## 4. Arquivos obrigatórios por artefato

Para cada artefato assinado, a release deve conter:

```text
<artefato>
<artefato>.sig
<artefato>.pem
```

O workflow também publica:

```text
<artefato>.bundle
```

O `.bundle` é complementar. Ele não substitui os arquivos obrigatórios `.sig` e `.pem`.

Exemplo para Linux:

```text
assinatura-1.0.0-linux-amd64.AppImage
assinatura-1.0.0-linux-amd64.AppImage.sig
assinatura-1.0.0-linux-amd64.AppImage.pem
assinatura-1.0.0-linux-amd64.AppImage.bundle
```

## 5. Comando de assinatura usado no pipeline

O trecho central do workflow é:

```bash
cosign sign-blob --yes --tlog-upload=true \
  --bundle "$f.bundle" \
  --output-signature "$f.sig" \
  --output-certificate "$f.pem" \
  "$f"
```

Esse comando gera automaticamente a assinatura `.sig`, o certificado `.pem` e o bundle `.bundle` para cada artefato.

## 6. Verificação pelos usuários

Depois de baixar o artefato, a assinatura e o certificado da release, o usuário pode verificar a autenticidade com:

```bash
cosign verify-blob \
  --certificate assinatura-1.0.0-linux-amd64.AppImage.pem \
  --signature assinatura-1.0.0-linux-amd64.AppImage.sig \
  assinatura-1.0.0-linux-amd64.AppImage
```

Se a verificação for bem-sucedida, o Cosign indicará que a assinatura é válida.

Também é possível usar a verificação com bundle, quando desejado:

```bash
cosign verify-blob \
  --bundle assinatura-1.0.0-linux-amd64.AppImage.bundle \
  assinatura-1.0.0-linux-amd64.AppImage
```

## 7. Automação

A assinatura não depende de uma etapa manual. Ela ocorre automaticamente quando uma tag SemVer é enviada ao GitHub:

```bash
git tag v1.0.0
git push origin v1.0.0
```

O workflow valida o padrão SemVer `vMAJOR.MINOR.PATCH`, executa os testes, compila os artefatos, gera os checksums, assina os arquivos com Cosign e publica tudo no GitHub Releases.

## 8. Justificativa

A assinatura dos artefatos distribuídos proporciona:

- verificação da autenticidade dos binários;
- proteção contra adulteração dos artefatos;
- rastreabilidade da origem do software;
- maior segurança para usuários e integradores;
- melhoria da segurança da cadeia de suprimentos de software.

## 9. Validação local do requisito

Para verificar localmente se o projeto contém a automação esperada, execute:

```bash
./scripts/check-release-artifacts.sh
```

Esse script confere se o workflow contém geração dos artefatos, publicação via GitHub Releases, SHA256, Cosign, OIDC, transparency log, `.sig` e `.pem`.
