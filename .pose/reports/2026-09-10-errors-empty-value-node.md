# Node de valor vazio presente

Spec: go-foundations-packs, in-progress.

O desafio de encadeamento de erro sentinela recebeu uma verificação para
distinguir chave presente com valor vazio de chave ausente. O node reforça a
necessidade de usar o segundo retorno da consulta ao mapa antes de construir
`ErrNotFound`; a fixture existente cobre os casos de sucesso e ausência.

`go run ./cmd/codinho catalog validate --checks --json` passou com diagnostics
null e editorial null. `git diff --check` passou; revisão humana e playtest
continuam pendentes.
