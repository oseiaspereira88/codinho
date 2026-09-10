# Node de parada rotulada

Spec: go-foundations-packs, in-progress.

O desafio de `labeled break` recebeu um node de verificação para a fronteira
observável após o primeiro comando `stop`: nenhum comando posterior pode ser
processado. A fixture já exercita a sequência `a, b, stop, c`; o novo node
torna esse comportamento explícito no percurso pedagógico.

`go run ./cmd/codinho catalog validate --checks --json` passou com diagnostics
null e editorial null. `git diff --check` passou. O contrato, autoria e
publicação permanecem inalterados; revisão humana e playtest continuam pendentes.
