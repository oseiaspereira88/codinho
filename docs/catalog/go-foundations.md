---
title: Catálogo fundamental de Go
doc_type: reference
---

# Catálogo fundamental de Go

Inventário de autoria em 2026-09-09. A spec `go-foundations-packs` permanece
in-progress: estes números não atestam publicação nem playtest humano.
A = atômico; C = combinado; F = fatia funcional. Nodes contam macro/meso/micro
recursivamente, sem incluir o nó de camada.

| Pack | Estado | Desafios | Conceitos | Competências | Nodes |
|---|---|---|---|---|---|
| go-first-steps | draft | 7 (7 A / 0 C / 0 F) | 17 | 8 | 14 |
| go-core | draft | 10 (8 A / 2 C / 0 F) | 9 | 8 | 23 |
| go-data-text | draft | 10 (8 A / 2 C / 0 F) | 8 | 8 | 22 |
| go-type-design | draft | 8 (6 A / 2 C / 0 F) | 6 | 6 | 16 |
| go-errors | draft | 5 (4 A / 1 C / 0 F) | 4 | 4 | 10 |
| go-io | draft | 2 (0 A / 1 C / 1 F) | 10 | 6 | 12 |
| go-testing | draft | 1 (0 A / 1 C / 0 F) | 8 | 5 | 7 |

Total atual: 33 atômicos, 9 combinados e 1 fatia funcional;
62 conceitos, 45 competências e 104 step nodes.
A meta segue 32 + 10 + 2 desafios, 100 conceitos, 60 competências e 300 nodes.

## Lote go-io

- `go-io.decode-strict-config`: documento JSON único, campos conhecidos,
  validação semântica, leitura fragmentada e preservação de erros da entrada.
- `go-io.import-csv-inventory`: cabeçalho, CSV com aspas/UTF-8, nomes únicos,
  quantidades positivas, serialização ordenada e propagação de erros. Nenhuma
  escrita ocorre antes da validação completa; uma falha do Writer pode deixar
  saída parcial, pois uma interface Writer genérica não oferece rollback.

Os micropassos intermediários usam avaliação qualitativa da parte alterada,
com observação e referência de rubrica explícitas via evidence_record. O último
passo usa o check completo; os testes finais não bloqueiam artificialmente
uma implementação ainda parcial nos passos anteriores.

Cada desafio tem seis níveis de apoio declarados, reflexões, decomposição
macro/meso/micro e fixture inicial que falha em testes comportamentais. A
referência privada deve passar pelos mesmos testes. Não há rede nem dependência
externa; a versão mínima da fixture é Go 1.25. A autoria deste lote é `codex`,
sem reviewed_by nem playtested. As quatro variantes atômicas planejadas para
I/O ainda não foram autoradas; não estão sendo contabilizadas como entregues.

## Pendências para a entrega fundamental

- Resolver a divergência de `go-first-steps`: sete atômicos existentes ante seis
  previstos. O conteúdo extra foi preservado; não se alterou sua classificação
  para ajustar artificialmente a meta. Concluir os demais lotes como planejados
  produziria 45 desafios, então a distribuição precisa de decisão explícita.
- Completar um combinado e uma fatia funcional em `go-errors`; aprofundar conceitos/competências e decomposição existentes.
- Autorar variantes contextualizadas, completar as cinco trilhas previstas e
  conferir cobertura cruzada. O lote go-io não adiciona trilha isolada que
  substitua uma das cinco jornadas canônicas.
- Realizar pré-revisão pedagógica e playtest humano por desafio antes de marcar
  qualquer conteúdo como publicado. A aprovação mecânica das fixtures não
  comprova essa revisão.

Validação reproduzível: `go run ./cmd/codinho catalog validate --checks --json`.
O comando examina o catálogo completo do manifest e compara as expectativas
baseline/reference; execução real dos checks depende da autorização local.

## Lote go-testing

`go-testing.select-active-tokens-with-injected-clock` combina seleção temporal,
casos de fronteira e isolamento de slices. Um relógio injetado permite testar
antes, igualdade e depois sem esperar tempo real. A suíte inclui fuso diferente
para o mesmo instante, dependência ausente, entrada vazia, preservação da ordem,
contagem de chamadas e modificação da saída para detectar aliasing.

O aluno deve justificar e acrescentar um caso à tabela. Essa parte requer
observação qualitativa com rubrica: o check automatizado confirma o contrato,
mas não certifica a qualidade de novos testes nem a autoria do aluno.
Os quatro cenários variantes planejados ainda não foram entregues.
