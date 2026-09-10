# Complemento conceitual de erros — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

Go-errors 1.4.0 completa o conteúdo dos oito conceitos antes sem explicação.
Os 12 conceitos do pack agora têm conteúdo; os quatro conteúdos anteriores,
desafios, fixtures e relações anteriores foram preservados por comparação
estrutural com 1cd4ef1. Nenhuma contagem de conceitos ou desafios foi alterada.

Oito exemplos cobrem precedência explícita entre categorias, identidade distinta
com mensagem igual, busca através de dois wrappers, primeira correspondência
em árvore agregada, mapa inicializado na primeira escrita, rejeição sem mutação,
ordenação lexical e isolamento de structs por valor. A revisão local conferiu
os limites: sem concorrência, sem cópia profunda implícita e sem prioridade
de negócio inferida da ordem da árvore. Não houve avaliação humana ou publicação.

Reprodução: extrair content.example.code dos oito conceitos complementados para
example_test.go em diretórios isolados, criar go.mod (module example; Go 1.25.0)
e executar go test ./... -count=1. Oito exemplos passaram. Validação estrutural
do catálogo também passou. Todos usam biblioteca padrão, sem rede ou arquivos.

Próxima demanda: conteúdo de apoio em go-core. A spec permanece aberta para
aprofundamento, distribuição, revisão pedagógica e playtest.

Matriz strict: 23/23 checks aprovados, zero skips. Docs-check, knowledge-check,
ready-check, check strict e artifact-check strict passaram. Os avisos globais
de arquivos sem atribuição permanecem fora deste lote.
