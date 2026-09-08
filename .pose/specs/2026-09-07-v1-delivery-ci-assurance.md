---
slug: v1-delivery-ci-assurance
status: in-progress
created_at: 2026-09-07
completed_at:
supersedes:
depends_on: installation-documentation-ci, reliability-observability-compatibility
priority: 35
components: .github/workflows, internal/ciassurance, internal/mcpserver, cmd/ci-assurance, scripts, packs, docs, .pose/contracts
delivers: governance:v1-delivery-ci
---

# Spec: v1-delivery-ci-assurance

## 1. Intent

### Goal
Tornar os gates de entrega executáveis em CI e reconciliar proveniência, plataformas e documentação antes do aceite.

### Business value
Fechar uma lacuna verificável da V1 antes de ampliar a superfície que depende dela.

### Constraints
- Preservar autoria do aluno, operação local, JSONL e separação pedagógica.
- Usar o [relatório de auditoria](../reports/2026-09-07-doc-audit-auditoria-do-planejamento-v1.md) como baseline, não como evidência de entrega.
- Consumir knowledge:planning-audit-2026-09 antes de iniciar a implementação.

### Non-goals
- Reabrir ou reescrever atestações históricas de specs done.
- Publicar release, alterar código do aluno ou ampliar a V1 para nuvem/multiusuário.

## 2. Requirements

### Functional
- R1: Instalar versão fixada e verificada do POSE em CI e executar check/validate estritos, skills-check e resultados estruturados; ausência de ferramenta requerida falha o job.
- R2: Alinhar matriz, Makefile e CI para formatação com código de falha, race, vet, build, catálogo, contratos e fixtures; executar o gate quantitativo V1 apenas no candidato completo.
- R3: Executar testes e smoke-install nativos em Linux e macOS conforme RNF-005; Windows permanece desejável e não bloqueia V1 sem mudança de escopo. Distinguir build cruzado de execução por arquitetura.
- R4: Fixar actions e scanner por versão/digest verificado; tratar inputs de release como dados validados, sem interpolação direta em shell; manter build sem publicação automática.
- R5: Mapear componentes/contratos reais, habilitar políticas de artifacts/delivery com roots e evidências corretas, reconciliar rename ailearn→codinho sem falsificar histórico ou editar bundles imutáveis.
- R6: Registrar e validar critérios C1–C9 do roadmap por evidência atual; comprovar que stale, skipped, check inexistente e relatório ausente bloqueiam a entrega.
- R7: Governar documentação por manifest e revisar comandos/links; corrigir alegações de retomada e suporte até existirem evidências dos contratos correspondentes.
- R8: Produzir relatório de prontidão com matriz por requisito, plataforma, host, commit e resultado; obter revisão independente do threat model antes do aceite.

### Non-functional
- Manter determinismo e falhas explícitas; provar o caminho composto, não apenas helpers isolados.
- Identificar resultados por versão, commit e cenário; skips obrigatórios não são sucesso.

### Security
- Não copiar estado real do aluno para fixtures, logs ou relatórios.
- Aplicar confinamento de paths, redaction e consentimento nos novos caminhos.

### Compatibility
Preservar a política de release sem publish. Definir roots e equivalência de componentes com evidência; não desabilitar gates para fazer o closeout passar.

## 3. Technical Plan

### Affected areas
- ci
- governance
- distribution

### Artifacts
- created: cmd/ci-assurance/main_test.go
- modified: packs/go-first-steps.yaml
- modified: packs/go-core.yaml
- modified: packs/go-data-text.yaml
- modified: packs/go-type-design.yaml
- modified: packs/go-errors.yaml
- created: .pose/contributions/20260908-root-module-metadata-normalization.md
- created: internal/mcpserver/governance_contract_test.go
- created: .pose/contracts/historical-renames.json
- modified: .pose/indexes/repo-map.json
- modified: .github/workflows/ci.yml
- modified: .github/workflows/release.yml
- modified: Makefile
- modified: .pose/indexes/validation-matrix.json
- modified: .pose/indexes/module-metadata.json
- modified: .pose/policy/artifacts.json
- modified: .pose/policy/delivery.json
- modified: .pose/policy/docs.json
- modified: .pose/roadmaps/codinho-v1.md
- modified: docs/compatibility.md
- modified: docs/install.md
- modified: docs/quickstart.md
- created: docs/acceptance/v1-release-readiness.md

