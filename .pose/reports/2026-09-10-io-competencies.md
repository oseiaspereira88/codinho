# Competências de I/O — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

O pack `go-io` ganhou cinco competências explícitas derivadas da árvore já
existente: consumir entrada fragmentada, distinguir EOF, validar a fronteira
do documento, interpretar campos citados e codificar saída estruturada. Cada
uma foi associada ao desafio que já exercita o comportamento; nenhum desafio
canônico novo foi criado.

`go run ./cmd/codinho catalog validate --checks --json` passou com
`diagnostics: null`. O conteúdo, fixtures e competências anteriores foram
preservados.
