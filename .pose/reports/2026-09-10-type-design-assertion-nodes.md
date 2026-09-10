# Step nodes de assertion segura — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

O desafio `describe-if-circle-safely` ganhou dois nodes: verificação de Circle
e tipo alternativo, e explicação da forma `comma-ok`. A árvore agora torna
observável o caminho sem panic para assertions inesperadas.

`go run ./cmd/codinho catalog validate --checks --json` passou com
`diagnostics: null`.