- modified: README.md
- modified: docs/configuration.md
- modified: docs/content-authoring.md
- modified: docs/content-review-checklist.md
- modified: docs/agent-review-workflow.md
- modified: docs/troubleshooting.md
- modified: docs/security/privacy.md
- modified: docs/security/executor-limitations.md
- modified: docs/security/threat-model.md
- created: scripts/ci/install-tools.sh
- created: scripts/ci/tools.env
- created: scripts/ci/check-format.sh
- created: scripts/ci/validate.sh
- created: scripts/ci/build-release.sh
- created: cmd/ci-assurance/main.go
- created: internal/ciassurance/evidence.go
- created: internal/ciassurance/evidence_test.go
- created: internal/ciassurance/workflows_test.go
- created: .pose/docs.json
- created: .pose/contracts/mcp-stdio.json
- created: .pose/adr/2026-09-08-native-ci-evidence-and-prospective-delivery-governance.md
- created: .pose/knowledge/2026-09-08-decision-log-adr-ci-assurance-review.md
- created: .pose/changelogs/unreleased/v1-delivery-ci-assurance.md

### Delivery targets
- governance:v1-delivery-ci module:.github/workflows profile:release-governance entrypoint:.github/workflows/ci.yml

### API/contract changes
Preservar a política de release sem publish. Definir roots e equivalência de componentes com evidência; não desabilitar gates para fazer o closeout passar.

### Data/storage changes
Documentar os campos e migração exigidos pelos requisitos antes do primeiro incremento. Nenhuma migração é executada nesta rodada de planejamento.

### Technical risks
Renames históricos e evidências de módulos divergentes bloqueiam roll-up. CI remota e revisão humana precisam de evidência real, não de YAML existente.

## 4. Tasks

### Planning
- [x] Reproduzir o achado e revisar contratos/ADRs aplicáveis.
- [x] Completar decisões de formato e plano de testes negativos antes de modificar código.
- [x] Reconciliar esta lista de artefatos com os arquivos efetivos; declarar arquivos adicionais antes de alterá-los.

### Implementation
- [x] Completar validation privado dos rascunhos com checks: preservar fixture inicial, reutilizar seus testes e fornecer referência executável; não mudar publication.
- [x] Implementar primeiro o menor fluxo que fecha a lacuna.
- [x] Integrar entradas reais, persistência/compatibilidade e diagnósticos.
- [x] Atualizar documentação e checks declarativos junto com o contrato.

### Validation
- [ ] Executar os cenários de cada R-ID, incluindo negativos.
- [ ] Executar pose assess integrate e validação estruturada no candidato.
- [ ] Reconciliar artifacts, surface e revisão independente antes do closeout.

## 5. Decisions

### Decision 3
- Date: 2026-09-08
- Decision: incluir provas privadas dos checks já declarados nos cinco packs existentes para satisfazer R2. Reutilizar os testes do desafio na referência e registrar a falha de baseline observada; o protótipo sem fixture usa baseline_fixture com justificativa.
- Rationale: o gate --checks exige prova de rascunhos também; placeholders não constituem referência executável. Alterar o gate esconderia a lacuna.
- Consequences: escopo inclui somente validation privado; autoria/revisão humana, status de publicação e gate quantitativo V1 permanecem separados. Validar que as fixtures distribuídas e contratos públicos foram preservados.

### Decision 2
- Date: 2026-09-08
- Decision: registrar ADR native-ci-evidence-and-prospective-delivery-governance antes da implementação. Consumidos knowledge:planning-audit-2026-09 e knowledge:adr-publication-review.
- Rationale: jobs locais não provam macOS; evidência histórica não deve ser reescrita para satisfazer regras novas. Versões externas precisam ser fixadas e conferidas por origem/hashes.
- Consequences: CI executa e arquiva provas nativas por plataforma; publicação de release permanece fora do workflow. Políticas prospectivas e reconciliação explícita preservam registros imutáveis. C9 já existente no roadmap também deve entrar na matriz de prontidão.

### Decision 1
- Date: 2026-09-07
- Context: ci.yml não instala POSE nem executa pose check/validate; skills-check é pulado sem binário. A matriz não inclui catálogo, checks de fixtures, formatação com falha nem race. delivery.json e artifacts.json estão desativados; roadmap não tinha Cut criteria; o bundle esbarra em cmd/ailearn/main.go histórico inexistente.
- Options considered: deixar implementação implícita no aceite; criar remediação focada.
- Decision: planejar remediação própria e exigir sua entrega antes do aceite.
- Rationale: v1-integrated-acceptance tem non-goal de adicionar features ou correções.
- Consequences: esta spec permanece draft; decisões estruturais novas exigem ADR no início da implementação.

## 6. Validation

