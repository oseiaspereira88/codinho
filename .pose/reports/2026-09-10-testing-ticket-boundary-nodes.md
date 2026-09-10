# Step nodes de fronteira temporal — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

A variante `check-ticket-boundary` ganhou dois nodes de verificação: igualdade
como ingresso expirado e equivalência de representações em fusos distintos.
O contrato compara instantes absolutos sem depender do relógio do sistema.

`go run ./cmd/codinho catalog validate --checks --json` passou com
`diagnostics: null`.
