---
slug: go-production-architecture-packs
status: in-progress
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: curriculum-graph-path-recommendation, catalog-authoring-quality, go-backend-packs
priority: 200
components: curriculum-content
delivers:
---

# Spec: go-production-architecture-packs

## 1. Intent

### Goal
Publicar packs de produção, arquitetura e depuração que completem a taxonomia e o volume curricular da V1.

### Business value
Praticar decisões seniores, investigação, operabilidade e evolução de sistemas em vez de apenas aumentar o número de arquivos.

### Constraints
- Publicar exatamente 13 desafios: 2 atômicos, 2 combinados, 2 fatias funcionais, 4 depuração/refatoração/revisão e 3 missões.
- Adicionar pelo menos 20 conceitos, 15 competências e 80 step nodes, atingindo 160 conceitos, 100 competências e 500 nodes cumulativos.
- Missões devem conter requisito incompleto, falhas e trade-offs.

### Non-goals
- Cobertura exaustiva de todo broker, banco, cloud ou framework.
- Ensinar padrões por nome sem força de design.

## 2. Requirements

### Functional
- R1: Publicar go-production e go-architecture-debugging.
- R2: Cobrir mensageria e gRPC conceituais, observabilidade, performance, runtime, segurança, arquitetura e depuração.
- R3: Entregar 2 atômicos sobre profiling e correlação observável.
- R4: Entregar 2 combinados sobre shutdown observável e cliente idempotente.
- R5: Entregar 2 fatias sobre worker operável e serviço instrumentado.
- R6: Entregar 4 desafios de depuração, refatoração ou revisão sobre race, retenção, overengineering e incidente.
- R7: Entregar 3 missões sobre backpressure, migração compatível e resiliência operacional.
- R8: Publicar duas trilhas: depuração e performance; arquitetura e produção.
- R9: Exigir justificativa de custo, falha, operabilidade e alternativa em toda missão.
- R10: Comprovar as metas cumulativas de 160 conceitos, 100 competências e 500 step nodes.

### Non-functional
- Desafios seniores podem ter pouco código e devem avaliar decisão.
- Checks de performance usam baseline local e tolerâncias documentadas.

### Security
- Incluir input limits, secrets, TLS, auth, supply chain e disponibilidade em casos pertinentes.
- Não usar serviços reais ou credenciais.

### Compatibility
- Contratos e migrações devem incluir backward compatibility e rollout progressivo.

## 3. Technical Plan

### Affected areas
- `packs/go-production.yaml`, `packs/go-architecture-debugging.yaml`,
  `packs/manifest.yaml`, auditorias de catálogo em Go, documentação do catálogo
  e inventário editorial.

### Artifacts
- modified: packs/manifest.yaml
- created: packs/go-production.yaml
- created: packs/go-architecture-debugging.yaml
- created: testdata/packs/production/README.md
- created: docs/catalog/go-production-architecture.md
- created: docs/editorial/go-production-architecture-inventory.md
- created: docs/editorial/go-production-architecture-queue.json
- modified: .pose/docs.json
- modified: cmd/codinho/foundation_tracks_integration_test.go
- modified: internal/curriculum/foundation_variants_test.go
- modified: .pose/assessments/README.md
- modified: .pose/assessments/cmd-codinho.md
- modified: .pose/assessments/consolidated.md
- modified: .pose/assessments/docs-catalog.md
- modified: .pose/assessments/docs-editorial.md
- modified: .pose/assessments/internal-curriculum.md
- modified: .pose/assessments/packs.md
- modified: .pose/state/components/cmd-codinho.json
- modified: .pose/state/components/docs-catalog.json
- modified: .pose/state/components/docs-editorial.json
- modified: .pose/state/components/internal-curriculum.json
- modified: .pose/state/components/packs.json
- modified: .pose/indexes/delivery-integrity.json

### Delivery targets
Nenhum tipado; conteúdo do catálogo.

### API/contract changes
- Consumir schemas e checks existentes; extensão exige spec própria.

### Data/storage changes
- Fixtures sintéticas de logs, traces, perfis, incidentes e migrações.

### Technical risks
- Missões podem ficar subjetivas demais para avaliar.
- Benchmarks podem variar por hardware.

## 4. Tasks

### Planning
- [x] Definir matriz dos 13 desafios, IDs, temas, competências e distribuição no inventário editorial.
- [x] Definir fixtures herméticas, cenários negativos e limites de benchmark antes da autoria.
- [x] Definir fila serializável com paths e aceite executável por pack e integração.

