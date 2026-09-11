# Node de prefixo UTF-8

Spec: go-foundations-packs, in-progress.

O desafio de truncamento por bytes recebeu uma verificação para confirmar que
os resultados são prefixos UTF-8 válidos e respeitam o orçamento, incluindo
cortes próximos a runes multibyte. A fixture já cobre casos ASCII, 2 bytes e
3 bytes; o node explicita essa propriedade sem alterar o contrato.

`go run ./cmd/codinho catalog validate --checks --json` passou com diagnostics
null e editorial null. `git diff --check` passou; revisão humana e playtest
permanecem pendentes.
