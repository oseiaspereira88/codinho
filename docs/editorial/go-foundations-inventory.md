---
title: Inventário editorial dos fundamentos de Go
doc_type: reference
---

# Inventário editorial dos fundamentos de Go

Use a [fila estruturada](go-foundations-queue.json) para distribuir trabalho da
[spec go-foundations-packs](../../.pose/specs/2026-08-22-go-foundations-packs.md).
Este documento registra planejamento e baseline, não conclusão de autoria ou
publicação. Aplique o [checklist de revisão](../content-review-checklist.md) e
as [regras de autoria](../content-authoring.md) a cada tarefa.

## Baseline de 2026-09-11

Medição no commit `9caef5c1a585813076d4a4bc3957d71c0e1a2fe2`: leitura YAML
com `gopkg.in/yaml.v3` dos arquivos do manifest, contando listas de conceitos,
competências e desafios e percorrendo recursivamente `layers` para contar
nodes com `kind: macro`, `meso` ou `micro`. O nome da propriedade
`macro_steps` não transforma um node declarado `micro` em macro.

| Pack | Conceitos | Competências | Entradas | Nodes de bases | Nodes de variantes |
|---|---:|---:|---:|---:|---:|
| go-first-steps | 17 | 8 | 7 | 22 | 0 |
| go-core | 9 | 8 | 10 | 33 | 0 |
| go-data-text | 8 | 8 | 10 | 25 | 0 |
| go-type-design | 6 | 6 | 8 | 30 | 0 |
| go-errors | 12 | 10 | 7 | 35 | 0 |
| go-io | 10 | 11 | 6 | 12 | 24 |
| go-testing | 8 | 5 | 5 | 10 | 29 |
| Sete packs fundamentais | 70 | 56 | 53 | 167 | 53 |
| go-debugging, fora do escopo | 1 | 4 | 1 | 6 | 0 |
| Manifest completo | 71 | 60 | 54 | 173 | 53 |

O total global é 226 nodes: 14 macro, 14 meso e 198 micro. Os sete packs
somam 220 nodes. Excluir o protótipo de filtro deixa 219 nodes nos desafios
fundamentais e variantes: 166 nas bases e 53 nas oito variantes. Esses
números não comprovam profundidade pedagógica; muitas bases ainda declaram
apenas micro, sem decomposição macro/meso real.

Há 45 entradas de base nos sete packs, incluindo
`go-data.slice-filter-preserve-input`, exemplo histórico de um node em
go-first-steps. Excluí-lo conceitualmente deixa exatamente os 44 desafios
previstos: 32 atômicos, 10 combinados e duas fatias funcionais. O outro
protótipo é `go-debug.slice-off-by-one`, em go-debugging. A classificação
persistida e o destino de ambos precisam de disposição explícita; não
apague IDs, não modifique sessões e não reclassifique conteúdo silenciosamente
para ajustar a contagem. A tarefa `disposition-count-prototypes` prepara a
decisão humana e registra seu resultado quando existir.

## Cobertura da fila

A fila contém 70 tarefas identificadas: cinco correções do piloto, 44 auditorias
de bases fundamentais, oito auditorias de variantes, um planejamento de
expansão, sete lotes conceituais por pack e cinco tarefas de consolidação.
Cada auditoria trata o desafio integralmente, incluindo fixtures, árvore,
pistas, reflexões, conceitos, critérios e evidência. Tarefas de auditoria
podem confirmar conteúdo correto; não exigem mudanças cosméticas.

Execute primeiro o piloto, na ordem da fila:

1. Prove chave presente com valor vazio em `LoadValue`.
2. Prove zero encontrado com `found: true` em `FindFirstAtLeast`.
3. Prove consulta sem panic em mapa nil em `Score`.
4. Prove validação depois da aplicação de todas as opções em `NewServer`.
5. Alinhe a referência `RunCommands` ao contrato de switch e break rotulado.

Consuma a integração aprovada do piloto antes de iniciar auditorias que
toquem os mesmos desafios. Conclua `plan-concept-depth` antes de distribuir
`concepts-*`. Use o plano para alocar pelo menos 30 conceitos adicionais
nos sete packs, quatro competências fundamentais e as decisões pedagógicas
necessárias à meta de nodes; recalcule o déficit depois das auditorias.
Não use as quatro competências de go-debugging para satisfazer uma meta
dos sete packs. A expansão deve elevar go-type-design de seis para pelo
menos oito conceitos, respeitando o piso registrado na spec.

Revise todos os lotes antes de executar `audit-global-gates`. Execute
`prepare-human-review-playtest` depois da consolidação técnica e
`prepare-publication-closeout` somente com a disposição explícita dos gates.
Essas dependências são de planejamento; a fila não afirma que sua ordem
sozinho as impõe. O coordenador deve despachar somente tarefas prontas.

## Regras de execução e aceite

Configure tamanho de lote, paralelismo, modelo e reasoning ao iniciar o
mecanismo. Use cinco tarefas e `gpt-5.6-luna` com reasoning `high` no piloto.
Considere `paths` o limite de escrita; leia referências sem ampliar esse
limite. Várias tarefas editam desafios distintos no mesmo YAML: use isolamento
por lote e integre serialmente com revisão das diferenças para evitar perda
de conteúdo. Não autorize publicação nem commits pelo autor apenas por
constarem como etapas futuras nesta fila.

Registre tarefa, revisão, comandos executados, resultado e commit integrado.
Devolva achados corrigíveis à mesma sessão do autor e preserve pendências
externas como pendências. Não confunda saída bem-sucedida do agente com
aprovação do revisor primário. Faça checks do catálogo e das fixtures em cada
lote aprovado e execute a matriz POSE aplicável na integração.

Mantenha a spec `in-progress` enquanto faltar qualquer requisito: 44 desafios
com disposição documentada, pelo menos 100 conceitos, 60 competências e 300
nodes contextualizados, cinco trilhas, cobertura dos temas A–M, gates de
conteúdo/publicação e revisão/playtest reais. Reconciliar contagens não
dispensa revisão de duplicações, hints 1–6, leaks ou contextos insuficientes.

A pré-revisão automatizada prepara o julgamento humano. Não preencha
`reviewed_by`, `playtested: true` ou `status: published` com base na execução
dos agentes. A revisão humana por pessoa diferente do autor e o playtest de
ponta a ponta seguem a decisão vigente na spec e no checklist; registre a
ausência de evidência honestamente.
