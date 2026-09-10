# Conteúdo conceitual de testes — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

O pack go-testing 1.2.0 acrescenta conteúdo aos oito conceitos existentes:
tabelas, subtestes, injeção de relógio, dublês, fronteiras, fusos, isolamento e
repetibilidade. Cada conteúdo contém explicação, exemplo Go completo, leitura
do exemplo e referência a outro conceito. Oito relações correspondentes foram
adicionadas ao grafo. Desafios, fixtures e relações anteriores foram preservados
por comparação estrutural do YAML com o commit 43dba55.

Os exemplos usam contextos diferentes dos desafios e não incluem suas soluções.
A revisão local conferiu limites explícitos de cada exemplo: ASCII, relógio não
nil, dublê sequencial, fronteira inclusiva, fuso fixo e slice de elementos int.
O conteúdo permanece draft e não altera contagens, revisão humana ou playtest.

Validação dos exemplos: extrair cada concepts[].content.example.code para um
example_test.go em diretório isolado, criar go.mod com module example e Go
1.25.0, executar go test ./... -count=1 em cada diretório. Resultado: 8/8 passaram.

A primeira matriz detectou relation_refs sem arestas correspondentes. Foram
adicionadas as oito arestas e a validação estrutural passou. O filtro --module
packs não executa checks nesta matriz; ele não foi usado como evidência de
validação do conteúdo. Consulte a matriz completa em delivery-validation.json.

Próximo lote: conteúdo conceitual em go-io. Manter a spec aberta até satisfazer
aprofundamento, distribuição e avaliação pedagógica. Nenhuma issue externa
foi submetida.

Matriz final: 23/23 checks aprovados, zero skips; docs-check, knowledge-check
strict e ready-check passaram.
