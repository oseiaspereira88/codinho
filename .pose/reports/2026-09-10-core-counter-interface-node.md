# Node de isolamento do contador por interface

Spec: go-foundations-packs, in-progress.

O desafio `go-core.hide-counter-type-behind-interface` recebeu uma verificação
para confirmar que dois valores retornados por `NewCounter` mantêm contagens
independentes. O node reforça o encapsulamento do tipo concreto e reutiliza a
fixture existente, que já exercita duas instâncias.

`go run ./cmd/codinho catalog validate --checks --json` passou com diagnostics
null e editorial null. `git diff --check` passou; a matriz strict, revisão
humana e playtest permanecem pendentes.
