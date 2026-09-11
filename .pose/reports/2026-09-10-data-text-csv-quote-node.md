# Node de caracteres especiais CSV

Spec: go-foundations-packs, in-progress.

O desafio de formatação de campo CSV recebeu um node para verificar que
vírgulas, aspas, newline e retorno de carro acionam citação e preservam o
conteúdo escapado. A fixture já cobre cada caso isolado e a combinação de
vírgula com aspas; o node reúne a propriedade de forma observável.

`go run ./cmd/codinho catalog validate --checks --json` passou com diagnostics
null e editorial null. `git diff --check` passou; revisão humana e playtest
continuam pendentes.