### Implementation
- [ ] Autorar os dois packs e seus 13 desafios, conceitos, competências, relações e fixtures.
- [ ] Integrar os dois packs no manifest, documentar as trilhas e reservar testdata sem duplicar fixtures.
- [ ] Atualizar regressões de inventário global e verificar alcance MCP das duas novas trilhas.
- [ ] Realizar playtests com revisão sênior.

### Validation
- [ ] Executar validação estrutural/editorial e os checks de todas as fixtures, inclusive race e benchmark local.
- [ ] Auditar tipos, IDs, conceitos, competências, nodes, trilhas e metas cumulativas sem promoção.
- [ ] Executar testes Go, race, vet, build, validação POSE e revisão de segurança/compatibilidade.
- [ ] Completar revisão pedagógica e playtests humanos antes de qualquer publicação.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Senioridade não é proporcional ao tamanho do projeto.
- Options considered: projetos maiores; algoritmos difíceis; decisões sob restrições.
- Decision: usar investigações e missões com falhas e trade-offs.
- Rationale: aproxima competências reais de engenharia.
- Consequences: rubricas e revisão humana têm papel maior.

### Decision 2
- Date: 2026-09-16
- Context: tempo absoluto de benchmark varia por hardware, toolchain e carga local.
- Options considered: limite absoluto global; benchmark apenas descritivo; comparação pareada local com tolerância.
- Decision: usar checks de benchmark apenas para provar que a medição é executável e comparar medianas locais repetidas; diferenças abaixo de 15% são inconclusivas e não se tornam gate global.
- Rationale: preserva evidência prática sem converter ruído de máquina em avaliação de correção.
- Consequences: nenhuma comparação entre hosts será usada como aprovação; interpretação final ainda requer contexto do desafio.

## 6. Validation

### Strategy
Combinar checks determinísticos, rubricas, fixtures de incidente e playtest independente.

### Test plan

| Cenário | Classificação | Comando/check | Evidência esperada |
|---|---|---|---|
| Inventário e schema | Obrigatório | `./codinho catalog validate --json` | Zero diagnostics/editorial bloqueantes; 13 desafios novos, IDs resolvidos e tudo em draft. |
| Baseline/reference de todas as fixtures | Obrigatório | `./codinho catalog validate --checks --json` | Cada check declarado executa; baseline corresponde ao resultado esperado e toda referência passa; nenhuma fixture usa rede. |
| Corrida e encerramento sob concorrência | Obrigatório | `go_test_race` nas fixtures de race/backpressure e `go test -race -count=1 ./...` | Sem data race na referência; baseline reproduz o defeito previsto; cancelamento e drain têm limite explícito. |
| Perfis/benchmark | Obrigatório | Checks `go_benchmark` de `go-production.profile-hot-path` e `go-production.correlate-request-signals` | Benchmark compila e emite medição local; nenhum limite absoluto de ns/op é inferido do runner. |
| Comparação de desempenho no mesmo host | Opcional | `./codinho workspace prepare go-production.profile-hot-path --dest /tmp/codinho-profile-hot-path && (cd /tmp/codinho-profile-hot-path && go test -run '^$' -bench '^BenchmarkHotPath$' -benchmem -count=5 -benchtime=1s)`; `./codinho workspace prepare go-production.correlate-request-signals --dest /tmp/codinho-request-correlation && (cd /tmp/codinho-request-correlation && go test -run '^$' -bench '^BenchmarkRequestCorrelation$' -benchmem -count=5 -benchtime=1s)` | Cinco amostras e medianas sob o mesmo toolchain/carga; diferenças abaixo de 15% são inconclusivas. |
| Entrada inválida/oversized, identidade e segredo | Obrigatório | Checks locais em `go-production.instrumented-service-boundary` | Rejeição limitada, ausência de segredo em resposta/log e identidade local validada; sem credenciais reais. |
| Dependência indisponível, timeout, retry e fila saturada | Obrigatório | Fixtures locais de `idempotent-retry-client`, `operable-message-worker` e missão backpressure | Retry limitado/idempotente, contexto cancelado, fila bounded e ack/redelivery explícitos; sem serviço remoto. |
| Compatibilidade de schema e rollback | Obrigatório | Fixture sintética de `go-production.compatible-schema-rollout` | Leitor antigo e novo coexistem no estágio expand; migração não destrutiva e retorno de versão demonstrado. |
| Segredos no diff | Obrigatório | `git diff --cached --unified=0 | rg --color never -n -i '(-----BEGIN [A-Z ]*PRIVATE KEY-----|AKIA[0-9A-Z]{16}|gh[pousr]_[A-Za-z0-9]{20,}|sk-[A-Za-z0-9]{20,}|Bearer [A-Za-z0-9._~+/=-]{16,})'` | Nenhuma correspondência; saída vazia e status 1 do `rg` significam que os padrões não foram encontrados. Valores de fixture são sintéticos. |
| Verificações estáticas e regressão do módulo | Obrigatório | `bash scripts/ci/check-format.sh`; `go vet ./...`; `go build ./cmd/codinho`; `go test -race -count=1 ./...` | Formato limpo, sem diagnósticos de vet, build e suite completos aprovados. |
| Matriz nativa POSE | Obrigatório | `pose validate --strict --module . --report` | Relatório atualizado; limitações externas, como indisponibilidade de DNS para govulncheck, permanecem explicitamente não aprovadas. |
| Comparar benchmark entre hosts ou julgar pedagogia/playtest | Fora da automação | Revisão independente humana | Nunca converter em gate automático nem declarar satisfeito nesta execução. |

