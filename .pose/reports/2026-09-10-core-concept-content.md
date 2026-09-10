# Conteúdo conceitual de go-core — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

## Primeiro lote

Go-core 1.2.0 acrescenta conteúdo a defer-execution-order,
variadic-parameters, labeled-break-scope e pointer-vs-value-receiver.
Os quatro exemplos Go foram formatados e executados individualmente em
example_test.go, com module example e Go 1.25.0: go test ./... -count=1 passou.
Desafios e competências foram comparados estruturalmente com o conteúdo anterior.
As referências de conteúdo receberam relações explícitas no grafo.

A primeira validação detectou relações ausentes, corrigidas antes do commit.
A execução no sandbox também encontrou cache Go somente leitura. A matriz
foi repetida fora do sandbox, com /home/go/go/bin no PATH para govulncheck.

A revisão local conferiu avaliação antecipada de argumentos de defer,
compartilhamento de slices, alvo de break e diferença entre chamada automática
por endereço e satisfação de interface. Não houve revisão humana ou publicação.

## Próximo lote

Completar panic/recover, campos não exportados, curto-circuito e construtores
retornando interface. Preservar o exemplo histórico de opções funcionais.

Primeiro lote: matriz strict aprovada, 23/23 checks; pose check --strict e
pose lint-spec go-foundations-packs --ready-check passaram.
