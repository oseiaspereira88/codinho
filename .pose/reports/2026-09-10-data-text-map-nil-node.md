# Node de consulta em mapa nil

Spec: go-foundations-packs, in-progress.

O desafio de consulta com comma-ok recebeu um node para a fronteira de mapa
nil. A consulta deve produzir o valor zero e `scored: false` sem panic, como
uma ausência comum. O caso complementa a distinção entre zero presente e zero
ausente já coberta pela fixture.

`go run ./cmd/codinho catalog validate --checks --json` produziu diagnostics
null e editorial null. `git diff --check` passou. A matriz strict, revisão
humana e playtest continuam pendentes.