### Deterministic checks
- Test: checks de todas as fixtures, race detector e benchmarks controlados.
- Lint: go run ./cmd/codinho catalog validate --json (executar na raiz; catálogo resolvido pelo manifest).
- Typecheck: go vet ./...
- Build: go build ./cmd/codinho
- Security / Contract: threat cases, compatibility fixtures e secret scan.

Não alterar `.pose/indexes/validation-matrix.json`: o check `catalog` já executa
as fixtures declaradas e a nova cobertura pertence ao catálogo, não a um novo
gate global.

### Execution log
- 2026-09-16: baseline conferido em `1fde57a`: `./codinho catalog validate --json` reporta 145 conceitos, 89 competências, 77 desafios, 501 nodes e oito trilhas, todos draft.
- 2026-09-16: `pose assess discover` passou para `packs`, `cmd/codinho`, `internal/curriculum`, `docs/catalog` e `docs/editorial`; LOC/debt metrics foram registrados. Sem TODO/FIXME nos componentes avaliados.
- 2026-09-16: continuação sob instrução explícita do usuário. `go-backend-packs` permanece `in-progress` com automação feita e gates humanos pendentes; `go-foundations-packs` permanece `in-progress` com somente gates humanos pendentes. Nenhuma dependência foi falsamente marcada como `done`.
- 2026-09-16: matriz e fila criadas antes da autoria. Alvos da política confirmados em `packs/distribution.json`: `go-production` 2/2/2/0/3 e `go-architecture-debugging` 0/0/0/4/0 (atomic/combined/functional/investigation/mission).

### Results summary
- Planejamento editorial definido; nenhum pack de produção autorado ou publicado ainda.

### Requirement trace
- R1: IDs/paths reservados para `go-production` e `go-architecture-debugging` na matriz editorial; integração ainda pendente.
- R2–R8: domínios, 13 desafios e duas trilhas distribuídos por ID na matriz editorial; autoria pendente.
- R9: cada uma das três missões tem requisitos de trade-off/custo/falha/operabilidade/alternativa no aceite da fila.
- R10: baseline e soma requerida (160/100/500) declaradas; cobertura resultante será medida após integração.

### Known gaps
- Tolerâncias de performance serão específicas por desafio, não globais.
- A spec depende editorialmente de `go-backend-packs`, ainda `in-progress` porque gates humanos foram pulados por instrução; o conteúdo backend e suas fixtures determinísticas estão integrados como drafts.
- As comparações locais de profiling são evidência auxiliar e não substituem julgamento sobre a adequação da otimização.

## 7. Final Report

### Delivered scope
Parcial: matriz e fila autoral estão definidas; os packs ainda não foram materializados.

### Files and modules changed
- Planejados em dois packs, documentação de catálogo/editorial, fila, fixtures e regressões de inventário/roteamento.

### Validation executed
- Command: `./codinho catalog validate --json`
- Result: baseline limpo, 145 conceitos, 89 competências, 77 desafios, 501 nodes e oito trilhas; todo o conteúdo está em draft.
- Command: `pose lint-spec go-production-architecture-packs --ready-check`
- Result: pass antes da autoria; repetir após cada integração.

### Residual risks
- Critérios semânticos precisam de calibração interavaliador.

### Follow-ups
- [covered: v1-integrated-acceptance] Auditar metas cumulativas e uma missão completa.
