# Step nodes de shadowing — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

O desafio `avoid-shadowing-named-returns` ganhou dois nodes após a implementação:
verificação dos casos válidos e inválidos e explicação da diferença entre
atribuição e redeclaração dentro do bloco. Fixture e desafio canônico foram
preservados.

`go run ./cmd/codinho catalog validate --checks --json` passou com
`diagnostics: null`.
