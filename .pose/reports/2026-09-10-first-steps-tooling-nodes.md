# Step nodes de tooling em go-first-steps — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

Os desafios de módulo mínimo e mensagem de limite de taxa ganharam quatro
nodes: inspeção do module path, confirmação da versão suportada, execução do
verificador de formatação e comparação da mensagem com o contrato. A mudança
separa implementação de evidência e preserva fixtures e desafios canônicos.

`go run ./cmd/codinho catalog validate --checks --json` passou com
`diagnostics: null`.
