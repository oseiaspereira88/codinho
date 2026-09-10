# Step nodes de classificação de erros — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

O desafio `classify-wrapped-errors` ganhou dois nodes de verificação: distinguir
identidade de mensagem igual e confirmar precedência de `denied` em erros
agregados independentemente da ordem. Fixture e desafio canônico foram
preservados.

`go run ./cmd/codinho catalog validate --checks --json` passou com
`diagnostics: null`.