### Strategy
Risco alto: CI, execução de ferramentas e evidência de entrega. Plano antes do código:

| Cenário | Comando | Evidência esperada |
|---|---|---|
| Contratos de workflows, pins, inputs e evidência negativa | go test ./internal/ciassurance/... -count=1 | Rejeita versão com metacaracteres, relatório ausente/stale, skip e check inexistente; versões/hashes fechados. |
| Matriz local composta | make check | POSE obrigatório, formatação/race/vet/build/catálogo/MCP/smoke e JSON; ferramenta ausente falha. |
| Documentos e referências | pose docs-check; pose check --strict | Manifest cobre docs e links; nenhum erro. |
| Plataformas nativas | GitHub Actions CI em Linux/macOS | Jobs reais, commit/OS/arquitetura/toolchain e resultados completos, sem equiparar cross-build a execução. |
| Release somente build | scripts/ci/build-release.sh v0.0.0-ci | Binário e checksum; input inválido falha antes de build, sem tag/release. |
| Integridade do candidato | pose artifact-check --spec v1-delivery-ci-assurance --strict; pose surface-check --spec v1-delivery-ci-assurance --strict | Origem atribuída, entrada de produção e evidência atual. |

Comandos são obrigatórios. Windows é checagem cruzada explicitamente distinta;
aceite quantitativo/editorial V1 e hosts de tutor continuam no aceite integrado.
Gates do roadmap são exercitados como prontidão (falhas pendentes esperadas),
sem declarar V1 pronta. Revisão independente cobre também threat model.


### Deterministic checks
- Test: pose check --strict; pose validate --strict; pose skills-check --strict; pose docs-check
- Lint: pose check --strict; pose lint-spec v1-delivery-ci-assurance --ready-check
- Typecheck: go vet ./...
- Build: go build ./cmd/codinho
- Security / Contract: pose assess integrate; pose validate --strict --json .pose/results/delivery-validation.json; pose surface-check --spec v1-delivery-ci-assurance --strict

### Execution log
- 2026-09-08 UTC (CI nativa): run 34273134508 executou o commit 649725999449ebbdb77dea0a7880fb60fe471f54 com sucesso em Linux/amd64 e macOS/arm64. Artefatos baixados e conferidos: Go 1.27.1, host github-actions, 19 checks pass, source_modified=false. Relatório de prontidão inclui a matriz e o link do provider. Go 1.27.1 conferido em https://go.dev/VERSION (acesso 2026-09-08).
- 2026-09-08 UTC (revisão local): threat model distingue controle de download de módulos de bloqueio de sockets, registra corrida residual de symlinks e explicita limites de subprocessos/CI. Mapeamento por paths reais eliminou contratos não classificados do bundle. Revisão independente solicitada, ainda não atestada.
- 2026-09-08 UTC (provas editoriais): validation privado adicionado a 37 desafios em cinco packs, reutilizando os testes existentes. O protótipo sem fixture recebeu baseline_fixture privada com justificativa. Comparação estrutural antes/depois confirmou fixture distribuída, checks, acceptance e publication idênticos. `catalog validate --checks --json` passou: 38/38 checks, 76/76 cenários, zero achados blocking. `make check` passou 19/19 checks. Não houve publicação ou atestado pedagógico.
- 2026-09-08 UTC (proveniência): commit dbd848e atribuiu o primeiro incremento à spec; artifact-check estrito passou sem erros (avisos de órfãos históricos). Novas provas privadas serão atribuídas no próximo commit. Relatório nativo passa a declarar source_modified e a rejeitar fonte alterada durante GitHub Actions; teste cobre fonte tracked/untracked e separa relatórios gerados.
- 2026-09-08 UTC (continuação): policies artifacts/delivery habilitadas com roots reais; metadados de CLI, MCP, governança e distribuição registrados. Inventário de 33 tools comparado com tools/list do servidor; rename 9d88c68f7a1a3a2d1a84730a724fa28f99bb9821 verificado por git show --find-renames. Bundles históricos preservados. `go test -race ./internal/ciassurance ./internal/mcpserver -count=1` passa.
- 2026-09-08 UTC (continuação): pose index não associa chave raiz "." ao path vazio emitido; alias equivalente aplicado e projeção passou a declared/isComplete=true. Limitação sanitizada registrada localmente em .pose/contributions/20260908-root-module-metadata-normalization.md. Nenhum envio upstream.
- 2026-09-08 UTC: retomada consumiu knowledge:adr-ci-assurance-review e confirmou a spec in-progress. Discover executado antes de editar; tools POSE 1.7.12 e govulncheck 1.6.0 instaladas em diretório temporário com hashes verificados.
- 2026-09-08 UTC: testes em internal/ciassurance passam para relatório inválido/ausente, stale/futuro, skip, comando divergente, check desconhecido/duplicado, matriz inválida e labels inseguros. Corrigida rejeição de checks duplicados na matriz e SHA não hexadecimal.
- 2026-09-08 UTC: corrigido SIGPIPE de `pose version | head` sob pipefail; teste executa o script real com ferramentas sintéticas e verifica propagação da falha do gate. Formatação inclui arquivos novos ainda não staged.
- 2026-09-08 UTC: `pose docs-check` passa com 13 documentos, zero erros/avisos; corrigidos diretório de instalação, remoção do binário e alegações de plataforma. `pose check --strict` e `pose lint-spec v1-delivery-ci-assurance --ready-check` passam.
- 2026-09-08 UTC: `PATH=/tmp/codinho-ci-tools:$PATH make check` fora do sandbox executou 19 checks: 18 pass, catalog fail. Race completo, vet, build, govulncheck, MCP, CLI, startup, smoke e ci-assurance passaram. O catálogo tem 38 checks declarados, zero verificados; rascunhos sem validation/reference_fixture bloqueiam `catalog validate --checks`. Gate preservado; relatório nativo de sucesso não foi emitido.
- 2026-09-08 UTC: `pose assess integrate` detecta zero contratos (limitação já registrada no handoff de auditoria); isso não comprova integração MCP. `pose assess tech-debt` retorna zero marcadores. `pose surface-check` e `pose artifact-check` falham por proveniência/artefatos pendentes. `pose roadmap-check codinho-v1 --strict` mantém terminal=false, nove critérios e blockers de documentos de aceite ausentes e specs não terminais.
- 2026-09-08 UTC: estado lido e discover executado antes das alterações:1módulo,12.127LOCprod/11.524LOCtest,zero marcadores. CI atual pula skills sem POSE, usa actions por tags e scanner latest; release interpola input em shell. Último job remoto34168152576 falhou em govulncheck.
- 2026-09-08 UTC: origem oficial verificada via GitHub API: POSE v1.7.12 commit cd9a050f4de9d74f8b695f5f0c5bdeabda96a9d4, digests dos assets; checkout/setup-go/upload-artifact resolvidos para commits dos respectivos repositórios. Referência de segurança: https://docs.github.com/en/actions/reference/security/secure-use (acesso2026-09-08).
- 2026-09-07 UTC: criada em revisão de planejamento; implementação e gates de entrega não executados.

