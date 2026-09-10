# Step nodes de contagem de prazos — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

A variante `count-deadline-clock-calls` ganhou dois nodes: lote vazio com uma
única leitura do relógio e exclusão de prazo exatamente no instante observado.
O contrato temporal fica explícito sem alterar fixture ou desafio canônico.

`go run ./cmd/codinho catalog validate --checks --json` passou com
`diagnostics: null`.
