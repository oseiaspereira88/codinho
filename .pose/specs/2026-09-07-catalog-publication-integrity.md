---
slug: catalog-publication-integrity
status: in-progress
created_at: 2026-09-07
completed_at:
supersedes:
depends_on: catalog-authoring-quality, catalog-schema-loader, safe-check-executor
priority: 30
components: curriculum, cli, mcp-server
delivers: capability:catalog-publication-integrity
---

# Spec: catalog-publication-integrity

## 1. Intent

### Goal
Separar inventário autorado de conteúdo publicável e ligar os gates editoriais à CLI e à CI.

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
- R1: Projetar separadamente inventário, rascunhos, publicações e conteúdo elegível aos gates; contar nos 84 apenas desafios curados, sem duplicar variantes ou protótipos.
- R2: Aplicar distribuição declarativa por pack e global na CLI, incluindo tipos inesperados; ligar CheckTypeDistribution ao caminho de produção e testar que uma contagem divergente falha.
- R3: Rejeitar publicação inválida antes do uso normal pelo MCP; busca/recomendação padrão só mostram conteúdo publicado, com modo explícito de autoria/playtest para rascunhos locais.
- R4: Executar checks com expectativa declarada de baseline e solução de referência; falha pedagógica prevista é distinta de timeout, skip, teste inexistente e erro de compilação inesperado. Gabaritos ficam fora da projeção pública.
- R5: Detectar check sem fixture ou evidência de execução alternativa e registrar justificativa por desafio; comprovar 100% dos checks publicáveis em CI ou por fixture.
- R6: Validar um corpus compartilhado de documentos válidos/inválidos contra schemas JSON e loader Go, incluindo publication, relations, fixture, variantes e campos desconhecidos.
- R7: Definir comportamento para metadados legados ausentes sem promovê-los a published; preservar IDs e sessões existentes, sem excluir trabalho autorado.
- R8: Provar pela CLI real que catálogo numericamente suficiente mas sem revisão/playtest falha o gate V1; manter validação estrutural incremental utilizável.

### Non-functional
- Manter determinismo e falhas explícitas; provar o caminho composto, não apenas helpers isolados.
- Identificar resultados por versão, commit e cenário; skips obrigatórios não são sucesso.

### Security
- Não copiar estado real do aluno para fixtures, logs ou relatórios.
- Aplicar confinamento de paths, redaction e consentimento nos novos caminhos.

### Compatibility
Registrar ADR para visibilidade, distribuição e expectativas de checks antes do código. Reusar publication e o gate humano existentes; geração dinâmica/quarentena fica em agent-authored-catalog-drafts.

## 3. Technical Plan

### Affected areas
- curriculum
- cli
- mcp-server

### Artifacts
- modified: internal/curriculum/coverage.go
- modified: internal/curriculum/coverage_test.go
- modified: internal/curriculum/model.go
- modified: internal/curriculum/index.go
- modified: internal/curriculum/loader.go
- modified: internal/curriculum/playtest.go
- modified: internal/curriculum/editorial.go
- modified: internal/checks/executor.go
- modified: internal/cli/catalog.go
- modified: internal/cli/editorial.go
- modified: internal/cli/editorial_test.go
- modified: internal/cli/cli_test.go
- modified: cmd/codinho/main.go
- modified: cmd/codinho/integration_test.go
- modified: cmd/codinho/recovery_integration_test.go
- modified: cmd/codinho/tree_progression_integration_test.go
- modified: cmd/codinho/evaluation_evidence_integration_test.go
- modified: schemas/catalog.schema.json
- modified: schemas/pack.schema.json
- modified: schemas/challenge.schema.json
- modified: docs/content-authoring.md
- modified: docs/compatibility.md
- modified: .github/workflows/ci.yml
- modified: .pose/indexes/validation-matrix.json
- modified: go.mod
- created: internal/curriculum/publication.go
- created: internal/curriculum/publication_test.go
- created: internal/curriculum/distribution.go
- created: internal/curriculum/schema_contract_test.go
- created: internal/checks/editorial.go
- created: internal/checks/editorial_test.go
- created: internal/cli/publication_test.go
- created: cmd/codinho/publication_integration_test.go
- created: testdata/catalog-quality/distribution.json
- created: testdata/catalog-quality/schema-corpus.json
- created: packs/distribution.json
- created: .pose/adr/2026-09-08-publication-projections-and-executable-editorial-proof.md
- created: .pose/knowledge/2026-09-08-decision-log-adr-publication-review.md
- created: .pose/changelogs/unreleased/catalog-publication-integrity.md