### Results summary
Matriz local passou 19/19 checks após completar as provas editoriais. Proveniência do primeiro incremento passou; evidências do candidato final e revisão independente ainda governam closeout.

### Requirement trace
- R1: ferramentas verificadas e gates obrigatórios implementados; make check passou localmente após completar as provas privadas.
- R2: matriz compartilhada executa os checks; 38 checks editoriais verificados em 76 cenários.
- R3: run 34273134508 passou em Linux/amd64 e macOS/arm64; veja matriz por commit no relatório de prontidão.
- R4: labels inseguros rejeitados em teste; build local v0.0.0-ci gerou binário/checksum sem publicar.
- R5: policies, metadados, inventário MCP e mapa histórico implementados/testados; reconciliação Git e gate de proveniência em andamento.
- R6: negativos do validador passam; roadmap C1–C9 permanece não terminal.
- R7: docs-check passa para 13 documentos; comandos e alegações de plataforma corrigidos.
- R8: relatório de prontidão registra pendências; evidências remotas conferidas; revisão independente pendente.

### Known gaps
Renames históricos e evidências de módulos divergentes bloqueiam roll-up. CI remota e revisão humana precisam de evidência real, não de YAML existente.

## 7. Final Report

### Delivered scope
Workflows, instalação verificada, validação de evidência e manifest documental implementados no worktree. Testes negativos e checks locais registrados acima. Spec permanece in-progress: evidência final de proveniência, execução nativa remota e revisão independente pendentes.

### Files and modules changed
- CI, scripts de ferramentas/validação/build, internal/ciassurance, cmd/ci-assurance, matriz, manifest e documentação declarados nos artefatos.

### Validation executed
- Consulte Execution log: validação composta executada com falha editorial explícita, checks de docs/estrutura e testes negativos passando.

### Residual risks
Validar o comportamento implementado em execução independente; não reutilizar resultado histórico como aprovação do código futuro.

### Follow-ups
- [covered: v1-integrated-acceptance] Reexecutar os requisitos desta spec no candidato composto da V1.
