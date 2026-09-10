# Competências de investigação em debugging — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

O desafio `go-debug.slice-off-by-one` agora explicita as competências de
reproduzir a falha, localizar a divergência de índice e formular uma hipótese.
Elas já correspondiam aos nodes de investigação existentes; a mudança torna a
progressão observável no catálogo sem criar outro desafio canônico.

`go run ./cmd/codinho catalog validate --checks --json` passou com
`diagnostics: null`. O desafio, fixture e conceito foram preservados.
