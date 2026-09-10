# Conteúdo conceitual de go-debugging — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

O conceito `slice-bounds` agora explica a faixa válida de índices, o erro de
usar `len` como índice e a guarda necessária para slices vazias. O exemplo
isolado cobre lista preenchida e entrada nil. O desafio de depuração e sua
fixture foram preservados.

`go run ./cmd/codinho catalog validate --checks --json` passou com
`diagnostics: null`. Não houve revisão humana ou publicação.

Próxima demanda: medir competências e step nodes por pack e iniciar a expansão
qualitativa sem duplicar os desafios canônicos.
