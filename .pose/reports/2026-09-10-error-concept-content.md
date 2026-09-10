# Conteúdo conceitual de erros — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

O pack go-errors 1.3.0 acrescenta explicações e exemplos a errors.As,
errors.Join e tradução na fronteira, preservando o conteúdo existente de %w.
As referências conceituais têm relações correspondentes no grafo. Comparação
estrutural com 2171d72 confirmou a preservação dos desafios, fixtures, relações
anteriores e campos dos conceitos existentes. Nenhuma contagem mudou.

Os três exemplos usam somente a biblioteca padrão, sem rede ou arquivos.
Para reproduzir, extrair content.example.code de cada conceito novo para
example_test.go em diretório isolado, criar go.mod com module example e Go
1.25.0, executar go test ./... -count=1. Resultado: três exemplos aprovados.

A revisão local conferiu o destino **SizeError, a agregação de nil e a
preservação da categoria pública sem expor a identidade interna. O mapeamento
ilustrativo usa uma única categoria e declara esse limite. Não constitui
revisão humana, playtest nem publicação do conteúdo.

A validação estrutural passou após corrigir a indentação das duas relações
adicionadas. A relação de errors.As já existia e foi reutilizada.

Matriz strict: 23/23 checks aprovados, zero skips. Docs-check, ready-check,
check strict, knowledge-check strict e artifact-check strict passaram.
Avisos globais de atribuição histórica permanecem fora do lote.
