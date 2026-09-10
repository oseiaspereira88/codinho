# Nodes de saída da seleção temporal — 2026-09-10 UTC

Spec: `go-foundations-packs`, in-progress.

O desafio `go-testing.select-active-tokens-with-injected-clock` recebeu verificações para garantir slice vazio não nulo e preservação da ordem dos tokens ativos em uma entrada mista. Os nodes complementam as fronteiras temporais já existentes sem alterar a fixture.

Validação: `go run ./cmd/codinho catalog validate --checks --json` passou com `diagnostics: null`; `git diff --check` passou.
