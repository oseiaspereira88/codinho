---
slug: go-foundations-packs
status: in-progress
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: curriculum-graph-path-recommendation, catalog-authoring-quality
priority: 180
components: curriculum-content
delivers:
---

# Spec: go-foundations-packs

## 1. Intent

### Goal
Publicar os sete packs fundamentais de Go com cobertura profunda de linguagem, dados, tipos, erros, I/O e testes.

### Business value
Oferecer prática desde o primeiro contato até a base necessária para backend, com microdecomposição real.

### Constraints
- Publicar exatamente 44 desafios: 32 atômicos, 10 combinados e 2 fatias funcionais.
- Entregar pelo menos 100 conceitos, 60 competências e 300 step nodes únicos ou contextualizados.
- Todo conteúdo passa pelos gates editoriais e playtest.

### Non-goals
- Concorrência, HTTP, banco, produção e entrevistas.
- Exercícios que imponham arquitetura em camadas.

## 2. Requirements

### Functional
- R1: Publicar go-first-steps, go-core, go-data-text, go-type-design, go-errors, go-io e go-testing.
- R2: Cobrir temas A–M aplicáveis: tooling, tipos, fluxo, funções, coleções, texto, structs, interfaces, generics, erros, packages, I/O e testes.
- R3: Entregar 32 desafios atômicos com uma competência principal cada.
- R4: Entregar 10 combinados que integrem duas a quatro competências.
- R5: Entregar 2 fatias funcionais: importador validado e pequeno catálogo em memória testado.
- R6: Fornecer macro, meso e micro suficientes para iniciante, sem duplicar desafios por nível.
- R7: Incluir pistas 1–6, reflexões, variantes e checks seguros onde aplicável.
- R8: Publicar cinco trilhas: Go do zero, transição orientada a objetos, fluência prática, dados e tipos, testes e design.
- R9: Demonstrar cobertura de nil, zero values, aliasing, UTF-8, method sets, interface nil, wrapping e testabilidade.
- R10: Executar cada desafio por reviewer diferente do autor.

### Non-functional
- Distribuição por tema não pode deixar uma competência essencial com apenas um contexto.
- Briefings devem ser curtos e não conter assinatura completa salvo política.

### Security
- Fixtures não contêm rede, secrets ou paths externos.

### Compatibility
- Conteúdo declara versão mínima de Go e evita dependência externa sem propósito pedagógico.

## 3. Technical Plan

### Affected areas
- packs/go-first-steps/, packs/go-core/, packs/go-data-text/, packs/go-type-design/, packs/go-errors/, packs/go-io/, packs/go-testing/

### Artifacts
- created: packs/go-first-steps/
- created: packs/go-core/
- created: packs/go-data-text/
- created: packs/go-type-design/
- created: packs/go-errors/
- created: packs/go-io/
- created: packs/go-testing/
- created: testdata/packs/foundations/
- created: docs/catalog/go-foundations.md

### Delivery targets
Nenhum tipado; conteúdo consumido pelo catálogo V1.

### API/contract changes
- Popular schema v1 sem alterá-lo; mudanças necessárias exigem retorno à spec de schema.

### Data/storage changes
- Adicionar conteúdo YAML, fixtures, testes e evidência editorial.

### Technical risks
- Volume pode produzir micropassos mecânicos.
- Conteúdo introdutório pode ser condescendente para profissionais experientes.

## 4. Tasks

### Planning
- [x] Criar matriz 100 conceitos versus 60 competências e desafios.
- [x] Distribuir contagens por pack antes da autoria.

Matriz de distribuição (temas A–M de `PROJECT.md` §14.2 por pack):

| Pack | Temas A–M | Desafios previstos |
|---|---|---|
| go-first-steps | A (tooling), B (declarações e tipos) | 6 atômicos |
| go-core | C (fluxo), D (funções e métodos), K (packages e APIs) | 8 atômicos, 2 combinados |
| go-data-text | E (slices/arrays/maps), F (texto e dados binários) | 8 atômicos, 2 combinados |
| go-type-design | G (structs/ponteiros), H (interfaces), I (generics) | 6 atômicos, 2 combinados |
| go-errors | J (erros) | 4 atômicos, 2 combinados, 1 fatia funcional |
| go-io | L (I/O e serialização) | 4 atômicos, 1 combinado, 1 fatia funcional |
| go-testing | M (testes) | 4 atômicos, 1 combinado |

