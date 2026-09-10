# Step nodes de catálogo em memória — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

O desafio `maintain-memory-catalog` ganhou dois nodes de verificação: rejeição
de produto inválido sem alterar estado e listagem ordenada com resultados
isolados. A mudança explicita invariantes de estado e snapshots sem alterar a
fixture ou criar desafio canônico.

`go run ./cmd/codinho catalog validate --checks --json` passou com
`diagnostics: null`.
