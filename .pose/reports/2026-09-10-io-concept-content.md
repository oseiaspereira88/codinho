# Conteúdo conceitual de I/O — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

O pack go-io 1.2.0 aprofunda seus dez conceitos existentes com explicação,
exemplo Go completo e referência a outro conceito, acompanhada da relação no
grafo. O lote cobre Reader, leituras parciais, EOF, Decoder JSON, campos
estritos, CSV com aspas, falhas de Writer, validação antes da escrita, largura
CSV e Encoder JSON. Desafios, fixtures e relações anteriores foram comparados
estruturalmente com 237c3bc e preservados.

A revisão local separada da geração conferiu os limites dos exemplos e o
tratamento de bytes junto com EOF. O Reader do exemplo mantém posição para
preservar bytes entre chamadas. Os exemplos não reproduzem as soluções dos
exercícios e não usam arquivos externos, rede ou dependências adicionais.

Validação dos exemplos: extrair cada concepts[].content.example.code para
example_test.go em diretório isolado, criar go.mod com module example e Go
1.25.0 e executar go test ./... -count=1. Resultado: dez exemplos aprovados.
O catálogo estrutural também passou. A avaliação mecânica não constitui
revisão humana ou playtest; nenhum conteúdo foi publicado.

Próximo lote: conceitos de erros e núcleo de Go, sem aumentar contagens com
itens redundantes. A meta de distribuição permanece pendente de decisão.

Matriz final: 23/23 checks aprovados, zero skips. Docs-check, knowledge-check
strict, check strict e ready-check passaram.

A auditoria artifact-check da spec inteira ainda falha na atribuição histórica:
criação de go-core, go-data-text, go-errors e go-type-design não reconhecida no
intervalo atribuído, e dois documentos observados sem declaração
(content-review-checklist.md e agent-review-workflow.md). Isso não é falha dos
exemplos; requer reconciliação das revisões reais antes de fechar a spec.
