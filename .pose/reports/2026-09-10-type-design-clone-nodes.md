# Step nodes de clonagem de struct — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

O desafio `clone-inventory-independently` ganhou dois nodes de verificação:
capacidade excedente em slice vazia e preservação do campo `Owner`. O fluxo
agora evidencia tanto o isolamento profundo quanto a cópia dos campos escalares,
sem alterar fixture ou criar desafio canônico.

`go run ./cmd/codinho catalog validate --checks --json` passou com
`diagnostics: null`.
