# Node de limites do invariante Percentage

Spec: go-foundations-packs, in-progress.

O desafio `go-core.guard-invariant-with-unexported-field` recebeu uma
verificação explícita dos limites inclusivos (0 e 100) e dos valores adjacentes
inválidos (-1 e 101). A fixture já cobre a tabela completa; o node torna o
invariante observável no percurso.

`go run ./cmd/codinho catalog validate --checks --json` passou com diagnostics
null e editorial null. `git diff --check` passou; revisão humana e playtest
permanecem pendentes.
