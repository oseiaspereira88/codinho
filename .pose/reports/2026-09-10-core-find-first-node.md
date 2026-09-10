# Node de flag de busca

Spec: go-foundations-packs, in-progress.

O desafio `go-core.find-first-value-at-least` recebeu uma verificação para
distinguir um valor zero encontrado (`found: true`) de uma busca sem
correspondência (`found: false`). O node reforça o contrato do retorno duplo e
mantém a proteção contra acesso a nodes nil.

`go run ./cmd/codinho catalog validate --checks --json` passou com diagnostics
null e editorial null. `git diff --check` passou; a fixture, autoria e
publicação não foram alteradas. Revisão humana e playtest continuam pendentes.