Total: 32 atômicos + 10 combinados + 2 fatias funcionais = 44 (R3–R5).
Conceitos/competências/step nodes seguem a mesma proporção por pack, sem
prever menos de 8 conceitos e 5 competências por pack (piso para não
deixar nenhum tema raso), ajustado durante a autoria para fechar em ≥100
conceitos / ≥60 competências / ≥300 step nodes (R1's threshold em
`catalog-authoring-quality` R9 `--v1-gate`).

### Implementation
- [ ] Autorar conceitos, competências e relações.
- [ ] Autorar os 44 desafios e 300 step nodes.
- [ ] Autorar pistas, reflexões, variantes e checks.
- [ ] Criar cinco trilhas.
- [ ] Realizar revisão, playtest e correções.

### Validation
- [ ] Executar codinho catalog validate packs.
- [ ] Executar todos os checks de fixtures.
- [ ] Auditar contagens, duplicação, leaks e cobertura.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Conteúdo fundamental precisa atender iniciante e experiente.
- Options considered: packs por senioridade; desafios duplicados; árvore completa com profundidade variável.
- Decision: uma árvore reutilizada por granularidade e tracks distintas.
- Rationale: evita drift e adapta apoio sem infantilizar.
- Consequences: cada desafio exige decomposição de alta qualidade.

### Decision 2
- Date: 2026-08-23
- Context: `catalog-authoring-quality` R7/R10 exige `author`, `reviewed_by`
  diferente de `author` e `playtested: true` antes de `status: published`;
  sem efetivo de autoria dedicado, a spec ficou bloqueada esperando quem
  assumiria esses papéis (achado ao final da sessão autônoma anterior).
- Options considered: (a) aguardar autoria humana externa; (b) fase 1 de
  autoria com `author: claude` (este agente) e `reviewed_by` = usuário,
  seguida por curadoria em duas camadas na fase 2 quando houver volume
  vindo do uso real.
- Decision: (b), a pedido explícito do usuário.
- Rationale: destrava a autoria sem contornar o gate — `playtested: true`
  continua exigindo um playtest real do usuário via `workspace prepare` +
  sessão MCP, nunca autocertificado; `status: published` só é setado após
  esse playtest realmente acontecer.
- Consequences: todo desafio novo autorado nesta spec declara
  `publication.author: claude` desde a criação; `reviewed_by` e
  `playtested: true` ficam pendentes até o usuário revisar/playtestar cada
  lote. Ver `.pose/adr/2026-08-23-agent-generated-catalog-content-as-a-first-class-mode.md`
  para o modelo mais amplo de autoria assistida por agente.

### Decision 3
- Date: 2026-08-23
- Context: sem efetivo de autoria, cada lote precisava do usuário revisando
  tudo em detalhe antes de qualquer outra evidência de qualidade existir;
  o usuário pediu explicitamente uma pré-revisão automatizada por um
  agente diferente (Codex CLI, `codex exec`, sandbox `workspace-write`,
  não-interativo) como padrão, ficando ele como segundo revisor mais fora
  do fluxo de criação.
- Options considered: (a) usuário revisa cada lote sozinho, em detalhe;
  (b) Codex como pré-revisor padrão (adversarial, independente do autor
  por ser outro fornecedor/modelo), usuário como revisor humano final,
  mais leve; (c) Codex substituindo `reviewed_by`/`playtested` (rejeitada
  sem chegar a ser implementada).
