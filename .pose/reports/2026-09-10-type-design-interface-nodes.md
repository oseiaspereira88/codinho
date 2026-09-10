# Step nodes de interfaces — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

O desafio `select-validator-true-nil` ganhou dois nodes: verificação da
identidade da interface em modos conhecido e desconhecido, e explicação do
estado zero de uma interface. O fluxo torna observável o risco de encaixotar
um ponteiro nil sem alterar fixture ou desafio canônico.

`go run ./cmd/codinho catalog validate --checks --json` passou com
`diagnostics: null`.
