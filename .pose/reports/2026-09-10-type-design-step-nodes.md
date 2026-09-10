# Step nodes de transformação genérica — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

O desafio `map-generic-transform` ganhou dois nodes após a implementação:
verificação de transformação completa e ordenada, e verificação de isolamento
entre entrada e saída. A mudança torna observáveis os casos de slice longa e
aliasing sem criar desafio canônico novo.

`go run ./cmd/codinho catalog validate --checks --json` passou com
`diagnostics: null`.