- Decision: (b).
- Rationale: `catalog-authoring-quality` Decision 1 já fixou que
  "automação não deve fingir compreender pedagogia" — isso não muda. Uma
  pré-revisão automatizada por um agente independente do autor pode achar
  problemas objetivos (vazamento de solução no texto disclosado,
  micropasso com duas intenções, critério de aceite não verificável)
  antes do humano gastar tempo nisso, mas não pode atestar honestidade de
  playtest nem substituir julgamento pedagógico humano. Validado nesta
  própria sessão: rodada 1 do Codex rejeitou
  `go-first-steps.convert-celsius-to-fahrenheit` por vazar `type Celsius
  float64` e a fórmula de conversão no `objective` (texto sempre
  disclosado ao aluno) e por um critério de aceite vago; rodada 2 pegou
  um vazamento semântico residual na correção; rodada 3 aprovou após o
  ajuste final.
- Consequences: todo lote futuro desta spec passa por `codex exec -s
  workspace-write --skip-git-repo-check` com um prompt de revisão
  referenciando `docs/content-review-checklist.md`/
  `docs/content-authoring.md` antes de ser proposto ao usuário. O veredito
  do Codex nunca preenche `publication.reviewed_by`/`playtested` — só o
  usuário faz isso, após playtest real via `workspace prepare` + sessão
  MCP. Nenhum arquivo é alterado pela pré-revisão (papel só de leitura +
  comandos de validação).
- Follow-up (2026-08-23): o padrão foi extraído para artefatos
  reutilizáveis entre sessões e independentes de par de agentes —
  `scripts/agent-review.sh` (primitiva new/resume), `docs/agent-review-
  workflow.md` (papéis, templates de prompt, por que retomar sessão em vez
  de recomeçar do zero a cada rodada) e a skill `agent-batch-review`. Ver
  `.pose/knowledge/2026-08-24-note-agent-batch-review-pattern.md`.

## 6. Validation

### Strategy
Validar estrutura, distribuição exata, checks, cobertura e playtest humano.

### Deterministic checks
- Test: go test ./internal/curriculum/... e checks de cada fixture.
- Lint: codinho catalog validate packs/go-first-steps packs/go-core packs/go-data-text packs/go-type-design packs/go-errors packs/go-io packs/go-testing
- Typecheck: go vet ./...
- Build: go build ./cmd/codinho
- Security / Contract: secret scan, path confinement e leak detector.

### Execution log
- 2026-08-23: destravada (Decision 2) — `author: claude` liberado como
  identidade de autoria da fase 1.
- 2026-08-23: primeiro lote autorado em `packs/go-first-steps.yaml` — 2
  desafios atômicos novos (`enumerate-weekdays-with-iota`,
  `convert-celsius-to-fahrenheit`), tema `declarations-and-types`, 4
  conceitos e 2 competências novos, relations `relates_to` para os itens
  novos.
- `go run ./cmd/codinho catalog validate` → exit 0, só avisos
  `relation_isolated` esperados (itens pré-existentes e conceitos novos
  seguindo o mesmo padrão já presente no pack).
- `go build ./...`, `gofmt -l .`, `go vet ./...`,
  `go test ./internal/curriculum/... ./internal/cli/... -race` → todos ok.
- Pré-revisão Codex rodada 1 (`codex exec -s workspace-write`): **rejeitou**
  `convert-celsius-to-fahrenheit` (vazamento literal de `type Celsius
  float64` e da fórmula `c*9/5+32` no `objective`; critério de aceite
  "preserva precisão" não verificável) e aprovou `enumerate-weekdays-with-
  iota` com ressalva menor (micropasso com duas intenções).
- Correções aplicadas: micropasso de weekdays separado em dois; objectives
  de Celsius reescritos sem tipo/fórmula literal; critério trocado por
  `CelsiusToFahrenheit(36.6)` com faixa numérica concreta.
- Pré-revisão Codex rodada 2: achou vazamento semântico residual
  ("ponto flutuante de 64 bits" ainda entregava `float64` combinado ao
  nome `Celsius`); `enumerate-weekdays-with-iota` **aprovado sem
  ressalvas**.
- Correção final: objective trocado para "tipo subjacente que suporte
  casas decimais", deixando a restrição "não usar float32" já existente
  fazer a inferência.
- Pré-revisão Codex rodada 3: **aprovado sem ressalvas** os dois
  desafios. `go run ./cmd/codinho catalog validate` → exit 0 em todas as
  rodadas.
