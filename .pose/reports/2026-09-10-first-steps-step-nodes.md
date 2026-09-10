# Step nodes de go-first-steps — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

O desafio `clamp-int-to-byte` ganhou dois nodes: verificação dos limites
inferior, interno e superior, e reflexão sobre o risco de conversão estreitante.
Eles tornam explícita a validação da ordem das operações sem alterar fixture,
assinatura ou desafio canônico.

`go run ./cmd/codinho catalog validate --checks --json` passou com
`diagnostics: null`.