### Delivery targets
- capability:catalog-publication-integrity module:cmd/codinho profile:composed-capability entrypoint:cmd/codinho/main.go

### API/contract changes
Registrar ADR para visibilidade, distribuição e expectativas de checks antes do código. Reusar publication e o gate humano existentes; geração dinâmica/quarentena fica em agent-authored-catalog-drafts.

### Data/storage changes
Documentar os campos e migração exigidos pelos requisitos antes do primeiro incremento. Campos novos opcionais usam omitempty ou json:"-"; não reescrever JSONL. Publicação de packs governa entidades auxiliares, publicação individual governa desafios; canonical opt-in exclui protótipos e variant_of exclui derivados da contagem.

### Technical risks
O catálogo real ainda não tem conteúdo publicado. Autoria usa serve --authoring; dados sintéticos de teste não constituem revisão/playtest do catálogo real.

## 4. Tasks

### Planning
- [x] Reproduzir o achado e revisar contratos/ADRs aplicáveis.
- [x] Completar decisões de formato e plano de testes negativos antes de modificar código.
- [x] Reconciliar esta lista de artefatos com os arquivos efetivos; declarar arquivos adicionais antes de alterá-los.

### Implementation
- [x] Implementar primeiro o menor fluxo que fecha a lacuna.
- [x] Integrar entradas reais, persistência/compatibilidade e diagnósticos.
- [x] Atualizar documentação e checks declarativos junto com o contrato.

### Validation
- [x] Executar os cenários de cada R-ID, incluindo negativos.
- [ ] Executar pose assess integrate e validação estruturada no candidato.
- [ ] Reconciliar artifacts, surface e revisão independente antes do closeout.

## 5. Decisions

### Decision 2
- Date: 2026-09-08
- Decision: aplicar [ADR de publicação](../adr/2026-09-08-publication-projections-and-executable-editorial-proof.md). Consumidos knowledge:planning-audit-2026-09 e knowledge:adr-local-versioned-catalog-and-event-state-review.
- Rationale: inventário e visibilidade não são evidência de curadoria; prova executável não substitui revisão humana.
- Consequences: política versionada externa ao YAML de conteúdo; modo de autoria explícito; corpus testa o contrato estrutural compartilhado, sem alegar equivalência de regras semânticas globais com JSON Schema.

### Decision 1
- Date: 2026-09-07
- Context: ProjectCoverage conta todos os desafios, incluindo 41 itens sem publicação. CheckTypeDistribution é chamado apenas por testes. runChecksAgainstFixtures ignora desafios sem fixture e aceita fail/skipped como sucesso de infraestrutura; schemas JSON e loader Go não têm teste de equivalência.
- Options considered: deixar implementação implícita no aceite; criar remediação focada.
- Decision: planejar remediação própria e exigir sua entrega antes do aceite.
- Rationale: v1-integrated-acceptance tem non-goal de adicionar features ou correções.
- Consequences: esta spec permanece draft; decisões estruturais novas exigem ADR no início da implementação.

## 6. Validation

### Strategy
Risco médio/alto: composição MCP, visibilidade de gabaritos e execução de código local.
Antes do código: R1/R2 comparar inventário/publicado/elegível, variante/protótipo,
contagem global e por pack, tipo inesperado; R3/R7 subprocesso stdio padrão vs
serve --authoring e replay legado. R4/R5 executar baseline e referência em roots
isoladas: pass, falha de teste prevista, compilação inesperada, skip, zero testes,
timeout, fixture ausente, alternativa sem justificativa, referência incorreta.
R6 corpus JSON sintético compartilhado com loader YAML, campos desconhecidos e
metadados inválidos; regras semânticas de referências/revisão testadas à parte.
R8 CLI real com catálogo suficiente de rascunhos deve falhar V1 e passar validação
incremental. Suites completas, race, vet/build, CI e revisão independente no
candidato final. Nenhum metadado humano real será criado pela automação.

### Deterministic checks
- Test: go test -race ./internal/curriculum/... ./internal/cli/... ./internal/mcpserver/...
- Lint: pose check --strict; pose lint-spec catalog-publication-integrity --ready-check
- Typecheck: go vet ./...
- Build: go build ./cmd/codinho
- Security / Contract: pose assess integrate; pose validate --strict --json .pose/results/delivery-validation.json; pose surface-check --spec catalog-publication-integrity --strict