- 2026-08-24: checkpoint 2 autorado em `packs/go-first-steps.yaml` — 2
  desafios atômicos novos (`declare-a-minimal-module`, tema `tooling`;
  `clamp-int-to-byte`, tema `declarations-and-types`), 3 conceitos e 1
  competência novos, fixture real (`go.mod`/`clamp.go`/`clamp_test.go`)
  com check `go_test` executável para `clamp-int-to-byte`.
- Pré-revisão Codex rodada 1 (`new`): **rejeitou os dois** —
  `declare-a-minimal-module` por acceptance vago ("module path específico
  ao exercício", "go 1.25 ou superior" sem versão exata); `clamp-int-to-
  byte` por objective vazando o algoritmo completo (intervalo 0-255,
  substituição pelo limite mais próximo, ordem clamp-antes-de-converter) e
  por um micropasso com duas intenções (clamp + conversão juntos).
- Correções: valores exatos e concretos no `declare-a-minimal-module`
  (`module codinho-practice/declare-a-minimal-module`, `go 1.25.0`);
  micropasso de clamp dividido em dois, objective generalizado sem citar
  0-255 explicitamente, acceptance ganhou os casos de borda 0/255.
- Pré-revisão Codex rodada 2 (`new`): `declare-a-minimal-module`
  **aprovado**; `clamp-int-to-byte` **rejeitado** — os critérios
  comportamentais eram `kind: structural` sem nenhum check real por trás
  (achado técnico correto: `internal/cli/editorial.go` só executa checks
  quando há `fixture` + `checks` juntos; sem fixture, o check nem roda).
- Correção: adicionado fixture real (go.mod, stub com `panic("not
  implemented")`, `clamp_test.go` com os cinco casos do acceptance) mais
  o bloco `checks` apontando para `TestClampToByte`. Verificado
  manualmente antes de qualquer nova rodada: `workspace prepare` +
  `go test` reproduz panic contra o stub e passa contra uma implementação
  correta de referência.
- Pré-revisão Codex rodada 3 (`new`, pedindo pra ele mesmo materializar e
  rodar o teste): travou mais de 1h30 sem concluir — morta
  (`TaskStop`). **Lição que já entrou em `docs/agent-review-workflow.md`**:
  a partir daqui, toda rodada seguinte do mesmo lote usa `codex exec
  resume --last` em vez de `new`.
- Pré-revisão Codex rodada 4 (`resume --last`, salvando a sessão travada):
  confirmou `declare-a-minimal-module` aprovado; para `clamp-int-to-byte`
  achou que seu próprio sandbox tinha `GOCACHE` apontando para
  `~/.cache/go-build` (read-only fora de `workdir`/`/tmp`/`$TMPDIR`) — não
  um problema do conteúdo.
- Pré-revisão Codex rodada 5 (`resume --last`, pedindo `GOCACHE=/tmp/...`):
  rodou `workspace prepare` + `go test` de verdade, confirmou panic contra
  o stub, **aprovado sem ressalvas** `clamp-int-to-byte`. Custo: ~6,6 mil
  tokens, segundos de execução — contra ~60-100 mil tokens e minutos das
  rodadas `new`.
- `go run ./cmd/codinho catalog validate` e `catalog validate --checks` →
  exit 0 em todas as rodadas do checkpoint 2. `go build ./...`, `gofmt -l
  .`, `go vet ./...` → ok.
- 2026-08-24: checkpoint 3 autorado em `packs/go-first-steps.yaml` — 2
  desafios atômicos novos (`format-rate-limit-message`, tema `tooling`,
  fixture + checks `go_test`+`go_vet`; `avoid-shadowing-named-returns`,
  tema `declarations-and-types`, fixture + check `go_test` reproduzindo
  um bug real de shadowing de retorno nomeado dentro de bloco condicional
  — verificado manualmente: versão com `:=` no bloco falha 2 dos 3 casos
  silenciosamente, versão com `=` passa todos).
- Antes da pré-revisão: uma rodada `new` travou (>1h30, morta) e o
  processo `codex` ficou **órfão** — `TaskStop` matou só o wrapper, não o
  processo real. A rodada seguinte travou de novo (>30min) por disputa de
  lock de sessão com o órfão. Matar o PID do `codex` direto resolveu.
  Lição registrada em `docs/agent-review-workflow.md` e na nota de
  knowledge.
- Pré-revisão Codex rodada 1 (`new`, só estática): achou um **bug real**
  — `packs/go-first-steps.yaml` linha 1 estava corrompida para
  `/schema_version: 1` (barra a mais), fazendo `catalog validate` falhar
  com exit 1 (schema_version lido como 0). Corrigido para `schema_version:
  1`. Também: (a) `format-rate-limit-message` aprovado com ressalvas
  menores — escopo não cobria imports, e o check permitia passar por
  concatenação/hardcode sem exercitar `fmt.Sprintf`; (b) `avoid-shadowing-
  named-returns` rejeitado — a constraint entregava a correção exata ("use
  ="), a ausência de shadowing não era um critério observável (só um
  efeito colateral do teste), micropasso com 3 intenções, escopo sem
  imports; (c) ambos tinham um micropasso redundante de "declarar
  assinatura" já satisfeito pela fixture.
- Correções: constraint reescrita para vetar `:=` sem prescrever `=`;
  critério novo `no-redeclaration-in-conditional` (structural,
  source_inspection) tornando shadowing observável; micropasso de
  implementação dividido em dois (validar vazio; converter e atribuir);
  micropassos redundantes de assinatura removidos nos dois desafios;
  escopos passaram a incluir "os imports necessários"; constraint +
  critério `uses-sprintf` adicionados em `format-rate-limit-message`
  fechando a brecha do check. Reverificado manualmente: stub falha, versão
  com bug de shadowing falha 2/3 casos, versão correta passa 4/4 (com o
  caso de borda novo).
- Pré-revisão Codex rodada 2 (`resume --last`): `format-rate-limit-
  message` **aprovado sem ressalvas**; `avoid-shadowing-named-returns`
  **aprovado com ressalvas menores** — primeiro micropasso (validar vazio)
  ainda sem "imports necessários" no escopo, e sugeriu caso de teste para
  string só com espaços.
- Correções finais: escopo do primeiro micropasso ajustado; caso
  `ParseAge("   ")` adicionado ao acceptance e ao `parseage_test.go`.
  Reverificado manualmente: passa com implementação correta.
- Pré-revisão Codex rodada 3 (`resume --last`): **aprovado sem
  ressalvas** os dois desafios. Custo da rodada final: ~8,4 mil tokens.
- `go run ./cmd/codinho catalog validate` e `catalog validate --checks` →
  exit 0. `go build ./...`, `gofmt -l .`, `go vet ./...` → ok.
- 2026-08-25: primeiro lote autorado em `packs/go-core.yaml` (pack novo) —
  2 desafios atômicos novos (`close-resources-in-defer-order`, tema
  `control-flow`; `sum-variadic-numbers`, tema `functions-and-methods`), 2
  conceitos e 2 competências novos, fixture + check `go_test` executável
  para cada um. Verificado manualmente antes da pré-revisão: stub de cada
  fixture falha (panic / erro de compilação), implementação de referência
  passa em todos os casos.
- Pré-revisão Codex rodada 1 (`new`): **rejeitou**
  `close-resources-in-defer-order` — `objective` do único macro_step
  vazava a técnica de implementação inteira ("registrar com defer que
  será acrescentado ao retorno nomeado order"), o mesmo passo reunia 5
  intenções (percorrer, agendar closure, capturar recurso correto,
  acumular, retornar), e o critério `uses-defer` (`kind: structural`) não
  tinha nenhuma evidência executável por trás — um slice invertido sem
  qualquer defer passaria no `TestClosingOrder`. `sum-variadic-numbers`
  **aprovado com ressalvas menores**: reflexão pedia "o quê" em vez de
  "por quê", e a cobertura de teste não incluía sinais mistos/negativos.
- Correções: `close-resources-in-defer-order` dividido em dois
  macro_steps — agendar a ação diferida por recurso (critério
  `uses-defer`, sem tratar retorno ainda) e formar o retorno a partir da
  execução real dos defers (critério `order-is-lifo`, com evidência
  comportamental real); nenhum dos dois objectives cita `append` na
  closure. Verificado em `internal/checks/model.go` e
  `internal/checks/executor.go` que o único runner estrutural
  (`internal_ast`) só valida parse válido, sem inspeção de palavra-chave
  — não existe hoje um jeito de dar evidência executável dedicada a "usa
  defer" sem mudar o motor (fora do escopo desta spec de conteúdo);
  `uses-defer` manteve `source_inspection`, mesmo padrão já usado em
  `uses-sprintf`/`signature-is-variadic`. `sum-variadic-numbers`: reflexão
  reescrita para "por que Sum() sem argumentos retorna 0 naturalmente",
  caso `Sum(-1, 5, -4)` (sinais mistos) adicionado ao acceptance e ao
  `sum_test.go`, reverificado manualmente (passa).
- Pré-revisão Codex rodada 2 (`resume --last`): **aprovado sem
  ressalvas** os dois desafios, incluindo concordância explícita com o
  argumento de que `uses-defer` via `source_inspection` segue o mesmo
  contrato já aceito para `uses-sprintf`/`signature-is-variadic` — a
  ausência de um runner dedicado a palavras-chave é uma limitação
  sistêmica do motor, não um bloqueio deste conteúdo.
- `go run ./cmd/codinho catalog validate` e `catalog validate --checks` →
  exit 0 nas duas rodadas. `go build ./...`, `gofmt -l .`, `go vet ./...`
  → ok.

### Results summary
Spec destravada. Quatro checkpoints de conteúdo real autorado (8 de 44
desafios): três em go-first-steps (checkpoint 1 declarações/tipos,
checkpoint 2 tooling + conversão numérica, checkpoint 3 tooling +
shadowing) e o primeiro checkpoint de um segundo pack, go-core (defer +
parâmetros variádicos) — todos com fixture/checks executáveis reais e
pré-revisão automatizada aprovada sem ressalvas (Decision 3), aguardando
revisão humana final e playtest do usuário antes de continuar os próximos
lotes. Nenhum desafio está `status: published` — `author: claude` já
preenchido, `reviewed_by`/`playtested` pendentes do playtest real. O
checkpoint 3 também validou o padrão de retomar sessão sob condições reais
adversas (rodada travada, processo órfão) e achou um bug estrutural real
no pack (schema_version corrompido) — evidência de que a pré-revisão
compensa mesmo quando a mudança parece só de conteúdo. O checkpoint 1 de
go-core achou e corrigiu um vazamento de solução real (objective
prescrevendo a técnica de implementação completa) e deixou registrado um
limite sistêmico do motor (nenhum runner allowlisted verifica presença de
palavra-chave como `defer`; critérios sintáticos ficam por
`source_inspection`, mesmo padrão já aceito alhures).

### Requirement trace
Pendente — spec em progresso (checkpoint 1 de N; ver Execution log e
Results summary). O trace por R-ID será preenchido no fechamento, quando
os sete packs, as cinco trilhas e o playtest real estiverem completos
(Definition of Done desta spec).

### Known gaps
- Retenção e transferência serão comprovadas no aceite integrado.
- `go-first-steps.clamp-int-to-byte` (checkpoint 2, já aprovado/comitado)
  tem o mesmo micropasso redundante de "declarar assinatura" identificado
  e corrigido no checkpoint 3 — não foi revisitado; considerar corrigir
  num lote futuro antes do playtest humano.

## 7. Final Report

### Delivered scope
Nenhum; conteúdo planejado.

### Files and modules changed
- Planejados nos sete packs, testdata e documentação de catálogo.

### Validation executed
- Command: pose lint-spec go-foundations-packs --ready-check
- Result: registrar após gate.

### Residual risks
- A meta quantitativa não substitui revisão editorial.

### Follow-ups
- [covered: v1-integrated-acceptance] Validar sessões reais com amostra dos sete packs.
