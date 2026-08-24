# ADR: Agent-generated catalog content as a first-class mode

## Status
Accepted

## Context

`curriculum-graph-path-recommendation` declara non-goal explícito "Gerar
desafios" e `agent-mcp-and-core-boundaries` fixa que o servidor MCP "não
chama LLM nem redige feedback aberto por conta própria" — nenhuma das duas
decisões proíbe geração de conteúdo para sempre; elas só nunca a colocaram
dentro do núcleo/servidor, que é determinístico por desenho.

`catalog-authoring-quality` (selada) já define o único caminho legítimo para
um item entrar publicado no catálogo: `author`, `reviewed_by` (pessoa
diferente do autor) e `playtested: true`, verificado por
`published_same_reviewer`/`published_without_playtest`. A decisão de fase 1
já tomada nesta conversa é: eu (agente, `author: claude`) autoro, o usuário
revisa (`reviewed_by`) — isso resolve quem preenche os campos, mas não onde
o conteúdo gerado por um agente-tutor **durante uma sessão real** (não numa
sessão de autoria offline) deveria morar antes de passar por esse gate.

`learning-track-composition` (ADR irmã, `multi-subject-selection-and-path-
composition`) cobre os três modos determinísticos de seleção sobre conteúdo
já existente: trilha completa, N-temas compostos, desafio único. O usuário
pediu explicitamente um quarto modo — geração assistida pelo agente
conectado ao MCP — disponível na mesma granularidade (trilha inteira nova,
trilha de N desafios, ou um desafio único mesmo sobre assunto já coberto) e
**escolhido deliberadamente**, não como fallback só quando o catálogo falha.

Módulos afetados: `internal/mcpserver` (nova superfície de tool), `internal/
curriculum` (armazenamento de rascunho), `internal/session` (proveniência de
evidência), `internal/mastery` (integridade de evidência), `.agents/skills/
codinho`.

## Decision

1. Geração por agente é um **modo de primeira classe**, oferecido lado a
   lado com os três modos determinísticos de `learning-track-composition` —
   nunca condicionado a "catálogo não cobre isso". O aluno pode pedir
   geração mesmo para um assunto já coberto por conteúdo existente.
2. O núcleo/servidor MCP continua sem LLM embutido (mantém
   `agent-mcp-and-core-boundaries`): quem gera o conteúdo é o agente-tutor do
   lado do host (Codex CLI/IDE), que já tem modelo; o MCP ganha uma tool nova
   só para **receber e persistir** o rascunho produzido pelo agente
   (`content_draft_submit` ou nome equivalente), validando-o com as mesmas
   regras estruturais/editoriais de `catalog-authoring-quality`
   (`internal/curriculum.Validate`/`RunEditorialChecks`) antes mesmo de
   aceitar o rascunho.
3. Todo conteúdo gerado por agente entra como
   `publication.status: draft`, em uma área de quarentena do catálogo,
   nunca como `published` direto — segue exatamente o mesmo funil de
   `author`/`reviewed_by`/`playtested` já definido, com `author` marcando a
   identidade do agente que gerou.
4. Uma sessão pode rodar **sobre um rascunho ainda não revisado** se o aluno
   pedir explicitamente (ele sabe que é conteúdo não curado) — mas a
   evidência dessa sessão é registrada com proveniência distinta
   (`content_provenance: draft`, campo novo no evento de mastery) e **nunca
   é promovida à ladder de maestria revisada** (`mastery-review-scheduling`
   R1–R2) até o rascunho subjacente virar `published`. Isso preserva o
   requisito não-funcional de `mastery-review-scheduling` ("cada mudança de
   domínio deve listar evidências causais") sem contaminar evidência
   confiável com conteúdo ainda não revisado por humano.
5. Rascunhos gerados alimentam o mesmo funil de curadoria já combinado:
   ficam disponíveis para eu promover a `reviewed_by` pendente e o usuário
   confirmar como segundo revisor quando a fase 2 (curadoria em duas
   camadas) entrar em vigor.

### Alternativas rejeitadas

- **Geração só como fallback quando busca/composição não encontram nada**:
  rejeitado a pedido explícito do usuário — reduz geração a caso de exceção
  em vez de opção deliberada, e não cobre o caso de gerar um desafio novo
  sobre assunto já coberto.
- **Aceitar rascunho gerado direto como `published`, confiando no julgamento
  do agente**: rejeitado porque recria exatamente o autocertificação que
  `catalog-authoring-quality` Decision 1 já rejeitou ("automação não deve
  fingir compreender pedagogia") — vale tanto para heurística quanto para
  geração por LLM.
- **Misturar evidência de rascunho na mesma ladder de maestria revisada**:
  rejeitado porque um card `retained`/`transferred` deixaria de significar
  "sobre conteúdo revisado", quebrando a garantia central de
  `mastery-review-scheduling`.
- **Embutir o modelo de geração dentro do servidor MCP** (em vez do
  agente-tutor do host): rejeitado por violar diretamente
  `agent-mcp-and-core-boundaries` Decision 2 e sua alternativa já rejeitada
  "embutir um LLM no servidor MCP".

## Consequences

- Nova tool MCP de submissão de rascunho e novo estado de quarentena no
  catálogo — spec própria (`agent-authored-catalog-drafts`), depende de
  `catalog-authoring-quality`, `learning-track-composition`,
  `tutor-skill-host-integration`.
- `internal/mastery` ganha um campo de proveniência por evidência; specs que
  o tocam precisam declarar essa extensão como aditiva (compatibilidade já
  exigida por `mastery-review-scheduling`).
- A skill `codinho` precisa expor os quatro modos (trilha completa, N-temas
  composto, desafio único, geração por agente) como escolha explícita no
  fluxo de início de sessão, não como texto livre implícito.
- O funil de autor/revisor combinado nesta conversa (fase 1: eu autor, você
  revisor; fase 2: eu primeiro revisor, você segundo) passa a valer também
  para rascunhos nascidos em sessão real, não só para autoria offline dos
  packs `go-*`.

### Gatilho de revisão

Revisar esta decisão se o volume de rascunhos gerados em sessão real
crescer além do que a fase 1 (revisor único) consegue processar, exigindo
adiantar a curadoria em duas camadas da fase 2; ou se `playtest` real
mostrar que sessões sobre rascunho não revisado confundem o aluno sobre o
que é "conteúdo oficial" do catálogo.
