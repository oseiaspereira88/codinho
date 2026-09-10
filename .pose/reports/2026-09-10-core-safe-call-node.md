# Node de conversão de panic

Spec: go-foundations-packs, in-progress.

O desafio `go-core.recover-from-panic-in-safe-call` recebeu um node que reúne
os casos de panic string, inteiro e error, verificando que todos retornam erro
não nulo e informação útil quando aplicável. A fixture já cobre esses três
valores; o node explicita a fronteira pedagógica sem alterar o contrato.

`go run ./cmd/codinho catalog validate --checks --json` passou com diagnostics
null e editorial null. `git diff --check` passou; revisão humana e playtest
continuam pendentes.
