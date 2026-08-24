# ADR: Multi-subject selection and path composition

## Status
Accepted

## Context

`curriculum-graph-path-recommendation` (selada) modela seleção como um valor
único: `curriculum.Query.Theme`/`Query.Competency` filtram por um tema ou uma
competência por chamada, e `recommendation.Objective` carrega exatamente um
`CompetencyID` ou `ThemeID`. `session_start` (`mcp-stdio-foundation`,
`session-orchestration-disclosure`) exige `challenge_id` de um item já
existente; a spec de sessão promete "desafio ou trilha" (R1), mas nenhuma
noção de trilha (`track_id`) chegou a ser implementada.

Na prática, o aluno quer escolher entre modos de granularidade bem
diferentes: (a) uma trilha completa pré-existente de uma linguagem (V1: Go),
percorrendo todos os seus desafios em ordem; (b) um ou mais temas/competências
específicos que quer praticar; (c) um único desafio já conhecido. Quando (b)
aponta para múltiplos assuntos e nenhum desafio único do catálogo cobre todos
eles ao mesmo tempo, o sistema hoje não tem como sinalizar isso — devolve
poucos resultados ou nenhum, sem declarar que a cobertura pedida é parcial.

Módulos afetados: `internal/curriculum` (search, graph), `internal/
recommendation`, `internal/session`, `internal/mcpserver` (session_tools,
catalog_search, learning_path_recommend).

## Decision

1. `curriculum.Query` e `recommendation.Objective` passam a aceitar um
   **conjunto** de temas/competências (`ThemeIDs`/`CompetencyIDs`), não mais
   um valor único; o campo singular anterior é removido, não mantido em
   paralelo (sem shim de compatibilidade).
2. Toda resposta de busca/recomendação multi-assunto declara cobertura
   explícita — `total` (um item cobre todos os assuntos pedidos), `parcial`
   (nenhum item cobre tudo, mas existe um caminho de N itens que cobre a
   união) ou `nenhuma` — nunca silenciosamente incompleta.
3. `session_start` ganha um modo de trilha completa pré-existente
   (`track_id`), fechando a lacuna entre a spec de sessão (R1) e o código: o
   aluno percorre todos os desafios de uma trilha autorada, em ordem, como um
   lifecycle de sessão só de leitura sobre uma sequência já fixada — igual ao
   documento raiz da funcionalidade V1 sugerir (trilha de Go completa).
4. Quando a cobertura é parcial, a recomendação **compõe uma trilha com os N
   desafios existentes** que juntos cobrem os assuntos pedidos (ordenados
   pelo grafo de precedência), em vez de tentar sintetizar um único item
   artificial — cada desafio da composição continua sendo conteúdo real,
   revisado e versionado, preservando a integridade da evidência de maestria
   (`mastery-review-scheduling` depende de evidência ligada a conteúdo
   revisado, não inventado).
5. A escolha entre trilha completa, N-temas compostos ou desafio único
   continua sendo sempre do aluno — nenhum modo é selecionado
   automaticamente no lugar de outro; a recomendação apenas explica as
   opções disponíveis e sua cobertura (mantém o caráter consultivo de
   `curriculum-graph-path-recommendation` Constraint "sugestões são
   consultivas").

### Alternativas rejeitadas

- **Manter seleção só por valor único e deixar o cliente MCP intersectar
  múltiplas buscas por conta própria**: rejeitado porque empurra lógica de
  composição e o sinal de cobertura para cada agente-tutor reimplementar,
  produzindo resultados inconsistentes entre hosts.
- **Sintetizar um único item "virtual" concatenando trechos de vários
  desafios reais**: rejeitado porque quebraria a granularidade micro/meso/
  macro definida em `session-orchestration-disclosure` e misturaria
  evidência de proveniências diferentes num único nó de maestria.
- **Adicionar `track_id` sem sinal de cobertura para buscas multi-tema**:
  rejeitado por resolver só o caso (a) e deixar o caso (b)/(c) tão silencioso
  quanto hoje.

## Consequences

- `internal/recommendation.Objective` e `internal/curriculum.Query` mudam de
  assinatura (breaking, mas internos — nenhum consumidor externo além do MCP
  server); specs que os tocam (`curriculum-graph-path-recommendation`,
  `session-orchestration-disclosure`) precisam de uma spec de extensão
  própria (`learning-track-composition`) em vez de reabrir o fechamento
  selado.
- O contrato MCP ganha `track_id` em `session_start` e um campo de cobertura
  em `catalog_search`/`learning_path_recommend`; hosts existentes (Codex CLI,
  skill) precisam ser atualizados para exibir os três modos como escolhas
  explícitas.
- Determinismo é preservado: cobertura e composição de trilha continuam
  sendo cálculo de grafo, sem LLM nem heurística não determinística — mantém
  a fronteira já registrada em `agent-mcp-and-core-boundaries`.

### Gatilho de revisão

Revisar esta decisão se a composição de trilha por N assuntos produzir
sequências pedagogicamente ruins em playtest real (assuntos numa ordem que
não respeita pré-requisitos percebidos), ou se o número de assuntos
simultâneos pedidos em uso real ultrapassar o que uma trilha de tamanho
razoável consegue cobrir sem virar maratona.
