# Artefatos executáveis

A distribuição é automatizada por `.github/workflows/release.yml`, localizado na raiz do repositório. Como o módulo Go também está na raiz, os comandos do workflow executam diretamente `go test`, `go build`, `make -C assinador` e scripts em `scripts/`.

Para uma tag `v1.0.5`, os principais artefatos esperados são:

```text
assinatura-1.0.5-linux-amd64.AppImage
assinatura-1.0.5-windows-amd64.exe
assinatura-1.0.5-macos-amd64.dmg
simulador-1.0.5-linux-amd64.AppImage
simulador-1.0.5-windows-amd64.exe
simulador-1.0.5-macos-amd64.dmg
assinador-1.0.5.jar
checksums.txt
```

Cada artefato primário recebe também arquivos de assinatura Cosign:

```text
<artefato>.sig
<artefato>.pem
<artefato>.bundle
```

Os binários oficiais não são commitados no Git; eles são publicados na aba GitHub Releases.
