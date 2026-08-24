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

### Results summary
Spec destravada. Primeiro lote de conteúdo real autorado (2 de 44
desafios, pack go-first-steps) como checkpoint 1, aguardando revisão e
playtest do usuário antes de continuar os próximos lotes. Nenhum desafio
está `status: published` — `author: claude` já preenchido,
`reviewed_by`/`playtested` pendentes do playtest real.

### Requirement trace
Pendente — spec em progresso (checkpoint 1 de N; ver Execution log e
Results summary). O trace por R-ID será preenchido no fechamento, quando
os sete packs, as cinco trilhas e o playtest real estiverem completos
(Definition of Done desta spec).

### Known gaps
- Retenção e transferência serão comprovadas no aceite integrado.

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
