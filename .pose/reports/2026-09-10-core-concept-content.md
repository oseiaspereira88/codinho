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

## Segundo lote

Go-core 1.3.0 completa os quatro conceitos restantes. Os nove conceitos agora
têm conteúdo, incluindo o exemplo histórico de opções funcionais preservado.
Os quatro novos exemplos passaram isoladamente com go test ./... -count=1.
A comparação estrutural com d688816 confirmou preservação dos dez desafios,
competências, temas, trilhas, conteúdo anterior e relações históricas.

A revisão local verificou que recover é direto no defer e limitado à mesma
goroutine; o zero value de Balance é válido; rejeição não altera estado;
curto-circuito distingue nil de valor inativo; o construtor retorna nil real
e demonstra separadamente a interface contendo ponteiro nil.

Próxima demanda autônoma: aprofundar os conceitos sem conteúdo de go-data-text.
Contagens e publicação permanecem inalteradas. Revisão e playtest humanos,
distribuição das bases e profundidade curricular continuam pendentes na spec.

Segundo lote: matriz strict aprovada, 23/23 checks, incluindo race, build,
govulncheck e fixtures. Docs-check, ready-check e artifact-check strict passaram;
os avisos globais de proveniência fora do lote continuam registrados pelo gate.
