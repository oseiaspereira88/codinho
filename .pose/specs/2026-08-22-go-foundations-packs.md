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
- packs/*.yaml, fixtures embutidas nos desafios e docs/catalog/go-foundations.md

### Artifacts
- modified: docs/agent-review-workflow.md
- modified: docs/content-review-checklist.md
- modified: packs/go-first-steps.yaml
- created: packs/go-core.yaml
- created: packs/go-data-text.yaml
- created: packs/go-type-design.yaml
- created: packs/go-errors.yaml
- created: packs/go-io.yaml
- created: packs/go-testing.yaml
- modified: packs/manifest.yaml
- created: docs/catalog/go-foundations.md
- modified: .pose/docs.json
- modified: .pose/indexes/validation-matrix.json
- created: cmd/codinho/foundation_tracks_integration_test.go
- created: internal/curriculum/foundation_variants_test.go

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
| go-io | L (I/O e serialização) | 1 combinado, 1 fatia funcional; 4 variantes atômicas sem nova contagem |
| go-testing | M (testes) | 1 combinado; 4 variantes atômicas sem nova contagem |

Total: 32 atômicos + 10 combinados + 2 fatias funcionais = 44 (R3–R5).
Correção de planejamento em 2026-09-07 UTC: a tabela anterior somava 40
atômicos e 52 desafios, em contradição com R3 e com a distribuição global.
Os 32 atômicos já previstos nos cinco primeiros packs são preservados;
as oito práticas adicionais de I/O/testes passam a variantes contextualizadas,
sem contar como desafios canônicos novos. Validar equivalência de competência
e relações na matriz editorial antes de autorar; se precisarem ser desafios
independentes, propor revisão explícita da distribuição, sem renomear tipos
artificialmente nem apagar conteúdo existente.
Conceitos/competências/step nodes seguem a mesma proporção por pack, sem
prever menos de 8 conceitos e 5 competências por pack (piso para não
deixar nenhum tema raso), ajustado durante a autoria para fechar em ≥100
conceitos / ≥60 competências / ≥300 step nodes (R1's threshold em
`catalog-authoring-quality` R9 `--v1-gate`).

### Implementation
- [x] Completar os oito conceitos sem conteúdo em go-core em dois lotes (fluxo/funções e fronteiras/APIs), com exemplos isolados, relações e preservação estrutural dos desafios; executar exemplos e matriz strict. Consumir knowledge:go-foundations-io-batch.
- [x] Aprofundar os oito conceitos existentes de go-data-text com explicações, exemplos Go isolados e relações conceituais válidas; preservar desafios e competências e validar catálogo e matriz strict. Consumir knowledge:go-foundations-io-batch.
- [x] Aprofundar os seis conceitos existentes de go-type-design com explicações, exemplos Go isolados e relações conceituais válidas; preservar desafios e competências e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Aprofundar quatro conceitos iniciais de go-first-steps (declaração curta, slice, tipos nomeados e zero values), com exemplos isolados, relações válidas e preservação dos desafios; validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Aprofundar quatro conceitos de go-first-steps sobre constantes, iota, conversões e módulos, com exemplos isolados e relações válidas; preservar desafios e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Aprofundar quatro conceitos de go-first-steps sobre versão, overflow, formatação e retornos nomeados, com exemplos isolados e relações válidas; preservar desafios e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Completar os quatro conceitos restantes de go-first-steps (`range`, structs, tipos de campo e escopo), com exemplos isolados e relações válidas; preservar desafios e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Aprofundar o conceito de limites de slice em go-debugging com exemplo isolado e preservar o desafio de depuração; validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Tornar explícitas três competências de investigação já cobertas pelos nodes de go-debugging e associá-las ao desafio existente, sem criar desafio canônico; validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Adicionar cinco competências derivadas dos nodes existentes de go-io e associá-las aos desafios correspondentes, sem criar desafios canônicos; validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir o fluxo de go-debugging com nodes de verificação da saída e análise da fronteira vazia, preservando o desafio canônico; validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir a variante `go-io.read-fragmented-note` com nodes de verificação para leituras parciais e fronteira EOF; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir a variante `go-io.decode-weather-fields` com nodes de verificação para documento nulo e valor posterior; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir a variante `go-io.read-quoted-attendees` com nodes de verificação para vírgulas entre aspas e atomicidade em registro inválido; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir a variante `go-io.write-result-report` com nodes de verificação para newline final e identidade do erro do Writer; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir o desafio `go-testing.select-active-tokens-with-injected-clock` com nodes de verificação para resultado vazio não nulo e preservação da ordem; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir a variante `go-testing.count-deadline-clock-calls` com nodes de verificação para relógio ausente e observação única do instante; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir a variante `go-testing.measure-reservation-remaining` com node de verificação para relógio nil e sentinela de erro; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir a variante `go-testing.check-ticket-boundary` com node de verificação para conjunto misto de expirações; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir o desafio `go-testing.select-active-tokens-with-injected-clock` com node de verificação para equivalência de fusos; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir `go-core.close-resources-in-defer-order` com node de verificação para entrada vazia e retorno nil; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Revisar o lote de aprofundamento até 13e9dc1: corrigir instruções e referências conceituais e implementar regressões para bytes com EOF, duração com relógio nil e entrada vazia alocada em defer. Relatório: `.pose/reports/2026-09-10-catalog-luna-review.md`. Consumir knowledge:go-foundations-io-batch.
- [x] Corrigir os 51 avisos de passos avaliados sem hints e separar a verificação bidirecional de snapshot em duas intenções, com critérios próprios. Verificar catálogo e matriz strict. Consumir knowledge:go-foundations-io-batch.
- [x] Conectar os 32 itens isolados do grafo com 39 relações applies_in fundamentadas nas competências e conceitos usados pelos desafios; validar catálogo e matriz strict. Consumir knowledge:go-foundations-io-batch.
- [x] Aprofundar Sum variádico com comparação entre chamadas diretas e expansão de slice, fixtures para ambas as formas e reflexão sobre o uso de reticências; validar catálogo e matriz strict. Consumir knowledge:go-foundations-io-batch.
- [x] Aprofundar `go-core.stop-processing-commands-with-labeled-break` com verificação de que comandos após o primeiro `stop` são ignorados; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Aprofundar `go-core.increment-counter-with-pointer-receiver` com verificação de estado independente entre instâncias; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Aprofundar `go-core.recover-from-panic-in-safe-call` com verificação da conversão de valores panic em mensagens de erro úteis; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Aprofundar `go-core.guard-invariant-with-unexported-field` com verificação dos limites inclusivos e valores adjacentes inválidos; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Aprofundar `go-core.find-first-value-at-least` com verificação da distinção entre valor zero encontrado e ausência de correspondência; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Aprofundar `go-core.hide-counter-type-behind-interface` com verificação de isolamento entre contadores construídos separadamente; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir o desafio de conversão estreitante de go-first-steps com nodes de verificação de limites e reflexão sobre a ordem das operações; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir os desafios de módulo e formatação de go-first-steps com nodes de inspeção e verificação observável; preservar fixtures e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir o desafio de shadowing e retornos nomeados com nodes de verificação de casos e reflexão sobre escopo; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir o desafio genérico de transformação de go-type-design com nodes de verificação de ordem e isolamento; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir o desafio de clonagem de struct de go-type-design com nodes de capacidade excedente e preservação de campos; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir o desafio de seleção de validador de go-type-design com nodes de identidade de interface e reflexão sobre estado nil; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir o desafio de receptor nil-safe de go-type-design com nodes de limites de cadeia e reflexão sobre a guarda; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir o desafio de assertion segura de go-type-design com nodes de verificação de tipos e reflexão sobre comma-ok; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir o desafio genérico de pertencimento de go-type-design com nodes de instanciações comparáveis e fronteiras de busca; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir o desafio combinado de clonagem e primeiro Circle de go-type-design com nodes de ordem de busca e isolamento do Registry; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir o desafio de classificação de erros de go-errors com nodes de identidade e precedência agregada; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir o desafio de agregação de erros de go-errors com nodes de identidade das causas e sucesso nil; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir o desafio de tradução de erros de go-errors com nodes de identidade e mensagem públicas; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir o desafio de extração tipada de go-errors com nodes de wrapping profundo e ausência de tipo; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir o catálogo em memória de go-errors com nodes de rejeição sem mutação e snapshots ordenados; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir a variante temporal de go-testing com nodes de duração não negativa e relógio chamado uma vez; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir a variante de snapshot de go-testing com nodes de isolamento bidirecional e vazio não nil; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir a variante de contagem de prazos de go-testing com nodes de lote vazio e limite estrito; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Expandir a variante de fronteira temporal de go-testing com nodes de igualdade e equivalência entre fusos; preservar fixture e validar catálogo. Consumir knowledge:go-foundations-io-batch.
- [x] Completar conteúdo dos oito conceitos restantes de go-errors com exemplos isolados, relações e preservação de desafios/fixtures; executar exemplos e matriz strict. Consumir knowledge:go-foundations-io-batch.
- [x] Complementar errors.As, errors.Join e tradução na fronteira em go-errors, preservando o conteúdo de %w e desafios; testar três exemplos isolados e matriz strict. Consumir knowledge:go-foundations-io-batch.
- [x] Aprofundar os dez conceitos existentes de go-io com explicações, exemplos Go isolados e referências ligadas ao grafo; preservar desafios/fixtures e verificar exemplos e matriz strict. Consumir knowledge:go-foundations-io-batch.
- [x] Aprofundar os oito conceitos existentes de go-testing com explanation/example e referências conceituais válidas; preservar desafios e fixtures. Validar exemplos Go isolados, catálogo e matriz strict. Consumir knowledge:go-foundations-io-batch.
- [x] Lote de variantes: quatro atômicas em go-io e quatro em go-testing, com variant_of, canonical false, uma competência compartilhada com a origem e fixtures próprias. Manter 45 bases de autoria; variantes não suprem o déficit de profundidade das bases. Consumir knowledge:go-foundations-io-batch.
- [x] Compor as cinco trilhas previstas usando desafios existentes e verificar início/travessia por MCP real em modo de autoria; não publicar conteúdo nem alterar a contagem canônica.
- [x] Lote final de autoria go-errors: classificação de erros encadeados (combinado) e catálogo em memória testado (fatia funcional), preservando fontes históricas e autoria; fixture baseline/reference, ordem determinística e isolamento.
- [x] Lote go-testing: desafio combinado de seleção temporal determinística, clock injetado, casos de fronteira e isolamento de slices; fixtures baseline/reference, avaliação qualitativa intermediária e check final. Consumir knowledge:go-foundations-io-batch.
- [x] Lote 2026-09-09: criar go-io com decodificação JSON estrita (combinado) e importação CSV validada para JSON (fatia funcional), fixtures baseline/reference e checks sem rede. Autoria nova identifica codex; conteúdo histórico mantém seu autor.
- [ ] Autorar conceitos, competências e relações.
- [ ] Autorar os 44 desafios e 300 step nodes.
- [ ] Autorar pistas, reflexões, variantes e checks.
- [ ] Criar cinco trilhas.
- [ ] Realizar revisão, playtest e correções.
- [ ] Integrar conteúdo canônico de concept-content-authoring e gates de catalog-publication-integrity antes de publicar os sete packs.
- [ ] Distinguir desafios canônicos, variantes e protótipos; decidir o destino dos dois exemplos históricos sem quebrar IDs de sessões.
- [ ] Revisar cobertura real de macro/meso/micro e critérios observáveis, sem inflar a contagem de nodes com passos redundantes.

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
- Lint: go run ./cmd/codinho catalog validate --json (catálogo completo resolvido pelo manifest; a CLI atual não seleciona packs por argumentos posicionais).
- Typecheck: go vet ./...
- Build: go build ./cmd/codinho
- Security / Contract: secret scan, path confinement e leak detector.

### Execution log
- 2026-09-10 UTC, trilhas fundamentais: cinco percursos explícitos (47 posições reutilizando desafios), início e travessia via MCP real em --authoring. Matriz ampliada para 23/23 checks aprovados; 43/43 checks editoriais verificados. IDs/conteúdo anteriores preservados. O avanço explícito no teste verifica alcance, não playtest humano.
- 2026-09-10 UTC, complemento go-errors: um combinado e uma fatia funcional draft, oito conceitos, seis competências e 14 nodes. Matriz 22/22 passou, com 43/43 checks editoriais verificados. Comparação dos objetos YAML confirma preservação integral dos desafios, conceitos, competências e relações históricos; versão do pack incrementada para 1.2.0.
- 2026-09-10 UTC, lote go-testing: um combinado draft, oito conceitos, cinco competências e sete nodes. Matriz completa 22/22 passou; 41/41 checks editoriais verificados, incluindo baseline com falha e referência aprovada. Docs-check sem erros/warnings. Revisão qualitativa e playtest humano permanecem pendentes.
- Retomada: knowledge:go-foundations-io-batch registra o lote validado e as demandas seguintes.
- 2026-09-09, lote go-io: dois desafios draft, 10 conceitos, 6 competências e 12 nodes adicionados. Matriz completa 22/22 passou; após revisão dos critérios intermediários, catalog validate --checks confirmou 40/40 checks do catálogo, incluindo baseline=test_failure/reference=pass nos dois desafios novos. docs-check sem warnings de inventário. Publicação e playtest permanecem pendentes.
- 2026-09-09: retomada após agent-authored-catalog-drafts; assess discover --component packs executado. Inventário aponta 33 atômicos fundamentais (go-first-steps tem 7 ante 6 previstos); preservar o conteúdo extra e registrar a divergência no relatório antes de qualquer alteração da distribuição.
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
- 2026-08-25: checkpoint 2 autorado em `packs/go-core.yaml` — 2 desafios
  atômicos novos (`stop-processing-commands-with-labeled-break`, tema
  `control-flow`, ensina o gotcha real de que `break` dentro de um
  `switch` aninhado em `for` só sai do switch, exigindo label;
  `increment-counter-with-pointer-receiver`, tema `functions-and-methods`,
  ensina receptor de ponteiro vs. valor para mutação de estado), 2
  conceitos e 2 competências novos, fixture + check `go_test` executável
  para cada um. Verificado manualmente antes da pré-revisão: stub de cada
  fixture falha; para o primeiro, uma implementação com `break` não
  rotulado também falha o teste (reproduz o bug real ensinado); para o
  segundo, um receptor de valor também falha o teste (idem).
- Tentativa de pré-revisão Codex rodada 1 (`new`): a primeira invocação
  falhou silenciosamente — `codex exec` saiu lendo stdin vazio sem gerar
  relatório, sem processo órfão. Reexecutada do zero.
- Pré-revisão Codex rodada 1 (`new`, reexecutada): **rejeitou os dois**.
  `stop-processing-commands-with-labeled-break`: as constraints exigiam
  só que "stop" encerrasse o laço externo, sem exigir literalmente break
  rotulado — uma implementação com `return` ou uma flag auxiliar passaria
  no teste e nos critérios sem nunca usar label, quebrando a
  correspondência entre a técnica ensinada e o contrato verificável.
  `increment-counter-with-pointer-receiver`: o primeiro macro_step reunia
  declarar o tipo Counter, o campo Count e a assinatura de Increment — a
  forma de Counter já estava inteiramente decidida pelo teste, sem
  decisão de modelagem real para o aluno; o critério
  `increment-has-mutating-receiver` também era interpretativo demais
  frente à competência (receptor de ponteiro).
- Correções: `stop-processing-commands-with-labeled-break` ganhou a
  constraint "usar break rotulado" e um critério novo
  `uses-labeled-break` (`source_inspection`, mesmo padrão de
  `uses-defer`/`uses-sprintf`) no segundo macro_step; primeiro macro_step
  teve o objective levemente reformulado (ressalva menor). Em
  `increment-counter-with-pointer-receiver`, `Counter` foi movido
  inteiramente para o fixture (já com o campo `Count`), o primeiro
  macro_step passou a exigir só a assinatura de `Increment` com receptor
  de ponteiro (`increment-uses-pointer-receiver`), e a constraint do
  desafio nomeia "receptor de ponteiro" objetivamente em vez de "receptor
  que permita alterar". Reverificado manualmente: os dois stubs continuam
  falhando (agora `Counter` compila mas `Increment` é undefined).
- `go run ./cmd/codinho catalog validate --checks` → exit 0 após as
  correções, sem novos avisos. `go build ./...`, `gofmt -l .`,
  `go vet ./...` → ok.
- Pré-revisão Codex rodada 2 (`resume --last`): **aprovado sem
  ressalvas** os dois desafios — `uses-labeled-break` e
  `increment-uses-pointer-receiver` aceitos como critérios objetivos e
  verificáveis por `source_inspection`, coerentes com a técnica exigida
  pelas constraints.
- 2026-08-25: checkpoint 3 autorado em `packs/go-core.yaml` — 2 desafios
  atômicos novos (`recover-from-panic-in-safe-call`, tema
  `functions-and-methods`, ensina converter panic em erro na fronteira de
  uma chamada com defer+recover; `guard-invariant-with-unexported-field`,
  tema novo `packages-and-apis`, ensina proteger um invariante com campo
  não exportado e construtor validador), 2 conceitos e 2 competências
  novos, fixture + check `go_test` executável para cada um. Verificado
  manualmente antes da pré-revisão: stub de cada fixture falha (panic /
  `undefined`), implementação de referência escrita à mão passa em todos
  os casos. Reviewer padrão passou a ser fixado explicitamente em
  `scripts/agent-review.sh` como `gpt-5.6-luna`/`model_reasoning_effort=
  high` (antes dependia implicitamente do default de
  `~/.codex/config.toml`) — ver `docs/agent-review-workflow.md`.
- Pré-revisão Codex rodada 1 (`new`): **rejeitou os dois**.
  `recover-from-panic-in-safe-call`: cobertura de "qualquer valor de
  panic" insuficiente (só `string`/`int` testados) e o `objective` do
  segundo macro_step nomeava a técnica inteira (`defer`+`recover`) além da
  constraint, duplicando o vazamento. `guard-invariant-with-unexported-
  field`: primeiro macro_step reunia quatro decisões (tipo, visibilidade
  do campo, accessor, assinatura de `NewPercentage`) sob `target:
  function` (incorreto para declaração de tipo) e o teste estava no mesmo
  pacote do código, não provando proteção via API pública.
- Correções: `recover-from-panic-in-safe-call` ganhou um quarto caso de
  panic com `error` no acceptance/teste; objective reescrito para
  descrever o efeito observável em vez da técnica literal (mesmo padrão de
  `close-resources-in-defer-order`). `guard-invariant-with-unexported-
  field` dividido em dois macro_steps de intenção única — tipo+accessor
  (`target: named_type`, valor já usado em `go-first-steps.yaml`) e
  `NewPercentage` declarado+implementado junto — e `percentage_test.go`
  reescrito como `package percentage_test` externo, usando só a API
  pública. Os critérios estruturais restantes (`source_inspection`/
  `compile` sem runner dedicado a palavra-chave/visibilidade) foram
  mantidos citando o precedente já aceito nos checkpoints 1–2 deste pack
  (`uses-defer`, `uses-labeled-break`, `increment-uses-pointer-receiver`)
  em vez de tratados como achado novo.
- `go run ./cmd/codinho catalog validate --checks` → exit 0 após as
  correções, sem novos avisos. `go build ./...`, `gofmt -l .`,
  `go vet ./...` → ok.
- Pré-revisão Codex rodada 2 (`resume --last`): **aprovado sem
  ressalvas** os dois desafios — reconheceu explicitamente o precedente
  citado para os critérios estruturais e o check único por desafio; único
  ponto residual anotado (não bloqueante) é que `go test` não consegue
  provar visibilidade de campo por si só, mesma limitação sistêmica já
  aceita.
- 2026-08-25: checkpoint 4 autorado em `packs/go-core.yaml` — completa os
  8 desafios atômicos previstos na matriz de distribuição para este pack.
  2 desafios atômicos novos: `go-core.find-first-value-at-least` (tema
  `control-flow`, curto-circuito `&&` para nunca acessar um campo de
  ponteiro nil ao filtrar uma lista) e `go-core.hide-counter-type-behind-
  interface` (tema `packages-and-apis`, construtor exportado retornando
  um tipo de interface em vez do tipo concreto). 2 conceitos e 2
  competências novos. Verificado manualmente: stub de cada fixture falha;
  para `find-first-value-at-least`, verificado também que uma
  implementação com a ordem da condição invertida (`node.Value >= min &&
  node != nil`) falha de verdade com nil pointer dereference contra o
  mesmo teste — evidência executável real por trás do critério
  estrutural, não só sintática.
- Pré-revisão Codex rodada 1 (`new`): **aprovado com ressalvas menores**
  os dois. `find-first-value-at-least`: objective do primeiro macro_step
  já nomeava "curto-circuito", reduzindo a descoberta esperada da técnica.
  `hide-counter-type-behind-interface`: o primeiro macro_step (mutação via
  receptor de ponteiro) estava com `concepts: [constructor-returns-
  interface]`, inconsistente — esse conceito só se aplica ao segundo
  passo (o construtor retornando a interface).
- Correções: objective do primeiro passo de `find-first-value-at-least`
  reescrito para descrever só o efeito observável, sem citar a técnica; a
  técnica permanece explícita apenas na constraint do segundo passo, mesmo
  padrão já usado no pack. Primeiro macro_step de `hide-counter-type-
  behind-interface` passou a usar `concepts: [pointer-vs-value-receiver]`
  (concept já existente no pack desde o checkpoint 2), refletindo o que
  esse passo realmente ensina.
- `go run ./cmd/codinho catalog validate --checks` → exit 0 após as
  correções, sem novos avisos. `go build ./...`, `gofmt -l .`,
  `go vet ./...` → ok.
- Pré-revisão Codex rodada 2 (`resume --last`): **aprovado sem
  ressalvas** `find-first-value-at-least`; **aprovado com uma ressalva
  menor não bloqueante, já conhecida do precedente**
  `hide-counter-type-behind-interface` — o teste externo prova uso sem
  nomear o tipo concreto, mas a garantia de que `NewCounter` retorna
  `Counter` (não o tipo concreto) continua dependendo de
  `source_inspection` da assinatura, mesma limitação sistêmica já aceita
  para critérios estruturais neste pack.
- 2026-08-25: checkpoint 5 autorado em `packs/go-core.yaml` — os 2
  desafios `kind: combined` previstos na matriz de distribuição,
  completando os 10 desafios de go-core (8 atômicos + 2 combinados).
  `go-core.run-commands-with-cleanup-and-recovery` (temas control-flow +
  functions-and-methods) integra switch+break rotulado, defer e recover
  reaproveitando os concepts já existentes dos checkpoints 1–3, sem
  concepts novos. `go-core.configure-server-with-functional-options`
  (temas functions-and-methods + packages-and-apis) integra parâmetros
  variádicos com validação de invariante via campo não exportado
  (padrão de opções funcionais), reaproveitando as competências já
  existentes `use-variadic-parameters`/`guard-invariant-with-unexported-
  field` e introduzindo um único concept novo (`functional-options-
  pattern`) para nomear o padrão em si — nenhuma competência nova foi
  criada, conforme a definição de "combinado" da taxonomia do projeto.
  Verificado manualmente: stub de cada fixture falha, referência passa;
  para `RunCommands`, a ordem exata de report (start/ok-ou-failed/closed
  por comando) foi conferida contra a implementação de referência.
- O validador determinístico (`catalog validate --checks`) pegou sozinho
  um `compound_micro_instruction` real antes mesmo da pré-revisão do
  Codex — um macro_step de `run-commands-with-cleanup-and-recovery`
  reunia "recuperar panic" e "registrar sucesso" na mesma instruction
  (separadas por `;`); dividido em dois macro_steps antes de enviar para
  revisão.
- Pré-revisão Codex rodada 1 (`new`): **aprovado com ressalvas menores**
  os dois. Único ponto acionável: o subteste de parada de
  `RunCommands` só conferia o `report` final, sem provar
  comportamentalmente que `run` nunca era chamado para "stop"/"c" depois
  da parada.
- Correção: subteste reescrito com um `spy` que registra cada comando com
  que `run` foi chamado, falhando se a lista não for exatamente `["a"]`
  no caso de parada — prova comportamental direta, não só inferida do
  report.
- Pré-revisão Codex rodada 2 (`resume --last`): **aprovado sem
  ressalvas** `run-commands-with-cleanup-and-recovery`; **aprovado com
  ressalvas menores não bloqueantes, já cobertas pelo precedente**
  `configure-server-with-functional-options` — parte da forma pública
  (assinatura variádica, campos não exportados) já vem scaffoldada na
  fixture, e os critérios estruturais continuam dependendo de
  `source_inspection`.
- `go run ./cmd/codinho catalog validate --checks` → exit 0 em todas as
  rodadas, sem `compound_micro_instruction` residual. `go build ./...`,
  `gofmt -l .`, `go vet ./...` → ok.
- 2026-08-25: pack novo `packs/go-data-text.yaml` criado (temas E
  slices/arrays/maps e F texto/dados binários) e registrado em
  `packs/manifest.yaml`. Checkpoint 1: 2 desafios atômicos —
  `go-data-text.filter-without-mutating-input` (filtra uma slice sem
  reaproveitar o backing array da entrada, um bug real de compactação
  in-place) e `go-data-text.truncate-bytes-at-rune-boundary` (trunca uma
  string por bytes recuando até a fronteira de rune UTF-8 mais próxima).
  Achado incidental fora do escopo: `go-first-steps.yaml` tem um item
  pré-existente quebrado (`go-data.slice-filter-preserve-input`, cópia
  literal do exemplo de PROJECT.md §14.8, sem fixture e com macro_step
  sobre um assunto completamente diferente) — registrado em Known gaps,
  não corrigido (fora do escopo deste checkpoint).
- Verificado manualmente antes da pré-revisão: para os dois desafios,
  stub falha e referência passa; adicionalmente, para cada um, uma
  implementação incorreta plausível foi escrita e confirmada como
  reprovada pelo teste real (compactação in-place para o filtro; corte
  ingênuo `s[:maxBytes]` para o truncamento).
- Pré-revisão Codex rodada 1 (`new`): **aprovado sem ressalvas**
  `truncate-bytes-at-rune-boundary`. **Rejeitado**
  `filter-without-mutating-input` por dois motivos reais: (1) objective
  do primeiro macro_step vazava a técnica ("partindo de nil"); (2) o
  teste só verificava que os valores visíveis da entrada não mudavam, sem
  provar ausência de aliasing — o revisor escreveu uma implementação
  "esperta" que reaproveita a capacidade excedente da entrada
  (`nums[len(nums):len(nums)]` seguido de append) sem mutar nenhum valor
  visível, e ela passava no teste original.
- Correções: objective reescrito para descrever só o efeito, técnica só
  na constraint. Teste reescrito para materializar a entrada com
  capacidade excedente deliberada (`make([]int, 6, 20)`) e verificar,
  via `unsafe.SliceData`/aritmética de ponteiro, se o endereço inicial da
  slice retornada cai dentro do intervalo completo do backing array da
  entrada (não só igualdade de ponteiro — uma checagem de igualdade
  simples não pegava a implementação esperta do revisor, porque o
  `append` nela desloca o ponteiro inicial para dentro da região de
  capacidade excedente). Reverificado manualmente contra três
  implementações: referência (passa), compactação in-place (falha) e a
  implementação esperta de capacidade excedente (falha).
- `go run ./cmd/codinho catalog validate --checks` → exit 0 após as
  correções, sem avisos novos. `go build ./...`, `gofmt -l .`,
  `go vet ./...` → ok.
- Pré-revisão Codex rodada 2 (`resume --last`): **aprovado sem
  ressalvas** `filter-without-mutating-input` — confirmou que as duas
  implementações incorretas (in-place e capacidade excedente) falham no
  teste novo.
- 2026-08-25: checkpoint 2 autorado em `packs/go-data-text.yaml` — 2
  desafios atômicos novos: `go-data-text.lookup-map-value-with-comma-ok`
  (tema slices-arrays-maps, distingue pontuação zero de nunca pontuado
  via comma-ok) e `go-data-text.quote-csv-field-when-needed` (tema
  text-and-binary-data, citação de campo CSV estilo RFC 4180
  simplificada). Verificado manualmente antes da pré-revisão: stub falha,
  referência passa, e uma implementação incorreta plausível para cada um
  falha (comparação com zero em vez de comma-ok; esquecer de dobrar
  aspas internas).
- Pré-revisão Codex rodada 1 (`new`): **aprovado com ressalva menor**
  `lookup-map-value-with-comma-ok` (acceptance não cobria chave vazia).
  **Rejeitado** `quote-csv-field-when-needed`: o teste cobria cada
  caractere especial isoladamente, mas não a combinação vírgula+aspas —
  o revisor escreveu uma implementação sutil que cita corretamente
  quando há vírgula/quebra de linha mas esquece de escapar aspas internas
  nesse mesmo caminho, e ela passava no teste original.
- Correções: caso de chave vazia adicionado ao teste de
  `lookup-map-value-with-comma-ok` (confirmado que a implementação com
  comparação por zero falha nesse caso). Para `quote-csv-field-when-
  needed`, adicionado um caso combinando vírgula e aspas internas
  (`a,"b"` → `"a,""b"""`) e um caso de retorno de carro (`\r`), já
  sinalizado como lacuna relacionada pelo revisor — a detecção de
  caracteres especiais na fixture e no macro_step passou a incluir `\r`
  explicitamente. Reverificado manualmente: a implementação sutil descrita
  pelo revisor agora falha no caso combinado.
- `go run ./cmd/codinho catalog validate --checks` → exit 0 após as
  correções, sem avisos novos. `go build ./...`, `gofmt -l .`,
  `go vet ./...` → ok.
- Pré-revisão Codex rodada 2 (`resume --last`): **aprovado sem
  ressalvas** os dois desafios — confirmou que os mutantes relevantes
  (comparação por zero; citação sem escapar aspas no caminho com vírgula)
  falham nos testes novos.
- 2026-08-25: checkpoint 3 autorado em `packs/go-data-text.yaml` — 2
  desafios atômicos novos: `go-data-text.dedupe-preserving-first-
  occurrence` (tema slices-arrays-maps, deduplica usando `map[string]
  struct{}` como set, preservando ordem de primeira ocorrência) e
  `go-data-text.sum-valid-integers-safely` (tema text-and-binary-data,
  soma inteiros de strings tratando o erro de `strconv.Atoi` em vez de
  ignorá-lo). Verificado manualmente com `-count=1` antes da
  pré-revisão: stub falha, referência passa (5x), mutantes incorretos
  plausíveis falham.
- Pré-revisão Codex rodada 1 (`new`): **rejeitou os dois**, e achou algo
  que minha própria verificação manual não tinha pego. `dedupe-
  preserving-first-occurrence`: o mutante que monta a saída iterando um
  map no final (perdendo a ordem original) só falhava intermitentemente
  contra o teste original — o revisor rodou 30 vezes e mediu 19 falhas
  contra 11 passagens, porque a ordem de iteração aleatória do map às
  vezes coincidia por acaso com a ordem esperada. Minhas 5 repetições
  manuais tinham, por sorte, sempre caído do lado da falha — evidência de
  que 5 repetições não bastam para descartar não-determinismo por
  iteração de map. `sum-valid-integers-safely`: faltava um caso de
  inteiro válido igual a zero (`"0"`) no teste; uma implementação mutante
  que trata `n == 0` como inválido (além do erro real de parsing) passava
  integralmente.
- Correções: teste de `dedupe-preserving-first-occurrence` reescrito para
  chamar `Dedupe` com a mesma entrada 30 vezes DENTRO do mesmo subteste,
  falhando na primeira chamada cuja ordem não bater exatamente — torna a
  chance de um mutante baseado em iteração de map passar por coincidência
  praticamente nula ((1/720)^30), em vez de depender de repetir o
  processo de teste várias vezes por fora. `Dedupe(nil)` passou a
  comparar explicitamente com `nil` (era "vazio", ambíguo). Caso
  `SumValid(["0","5"])` → `(5, 0)` adicionado a `sum-valid-integers-
  safely`. Reverificado manualmente com `-count=1`: mutante de dedupe
  falha 5/5 execuções completas do processo (cada uma já contendo as 30
  chamadas internas); mutante de zero-como-inválido falha 5/5.
- `go run ./cmd/codinho catalog validate --checks` → exit 0 após as
  correções, sem avisos novos. `go build ./...`, `gofmt -l .`,
  `go vet ./...` → ok.
- Pré-revisão Codex rodada 2 (`resume --last`): **aprovado com ressalva
  menor não bloqueante** `dedupe-preserving-first-occurrence` (primeiro
  macro_step ainda agrupa percorrer/consultar/registrar — não bloqueante);
  **aprovado sem ressalvas** `sum-valid-integers-safely`. O revisor
  confirmou o determinismo do novo teste rodando o processo 10 vezes
  contra referência e mutante.
- Achado de processo registrado em `docs/agent-review-workflow.md`: `go
  test` sem `-count=1` pode devolver resultado em cache mesmo trocando o
  arquivo de implementação entre execuções, mascarando um mutante
  quebrado; e, mesmo com `-count=1`, poucas repetições (5) não bastam
  para descartar flakiness por iteração de map — repetir a chamada dentro
  do próprio subteste (N=30) é a forma barata e confiável de tornar esse
  tipo de evidência comportamental verdadeiramente determinística.
- 2026-08-25: checkpoint 4 autorado em `packs/go-data-text.yaml` — 2
  desafios atômicos novos, completando os 8 atômicos previstos para o
  pack (faltam só os 2 combinados). `go-data-text.title-first-letters-
  preserving-spacing` (tema text-and-binary-data, capitaliza a primeira
  letra de cada palavra preservando exatamente o espaçamento original,
  em contraste com strings.Fields+strings.Join que colapsaria espaços
  múltiplos) aprovado sem ressalvas já na rodada 1, após reforçar o teste
  com casos de espaço final/tabulação/letra Unicode e reescrever os
  objectives para descrever efeito em vez de narrar o algoritmo passo a
  passo. `go-data-text.preallocate-slice-with-zero-length` (tema
  slices-arrays-maps) precisou de duas rodadas de correção real, não só
  polimento — ver achado abaixo.
- Achado de desenho de desafio (mais significativo que os anteriores): o
  desafio original, `BuildSquares(n)` (soma dos quadrados de 0 a n-1,
  sempre produzindo exatamente n elementos 1:1), tinha uma constraint
  estruturalmente infalsificável — como toda entrada gera exatamente n
  saídas, `make([]int, n)` seguido de escrita indexada produz o MESMO
  resultado observável (valores, comprimento, capacidade) que
  `make([]int, 0, n)` seguido de append; nenhum teste de caixa preta
  poderia distinguir as duas técnicas. Verificação manual anterior não
  pegou isso porque só testei a implementação óbvia (make(n)+append, que
  de fato quebra), não a alternativa por indexação (que não quebra).
  Redesenhado para `EvenSquares(n)` (só os quadrados dos números pares),
  onde o comprimento final é variável/dependente dos dados — agora
  make(n)+append produz um resultado genuinamente errado (zeros à
  esquerda espúrios). Mesmo assim, a rodada 2 achou uma segunda lacuna:
  o teste de capacidade exigia só `cap <= n`, que uma implementação SEM
  pré-alocação nenhuma (`var out []int` + append) também satisfazia por
  acidente, porque o crescimento amortizado do append para poucos
  elementos frequentemente fica abaixo de n. Corrigido exigindo
  `cap == n` exatamente (só pré-alocação deliberada garante isso).
  Lição para desafios futuros: antes de escrever o teste, perguntar se
  existe uma implementação mais simples ou sem a técnica exigida que
  produziria o mesmo resultado observável — inclusive "sem a otimização
  nenhuma", não só "com a técnica errada".
- `go run ./cmd/codinho catalog validate --checks` → exit 0 em todas as
  rodadas. `go build ./...`, `gofmt -l .`, `go vet ./...` → ok.
- Pré-revisão Codex: rodada 1 (`new`) rejeitou os dois pelos motivos
  acima; rodada 2 (`resume --last`) aprovou `title-first-letters-
  preserving-spacing` sem ressalvas e rejeitou de novo `preallocate-
  slice-with-zero-length` (a lacuna de `cap <= n`); rodada 3 (`resume
  --last`) aprovou `preallocate-slice-with-zero-length` sem ressalvas
  após a correção para `cap == n` exato.
- 2026-08-25: checkpoint 5 (ÚLTIMO de go-data-text) autorado — os 2
  desafios `kind: combined`, completando o pack em 10/10 desafios (8
  atômicos + 2 combinados). `go-data-text.build-scores-from-entries`
  (temas text-and-binary-data + slices-arrays-maps) integra
  parse-integers-with-error-handling + distinguish-zero-value-from-
  absence + dedupe-preserving-first-occurrence, sem competências novas.
  `go-data-text.format-names-as-csv-fields` (temas slices-arrays-maps +
  text-and-binary-data) integra preallocate-slice-with-zero-length +
  title-first-letters-preserving-spacing + quote-csv-field-when-needed,
  também sem competências novas.
- Pré-revisão Codex rodada 1 (`new`): **rejeitou os dois**, achando
  lacunas de teste reais análogas às dos checkpoints 3-4.
  `build-scores-from-entries`: faltava o caso de uma entrada inválida
  seguida de uma válida para o mesmo nome — um mutante que marca o nome
  como "já visto" antes de validar o placar (em vez de só depois da
  conversão ter sucesso) passava no teste publicado, mas descartaria
  incorretamente um placar válido que vem depois de uma tentativa
  inválida. `format-names-as-csv-fields`: faltavam casos de espaçamento
  múltiplo/pontas, nome só de espaços, e aspas internas após capitalizar
  — o mutante `strings.Fields`+`strings.Join` (já usado no checkpoint 4)
  passava de novo por não ter cobertura suficiente neste novo contexto
  combinado; o primeiro macro_step também combinava pré-alocar e ignorar
  nomes vazios.
- Correções: `build-scores-from-entries` ganhou o caso `["bob:not-an-
  int","bob:5"]` → `scores["bob"]==(5,true)`, e uma constraint explícita
  proibindo que uma entrada inválida marque o nome como já tendo placar.
  `format-names-as-csv-fields` ganhou os três casos que faltavam, uma
  constraint explícita definindo que só a string vazia "" é ignorada (um
  nome só de espaços não), e o primeiro macro_step foi dividido em dois
  (pré-alocar / ignorar só string vazia), totalizando 4 macro_steps.
  Reverificado manualmente com -count=1 (10x): os mutantes relevantes de
  cada desafio agora falham de forma consistente.
- `go run ./cmd/codinho catalog validate --checks` → exit 0 após as
  correções, sem avisos novos. `go build ./...`, `gofmt -l .`,
  `go vet ./...` → ok.
- Pré-revisão Codex rodada 2 (`resume --last`): **aprovado sem
  ressalvas** `build-scores-from-entries`; **aprovado com ressalva menor
  não bloqueante** `format-names-as-csv-fields` (comentário da fixture
  ainda dizia "skipping blanks", desalinhado com o contrato final de só
  ignorar a string vazia) — corrigido o comentário.
- **go-data-text.yaml está completo: 10/10 desafios da matriz de
  distribuição** (8 atômicos + 2 combinados), mesmo padrão de rigor de
  go-core.
- 2026-08-25: pack novo `packs/go-type-design.yaml` criado nesta sessão
  (temas G structs/ponteiros, H interfaces, I generics) e registrado em
  `packs/manifest.yaml`. Checkpoint 1: 2 desafios atômicos —
  `go-type-design.clone-inventory-independently` (tema
  structs-and-pointers, clona uma struct Inventory sem que o campo Items
  compartilhe backing array com o original) e `go-type-design.select-
  validator-true-nil` (tema interfaces, a armadilha clássica do "typed
  nil": devolver um ponteiro nil embrulhado em uma interface produz uma
  interface != nil).
- Pré-revisão Codex rodada 1 (`new`): **rejeitou os dois**, achando
  mutantes sutis reais. `clone-inventory-independently`: o teste só
  verificava aliasing quando `len(Items) > 0` (saída antecipada no caso
  vazio); um mutante `clone := inv; if len(inv.Items) > 0 { clone.Items =
  append(...) }` compartilha o backing array quando Items está vazio mas
  tem capacidade excedente, e passava despercebido. `select-validator-
  true-nil`: o `strictValidator` original não tinha nenhum campo/estado,
  então `Validate()` nunca dereferenciava o receptor — um mutante que
  retorna um `*strictValidator` nil DENTRO do próprio caso "strict"
  (variação do bug clássico, não só no caso desconhecido) produzia uma
  interface não nil que "funcionava por acidente", sem nenhum teste capaz
  de provar que o validador estava de fato configurado. Os objectives
  também vazavam a técnica exata (`&strictValidator{}`, `return nil`
  diretamente).
- Correções: teste de `clone-inventory-independently` ganhou o caso
  `Items: make([]string, 0, 5)` (vazio com capacidade excedente),
  verificado via `unsafe.SliceData` sem saída antecipada.
  `strictValidator` ganhou um campo real (`minLength int`) e `Validate()`
  passou a retornar erro quando não configurado (`minLength <= 0`) — um
  ponteiro nil ou zero-value agora falha de verdade (panic ou erro) ao
  ser usado, não só "aceita por acaso". Objectives reescritos para
  descrever efeito ("de fato configurado", "comparável a nil com ==") em
  vez da técnica exata. Reverificado manualmente com -count=1: os dois
  mutantes sutis agora falham de forma consistente (5/5).
- `go run ./cmd/codinho catalog validate --checks` → exit 0 após as
  correções, sem avisos novos. `go build ./...`, `gofmt -l .`,
  `go vet ./...` → ok.
- Pré-revisão Codex rodada 2 (`resume --last`): **aprovado sem
  ressalvas** os dois desafios.
- 2026-08-25: checkpoint 2 autorado em `packs/go-type-design.yaml` — 2
  desafios atômicos novos: `go-type-design.sum-list-with-nil-safe-
  receiver` (tema structs-and-pointers, soma valores de uma lista
  encadeada com método que funciona corretamente em um *Node nil) e
  `go-type-design.contains-generic-comparable` (tema generics, função
  genérica Contains[T comparable]). Verificado manualmente com -count=1
  antes da pré-revisão: stub falha, referência passa, mutante plausível
  falha.
- Pré-revisão Codex: cinco rodadas até aprovação final, a mais
  trabalhosa desta sessão. Rodada 1 (`new`) rejeitou os dois — faltava
  caso de valor negativo em Sum e caso de alvo no índice 0/tipos além de
  int-string em Contains, além de objectives vazando técnica. Rodada 2
  achou lacunas novas em cada correção: Sum aceitava um mutante que
  deduplicava valores repetidos (cadeia "2 → 2" somava 2 em vez de 4);
  Contains aceitava um mutante que só examinava os dois primeiros
  elementos. Rodada 3 achou mais uma classe de mutante "processa só um
  prefixo de tamanho fixo" em cada um (cadeia/slice de 3-4 elementos) —
  nesse ponto, em vez de continuar caçando um tamanho por rodada, troquei
  a estratégia para um caso de 20 elementos construído em loop; isso
  fechou de vez a classe de prefixos fixos em Sum (aprovado sem
  ressalvas na rodada 4) mas Contains ainda tinha um mutante esparso
  (só índices 0, 1 e o último) que só a checagem de um alvo (o último
  índice) não pegava. Rodada 5: generalizado para testar TODOS os 20
  índices como alvo em loop, fechando também a classe de subconjuntos
  esparsos — aprovado sem ressalvas.
- Lição de processo mais importante desta sessão: para bugs do tipo
  "processa só um subconjunto da coleção" (prefixo fixo, índices
  esparsos, etc.), testar exaustivamente TODOS os elementos/índices de
  uma coleção razoavelmente grande (ex.: 20) construída em loop fecha a
  classe inteira de uma vez — muito mais robusto que adicionar um caso
  pontual a cada mutante que a pré-revisão encontra, o que vira caça ao
  gato-e-rato. Vale aplicar preventivamente em qualquer desafio futuro
  que itere uma coleção.
- `go run ./cmd/codinho catalog validate --checks` → exit 0 em todas as
  rodadas. `go build ./...`, `gofmt -l .`, `go vet ./...` → ok.
- 2026-08-25: checkpoint 3 (último atômico) autorado em `packs/go-type-
  design.yaml` — 2 desafios atômicos novos, completando os 6 atômicos
  previstos para o pack. `go-type-design.describe-if-circle-safely`
  (tema interfaces, type assertion segura via comma-ok ou type switch,
  sem panicar quando o Shape não é um Circle) e `go-type-design.map-
  generic-transform` (tema generics, função genérica Map[T,U any],
  aplicando a lição do checkpoint 2: caso de 20 elementos já incluído
  desde a primeira versão).
- Pré-revisão Codex rodada 1 (`new`): **rejeitou os dois**.
  `describe-if-circle-safely`: o teste checava só `strings.Contains(desc,
  "2")` para o raio 2 — um mutante que descreve a ÁREA em vez do raio
  (área = 12.566 para raio 2) também contém "2" e passava por
  coincidência. `map-generic-transform`: o teste "não muta a entrada" só
  verificava igualdade sem nunca mutar o resultado nem checar
  compartilhamento de backing array — um mutante que reaproveita a
  capacidade excedente da entrada passava; o objective do segundo
  macro_step também prescrevia "acrescentar" (vazando append como única
  técnica válida, quando escrita indexada é igualmente correta).
- Correções: `describe-if-circle-safely` passou a comparar a descrição
  EXATA ("circle with radius 10", raio escolhido para não coincidir com
  a área) em vez de checar um dígito por coincidência; confirmado que
  type switch sem comma-ok é uma técnica alternativa igualmente válida
  (não é bug), e a constraint/reflexão foram ajustadas para não excluir
  essa alternativa. `map-generic-transform` ganhou o mesmo padrão de
  teste de aliasing já usado em outros packs (entrada com capacidade
  excedente, mutação do resultado, checagem via unsafe.SliceData) e o
  objective deixou de mencionar "acrescentar", descrevendo só o efeito.
- `go run ./cmd/codinho catalog validate --checks` → exit 0 após as
  correções, sem avisos novos. `go build ./...`, `gofmt -l .`,
  `go vet ./...` → ok.
- Pré-revisão Codex rodada 2 (`resume --last`): **aprovado sem
  ressalvas** `map-generic-transform`; **aprovado com ressalva menor não
  bloqueante** `describe-if-circle-safely` — a reflexão continha uma
  afirmação factualmente incorreta sobre type switch sem default
  panicar (não panica, só não executa nenhum case) — corrigida.
- 2026-08-25: checkpoint 4 (ÚLTIMO) autorado em `packs/go-type-
  design.yaml` — os 2 desafios `kind: combined`, completando o pack em
  8/8 (6 atômicos + 2 combinados). `go-type-design.first-circle-clone`
  (temas structs-and-pointers + interfaces) integra clone-struct-with-
  independent-slice + assert-type-without-panicking + return-true-nil-
  interface, sem competências novas. `go-type-design.generic-node-
  values-nil-safe` (temas structs-and-pointers + generics) integra
  handle-nil-receiver-safely + implement-generic-transform, também sem
  competências novas.
- Pré-revisão Codex rodada 1 (`new`): **rejeitou os dois**, aplicando a
  lição do checkpoint 2 de forma proativa (o autor já tinha incluído
  casos de 20 elementos desde a primeira versão, mas incompletos).
  `first-circle-clone`: um mutante que retorna o ÚLTIMO Circle (não o
  primeiro) passava, porque o caso de 20 elementos só tinha um Circle
  por vez, nunca provando "primeiro" vs "último"; um mutante que só
  inspeciona os índices 0, 1 e o último também passava, porque o caso de
  20 só testava o Circle no último índice, não em todos.
  `generic-node-values-nil-safe`: um mutante que só coleta valores para
  os tipos concretos presentes nos testes (int, string, via type switch
  em tempo de execução) passava, porque a fixture nunca testava um tipo
  definido pelo usuário.
- Correções: `first-circle-clone` ganhou um caso com dois Circles (radius
  2 antes de radius 5), provando que o primeiro é devolvido, e o loop de
  20 elementos passou a colocar o Circle em CADA um dos 20 índices, um de
  cada vez (não só no último), fechando de vez a classe de subconjuntos
  esparsos. `generic-node-values-nil-safe` ganhou um caso com um tipo
  struct definido no teste (`item{ID int; Name string}`).
  Reverificado manualmente com -count=1: os mutantes "retorna o último",
  "índices esparsos" e "só int/string" agora falham de forma consistente.
- `go run ./cmd/codinho catalog validate --checks` → exit 0 após as
  correções, sem avisos novos. `go build ./...`, `gofmt -l .`,
  `go vet ./...` → ok.
- Pré-revisão Codex rodada 2 (`resume --last`): **aprovado sem
  ressalvas** os dois desafios.
- **go-type-design.yaml está completo: 8/8 desafios da matriz de
  distribuição** (6 atômicos + 2 combinados). Lição reforçada nesta
  sessão: aplicar a técnica de "testar TODOS os índices/posições/tipos
  representativos" desde a primeira versão de um desafio que itera uma
  coleção ou é genérico reduz — mas não elimina — a chance de rejeição
  na primeira rodada; ainda vale a pena a pré-revisão adversarial
  específica, porque a primeira tentativa de aplicar a lição
  (colocar o item-alvo em só uma posição por vez do loop de 20) ainda
  deixou passar duas classes de mutante que só a iteração completa
  fechou de fato.
- 2026-08-25: pack novo `packs/go-errors.yaml` criado nesta sessão (tema
  J erros) e registrado em `packs/manifest.yaml`. Checkpoint 1: 2
  desafios atômicos — `go-errors.wrap-sentinel-error-with-context`
  (encadeia um erro sentinela com fmt.Errorf %w, preservando errors.Is
  através do encadeamento) e `go-errors.extract-field-with-errors-as`
  (extrai um *ValidationError de um erro encadeado usando errors.As).
- Pré-revisão Codex rodada 1 (`new`): **rejeitou os dois**.
  `wrap-sentinel-error-with-context`: a fixture só testava uma chave
  válida ("a"); um mutante hardcoded que só reconhece essa chave
  específica passava integralmente. `extract-field-with-errors-as`: a
  fixture só testava uma camada de wrapping; um mutante que faz só
  `errors.Unwrap(err)` uma vez (em vez de errors.As, que percorre toda a
  cadeia) passava nesse caso raso.
- Correções: `wrap-sentinel-error-with-context` ganhou um mapa com 5
  chaves válidas testadas em loop e uma segunda chave ausente diferente
  da original (evitando qualquer hardcode); o `concepts` do macro_step de
  caminho feliz (que não envolve wrapping) foi corrigido de
  `error-wrapping-with-percent-w` para vazio. `extract-field-with-
  errors-as` ganhou um caso com duas camadas externas de wrapping
  (`fmt.Errorf("...: %w", fmt.Errorf("...: %w", ValidateAge(-1)))`).
  Reverificado manualmente com -count=1: o mutante hardcoded falha nas
  chaves não previstas; o mutante de unwrap único falha no caso profundo.
- `go run ./cmd/codinho catalog validate --checks` → exit 0 após as
  correções, sem avisos novos. `go build ./...`, `gofmt -l .`,
  `go vet ./...` → ok.
- Pré-revisão Codex rodada 2 (`resume --last`): **aprovado sem
  ressalvas** os dois desafios.
- 2026-08-25: checkpoint 2 autorado em `packs/go-errors.yaml` — 2
  desafios atômicos novos, completando os 4 atômicos previstos para o
  pack. `go-errors.validate-user-joining-all-errors` (agrega todas as
  violações de validação com errors.Join, não só a primeira) e
  `go-errors.translate-error-at-boundary` (traduz uma falha interna
  simulada em um erro público de domínio, sem deixar o erro interno
  detectável nem por errors.Is nem textualmente na mensagem).
- Pré-revisão Codex rodada 1 (`new`): **aprovado com ressalva menor não
  bloqueante** `validate-user-joining-all-errors` — uma implementação
  sem errors.Join, usando dois %w (`fmt.Errorf("%w: %w", ...)`), produz
  o mesmo resultado observável; a exigência de usar especificamente
  errors.Join continua dependendo de source_inspection, mesma limitação
  sistêmica já aceita no projeto para outras técnicas (uses-defer,
  errors.As vs asserção direta, etc.) — mantido sem alteração.
  **Rejeitado** `translate-error-at-boundary`: o teste só verificava
  ausência de vazamento via errors.Is; um mutante que troca %w por %v
  (`fmt.Errorf("%w: %v", ErrServiceUnavailable, errDeadlineExceeded)`)
  escapa de errors.Is mas ainda expõe o texto "i/o timeout" na mensagem
  do erro — uma fronteira de tradução real precisa impedir os dois tipos
  de vazamento, não só o estrutural.
- Correção: `translate-error-at-boundary` ganhou uma verificação de
  vazamento textual (`strings.Contains(err.Error(), errDeadlineExceeded.
  Error())`) além da checagem via errors.Is, e a constraint/acceptance
  foram reforçadas para exigir explicitamente os dois tipos de proteção.
  Reverificado manualmente com -count=1: o mutante %v agora falha.
- `go run ./cmd/codinho catalog validate --checks` → exit 0 após a
  correção, sem avisos novos. `go build ./...`, `gofmt -l .`,
  `go vet ./...` → ok.
- Pré-revisão Codex rodada 2 (`resume --last`): **aprovado sem
  ressalvas** `translate-error-at-boundary`.
- **Os 4 desafios atômicos de go-errors.yaml estão completos**, faltando
  o checkpoint com os 2 combinados e a fatia funcional (3 desafios
  restantes de 7).

### Results summary
Baseline da auditoria de 2026-09-07 UTC: 38 desafios fundamentais no HEAD
d758744, mais um combinado em go-errors no worktree do usuário, totalizando
39 canônicos autorados de 44. Há ainda dois protótipos fora dessa contagem
(go-data.slice-filter-preserve-input e go-debug.slice-off-by-one). Os 41 itens
carregados no worktree têm zero publicação confirmada. Os checkpoints abaixo
são histórico de autoria/pré-revisão, não comprovação de publicação ou do gate.

Histórico de autoria: checkpoints de conteúdo real em cinco packs,
com contagem reconciliada no baseline acima. Três em go-first-steps (checkpoint 1 declarações/tipos,
checkpoint 2 tooling + conversão numérica, checkpoint 3 tooling +
shadowing), cinco em go-core (checkpoint 1 defer + parâmetros
variádicos, checkpoint 2 labeled break + receptor de ponteiro, checkpoint
3 panic/recover + invariante com campo não exportado, checkpoint 4
curto-circuito + construtor de interface, checkpoint 5 os 2 desafios
combinados — **go-core está completo**: 10/10 desafios previstos na
matriz de distribuição, 8 atômicos + 2 combinados), cinco checkpoints em
go-data-text (pack novo, criado nesta sessão: checkpoint 1 filtro de
slice sem aliasing + truncamento de string em fronteira de rune UTF-8,
checkpoint 2 comma-ok em map + citação de campo CSV, checkpoint 3
deduplicação com set + soma segura de inteiros, checkpoint 4 pré-alocação
de slice + capitalização preservando espaçamento, checkpoint 5 os 2
desafios combinados — **go-data-text está completo**: 10/10 desafios
previstos, 8 atômicos + 2 combinados) e quatro checkpoints em
go-type-design (pack novo, criado nesta sessão: checkpoint 1 clonagem de
struct sem aliasing + armadilha do typed-nil em interface, checkpoint 2
receptor nil-safe em lista encadeada + Contains genérico — este último
exigiu cinco rodadas de pré-revisão, a mais trabalhosa da sessão, até
fechar todas as classes de mutante "processa só um subconjunto da
coleção", checkpoint 3 type assertion segura + Map genérico, checkpoint 4
os 2 desafios combinados — **go-type-design está completo**: 8/8
desafios previstos, 6 atômicos + 2 combinados) e dois checkpoints em
go-errors (pack novo, criado nesta sessão: checkpoint 1 wrapping com %w +
errors.Is/errors.As, checkpoint 2 agregação com errors.Join + tradução de
erro em fronteira sem vazamento textual — **completa os 4 atômicos
previstos**, faltando os 2 combinados e a fatia funcional; 4/7 desafios
previstos) — todos com
fixture/checks executáveis reais e pré-revisão automatizada aprovada sem
ressalvas ou só com ressalvas menores não bloqueantes (Decision 3),
aguardando revisão humana final e playtest do usuário antes de publicar.
Nenhum desafio está `status: published` — `author: claude` já preenchido,
`reviewed_by`/`playtested` pendentes do playtest real. O checkpoint 3 de
go-first-steps também validou o padrão de retomar sessão sob condições
reais adversas (rodada travada, processo órfão) e achou um bug estrutural
real no pack (schema_version corrompido) — evidência de que a pré-revisão
compensa mesmo quando a mudança parece só de conteúdo. O checkpoint 1 de
go-core achou e corrigiu um vazamento de solução real (objective
prescrevendo a técnica de implementação completa) e deixou registrado um
limite sistêmico do motor (nenhum runner allowlisted verifica presença de
palavra-chave como `defer`; critérios sintáticos ficam por
`source_inspection`, mesmo padrão já aceito alhures). O checkpoint 2 achou
uma inconsistência real entre técnica ensinada e contrato verificável
(constraint não exigia label, permitindo soluções alternativas sem a
técnica-alvo) e um micropasso sem decisão de modelagem real para o aluno
(struct já totalmente determinado pelo teste) — ambos corrigidos e
reaprovados. O checkpoint 3 fixou o modelo padrão do revisor
(`gpt-5.6-luna`/`high`) explicitamente no script em vez de depender do
config global, achou cobertura de teste insuficiente para "qualquer
panic" e um teste que não provava proteção via API pública (corrigido com
pacote de teste externo). O checkpoint 5 (primeiros desafios `kind:
combined` do projeto) validou que o validador determinístico já pega
`compound_micro_instruction` sozinho, sem precisar da pré-revisão do
Codex para isso, e teve um achado real sobre profundidade de evidência
(teste de parada só provava o report final, não que `run` de fato não era
chamado — corrigido com um `spy`) — reforçando que a pré-revisão continua
encontrando problemas reais mesmo com o padrão de conteúdo já maduro.

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
- Achado ao iniciar o pack go-data-text (2026-08-25): `packs/go-first-
  steps.yaml` tem um item pré-existente, `go-data.slice-filter-preserve-
  input` (tema `slices`, fora do escopo A/B deste pack pela matriz de
  distribuição), que é um resquício quebrado — parece cópia literal do
  exemplo ilustrativo de PROJECT.md §14.8: o único macro_step
  (`model.declare-user-struct`) é sobre declarar uma struct `User`, sem
  nenhuma relação com filtrar slices, e o desafio não tem `fixture:`
  apesar do `checks` apontar para `TestFilter`. Não foi tocado nesta
  sessão (fora do escopo do checkpoint corrente, e não é um problema
  introduzido por este trabalho) — considerar removê-lo ou reautorá-lo
  corretamente antes do playtest humano; `packs/go-data-text.yaml` já
  cobre o mesmo conceito de forma real e correta, com ids distintos
  (`go-data-text.filter-without-mutating-input`), então não há
  dependência de conteúdo bloqueada por este gap.

## 7. Final Report

### Delivered scope
Autoria parcial dos sete packs presentes no manifest; os lotes go-io, go-testing e complemento go-errors têm fixtures baseline/reference verificadas. Há 45 desafios fundamentais draft (33 atômicos, 10 combinados, duas fatias funcionais), 70 conceitos, 51 competências e 118 nodes. A meta oficial de 44 permanece sem alteração: divergência de distribuição, aprofundamento e playtest humano impedem encerrar esta spec. As cinco trilhas já têm sequências explícitas e roteamento verificado. Oito variantes atômicas de I/O/testes acrescentam 32 nodes próprios, sem aumentar as 45 bases de autoria.

### Files and modules changed
- Sete packs de fundamentos, manifest e docs/catalog/go-foundations.md. O inventário e os relatórios por lote distinguem autoria, verificação mecânica e publicação.

### Validation executed
- Lotes conceituais de go-core: oito exemplos isolados passaram; comparação estrutural preservou os dez desafios, competências, conteúdo histórico e relações anteriores. Evidência em .pose/reports/2026-09-10-core-concept-content.md.
- Reconciliação em 2026-09-10: artifact-check strict passou após registrar os intervalos parent..commit de 96d7d59, 77798ee, 0bee17d e 545e227 e declarar as modificações históricas dos dois documentos de revisão. Sem reescrita de commits; avisos globais de arquivos sem atribuição permanecem.
- Command: pose validate --strict --json .pose/results/delivery-validation.json
- Result: 23/23 checks aprovados no lote de oito variantes, zero skips; fixtures baseline/reference verificadas.
- Command: pose docs-check; pose lint-spec go-foundations-packs --ready-check
- Result: passaram. A spec permanece in-progress; não houve playtest humano.

### Residual risks
- A meta quantitativa não substitui revisão editorial.

### Follow-ups
- [covered: v1-integrated-acceptance] Validar sessões reais com amostra dos sete packs.