### Execution log
- 2026-09-08 UTC: rodada 1 da revisão publication_review rejeitou f9ca9c8 por V1 sem checks, timeout genérico, fail+skip e pack ausente com meta zero. Corrigidos com regressões reais; achado de variant_of retirado após distinguir derivado publicado de variants privado, esclarecido no ADR/doc e testado no MCP.
- 2026-09-08 UTC: três casos adicionais do corpus reproduziram aceitação de publicação sem título, com tipo desconhecido ou schema futuro. Corrigida validação do desafio publicado sem bloquear rascunhos mínimos. go test -race ./internal/curriculum/... ./internal/checks/... ./internal/cli/... ./cmd/codinho/... -count=1 passou após todas as correções.
- 2026-09-08 UTC: matriz POSE inicial passou 16/16 em f9ca9c8 (incluindo govulncheck). Execução do revisor sem PATH do govulncheck sobrescreveu o resultado com falha ambiental; o autor reexecutará a matriz no novo candidato. Nenhuma falha ambiental foi contada como sucesso.
- 2026-09-08 UTC: pose assess discover antes das alterações; 11.335 LOC de produção, 10.887 LOC de testes, sem marcadores de dívida.
- 2026-09-08 UTC: reprodução do filtro ausente por inspeção e contratos negativos. ADR/plano registrados antes do código. Implementados projeções de publicação, distribuição, provas de execução, corpus e composição.
- 2026-09-08 UTC: go test ./internal/curriculum/... ./internal/cli/... ./internal/checks/... passou; go test ./cmd/codinho -run TestPublicationIntegrity -count=1 passou. go test -race ./... passou em todos os pacotes; último ajuste de saída truncada será revalidado no candidato.
- 2026-09-08 UTC: pose assess integrate executado (zero contratos detectados pelo scanner; testes MCP reais provam composição). pose assess tech-debt: zero marcadores. Revisão independente publication_review iniciada por agent-batch-review, ainda sem decisão final.
- 2026-09-07 UTC: criada em revisão de planejamento; implementação e gates de entrega não executados.

### Results summary
Implementação e testes compostos passam; validação estruturada no commit, reconciliação de artefatos e revisão independente ainda condicionam entrega.

### Requirement trace
- R1: TestPublicationCoverageVisibilityAndPrivateProof; TestPublicationIntegrityCLI/published-proof — inventário 2, publicado/elegível 1; variantes e protótipos excluídos.
- R2: TestDistributionProductionPolicyAndUnexpectedTypes; TestTypeDistributionIncludesUnexpectedKindsDeterministically; CLI real --distribution.
- R3: TestPublicationIntegrityOverRealStdio — busca/recomendação/relações sem rascunho, get/start recusados, autoria explícita; publicação inválida impede startup.
- R4: TestEditorialRealGoTestProof; TestEditorialOtherRunners; TestEditorialExpectedBaselineAndReference — baseline/ref real, fail/compile/skip/no-match/timeout, gabarito privado.
- R5: TestEditorialExpectedBaselineAndReference — falta de fixture/prova, alternativa sem justificativa, referência falhando; TestPublicationIntegrityCLI verifica 1/1 publicado; CI executa --published-checks.
- R6: TestPublicationSchemaCorpus — documentos compartilhados; TestLoaderRejectsUnknownManifestAndSymlinkEscape — manifesto desconhecido, versão futura, confinamento.
- R7: TestPublicationIntegrityOverRealStdio retoma sessão de autoria com catálogo público; suites de recovery/tree/evidence legadas passam em -race sem reescrever eventos.
- R8: TestPublicationIntegrityCLI/numerically-sufficient-drafts — inventário supera todos os mínimos, validação incremental passa, V1 falha com zero elegíveis; com publicação e distribuição suficientes mas zero checks, falha isolada v1_check_proof_incomplete.

### Known gaps
O catálogo real ainda não tem conteúdo publicado. Autoria usa serve --authoring; dados sintéticos de teste não constituem revisão/playtest do catálogo real.

## 7. Final Report

### Delivered scope
Implementados os caminhos R1–R8 e verificados por testes locais. Encerramento depende da validação estruturada e atestação independente no candidato.

### Files and modules changed
- Curriculum, checks, CLI e composição MCP; schemas, corpus, política, CI, documentos e ADR declarados em Artifacts.

### Validation executed
- Suites de módulo, CLI/MCP reais, go test -race ./..., pose lint-spec --ready-check e pose check --strict passaram; registro final seguirá a validação estruturada.

### Residual risks
Validar o comportamento implementado em execução independente; não reutilizar resultado histórico como aprovação do código futuro.

### Follow-ups
- [covered: v1-integrated-acceptance] Reexecutar os requisitos desta spec no candidato composto da V1.
