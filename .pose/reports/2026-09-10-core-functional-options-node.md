# Node de padrões e precedência em opções funcionais

Spec: go-foundations-packs, in-progress.

O desafio combinado de opções funcionais recebeu uma verificação para os
valores padrão de `NewServer()` e para a precedência da última `WithPort`.
Esses casos já fazem parte da fixture; o node os conecta explicitamente ao
critério de aplicação ordenada das opções.

`go run ./cmd/codinho catalog validate --checks --json` passou com diagnostics
null e editorial null. `pose validate --strict`, `pose check --strict` e
`git diff --check` passaram; revisão humana e playtest continuam pendentes.
